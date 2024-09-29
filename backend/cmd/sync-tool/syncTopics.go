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

func (s SyncHelper) GetAllTopics() (oTopic []entities.Topic, oErr error) {
	slog.Info("GetAllTopics")

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

func (s SyncHelper) createTopic(t entities.Topic) (result entities.Topic, oErr error) {
	slog.Debug("createTopic")

	// prepare request body
	body := &bytes.Buffer{}
	data := entities.NewInTopic(t.Name, t.Description)
	if err := json.NewEncoder(body).Encode(data); err != nil {
		return entities.Topic{}, fmt.Errorf("createTopic: encode body failed for topic %q: %w", t.Name, err)
	}

	apiURL, err := url.JoinPath(s.baseURL, "topics")
	if err != nil {
		return entities.Topic{}, fmt.Errorf("createTopic: join api url failed for topic %q: %w", t.Name, err)
	}
	slog.Debug("api url", "url", apiURL)

	req, err := http.NewRequest(http.MethodPost, apiURL, body)
	if err != nil {
		return entities.Topic{}, fmt.Errorf("createTopic: new requset failed for topic %q: %w", t.Name, err)
	}
	req.Header.Set("content-type", "application/json")
	req.Header.Set("Authorization", "Bearer "+s.token)

	res, err := httpClient.Do(req)
	if err != nil {
		return entities.Topic{}, fmt.Errorf("createTopic: requset failed for topic %q: %w", t.Name, err)
	}

	defer func() {
		oErr = errors.Join(oErr, drainAndClose(res.Body))
	}()

	// process response and send it through the channel
	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		return entities.Topic{}, fmt.Errorf("createTopic: read response body failed for topic %q: %w", t.Name, err)
	}
	if res.StatusCode >= 400 {
		return entities.Topic{}, fmt.Errorf("createTopic: status code %d for topic %q: %s", res.StatusCode, t.Name, string(resBody))
	}
	resData := entities.RetSuccess[entities.Topic]{}
	if err := json.Unmarshal(resBody, &resData); err != nil {
		return entities.Topic{}, fmt.Errorf("createTopic: parse response body failed for topic %q: %w", t.Name, err)
	}

	return resData.Msg, nil
}

func (s SyncHelper) CreateTopics(topics []entities.Topic) ([]entities.Topic, error) {
	slog.Info("CreateTopics", "count", len(topics))

	batchData := make(chan []entities.Topic, 1)
	go batch(topics, s.batchSize, batchData)

	result := []entities.Topic{}

	// seperate into batches
	for currentBatch := range batchData {
		requestErr := make(chan error, 1)
		response := make(chan entities.Topic, 1)
		responseCount := 0

		// there is no 'bulk api' for now, so we just create topics one by one
		for _, topic := range currentBatch {
			go func(t entities.Topic) {
				res, err := s.createTopic(t)
				if err != nil {
					requestErr <- err
					return
				}
				response <- res
			}(topic)
		}

		// wait for all requests to finish or if an error occurs
		for {
			if responseCount == len(currentBatch) {
				break
			}
			select {
			case newTopic := <-response:
				result = append(result, newTopic)
				responseCount++
			case err := <-requestErr:
				return []entities.Topic{}, err
			}
		}
	}

	slog.Info("created topics", "count", len(result))
	return result, nil
}

func (s SyncHelper) updateTopic(t entities.Topic) (result entities.Topic, oErr error) {
	slog.Debug("updateTopic")

	// prepare request body
	body := &bytes.Buffer{}
	data := entities.NewInTopic(t.Name, t.Description)
	if err := json.NewEncoder(body).Encode(data); err != nil {
		return entities.Topic{}, fmt.Errorf("updateTopic: encode body failed for topic %q: %w", t.Name, err)
	}

	apiURL, err := url.JoinPath(s.baseURL, "topics", strconv.Itoa(t.ID))
	if err != nil {
		return entities.Topic{}, fmt.Errorf("updateTopic: join api url failed for topic %q: %w", t.Name, err)
	}
	slog.Debug("api url", "url", apiURL)

	req, err := http.NewRequest(http.MethodPut, apiURL, body)
	if err != nil {
		return entities.Topic{}, fmt.Errorf("updateTopic: new request failed for topic %q: %w", t.Name, err)
	}
	req.Header.Set("content-type", "application/json")
	req.Header.Set("Authorization", "Bearer "+s.token)

	res, err := httpClient.Do(req)
	if err != nil {
		return entities.Topic{}, fmt.Errorf("updateTopic: requset failed for topic %q: %w", t.Name, err)
	}

	defer func() {
		oErr = errors.Join(oErr, drainAndClose(res.Body))
	}()

	// process response and send it through the channel
	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		return entities.Topic{}, fmt.Errorf("updateTopic: read response body failed for topic %q: %w", t.Name, err)
	}
	if res.StatusCode >= 400 {
		return entities.Topic{}, fmt.Errorf("updateTopic: status code %d for topic %q: %s", res.StatusCode, t.Name, string(resBody))
	}
	resData := entities.RetSuccess[entities.Topic]{}
	if err := json.Unmarshal(resBody, &resData); err != nil {
		return entities.Topic{}, fmt.Errorf("updateTopic: parse response body failed for topic %q: %w", t.Name, err)
	}

	return resData.Msg, nil
}

