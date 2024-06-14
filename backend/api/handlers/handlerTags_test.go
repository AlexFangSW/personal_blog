package handlers_test

import (
	"blog/api/handlers"
	"blog/entities"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
)

type DummyAuthHelper struct{}

func (d *DummyAuthHelper) Verify(r *http.Request) (bool, error) {
	token := r.Header.Get("Authorization")
	if token == "" {
		return false, errors.New("Verify: verify failed")
	}
	return true, nil
}

type DummyTagsRepo struct{}

func (d *DummyTagsRepo) Create(ctx context.Context, tag entities.Tag) (*entities.Tag, error) {
	newTag := entities.NewTag(tag.Name, "create")
	return newTag, nil
}
func (d *DummyTagsRepo) List(ctx context.Context) ([]entities.Tag, error) {
	newTag := []entities.Tag{{Name: "list"}}
	return newTag, nil
}
func (d *DummyTagsRepo) Get(ctx context.Context, id int) (*entities.Tag, error) {
	newTag := &entities.Tag{Name: "get"}
	return newTag, nil
}
func (d *DummyTagsRepo) Update(ctx context.Context, tag entities.Tag, id int) (*entities.Tag, error) {
	newTag := entities.NewTag(tag.Name, "update")
	return newTag, nil
}
func (d *DummyTagsRepo) Delete(ctx context.Context, id int) (int, error) {
	return 1, nil
}

func initTags() *handlers.Tags {
	repo := &DummyTagsRepo{}
	auth := &DummyAuthHelper{}
	return handlers.NewTags(repo, auth)
}

func TestHandlerTagsCreate(t *testing.T) {
	tags := initTags()

	// prepare request
	reqData := entities.NewInTag("dummy name", "dummy description")
	reqBody := bytes.Buffer{}
	if err := json.NewEncoder(&reqBody).Encode(reqData); err != nil {
		t.Fatalf("TestHandlerTagsCreate: encode request body failed: %f", err)
	}
	r := httptest.NewRequest(http.MethodGet, "/tags", &reqBody)
	r.Header.Set("Authorization", "Bearer aaa.bbb.ccc")

	// prepare response recorder
	w := httptest.NewRecorder()

	// call api
	if err := tags.CreateTag(w, r); err != nil {
		t.Fatalf("TestHandlerTagsCreate: create tag failed: %s", err)
	}

	// read result
	res := w.Result()
	defer res.Body.Close()

	resData := entities.RetSuccess[entities.Tag]{}
	if err := json.NewDecoder(res.Body).Decode(&resData); err != nil {
		t.Fatalf("TestHandlerTagsCreate: read response body failed: %s", err)
	}

	// check response
	if resData.Status != http.StatusOK {
		t.Fatalf("TestHandlerTagsCreate: status incorrect")
	}
	if resData.Msg.Name != reqData.Name {
		t.Fatalf("TestHandlerTagsCreate: tag name incorrect")
	}
	if resData.Msg.Description != "create" {
		t.Fatalf("TestHandlerTagsCreate: didn't call the correct repository method")
	}
}

func TestHandlerTagsCreateAuthFail(t *testing.T) {
	tags := initTags()

	// prepare request
	reqData := entities.NewInTag("dummy name", "dummy description")
	reqBody := bytes.Buffer{}
	if err := json.NewEncoder(&reqBody).Encode(reqData); err != nil {
		t.Fatalf("TestHandlerTagsCreate: encode request body failed: %f", err)
	}
	r := httptest.NewRequest(http.MethodGet, "/tags", &reqBody)

	// prepare response recorder
	w := httptest.NewRecorder()

	// call api
	if err := tags.CreateTag(w, r); err != nil {
		t.Fatalf("TestHandlerTagsCreate: create tag failed: %s", err)
	}

	// read result
	res := w.Result()
	defer res.Body.Close()

	resData := entities.RetFailed{}
	if err := json.NewDecoder(res.Body).Decode(&resData); err != nil {
		t.Fatalf("TestHandlerTagsCreate: read response body failed: %s", err)
	}

	// check response
	if resData.Status != http.StatusForbidden {
		t.Fatalf("TestHandlerTagsCreate: status incorrect")
	}
}

func TestHandlerTagsCreateNoBody(t *testing.T) {
	tags := initTags()

	// prepare request
	r := httptest.NewRequest(http.MethodGet, "/tags", nil)
	r.Header.Set("Authorization", "Bearer aaa.bbb.ccc")

	// prepare response recorder
	w := httptest.NewRecorder()

	// call api
	if err := tags.CreateTag(w, r); err != nil {
		t.Fatalf("TestHandlerTagsCreate: create tag failed: %s", err)
	}

	// read result
	res := w.Result()
	defer res.Body.Close()

	resData := entities.RetFailed{}
	if err := json.NewDecoder(res.Body).Decode(&resData); err != nil {
		t.Fatalf("TestHandlerTagsCreate: read response body failed: %s", err)
	}

	// check response
	if resData.Status != http.StatusBadRequest {
		t.Fatalf("TestHandlerTagsCreate: status incorrect")
	}
}
