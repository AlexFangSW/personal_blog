package handlers_test

import (
	"blog/api/handlers"
	"blog/entities"
	"context"
	"errors"
	"net/http"
	"testing"
)

type dummyUsersRepo struct {
	jwt string
}

func (d *dummyUsersRepo) Get(ctx context.Context) (*entities.User, error) {
	user := &entities.User{
		JWT: d.jwt,
	}
	return user, nil
}
func (d *dummyUsersRepo) UpdateJWT(ctx context.Context, jwt string) error {
	return nil
}
func (d *dummyUsersRepo) ClearJWT(ctx context.Context) error {
	return nil
}

func (d *dummyUsersRepo) Create(ctx context.Context, user entities.InUser) (*entities.User, error) {
	return &entities.User{}, nil
}
func (d *dummyUsersRepo) Update(ctx context.Context, user entities.InUser) (*entities.User, error) {

	return &entities.User{}, nil
}
func (d *dummyUsersRepo) Delete(ctx context.Context) error {
	return nil
}

type dummyJWTHelper struct {
	jwt string
}

func (d *dummyJWTHelper) GenJWT(name string) (string, error) {
	return d.jwt, nil
}
func (d *dummyJWTHelper) VerifyJWT(token string) error {
	if token == d.jwt {
		return nil
	}
	return errors.New("VerifyJWT: failed")
}

func TestAuthVerify(t *testing.T) {
	authHelper := handlers.NewAuthHelper(
		&dummyUsersRepo{jwt: "aaa.bbb.ccc"},
		&dummyJWTHelper{jwt: "aaa.bbb.ccc"},
	)

	// pass
	dummyRequestWithAuth := &http.Request{
		Header: map[string][]string{
			"Authorization": {"Bearer aaa.bbb.ccc"},
		},
	}
	pass, err := authHelper.Verify(dummyRequestWithAuth)
	if err != nil {
		t.Fatalf("TestAuthVerify: should have passed: %s", err)
	}
	if !pass {
		t.Fatalf("TestAuthVerify: should have passed: %s", err)
	}

	// fail by no auth header
	dummyRequestWithNoAuth := &http.Request{
		Header: map[string][]string{},
	}
	pass2, err2 := authHelper.Verify(dummyRequestWithNoAuth)
	if err2 == nil {
		t.Fatalf("TestAuthVerify: should not have passed")
	}
	if pass2 {
		t.Fatalf("TestAuthVerify: should not have passed")
	}

	// fail by wrong token
	authHelper2 := handlers.NewAuthHelper(
		&dummyUsersRepo{jwt: "aaa.ccc.ccc"},
		&dummyJWTHelper{jwt: "aaa.bbb.ccc"},
	)
	dummyRequestWithNoAuth2 := &http.Request{
		Header: map[string][]string{
			"Authorization": {"Bearer aaa.bbb.ccc"},
		},
	}
	pass3, err3 := authHelper2.Verify(dummyRequestWithNoAuth2)
	if err3 == nil {
		t.Fatalf("TestAuthVerify: should not have passed")
	}
	if pass3 {
		t.Fatalf("TestAuthVerify: should not have passed")
	}
}
