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
)

func (s SyncHelper) GetAllBlogs() (oBlog []entities.OutBlogSimple, oErr error) {
	slog.Info("GetAllBlogs")

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

func (s *SyncHelper) createBlog(inpt BlogInfo) (result FileIDMap, oErr error) {
	slog.Debug("createBlog")

	// load content
	targetFile := path.Join(s.sourcePath, "blogs", inpt.Filename)
	content, err := os.ReadFile(targetFile)
	if err != nil {
		return FileIDMap{}, fmt.Errorf("createBlog: load content failed for blog %q: %w", inpt.Filename, err)
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
		return FileIDMap{}, fmt.Errorf("createBlog: encode body failed for blog %q: %w", inpt.Filename, err)
	}

	// If we somehow lost our database, we will have to create blogs with their original ids
	apiURL := ""
	if inpt.Frontmatter.ID == 0 {
		apiURL, err = url.JoinPath(s.baseURL, "blogs")
		if err != nil {
			return FileIDMap{}, fmt.Errorf("createBlog: join api url failed for blog %q: %w", inpt.Filename, err)
		}
	} else {
		apiURL, err = url.JoinPath(s.baseURL, "blogs", strconv.Itoa(inpt.Frontmatter.ID))
		if err != nil {
			return FileIDMap{}, fmt.Errorf("createBlog: join api url failed for blog %q: %w", inpt.Filename, err)
		}
	}
	slog.Debug("api url", "url", apiURL)

	req, err := http.NewRequest(http.MethodPost, apiURL, body)
	if err != nil {
		return FileIDMap{}, fmt.Errorf("createBlog: new requset failed for blog %q: %w", inpt.Filename, err)
	}
	req.Header.Set("content-type", "application/json")
	req.Header.Set("Authorization", "Bearer "+s.token)

	res, err := httpClient.Do(req)
	if err != nil {
		return FileIDMap{}, fmt.Errorf("createBlog: requset failed for blog %q: %w", inpt.Filename, err)
	}

	defer func() {
		oErr = errors.Join(oErr, drainAndClose(res.Body))
	}()

	// process response and send it through the channel
	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		return FileIDMap{}, fmt.Errorf("createBlog: read response body failed for blog %q: %w", inpt.Filename, err)
	}
	if res.StatusCode >= 400 {
		return FileIDMap{}, fmt.Errorf("createBlog: status code %d for blog %q: %s", res.StatusCode, inpt.Filename, string(resBody))
	}
	resData := entities.RetSuccess[entities.OutBlog]{}
	if err := json.Unmarshal(resBody, &resData); err != nil {
		return FileIDMap{}, fmt.Errorf("createBlog: parse response body failed for blog %q: %w", inpt.Filename, err)
	}

	return NewFileIDMap(inpt.Filename, resData.Msg.ID), nil
}

// return a mapping of blog_filename and id
func (s SyncHelper) CreateBlogs(blogs []BlogInfo) (map[string]int, error) {
	slog.Info("CreateBlogs", "count", len(blogs))

	// load blog content by batch
	batchData := make(chan []BlogInfo, 1)
	go batch[BlogInfo](blogs, s.batchSize, batchData)

	result := map[string]int{}

	// seperate into batches
	for currentBatch := range batchData {
		requestErr := make(chan error, 1)
		response := make(chan FileIDMap, 1)
		responseCount := 0

		for _, blog := range currentBatch {
			go func(b BlogInfo) {
				ret, err := s.createBlog(b)
				if err != nil {
					requestErr <- err
					return
				}
				response <- ret
			}(blog)
		}

		// wait for all requests to finish or if an error occurs
		for {
			if responseCount == len(currentBatch) {
				break
			}
			select {
			case err := <-requestErr:
				return map[string]int{}, err
			case res := <-response:
				responseCount++
				result[res.Filename] = res.Id
			}
		}
	}

	slog.Info("created blogs", "count", len(result), "id mapping", result)
	return result, nil
}

func (s SyncHelper) updateBlog(inpt BlogInfo) (oErr error) {
	slog.Debug("updateBlog")

	// load content
	targetFile := path.Join(s.sourcePath, "blogs", inpt.Filename)
	content, err := os.ReadFile(targetFile)
	if err != nil {
		return fmt.Errorf("updateBlog: load content failed for blog %q: %w", inpt.Filename, err)
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
		return fmt.Errorf("updateBlog: encode body failed for blog %q: %w", inpt.Filename, err)
	}

	if inpt.Frontmatter.ID == 0 {
		return fmt.Errorf("updateBlog: blog id shouldn't be '0', blog %q", inpt.Filename)
	}

	apiURL, err := url.JoinPath(s.baseURL, "blogs", strconv.Itoa(inpt.Frontmatter.ID))
	if err != nil {
		return fmt.Errorf("updateBlog: join api url failed for blog %q: %w", inpt.Filename, err)
	}
	slog.Debug("api url", "url", apiURL)

	req, err := http.NewRequest(http.MethodPut, apiURL, body)
	if err != nil {
		return fmt.Errorf("updateBlog: new requset failed for blog %q: %w", inpt.Filename, err)
	}
	req.Header.Set("content-type", "application/json")
	req.Header.Set("Authorization", "Bearer "+s.token)

	res, err := httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("updateBlog: requset failed for blog %q: %w", inpt.Filename, err)
	}

	defer func() {
		oErr = errors.Join(oErr, drainAndClose(res.Body))
	}()

	// process response and send it through the channel
	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		return fmt.Errorf("updateBlog: read response body failed for blog %q: %w", inpt.Filename, err)
	}
	if res.StatusCode >= 400 {
		return fmt.Errorf("updateBlog: status code %d for blog %q: %s", res.StatusCode, inpt.Filename, string(resBody))
	}

	// just to make sure the response is what we expect
	resData := entities.RetSuccess[entities.OutBlog]{}
	if err := json.Unmarshal(resBody, &resData); err != nil {
		return fmt.Errorf("updateBlog: parse response body failed for blog %q: %w", inpt.Filename, err)
	}

	return nil
}

