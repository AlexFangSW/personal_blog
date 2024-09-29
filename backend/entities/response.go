package entities

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
)

type JWT struct {
	JWT string `json:"jwt"`
}

func NewJWT(jwt string) *JWT {
	return &JWT{
		JWT: jwt,
	}
}

type RowsAffected struct {
	AffectedRows int `json:"affectedRows"`
}

func NewRowsAffected(affectedRows int) *RowsAffected {
	return &RowsAffected{
		AffectedRows: affectedRows,
	}
}

type MsgType interface {
	RowsAffected | OutBlog | []OutBlog | []OutBlogSimple |
		Tag | []Tag | Topic | []Topic |
		~string | JWT
}

type RetSuccess[T MsgType] struct {
	Error  string `json:"error"`
	Status int    `json:"status"`
	Msg    T      `json:"msg"`
}

func NewRetSuccess[T MsgType](msg T) *RetSuccess[T] {
	return &RetSuccess[T]{
		Status: http.StatusOK,
		Msg:    msg,
	}
}

func (r *RetSuccess[T]) WriteJSON(w http.ResponseWriter) error {
	w.Header().Set("content-type", "application/json")
	w.WriteHeader(r.Status)
	if err := json.NewEncoder(w).Encode(r); err != nil {
		slog.Error("WriteJSON: RetSuccess encode error", "error", err.Error())
		return fmt.Errorf("WriteJSON: RetSuccess encode error: %w", err)
	}
	return nil
}

type RetFailed struct {
	Error      string `json:"error"`
	Status     int    `json:"status"`
	Msg        string `json:"msg"`
	httpStatus int
}

func NewRetFailed(err error, status int) *RetFailed {
	return &RetFailed{
		Error:      err.Error(),
		Status:     status,
		httpStatus: status,
	}
}

// When the http status code differs from the cusotm status code
func NewRetFailedCustom(err error, status, httpStatus int) *RetFailed {
	return &RetFailed{
		Error:      err.Error(),
		Status:     status,
		httpStatus: httpStatus,
	}
}

func (r *RetFailed) WriteJSON(w http.ResponseWriter) error {
	w.Header().Set("content-type", "application/json")
	w.WriteHeader(r.httpStatus)
	if err := json.NewEncoder(w).Encode(r); err != nil {
		slog.Error("WriteJSON: RetFailed encode error", "error", err.Error())
		return fmt.Errorf("WriteJSON: RetFailed encode error: %w", err)
	}
	return nil
}