func (s SyncHelper) UpdateTopics(topics []entities.Topic) ([]entities.Topic, error) {
	slog.Info("UpdateTopics", "count", len(topics))

	batchData := make(chan []entities.Topic, 1)
	go batch(topics, s.batchSize, batchData)

	result := []entities.Topic{}

	// seperate into batches
	for currentBatch := range batchData {
		requestErr := make(chan error, 1)
		response := make(chan entities.Topic, 1)
		responseCount := 0

		// there is no 'bulk api' for now, so we just create topics one by one
		for _, topic := range currentBatch {
			go func(t entities.Topic) {
				res, err := s.updateTopic(t)
				if err != nil {
					requestErr <- err
					return
				}
				response <- res
			}(topic)
		}

		// wait for all requests to finish or if an error occurs
		for {
			if responseCount == len(currentBatch) {
				break
			}
			select {
			case newTopic := <-response:
				result = append(result, newTopic)
				responseCount++
			case err := <-requestErr:
				return []entities.Topic{}, err
			}
		}
	}

	slog.Info("updated topics", "count", len(result))
	return result, nil
}

func (s SyncHelper) deleteTopic(t entities.Topic) (oErr error) {
	slog.Debug("deleteTopic")

	apiURL, err := url.JoinPath(s.baseURL, "topics", strconv.Itoa(t.ID))
	if err != nil {
		return fmt.Errorf("deleteTopic: join api url failed for topic %q: %w", t.Name, err)
	}
	slog.Debug("api url", "url", apiURL)

	req, err := http.NewRequest(http.MethodDelete, apiURL, nil)
	if err != nil {
		return fmt.Errorf("deleteTopic: new request failed for topic %q: %w", t.Name, err)
	}
	req.Header.Set("content-type", "application/json")
	req.Header.Set("Authorization", "Bearer "+s.token)

	res, err := httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("deleteTopic: requset failed for topic %q: %w", t.Name, err)
	}

	defer func() {
		oErr = errors.Join(oErr, drainAndClose(res.Body))
	}()

	// process response and send it through the channel
	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		return fmt.Errorf("deleteTopic: read response body failed for topic %q: %w", t.Name, err)
	}
	if res.StatusCode >= 400 {
		return fmt.Errorf("deleteTopic: status code %d for topic %q: %s", res.StatusCode, t.Name, string(resBody))
	}
	resData := entities.RetSuccess[entities.RowsAffected]{}
	if err := json.Unmarshal(resBody, &resData); err != nil {
		return fmt.Errorf("deleteTopic: parse response body failed for topic %q: %w", t.Name, err)
	}

	if resData.Msg.AffectedRows != 1 {
		return fmt.Errorf("deleteTopic: should only delete one topic %q", t.Name)
	}
	return nil
}

func (s SyncHelper) DeleteTopics(topics []entities.Topic) error {
	slog.Info("DeleteTopics", "count", len(topics))

	batchData := make(chan []entities.Topic, 1)
	go batch(topics, s.batchSize, batchData)

	totalCount := 0

	// seperate into batches
	for currentBatch := range batchData {
		requestErr := make(chan error, 1)
		finish := make(chan bool, 1)
		finishCount := 0
		// there is no 'bulk api' for now, so we just create topics one by one
		for _, topic := range currentBatch {
			go func(t entities.Topic) {
				if err := s.deleteTopic(t); err != nil {
					requestErr <- err
					return
				}
				finish <- true
			}(topic)
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

	slog.Info("deleted topics", "count", totalCount)
	return nil
}