func (s SyncHelper) UpdateBlogs(blogs []BlogInfo) error {
	slog.Info("UpdateBlogs", "count", len(blogs))

	// load blog content by batch
	batchData := make(chan []BlogInfo, 1)
	go batch[BlogInfo](blogs, s.batchSize, batchData)
	totalCount := 0

	// seperate into batches
	for currentBatch := range batchData {
		requestErr := make(chan error, 1)
		finish := make(chan bool, 1)
		finishCount := 0

		for _, blog := range currentBatch {
			go func(b BlogInfo) {
				if err := s.updateBlog(b); err != nil {
					requestErr <- err
					return
				}
				finish <- true
			}(blog)
		}

		// wait for all requests to finish or if an error occurs
		for {
			if finishCount == len(currentBatch) {
				break
			}
			select {
			case err := <-requestErr:
				return err
			case <-finish:
				finishCount++
				totalCount++
			}
		}
	}

	slog.Info("updated blogs", "count", totalCount)
	return nil
}

func (s SyncHelper) deleteBlog(b entities.OutBlogSimple) (oErr error) {
	slog.Debug("deleteBlog")
	if b.ID == 0 {
		return fmt.Errorf("deleteBlog: blog id shouldn't be '0', blog %q", b.Slug)
	}

	apiURL, err := url.JoinPath(s.baseURL, "blogs", "delete-now", strconv.Itoa(b.ID))
	if err != nil {
		return fmt.Errorf("deleteBlog: join api url failed for blog (id: %d) %q: %w", b.ID, b.Slug, err)
	}
	slog.Debug("api url", "url", apiURL)

	req, err := http.NewRequest(http.MethodDelete, apiURL, nil)
	if err != nil {
		return fmt.Errorf("deleteBlog: new requset failed for blog (id: %d) %q: %w", b.ID, b.Slug, err)
	}
	req.Header.Set("content-type", "application/json")
	req.Header.Set("Authorization", "Bearer "+s.token)

	res, err := httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("deleteBlog: requset failed for blog (id: %d) %q: %w", b.ID, b.Slug, err)
	}

	defer func() {
		oErr = errors.Join(oErr, drainAndClose(res.Body))
	}()

	// process response and send it through the channel
	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		return fmt.Errorf("deleteBlog: read response body failed for blog (id: %d) %q: %w", b.ID, b.Slug, err)
	}
	if res.StatusCode >= 400 {
		return fmt.Errorf("deleteBlog: status code %d for blog (id: %d) %q: %s", res.StatusCode, b.ID, b.Slug, string(resBody))
	}

	// just to make sure the response is what we expect
	resData := entities.RetSuccess[entities.RowsAffected]{}
	if err := json.Unmarshal(resBody, &resData); err != nil {
		return fmt.Errorf("deleteBlog: parse response body failed for blog (id: %d) %q: %w", b.ID, b.Slug, err)
	}

	if resData.Msg.AffectedRows != 1 {
		return fmt.Errorf("deleteBlog: should only delete one blog (id: %d) %q", b.ID, b.Slug)
	}

	return nil
}

func (s SyncHelper) DeleteBlogs(blogs []entities.OutBlogSimple) error {
	slog.Info("DeleteBlogs", "count", len(blogs))

	// load blog content by batch
	batchData := make(chan []entities.OutBlogSimple, 1)
	go batch[entities.OutBlogSimple](blogs, s.batchSize, batchData)
	totalCount := 0

	// seperate into batches
	for currentBatch := range batchData {
		requestErr := make(chan error, 1)
		finish := make(chan bool, 1)
		finishCount := 0

		for _, blog := range currentBatch {
			go func(b entities.OutBlogSimple) {
				if err := s.deleteBlog(b); err != nil {
					requestErr <- err
					return
				}
				finish <- true
			}(blog)
		}

		// wait for all requests to finish or if an error occurs
		for {
			if finishCount == len(currentBatch) {
				break
			}
			select {
			case err := <-requestErr:
				return err
			case <-finish:
				finishCount++
				totalCount++
			}
		}
	}

	slog.Info("deleted blogs", "count", totalCount)
	return nil
}

// update ids.json (blog filename to id mapping)
func updateIDMapping(blogs []BlogInfo, newBlogIDs map[string]int, targetFile string) error {
	slog.Info("updateIDMapping")

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

	slog.Info("finish updating ids.json", "total blogs", len(newMapping))

	return nil
}
