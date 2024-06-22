package main

import (
	"blog/entities"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
)

func (s SyncHelper) GetAllTopics() (oTag []entities.Topic, oErr error) {
	slog.Debug("GetAllTopics")

	res, err := httpClient.Get(s.baseURL + "/topics")
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
func (s SyncHelper) CreateTopics(topics []entities.Topic) error {
	return nil
}
func (s SyncHelper) UpdateTopics(topics []entities.Topic) error {
	return nil
}
func (s SyncHelper) DeleteTopics(topics []entities.Topic) error {
	return nil
}
