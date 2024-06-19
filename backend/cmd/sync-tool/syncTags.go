package main

import (
	"blog/entities"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
)

func (s SyncHelper) GetAllTags() (oTag []entities.Tag, oErr error) {
	slog.Debug("GetAllTags")

	res, err := httpClient.Get(s.baseURL + "/tags")
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

	return data.Msg, nil
}
func (s SyncHelper) CreateTags(tags []entities.Tag) error {
	return nil
}
func (s SyncHelper) UpdateTags(tags []entities.Tag) error {
	return nil
}
func (s SyncHelper) DeleteTags(tags []entities.Tag) error {
	return nil
}
