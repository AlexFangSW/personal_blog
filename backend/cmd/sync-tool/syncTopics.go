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

func (s SyncHelper) GetAllTopics() (oTag []entities.Topic, oErr error) {
	slog.Debug("GetAllTopics")

	apiURL, err := url.JoinPath(s.baseURL, "topics")
	if err != nil {
		return []entities.Topic{}, fmt.Errorf("GetAllTopics: join api url failed: %w", err)
	}
	slog.Debug("api url", "url", apiURL)

	res, err := httpClient.Get(apiURL)
	if err != nil {
		return []entities.Topic{}, fmt.Errorf("GetAllTopics: get failed: %w", err)
	}

	// cleanup
	defer func() {
		oErr = errors.Join(oErr, drainAndClose(res.Body))
	}()

	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		return []entities.Topic{}, fmt.Errorf("GetAllTopics: read body failed: %w", err)
	}

	if res.StatusCode >= 400 {
		return []entities.Topic{}, fmt.Errorf("GetAllTopics: status code %d, msg: %s", res.StatusCode, string(resBody))
	}

	data := entities.RetSuccess[[]entities.Topic]{}
	if err := json.Unmarshal(resBody, &data); err != nil {
		return []entities.Topic{}, fmt.Errorf("GetAllTopics: decode body failed: %w", err)
	}

	slog.Debug("got topics", "topics", data.Msg)
	return data.Msg, nil
}

