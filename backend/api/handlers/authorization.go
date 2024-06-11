package handlers

import (
	"fmt"
	"net/http"
)

type authHelper interface {
	Verify(r *http.Request) (bool, error)
}

type AuthHelper struct {
	repo usersRepository
	jwt  jwtHelper
}

func NewAuthHelper(repo usersRepository, jwt jwtHelper) *AuthHelper {
	return &AuthHelper{
		repo: repo,
		jwt:  jwt,
	}
}

func (a *AuthHelper) Verify(r *http.Request) (bool, error) {
	// get auth header
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return false, fmt.Errorf("Verify: auth header empty: %w", ErrorAuthorizationHeaderEmpty)
	}

	// verify token
	token := readToken(authHeader)
	err := a.jwt.VerifyJWT(token)
	if err != nil {
		return false, fmt.Errorf("Verify: verification failed: %w", err)
	}

	// match with cached token
	user, err := a.repo.Get(r.Context())
	if err != nil {
		return false, fmt.Errorf("Verify: get user failed: %w", err)
	}

	if token != user.JWT {
		return false, fmt.Errorf("Verify: where did this token come from ???")
	}

	return true, nil
}
