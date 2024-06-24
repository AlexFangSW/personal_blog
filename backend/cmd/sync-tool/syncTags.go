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
	slog.Debug("GetAllTags")

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
func (s SyncHelper) CreateTags(tags []entities.Tag) ([]entities.Tag, error) {
	slog.Debug("CreateTags")

	batchData := make(chan []entities.Tag, 1)
	go batch[entities.Tag](tags, s.batchSize, batchData)

	requestErr := make(chan error, 1)
	successResponse := make(chan entities.Tag, len(tags))

	// seperate into batches
	for currentBatch := range batchData {
		// there is no 'bulk api' for now, so we just create tags one by one
		for _, tag := range currentBatch {
			go func(t entities.Tag) {
				var oErr error
				defer func() {
					if oErr != nil {
						requestErr <- oErr
					}
				}()

				// prepare request body
				body := &bytes.Buffer{}
				data := entities.NewInTag(t.Name, t.Description)
				if err := json.NewEncoder(body).Encode(data); err != nil {
					oErr = errors.Join(oErr, fmt.Errorf("CreateTags: encode body failed for tag %q: %w", t.Name, err))
					return
				}

				apiURL, err := url.JoinPath(s.baseURL, "tags")
				if err != nil {
					oErr = errors.Join(oErr, fmt.Errorf("CreateTags: join api url failed for tag %q: %w", t.Name, err))
					return
				}
				slog.Debug("api url", "url", apiURL)

				req, err := http.NewRequest(http.MethodPost, apiURL, body)
				if err != nil {
					oErr = errors.Join(oErr, fmt.Errorf("CreateTags: new requset failed for tag %q: %w", t.Name, err))
					return
				}
				req.Header.Set("content-type", "application/json")
				req.Header.Set("Authorization", "Bearer "+s.token)

				res, err := httpClient.Do(req)
				if err != nil {
					oErr = errors.Join(oErr, fmt.Errorf("CreateTags: requset failed for tag %q: %w", t.Name, err))
					return
				}

				defer func() {
					oErr = errors.Join(oErr, drainAndClose(res.Body))
				}()

				// process response and send it through the channel
				resBody, err := io.ReadAll(res.Body)
				if err != nil {
					oErr = errors.Join(oErr, fmt.Errorf("CreateTags: read response body failed for tag %q: %w", t.Name, err))
					return
				}
				if res.StatusCode >= 400 {
					oErr = errors.Join(oErr, fmt.Errorf("CreateTags: status code %d for tag %q: %s", res.StatusCode, t.Name, string(resBody)))
					return
				}
				resData := entities.RetSuccess[entities.Tag]{}
				if err := json.Unmarshal(resBody, &resData); err != nil {
					oErr = errors.Join(oErr, fmt.Errorf("CreateTags: parse response body failed for tag %q: %w", t.Name, err))
					return
				}

				successResponse <- resData.Msg

			}(tag)
		}
	}

	// wait for all requests to finish or if an error occurs
	result := []entities.Tag{}
	for {
		if len(result) == len(tags) {
			slog.Debug("new tags", "count", len(result))
			return result, nil
		}
		select {
		case newTag := <-successResponse:
			result = append(result, newTag)
		case err := <-requestErr:
			return []entities.Tag{}, err
		}
	}

}
func (s SyncHelper) UpdateTags(tags []entities.Tag) ([]entities.Tag, error) {
	slog.Debug("UpdateTags")

	batchData := make(chan []entities.Tag, 1)
	go batch[entities.Tag](tags, s.batchSize, batchData)

	requestErr := make(chan error, 1)
	successResponse := make(chan entities.Tag, len(tags))

	// seperate into batches
	for currentBatch := range batchData {
		// there is no 'bulk api' for now, so we just create tags one by one
		for _, tag := range currentBatch {
			go func(t entities.Tag) {
				var oErr error
				defer func() {
					if oErr != nil {
						requestErr <- oErr
					}
				}()

				// prepare request body
				body := &bytes.Buffer{}
				data := entities.NewInTag(t.Name, t.Description)
				if err := json.NewEncoder(body).Encode(data); err != nil {
					oErr = errors.Join(oErr, fmt.Errorf("UpdateTags: encode body failed for tag %q: %w", t.Name, err))
					return
				}

				apiURL, err := url.JoinPath(s.baseURL, "tags", strconv.Itoa(t.ID))
				if err != nil {
					oErr = errors.Join(oErr, fmt.Errorf("UpdateTags: join api url failed for tag %q: %w", t.Name, err))
					return
				}
				slog.Debug("api url", "url", apiURL)

				req, err := http.NewRequest(http.MethodPut, apiURL, body)
				if err != nil {
					oErr = errors.Join(oErr, fmt.Errorf("UpdateTags: new request failed for tag %q: %w", t.Name, err))
					return
				}
				req.Header.Set("content-type", "application/json")
				req.Header.Set("Authorization", "Bearer "+s.token)

				res, err := httpClient.Do(req)
				if err != nil {
					oErr = errors.Join(oErr, fmt.Errorf("UpdateTags: requset failed for tag %q: %w", t.Name, err))
					return
				}

				defer func() {
					oErr = errors.Join(oErr, drainAndClose(res.Body))
				}()

				// process response and send it through the channel
				resBody, err := io.ReadAll(res.Body)
				if err != nil {
					oErr = errors.Join(oErr, fmt.Errorf("UpdateTags: read response body failed for tag %q: %w", t.Name, err))
					return
				}
				if res.StatusCode >= 400 {
					oErr = errors.Join(oErr, fmt.Errorf("UpdateTags: status code %d for tag %q: %s", res.StatusCode, t.Name, string(resBody)))
					return
				}
				resData := entities.RetSuccess[entities.Tag]{}
				if err := json.Unmarshal(resBody, &resData); err != nil {
					oErr = errors.Join(oErr, fmt.Errorf("UpdateTags: parse response body failed for tag %q: %w", t.Name, err))
					return
				}

				successResponse <- resData.Msg

			}(tag)
		}
	}

	// wait for all requests to finish or if an error occurs
	result := []entities.Tag{}
	for {
		if len(result) == len(tags) {
			slog.Debug("updated tags", "count", len(result))
			return result, nil
		}
		select {
		case newTag := <-successResponse:
			result = append(result, newTag)
		case err := <-requestErr:
			return []entities.Tag{}, err
		}
	}
}
func (s SyncHelper) DeleteTags(tags []entities.Tag) error {
	slog.Debug("DeleteTags")

	batchData := make(chan []entities.Tag, 1)
	go batch[entities.Tag](tags, s.batchSize, batchData)

	requestErr := make(chan error, 1)
	successResponse := make(chan bool, len(tags))

	// seperate into batches
	for currentBatch := range batchData {
		// there is no 'bulk api' for now, so we just create tags one by one
		for _, tag := range currentBatch {
			go func(t entities.Tag) {
				var oErr error
				defer func() {
					if oErr != nil {
						requestErr <- oErr
					}
				}()

				apiURL, err := url.JoinPath(s.baseURL, "tags", strconv.Itoa(t.ID))
				if err != nil {
					oErr = errors.Join(oErr, fmt.Errorf("DeleteTags: join api url failed for tag %q: %w", t.Name, err))
					return
				}
				slog.Debug("api url", "url", apiURL)

				req, err := http.NewRequest(http.MethodDelete, apiURL, nil)
				if err != nil {
					oErr = errors.Join(oErr, fmt.Errorf("DeleteTags: new request failed for tag %q: %w", t.Name, err))
					return
				}
				req.Header.Set("content-type", "application/json")
				req.Header.Set("Authorization", "Bearer "+s.token)

				res, err := httpClient.Do(req)
				if err != nil {
					oErr = errors.Join(oErr, fmt.Errorf("DeleteTags: requset failed for tag %q: %w", t.Name, err))
					return
				}

				defer func() {
					oErr = errors.Join(oErr, drainAndClose(res.Body))
				}()

				// process response and send it through the channel
				resBody, err := io.ReadAll(res.Body)
				if err != nil {
					oErr = errors.Join(oErr, fmt.Errorf("DeleteTags: read response body failed for tag %q: %w", t.Name, err))
					return
				}
				if res.StatusCode >= 400 {
					oErr = errors.Join(oErr, fmt.Errorf("DeleteTags: status code %d for tag %q: %s", res.StatusCode, t.Name, string(resBody)))
					return
				}
				resData := entities.RetSuccess[entities.RowsAffected]{}
				if err := json.Unmarshal(resBody, &resData); err != nil {
					oErr = errors.Join(oErr, fmt.Errorf("DeleteTags: parse response body failed for tag %q: %w", t.Name, err))
					return
				}

				if resData.Msg.AffectedRows != 1 {
					oErr = errors.Join(oErr, fmt.Errorf("DeleteTags: should only delete one tag %q", t.Name))
					return
				}
				successResponse <- true

			}(tag)
		}
	}

	// wait for all requests to finish or if an error occurs
	successCount := 0
	for {
		if successCount == len(tags) {
			slog.Debug("deleted tags", "count", successCount)
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