// creates topics and updates ID field in place
func (s SyncHelper) CreateTopics(topics []entities.Topic) ([]entities.Topic, error) {
	slog.Debug("CreateTopics")

	batchData := make(chan []entities.Topic, 1)
	go batch[entities.Topic](topics, s.batchSize, batchData)

	requestErr := make(chan error, 1)
	successResponse := make(chan entities.Topic, len(topics))

	// seperate into batches
	for currentBatch := range batchData {
		// there is no 'bulk api' for now, so we just create topics one by one
		for _, topic := range currentBatch {
			go func(t entities.Topic) {
				var oErr error
				defer func() {
					if oErr != nil {
						requestErr <- oErr
					}
				}()

				// prepare request body
				body := &bytes.Buffer{}
				data := entities.NewInTopic(t.Name, t.Description)
				if err := json.NewEncoder(body).Encode(data); err != nil {
					oErr = errors.Join(oErr, fmt.Errorf("CreateTopics: encode body failed for topic %q: %w", t.Name, err))
					return
				}

				apiURL, err := url.JoinPath(s.baseURL, "topics")
				if err != nil {
					oErr = errors.Join(oErr, fmt.Errorf("CreateTopics: join api url failed for topic %q: %w", t.Name, err))
					return
				}
				slog.Debug("api url", "url", apiURL)

				req, err := http.NewRequest(http.MethodPost, apiURL, body)
				if err != nil {
					oErr = errors.Join(oErr, fmt.Errorf("CreateTopics: new request failed for topic %q: %w", t.Name, err))
					return
				}
				req.Header.Set("content-type", "application/json")
				req.Header.Set("Authorization", "Bearer "+s.token)

				res, err := httpClient.Do(req)
				if err != nil {
					oErr = errors.Join(oErr, fmt.Errorf("CreateTopics: requset failed for topic %q: %w", t.Name, err))
					return
				}

				// cleanup
				defer func() {
					oErr = errors.Join(oErr, drainAndClose(res.Body))
				}()

				// process response and send it through the channel
				resBody, err := io.ReadAll(res.Body)
				if err != nil {
					oErr = errors.Join(oErr, fmt.Errorf("CreateTopics: read response body failed for topic %q: %w", t.Name, err))
					return
				}
				if res.StatusCode >= 400 {
					oErr = errors.Join(oErr, fmt.Errorf("CreateTopics: status code %d for topic %q: %s", res.StatusCode, t.Name, string(resBody)))
					return
				}
				resData := entities.RetSuccess[entities.Topic]{}
				if err := json.Unmarshal(resBody, &resData); err != nil {
					oErr = errors.Join(oErr, fmt.Errorf("CreateTopics: parse response body failed for topic %q: %w", t.Name, err))
					return
				}

				successResponse <- resData.Msg

			}(topic)
		}
	}

	// wait for all requests to finish or if an error occurs
	result := []entities.Topic{}
	for {
		if len(result) == len(topics) {
			slog.Debug("new topics", "count", len(result))
			return result, nil
		}
		select {
		case newTopic := <-successResponse:
			result = append(result, newTopic)
		case err := <-requestErr:
			return []entities.Topic{}, err
		}
	}
}
func (s SyncHelper) UpdateTopics(topics []entities.Topic) ([]entities.Topic, error) {
	slog.Debug("UpdateTopics")

	batchData := make(chan []entities.Topic, 1)
	go batch[entities.Topic](topics, s.batchSize, batchData)

	requestErr := make(chan error, 1)
	successResponse := make(chan entities.Topic, len(topics))

	// seperate into batches
	for currentBatch := range batchData {
		// there is no 'bulk api' for now, so we just create topics one by one
		for _, topic := range currentBatch {
			go func(t entities.Topic) {
				var oErr error
				defer func() {
					if oErr != nil {
						requestErr <- oErr
					}
				}()

				// prepare request body
				body := &bytes.Buffer{}
				data := entities.NewInTopic(t.Name, t.Description)
				if err := json.NewEncoder(body).Encode(data); err != nil {
					oErr = errors.Join(oErr, fmt.Errorf("UpdateTopics: encode body failed for topic %q: %w", t.Name, err))
					return
				}

				apiURL, err := url.JoinPath(s.baseURL, "topics", strconv.Itoa(t.ID))
				if err != nil {
					oErr = errors.Join(oErr, fmt.Errorf("UpdateTopics: join api url failed for topic %q: %w", t.Name, err))
					return
				}
				slog.Debug("api url", "url", apiURL)

				req, err := http.NewRequest(http.MethodPut, apiURL, body)
				if err != nil {
					oErr = errors.Join(oErr, fmt.Errorf("UpdateTopics: new request failed for topic %q: %w", t.Name, err))
					return
				}
				req.Header.Set("content-type", "application/json")
				req.Header.Set("Authorization", "Bearer "+s.token)

				res, err := httpClient.Do(req)
				if err != nil {
					oErr = errors.Join(oErr, fmt.Errorf("UpdateTopics: requset failed for topic %q: %w", t.Name, err))
					return
				}

				defer func() {
					oErr = errors.Join(oErr, drainAndClose(res.Body))
				}()

				// process response and send it through the channel
				resBody, err := io.ReadAll(res.Body)
				if err != nil {
					oErr = errors.Join(oErr, fmt.Errorf("UpdateTopics: read response body failed for topic %q: %w", t.Name, err))
					return
				}
				if res.StatusCode >= 400 {
					oErr = errors.Join(oErr, fmt.Errorf("UpdateTopics: status code %d for topic %q: %s", res.StatusCode, t.Name, string(resBody)))
					return
				}
				resData := entities.RetSuccess[entities.Topic]{}
				if err := json.Unmarshal(resBody, &resData); err != nil {
					oErr = errors.Join(oErr, fmt.Errorf("UpdateTopics: parse response body failed for topic %q: %w", t.Name, err))
					return
				}

				successResponse <- resData.Msg

			}(topic)
		}
	}

	// wait for all requests to finish or if an error occurs
	result := []entities.Topic{}
	for {
		if len(result) == len(topics) {
			slog.Debug("updated topics", "count", len(result))
			return result, nil
		}
		select {
		case newTopic := <-successResponse:
			result = append(result, newTopic)
		case err := <-requestErr:
			return []entities.Topic{}, err
		}
	}
}
func (s SyncHelper) DeleteTopics(topics []entities.Topic) error {
	slog.Debug("DeleteTopics")

	batchData := make(chan []entities.Topic, 1)
	go batch[entities.Topic](topics, s.batchSize, batchData)

	requestErr := make(chan error, 1)
	successResponse := make(chan bool, len(topics))

	// seperate into batches
	for currentBatch := range batchData {
		// there is no 'bulk api' for now, so we just create topics one by one
		for _, topic := range currentBatch {
			go func(t entities.Topic) {
				var oErr error
				defer func() {
					if oErr != nil {
						requestErr <- oErr
					}
				}()

				apiURL, err := url.JoinPath(s.baseURL, "topics", strconv.Itoa(t.ID))
				if err != nil {
					oErr = errors.Join(oErr, fmt.Errorf("DeleteTopics: join api url failed for topic %q: %w", t.Name, err))
					return
				}
				slog.Debug("api url", "url", apiURL)

				req, err := http.NewRequest(http.MethodDelete, apiURL, nil)
				if err != nil {
					oErr = errors.Join(oErr, fmt.Errorf("DeleteTopics: new request failed for topic %q: %w", t.Name, err))
					return
				}
				req.Header.Set("content-type", "application/json")
				req.Header.Set("Authorization", "Bearer "+s.token)

				res, err := httpClient.Do(req)
				if err != nil {
					oErr = errors.Join(oErr, fmt.Errorf("DeleteTopics: requset failed for topic %q: %w", t.Name, err))
					return
				}

				defer func() {
					oErr = errors.Join(oErr, drainAndClose(res.Body))
				}()

				// process response and send it through the channel
				resBody, err := io.ReadAll(res.Body)
				if err != nil {
					oErr = errors.Join(oErr, fmt.Errorf("DeleteTopics: read response body failed for topic %q: %w", t.Name, err))
					return
				}
				if res.StatusCode >= 400 {
					oErr = errors.Join(oErr, fmt.Errorf("DeleteTopics: status code %d for topic %q: %s", res.StatusCode, t.Name, string(resBody)))
					return
				}
				resData := entities.RetSuccess[entities.RowsAffected]{}
				if err := json.Unmarshal(resBody, &resData); err != nil {
					oErr = errors.Join(oErr, fmt.Errorf("DeleteTopics: parse response body failed for topic %q: %w", t.Name, err))
					return
				}

				if resData.Msg.AffectedRows != 1 {
					oErr = errors.Join(oErr, fmt.Errorf("DeleteTopics: should only delete one topic %q", t.Name))
					return
				}

				successResponse <- true

			}(topic)
		}
	}

	// wait for all requests to finish or if an error occurs
	successCount := 0
	for {
		if successCount == len(topics) {
			slog.Debug("deleted topics", "count", successCount)
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
