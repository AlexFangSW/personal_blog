package main

import (
	"fmt"
	"io"
)

var LimitReaderSize int64 = 10 * 1024 * 1024 // 10MB

func drainAndClose(body io.ReadCloser) error {
	reader := io.LimitReader(body, LimitReaderSize)
	_, drainErr := io.Copy(io.Discard, reader)
	if drainErr != nil {
		return fmt.Errorf("drainAndClose: drain failed: %w", drainErr)
	}
	return body.Close()
}

type SyncHelper struct {
	baseURL string
	token   string // jwt token
}

func NewSyncHelper(baseURL, token string) SyncHelper {
	return SyncHelper{
		baseURL: baseURL,
		token:   token,
	}
}
