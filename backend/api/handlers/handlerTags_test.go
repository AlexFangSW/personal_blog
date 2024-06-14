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
	"strconv"
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
func (d *DummyTagsRepo) ListByTopicID(ctx context.Context, topicID int) ([]entities.Tag, error) {
	newTag := []entities.Tag{{Name: "list by topic ID", Description: strconv.Itoa(topicID)}}
	return newTag, nil
}
func (d *DummyTagsRepo) Get(ctx context.Context, id int) (*entities.Tag, error) {
	newTag := &entities.Tag{Name: "get", Description: strconv.Itoa(id)}
	return newTag, nil
}
func (d *DummyTagsRepo) Update(ctx context.Context, tag entities.Tag, id int) (*entities.Tag, error) {
	newTag := entities.NewTag(tag.Name, "update"+strconv.Itoa(id))
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

/* ============ Create ============== */

func TestHandlerTagsCreate(t *testing.T) {
	tags := initTags()

	// prepare request
	reqData := entities.NewInTag("dummy name", "dummy description")
	reqBody := bytes.Buffer{}
	if err := json.NewEncoder(&reqBody).Encode(reqData); err != nil {
		t.Fatalf("TestHandlerTagsCreate: encode request body failed: %f", err)
	}
	r := httptest.NewRequest(http.MethodPut, "/tags", &reqBody)
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
	r := httptest.NewRequest(http.MethodPut, "/tags", &reqBody)

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
	r := httptest.NewRequest(http.MethodPut, "/tags", nil)
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

/* ============ List ============== */

func TestHandlerTagsList(t *testing.T) {
	tags := initTags()

	// prepare request
	r := httptest.NewRequest(http.MethodGet, "/tags", nil)

	// prepare response recorder
	w := httptest.NewRecorder()

	// call api
	if err := tags.ListTags(w, r); err != nil {
		t.Fatalf("TestHandlerTagsList: list tags failed: %s", err)
	}

	// read result
	res := w.Result()
	defer res.Body.Close()

	resData := entities.RetSuccess[[]entities.Tag]{}
	if err := json.NewDecoder(res.Body).Decode(&resData); err != nil {
		t.Fatalf("TestHandlerTagsList: read response body failed: %s", err)
	}

	// check response
	if resData.Status != http.StatusOK {
		t.Fatalf("TestHandlerTagsList: status incorrect")
	}
	if resData.Msg[0].Name != "list" {
		t.Fatalf("TestHandlerTagsList: didn't call the correct repo method")
	}
}

/* ============ Get ============== */

func TestHandlerTagsGet(t *testing.T) {
	tags := initTags()

	// prepare request
	// any thing here besides body isn't quite useful.
	// still not sure how to test handler with mux.
	r := httptest.NewRequest(http.MethodGet, "/tags/1", nil)
	// had to manually set path value
	// but I think it's acceptable for unit test
	r.SetPathValue("id", "1")

	// prepare response recorder
	w := httptest.NewRecorder()

	// call api
	if err := tags.GetTag(w, r); err != nil {
		t.Fatalf("TestHandlerTagsGet: get tag failed: %s", err)
	}

	// read result
	res := w.Result()
	defer res.Body.Close()

	resData := entities.RetSuccess[entities.Tag]{}
	if err := json.NewDecoder(res.Body).Decode(&resData); err != nil {
		t.Fatalf("TestHandlerTagsGet: read response body failed: %s", err)
	}

	// check response
	if resData.Status != http.StatusOK {
		t.Fatalf("TestHandlerTagsGet: status incorrect")
	}
	if resData.Msg.Name != "get" {
		t.Fatalf("TestHandlerTagsGet: didn't call the correct repo method")
	}
	if resData.Msg.Description != "1" {
		t.Fatalf("TestHandlerTagsGet: id isn't properly passed")
	}
}

/* ============ Update ============== */

func TestHandlerTagsUpdate(t *testing.T) {
	tags := initTags()

	// prepare request
	reqData := entities.NewInTag("dummy name", "aaaa")
	reqBody := bytes.Buffer{}
	if err := json.NewEncoder(&reqBody).Encode(reqData); err != nil {
		t.Fatalf("TestHandlerTagsUpdate: encode req body failed: %s", err)
	}
	r := httptest.NewRequest(http.MethodPut, "/tags/1", &reqBody)
	r.SetPathValue("id", "1")
	r.Header.Set("Authorization", "Bearer aaa.bbb.ccc")

	// prepare response recorder
	w := httptest.NewRecorder()

	// call api
	if err := tags.UpdateTag(w, r); err != nil {
		t.Fatalf("TestHandlerTagsUpdate: update tag failed: %s", err)
	}

	// read result
	res := w.Result()
	defer res.Body.Close()

	resData := entities.RetSuccess[entities.Tag]{}
	if err := json.NewDecoder(res.Body).Decode(&resData); err != nil {
		t.Fatalf("TestHandlerTagsUpdate: read response body failed: %s", err)
	}

	// check response
	if resData.Status != http.StatusOK {
		t.Fatalf("TestHandlerTagsUpdate: status incorrect")
	}
	if resData.Msg.Name != reqData.Name || resData.Msg.Description != "update1" {
		t.Fatalf("TestHandlerTagsUpdate: data isn't properly passed")
	}
}

func TestHandlerTagsUpdateAuthFail(t *testing.T) {
	tags := initTags()

	// prepare request
	reqData := entities.NewInTag("dummy name", "aaaa")
	reqBody := bytes.Buffer{}
	if err := json.NewEncoder(&reqBody).Encode(reqData); err != nil {
		t.Fatalf("TestHandlerTagsUpdateAuthFail: encode req body failed: %s", err)
	}
	r := httptest.NewRequest(http.MethodPut, "/tags/1", &reqBody)
	r.SetPathValue("id", "1")

	// prepare response recorder
	w := httptest.NewRecorder()

	// call api
	if err := tags.UpdateTag(w, r); err != nil {
		t.Fatalf("TestHandlerTagsUpdateAuthFail: update tag failed: %s", err)
	}

	// read result
	res := w.Result()
	defer res.Body.Close()

	resData := entities.RetFailed{}
	if err := json.NewDecoder(res.Body).Decode(&resData); err != nil {
		t.Fatalf("TestHandlerTagsUpdateAuthFail: read response body failed: %s", err)
	}

	// check response
	if resData.Status != http.StatusForbidden {
		t.Fatalf("TestHandlerTagsUpdateAuthFail: status incorrect")
	}
}

func TestHandlerTagsUpdateNoBody(t *testing.T) {
	tags := initTags()

	// prepare request
	r := httptest.NewRequest(http.MethodPut, "/tags/1", nil)
	r.Header.Set("Authorization", "Bearer aaa.bbb.ccc")
	r.SetPathValue("id", "1")

	// prepare response recorder
	w := httptest.NewRecorder()

	// call api
	if err := tags.UpdateTag(w, r); err != nil {
		t.Fatalf("TestHandlerTagsUpdateNoBody: update tag failed: %s", err)
	}

	// read result
	res := w.Result()
	defer res.Body.Close()

	resData := entities.RetFailed{}
	if err := json.NewDecoder(res.Body).Decode(&resData); err != nil {
		t.Fatalf("TestHandlerTagsUpdateNoBody: read response body failed: %s", err)
	}

	// check response
	if resData.Status != http.StatusBadRequest {
		t.Fatalf("TestHandlerTagsUpdateNoBody: status incorrect")
	}
}

/* ============ Delete ============== */

func TestHandlerTagsDelete(t *testing.T) {
	tags := initTags()

	// prepare request
	r := httptest.NewRequest(http.MethodPut, "/tags/1", nil)
	r.SetPathValue("id", "1")
	r.Header.Set("Authorization", "Bearer aaa.bbb.ccc")

	// prepare response recorder
	w := httptest.NewRecorder()

	// call api
	if err := tags.DeleteTag(w, r); err != nil {
		t.Fatalf("TestHandlerTagsDelete: update tag failed: %s", err)
	}

	// read result
	res := w.Result()
	defer res.Body.Close()

	resData := entities.RetSuccess[entities.RowsAffected]{}
	if err := json.NewDecoder(res.Body).Decode(&resData); err != nil {
		t.Fatalf("TestHandlerTagsDelete: read response body failed: %s", err)
	}

	// check response
	if resData.Status != http.StatusOK {
		t.Fatalf("TestHandlerTagsDelete: status incorrect")
	}
	if resData.Msg.AffectedRows != 1 {
		t.Fatalf("TestHandlerTagsDelete: not passed to the correct method")
	}
}

func TestHandlerTagsDeleteAuthFail(t *testing.T) {
	tags := initTags()

	// prepare request
	r := httptest.NewRequest(http.MethodPut, "/tags/1", nil)
	r.SetPathValue("id", "1")

	// prepare response recorder
	w := httptest.NewRecorder()

	// call api
	if err := tags.DeleteTag(w, r); err != nil {
		t.Fatalf("TestHandlerTagsDeleteAuthFail: update tag failed: %s", err)
	}

	// read result
	res := w.Result()
	defer res.Body.Close()

	resData := entities.RetFailed{}
	if err := json.NewDecoder(res.Body).Decode(&resData); err != nil {
		t.Fatalf("TestHandlerTagsDeleteAuthFail: read response body failed: %s", err)
	}

	// check response
	if resData.Status != http.StatusForbidden {
		t.Fatalf("TestHandlerTagsDeleteAuthFail: status incorrect")
	}
}
