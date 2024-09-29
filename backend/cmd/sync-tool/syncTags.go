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
	"strconv"
)

func (s SyncHelper) GetAllTags() (oTag []entities.Tag, oErr error) {
	slog.Info("GetAllTags")

	apiURL, err := url.JoinPath(s.baseURL, "tags")
	if err != nil {
		return []entities.Tag{}, fmt.Errorf("GetAllTags: join api url failed: %w", err)
	}
	slog.Debug("api url", "url", apiURL)

	res, err := httpClient.Get(apiURL)
	if err != nil {
		return []entities.Tag{}, fmt.Errorf("GetAllTags: get failed: %w", err)
	}

	// cleanup
	defer func() {
		oErr = errors.Join(oErr, drainAndClose(res.Body))
	}()

	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		return []entities.Tag{}, fmt.Errorf("GetAllTags: read body failed: %w", err)
	}

	if res.StatusCode >= 400 {
		return []entities.Tag{}, fmt.Errorf("GetAllTags: status code %d, msg: %s", res.StatusCode, string(resBody))
	}

	data := entities.RetSuccess[[]entities.Tag]{}
	if err := json.Unmarshal(resBody, &data); err != nil {
		return []entities.Tag{}, fmt.Errorf("GetAllTags: decode body failed: %w", err)
	}

	slog.Debug("got tags", "tags", data.Msg)

	return data.Msg, nil
}

func (s SyncHelper) createTag(t entities.Tag) (result entities.Tag, oErr error) {
	slog.Debug("createTag")

	// prepare request body
	body := &bytes.Buffer{}
	data := entities.NewInTag(t.Name, t.Description)
	if err := json.NewEncoder(body).Encode(data); err != nil {
		return entities.Tag{}, fmt.Errorf("createTag: encode body failed for tag %q: %w", t.Name, err)
	}

	apiURL, err := url.JoinPath(s.baseURL, "tags")
	if err != nil {
		return entities.Tag{}, fmt.Errorf("createTag: join api url failed for tag %q: %w", t.Name, err)
	}
	slog.Debug("api url", "url", apiURL)

	req, err := http.NewRequest(http.MethodPost, apiURL, body)
	if err != nil {
		return entities.Tag{}, fmt.Errorf("createTag: new requset failed for tag %q: %w", t.Name, err)
	}
	req.Header.Set("content-type", "application/json")
	req.Header.Set("Authorization", "Bearer "+s.token)

	res, err := httpClient.Do(req)
	if err != nil {
		return entities.Tag{}, fmt.Errorf("createTag: requset failed for tag %q: %w", t.Name, err)
	}

	defer func() {
		oErr = errors.Join(oErr, drainAndClose(res.Body))
	}()

	// process response and send it through the channel
	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		return entities.Tag{}, fmt.Errorf("createTag: read response body failed for tag %q: %w", t.Name, err)
	}
	if res.StatusCode >= 400 {
		return entities.Tag{}, fmt.Errorf("createTag: status code %d for tag %q: %s", res.StatusCode, t.Name, string(resBody))
	}
	resData := entities.RetSuccess[entities.Tag]{}
	if err := json.Unmarshal(resBody, &resData); err != nil {
		return entities.Tag{}, fmt.Errorf("createTag: parse response body failed for tag %q: %w", t.Name, err)
	}

	return resData.Msg, nil
}

func (s SyncHelper) CreateTags(tags []entities.Tag) ([]entities.Tag, error) {
	slog.Info("CreateTags", "count", len(tags))

	batchData := make(chan []entities.Tag, 1)
	go batch(tags, s.batchSize, batchData)

	result := []entities.Tag{}

	// seperate into batches
	for currentBatch := range batchData {
		requestErr := make(chan error, 1)
		response := make(chan entities.Tag, 1)
		responseCount := 0

		// there is no 'bulk api' for now, so we just create tags one by one
		for _, tag := range currentBatch {
			go func(t entities.Tag) {
				res, err := s.createTag(t)
				if err != nil {
					requestErr <- err
					return
				}
				response <- res
			}(tag)
		}

		// wait for all requests to finish or if an error occurs
		for {
			if responseCount == len(currentBatch) {
				break
			}
			select {
			case newTag := <-response:
				result = append(result, newTag)
				responseCount++
			case err := <-requestErr:
				return []entities.Tag{}, err
			}
		}
	}

	slog.Info("created tags", "count", len(result))
	return result, nil
}

