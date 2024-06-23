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
	baseURL   string
	token     string // jwt token
	batchSize int
}

func NewSyncHelper(baseURL, token string, batchSize int) SyncHelper {
	return SyncHelper{
		baseURL:   baseURL,
		token:     token,
		batchSize: batchSize,
	}
}

func batch[T any](inpt []T, size int, out chan<- []T) {
	inptLength := len(inpt)
	for i := 0; i < inptLength; i += size {
		upperBound := min(i+size, inptLength)
		out <- inpt[i:upperBound:upperBound]
	}
	close(out)
	return
}
