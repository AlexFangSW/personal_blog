package main

import (
	"blog/entities"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"os"
	"path"
	"strconv"
	"strings"

	"golang.org/x/exp/maps"
)

func (s SyncHelper) GetAllBlogs() (oBlog []entities.OutBlogSimple, oErr error) {
	slog.Debug("GetAllBlogs")

	apiURL, err := url.JoinPath(s.baseURL, "blogs")
	if err != nil {
		return []entities.OutBlogSimple{}, fmt.Errorf("GetAllBlogs: join api url failed: %w", err)
	}
	slog.Debug("api url", "url", apiURL)

	req, err := http.NewRequest(http.MethodGet, apiURL, nil)
	if err != nil {
		return []entities.OutBlogSimple{}, fmt.Errorf("GetAllBlogs: create new request failed: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+s.token)
	query := req.URL.Query()
	query.Set("all", "true")
	query.Set("simple", "true")
	req.URL.RawQuery = query.Encode()

	res, err := httpClient.Do(req)
	if err != nil {
		return []entities.OutBlogSimple{}, fmt.Errorf("GetAllBlogs: req failed: %w", err)
	}

	defer func() {
		oErr = errors.Join(oErr, drainAndClose(res.Body))
	}()

	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		return []entities.OutBlogSimple{}, fmt.Errorf("GetAllBlogs: read body failed: %w", err)
	}

	if res.StatusCode >= 400 {
		return []entities.OutBlogSimple{}, fmt.Errorf("GetAllBlogs: status code %d, msg: %s", res.StatusCode, string(resBody))
	}

	data := entities.RetSuccess[[]entities.OutBlogSimple]{}
	if err := json.Unmarshal(resBody, &data); err != nil {
		return []entities.OutBlogSimple{}, fmt.Errorf("GetAllBlogs: unmarshal failed: %w", err)
	}

	slog.Debug("got blogs", "blogs", data.Msg)
	return data.Msg, nil
}

type FileIDMap struct {
	Filename string
	Id       int
}

func NewFileIDMap(filename string, id int) FileIDMap {
	return FileIDMap{
		Filename: filename,
		Id:       id,
	}
}

func (s *SyncHelper) createBlogs(inpt BlogInfo) (result FileIDMap, oErr error) {
	// load content
	targetFile := path.Join(s.sourcePath, "blogs", inpt.Filename)
	content, err := os.ReadFile(targetFile)
	if err != nil {
		return FileIDMap{}, fmt.Errorf("createBlogs: load content failed for blog %q: %w", inpt.Filename, err)
	}

	// prepare request body
	newBlog := entities.NewBlog(
		inpt.Frontmatter.Title,
		strings.Split(string(content), "---")[2],
		inpt.Frontmatter.Description,
		inpt.Frontmatter.Pined,
		inpt.Frontmatter.Visible,
	)
	newInBlog := entities.NewInBlog(
		*newBlog,
		inpt.Frontmatter.TagIDs,
		inpt.Frontmatter.TopicIDs,
	)

	body := &bytes.Buffer{}
	if err := json.NewEncoder(body).Encode(newInBlog); err != nil {
		return FileIDMap{}, fmt.Errorf("createBlogs: encode body failed for blog %q: %w", inpt.Filename, err)
	}

	// If we somehow lost our database, we will have to create blogs with their original ids
	apiURL := ""
	if inpt.Frontmatter.ID == 0 {
		apiURL, err = url.JoinPath(s.baseURL, "blogs")
		if err != nil {
			return FileIDMap{}, fmt.Errorf("createBlogs: join api url failed for blog %q: %w", inpt.Filename, err)
		}
	} else {
		apiURL, err = url.JoinPath(s.baseURL, "blogs", strconv.Itoa(inpt.Frontmatter.ID))
		if err != nil {
			return FileIDMap{}, fmt.Errorf("createBlogs: join api url failed for blog %q: %w", inpt.Filename, err)
		}
	}
	slog.Debug("api url", "url", apiURL)

	req, err := http.NewRequest(http.MethodPost, apiURL, body)
	if err != nil {
		return FileIDMap{}, fmt.Errorf("createBlogs: new requset failed for blog %q: %w", inpt.Filename, err)
	}
	req.Header.Set("content-type", "application/json")
	req.Header.Set("Authorization", "Bearer "+s.token)

	res, err := httpClient.Do(req)
	if err != nil {
		return FileIDMap{}, fmt.Errorf("createBlogs: requset failed for blog %q: %w", inpt.Filename, err)
	}

	defer func() {
		oErr = errors.Join(oErr, drainAndClose(res.Body))
	}()

	// process response and send it through the channel
	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		return FileIDMap{}, fmt.Errorf("createBlogs: read response body failed for blog %q: %w", inpt.Filename, err)
	}
	if res.StatusCode >= 400 {
		return FileIDMap{}, fmt.Errorf("createBlogs: status code %d for blog %q: %s", res.StatusCode, inpt.Filename, string(resBody))
	}
	resData := entities.RetSuccess[entities.OutBlog]{}
	if err := json.Unmarshal(resBody, &resData); err != nil {
		return FileIDMap{}, fmt.Errorf("createBlogs: parse response body failed for blog %q: %w", inpt.Filename, err)
	}

	return NewFileIDMap(inpt.Filename, resData.Msg.ID), nil
}

