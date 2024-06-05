package handlers

import (
	"blog/entities"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
)

var (
	ErrorTargetNotFound = errors.New("Target not found")
)

func writeJSON(w http.ResponseWriter, err error, msg any, status int) error {
	w.Header().Set("content-type", "application/json")
	w.WriteHeader(status)
	var body *entities.Ret
	if err != nil {
		body = entities.NewRet(err.Error(), status, msg)
	} else {
		body = entities.NewRet(nil, status, msg)
	}
	if err := json.NewEncoder(w).Encode(body); err != nil {
		slog.Error("writeJSON: encode error", "error", err.Error())
		return fmt.Errorf("writeJSON: encode error: %w", err)
	}
	return nil
}

func affectedRowsResponse(affectedRows int) map[string]int {
	return map[string]int{
		"affactedRows": affectedRows,
	}
}
