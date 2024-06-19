package main

import (
	"blog/entities"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
)

func (s SyncHelper) GetAllBlogs() (oBlog []entities.OutBlogSimple, oErr error) {
	slog.Debug("GetAllBlogs")

	req, err := http.NewRequest(http.MethodGet, s.baseURL+"/blogs", nil)
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

	return data.Msg, nil
}
func (s SyncHelper) CreateBlogs(blogs []entities.InBlog) error {
	// load blog content by batch
	return nil
}
func (s SyncHelper) UpdateBlogs(blogs []entities.InBlog) error {
	// load blog content by batch
	return nil
}
func (s SyncHelper) DeleteBlogs(blogs []entities.InBlog) error {
	return nil
}