// return a mapping of blog_filename and id
func (s SyncHelper) CreateBlogs(blogs []BlogInfo) (map[string]int, error) {
	slog.Debug("CreateBlogs")

	// load blog content by batch
	batchData := make(chan []BlogInfo, 1)
	go batch[BlogInfo](blogs, s.batchSize, batchData)

	result := map[string]int{}

	// seperate into batches
	for currentBatch := range batchData {
		requestErr := make(chan error, 1)
		response := make(chan FileIDMap, len(currentBatch))
		responseBuffer := map[string]int{}

		for _, blog := range currentBatch {
			go func(b BlogInfo) {
				ret, err := s.createBlogs(b)
				if err != nil {
					requestErr <- err
					return
				}
				response <- ret
			}(blog)
		}

		// wait for all requests to finish or if an error occurs
		for {
			if len(responseBuffer) == len(currentBatch) {
				maps.Copy[map[string]int](result, responseBuffer)
				break
			}
			select {
			case err := <-requestErr:
				return map[string]int{}, err
			case res := <-response:
				responseBuffer[res.Filename] = res.Id
			}
		}
	}

	slog.Debug("created blogs", "count", len(blogs), "id mapping", result)
	return result, nil
}

func (s SyncHelper) updateBlogs(inpt BlogInfo) (oErr error) {

	// load content
	targetFile := path.Join(s.sourcePath, "blogs", inpt.Filename)
	content, err := os.ReadFile(targetFile)
	if err != nil {
		return fmt.Errorf("updateBlogs: load content failed for blog %q: %w", inpt.Filename, err)
	}

	// prepare request body
	newBlog := entities.NewBlog(
		inpt.Frontmatter.Title,
		strings.Split(string(content), "---")[2],
		inpt.Frontmatter.Description,
		inpt.Frontmatter.Pined,
		inpt.Frontmatter.Visible,
	)
	newInBlog := entities.NewInBlog(
		*newBlog,
		inpt.Frontmatter.TagIDs,
		inpt.Frontmatter.TopicIDs,
	)

	body := &bytes.Buffer{}
	if err := json.NewEncoder(body).Encode(newInBlog); err != nil {
		return fmt.Errorf("updateBlogs: encode body failed for blog %q: %w", inpt.Filename, err)
	}

	if inpt.Frontmatter.ID == 0 {
		return fmt.Errorf("updateBlogs: blog id shouldn't be '0', blog %q", inpt.Filename)
	}

	apiURL, err := url.JoinPath(s.baseURL, "blogs", strconv.Itoa(inpt.Frontmatter.ID))
	if err != nil {
		return fmt.Errorf("updateBlogs: join api url failed for blog %q: %w", inpt.Filename, err)
	}
	slog.Debug("api url", "url", apiURL)

	req, err := http.NewRequest(http.MethodPut, apiURL, body)
	if err != nil {
		return fmt.Errorf("updateBlogs: new requset failed for blog %q: %w", inpt.Filename, err)
	}
	req.Header.Set("content-type", "application/json")
	req.Header.Set("Authorization", "Bearer "+s.token)

	res, err := httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("updateBlogs: requset failed for blog %q: %w", inpt.Filename, err)
	}

	defer func() {
		oErr = errors.Join(oErr, drainAndClose(res.Body))
	}()

	// process response and send it through the channel
	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		return fmt.Errorf("updateBlogs: read response body failed for blog %q: %w", inpt.Filename, err)
	}
	if res.StatusCode >= 400 {
		return fmt.Errorf("updateBlogs: status code %d for blog %q: %s", res.StatusCode, inpt.Filename, string(resBody))
	}

	// just to make sure the response is what we expect
	resData := entities.RetSuccess[entities.OutBlog]{}
	if err := json.Unmarshal(resBody, &resData); err != nil {
		return fmt.Errorf("updateBlogs: parse response body failed for blog %q: %w", inpt.Filename, err)
	}

	return nil
}