func (s SyncHelper) updateTag(t entities.Tag) (result entities.Tag, oErr error) {
	slog.Debug("updateTag")

	// prepare request body
	body := &bytes.Buffer{}
	data := entities.NewInTag(t.Name, t.Description)
	if err := json.NewEncoder(body).Encode(data); err != nil {
		return entities.Tag{}, fmt.Errorf("updateTag: encode body failed for tag %q: %w", t.Name, err)
	}

	apiURL, err := url.JoinPath(s.baseURL, "tags", strconv.Itoa(t.ID))
	if err != nil {
		return entities.Tag{}, fmt.Errorf("updateTag: join api url failed for tag %q: %w", t.Name, err)
	}
	slog.Debug("api url", "url", apiURL)

	req, err := http.NewRequest(http.MethodPut, apiURL, body)
	if err != nil {
		return entities.Tag{}, fmt.Errorf("updateTag: new request failed for tag %q: %w", t.Name, err)
	}
	req.Header.Set("content-type", "application/json")
	req.Header.Set("Authorization", "Bearer "+s.token)

	res, err := httpClient.Do(req)
	if err != nil {
		return entities.Tag{}, fmt.Errorf("updateTag: requset failed for tag %q: %w", t.Name, err)
	}

	defer func() {
		oErr = errors.Join(oErr, drainAndClose(res.Body))
	}()

	// process response and send it through the channel
	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		return entities.Tag{}, fmt.Errorf("updateTag: read response body failed for tag %q: %w", t.Name, err)
	}
	if res.StatusCode >= 400 {
		return entities.Tag{}, fmt.Errorf("updateTag: status code %d for tag %q: %s", res.StatusCode, t.Name, string(resBody))
	}
	resData := entities.RetSuccess[entities.Tag]{}
	if err := json.Unmarshal(resBody, &resData); err != nil {
		return entities.Tag{}, fmt.Errorf("updateTag: parse response body failed for tag %q: %w", t.Name, err)
	}

	return resData.Msg, nil
}

func (s SyncHelper) UpdateTags(tags []entities.Tag) ([]entities.Tag, error) {
	slog.Info("UpdateTags", "count", len(tags))

	batchData := make(chan []entities.Tag, 1)
	go batch(tags, s.batchSize, batchData)

	result := []entities.Tag{}

	// seperate into batches
	for currentBatch := range batchData {
		requestErr := make(chan error, 1)
		response := make(chan entities.Tag, 1)
		responseCount := 0

		// there is no 'bulk api' for now, so we just create tags one by one
		for _, tag := range currentBatch {
			go func(t entities.Tag) {
				res, err := s.updateTag(t)
				if err != nil {
					requestErr <- err
					return
				}
				response <- res
			}(tag)
		}

		// wait for all requests to finish or if an error occurs
		for {
			if responseCount == len(currentBatch) {
				break
			}
			select {
			case newTag := <-response:
				result = append(result, newTag)
				responseCount++
			case err := <-requestErr:
				return []entities.Tag{}, err
			}
		}
	}

	slog.Info("updated tags", "count", len(result))
	return result, nil
}

func (s SyncHelper) deleteTag(t entities.Tag) (oErr error) {
	slog.Debug("deleteTag")

	apiURL, err := url.JoinPath(s.baseURL, "tags", strconv.Itoa(t.ID))
	if err != nil {
		return fmt.Errorf("DeleteTags: join api url failed for tag %q: %w", t.Name, err)
	}
	slog.Debug("api url", "url", apiURL)

	req, err := http.NewRequest(http.MethodDelete, apiURL, nil)
	if err != nil {
		return fmt.Errorf("DeleteTags: new request failed for tag %q: %w", t.Name, err)
	}
	req.Header.Set("content-type", "application/json")
	req.Header.Set("Authorization", "Bearer "+s.token)

	res, err := httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("DeleteTags: requset failed for tag %q: %w", t.Name, err)
	}

	defer func() {
		oErr = errors.Join(oErr, drainAndClose(res.Body))
	}()

	// process response and send it through the channel
	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		return fmt.Errorf("DeleteTags: read response body failed for tag %q: %w", t.Name, err)
	}
	if res.StatusCode >= 400 {
		return fmt.Errorf("DeleteTags: status code %d for tag %q: %s", res.StatusCode, t.Name, string(resBody))
	}
	resData := entities.RetSuccess[entities.RowsAffected]{}
	if err := json.Unmarshal(resBody, &resData); err != nil {
		return fmt.Errorf("DeleteTags: parse response body failed for tag %q: %w", t.Name, err)
	}

	if resData.Msg.AffectedRows != 1 {
		return fmt.Errorf("DeleteTags: should only delete one tag %q", t.Name)
	}
	return nil
}

func (s SyncHelper) DeleteTags(tags []entities.Tag) error {
	slog.Info("DeleteTags", "count", len(tags))

	batchData := make(chan []entities.Tag, 1)
	go batch(tags, s.batchSize, batchData)

	totalCount := 0

	// seperate into batches
	for currentBatch := range batchData {
		requestErr := make(chan error, 1)
		finish := make(chan bool, 1)
		finishCount := 0
		// there is no 'bulk api' for now, so we just create tags one by one
		for _, tag := range currentBatch {
			go func(t entities.Tag) {
				if err := s.deleteTag(t); err != nil {
					requestErr <- err
					return
				}
				finish <- true
			}(tag)
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

	slog.Info("deleted tags", "count", totalCount)
	return nil
}
