package handlers_test

import (
	"blog/api/handlers"
	"blog/entities"
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
	newTag := &entities.Tag{Name: "from create"}
	return newTag, nil
}
func (d *DummyTagsRepo) List(ctx context.Context) ([]entities.Tag, error) {
	newTag := []entities.Tag{{Name: "from list"}}
	return newTag, nil
}
func (d *DummyTagsRepo) Get(ctx context.Context, id int) (*entities.Tag, error) {
	newTag := &entities.Tag{Name: "from get"}
	return newTag, nil
}
func (d *DummyTagsRepo) Update(ctx context.Context, tag entities.Tag, id int) (*entities.Tag, error) {
	newTag := &entities.Tag{Name: "from update"}
	return newTag, nil
}
func (d *DummyTagsRepo) Delete(ctx context.Context, id int) (int, error) {
	return 111, nil
}

func initTags() *handlers.Tags {
	repo := &DummyTagsRepo{}
	auth := &DummyAuthHelper{}
	return handlers.NewTags(repo, auth)
}

func TestHandlerTagsCreate(t *testing.T) {
	tags := initTags()

	// how to add body and header ????
	r := httptest.NewRequest(http.MethodGet, "/tags", nil)
	w := httptest.NewRecorder()

	if err := tags.CreateTag(w, r); err != nil {
		t.Fatalf("TestHandlerTagsCreate: create tag failed: %s", err)
	}

	res := w.Result()
	defer res.Body.Close()

	body := &entities.Tag{}
	if err := json.NewDecoder(res.Body).Decode(body); err != nil {
		t.Fatalf("TestHandlerTagsCreate: read response body failed: %s", err)
	}

	if body.Name != "from create" {
		t.Fatalf("TestHandlerTagsCreate: weird response")
	}
}