func (s SyncHelper) UpdateBlogs(blogs []BlogInfo) error {
	slog.Debug("UpdateBlogs")

	// load blog content by batch
	batchData := make(chan []BlogInfo, 1)
	go batch[BlogInfo](blogs, s.batchSize, batchData)

	// seperate into batches
	for currentBatch := range batchData {
		requestErr := make(chan error, 1)
		finish := make(chan bool, 1)
		finishCount := 0

		for _, blog := range currentBatch {
			go func(b BlogInfo) {
				if err := s.updateBlogs(b); err != nil {
					requestErr <- err
					return
				}
				finish <- true
			}(blog)
		}

		// wait for all requests to finish or if an error occurs
		for {
			if finishCount == len(blogs) {
				return nil
			}
			select {
			case err := <-requestErr:
				return err
			case <-finish:
				finishCount++
			}
		}
	}
	return nil
}

func (s SyncHelper) DeleteBlogs(blogs []entities.OutBlogSimple) error {
	slog.Debug("DeleteBlogs")

	// load blog content by batch
	batchData := make(chan []entities.OutBlogSimple, 1)
	go batch[entities.OutBlogSimple](blogs, s.batchSize, batchData)

	requestErr := make(chan error, 1)
	successResponse := make(chan bool, len(blogs))

	// seperate into batches
	for currentBatch := range batchData {
		for _, blog := range currentBatch {
			go func(b entities.OutBlogSimple) {
				var oErr error
				defer func() {
					if oErr != nil {
						requestErr <- oErr
					}
				}()

				if b.ID == 0 {
					oErr = errors.Join(oErr, fmt.Errorf("DeleteBlogs: blog id shouldn't be '0', blog %q", b.Slug))
					return
				}

				apiURL, err := url.JoinPath(s.baseURL, "blogs", "delete-now", strconv.Itoa(b.ID))
				if err != nil {
					oErr = errors.Join(oErr, fmt.Errorf("DeleteBlogs: join api url failed for blog (id: %d) %q: %w", b.ID, b.Slug, err))
					return
				}
				slog.Debug("api url", "url", apiURL)

				req, err := http.NewRequest(http.MethodDelete, apiURL, nil)
				if err != nil {
					oErr = errors.Join(oErr, fmt.Errorf("DeleteBlogs: new requset failed for blog (id: %d) %q: %w", b.ID, b.Slug, err))
					return
				}
				req.Header.Set("content-type", "application/json")
				req.Header.Set("Authorization", "Bearer "+s.token)

				res, err := httpClient.Do(req)
				if err != nil {
					oErr = errors.Join(oErr, fmt.Errorf("DeleteBlogs: requset failed for blog (id: %d) %q: %w", b.ID, b.Slug, err))
					return
				}

				defer func() {
					oErr = errors.Join(oErr, drainAndClose(res.Body))
				}()

				// process response and send it through the channel
				resBody, err := io.ReadAll(res.Body)
				if err != nil {
					oErr = errors.Join(oErr, fmt.Errorf("DeleteBlogs: read response body failed for blog (id: %d) %q: %w", b.ID, b.Slug, err))
					return
				}
				if res.StatusCode >= 400 {
					oErr = errors.Join(oErr, fmt.Errorf("DeleteBlogs: status code %d for blog (id: %d) %q: %s", res.StatusCode, b.ID, b.Slug, string(resBody)))
					return
				}

				// just to make sure the response is what we expect
				resData := entities.RetSuccess[entities.RowsAffected]{}
				if err := json.Unmarshal(resBody, &resData); err != nil {
					oErr = errors.Join(oErr, fmt.Errorf("DeleteBlogs: parse response body failed for blog (id: %d) %q: %w", b.ID, b.Slug, err))
					return
				}

				if resData.Msg.AffectedRows != 1 {
					oErr = errors.Join(oErr, fmt.Errorf("DeleteBlogs: should only delete one blog (id: %d) %q", b.ID, b.Slug))
					return
				}
				successResponse <- true

			}(blog)
		}
	}

	// wait for all requests to finish or if an error occurs
	successCount := 0
	for {
		if successCount == len(blogs) {
			slog.Debug("deleted blogs", "count", successCount)
			return nil
		}

		select {
		case <-successResponse:
			successCount++
		case err := <-requestErr:
			return err
		}
	}
}

// update ids.json (blog filename to id mapping)
func updateIDMapping(blogs []BlogInfo, newBlogIDs map[string]int, targetFile string) error {
	slog.Debug("updateIDMapping")

	newMapping := newBlogIDs
	for _, blog := range blogs {
		newMapping[blog.Filename] = blog.Frontmatter.ID
	}
	slog.Debug("new id mapping", "mapping", newMapping)

	data, err := json.Marshal(newMapping)
	if err != nil {
		return fmt.Errorf("updateIDMapping: marshal failed: %w", err)
	}

	if err := os.WriteFile(targetFile, data, 0644); err != nil {
		return fmt.Errorf("updateIDMapping: write file failed: %w", err)
	}

	slog.Debug("finish updating ids.json", "total blogs", len(newMapping))

	return nil
}
