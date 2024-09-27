package handlers

import (
	"blog/entities"
	"context"
	"encoding/base64"
	"fmt"
	"log/slog"
	"net/http"
	"strings"
)

// Concrete implementations are at repository/<name>
type usersRepository interface {
	Get(ctx context.Context) (*entities.User, error)
	UpdateJWT(ctx context.Context, jwt string) error
	ClearJWT(ctx context.Context) error

	Create(ctx context.Context, user entities.InUser) (*entities.User, error)
	Update(ctx context.Context, user entities.InUser) (*entities.User, error)
	Delete(ctx context.Context) error
}

type Users struct {
	repo usersRepository
	jwt  jwtHelper
	auth authHelper
}

func NewUsers(repo usersRepository, jwt jwtHelper, auth authHelper) *Users {
	return &Users{
		repo: repo,
		jwt:  jwt,
		auth: auth,
	}
}

// Login
//
//	@Summary		Login
//	@Description	login to get jwt token
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Param			Authorization	header		string	true	"user credentials"
//	@Success		200				{object}	entities.RetSuccess[entities.JWT]
//	@Failure		400				{object}	entities.RetFailed
//	@Failure		412				{object}	entities.RetFailed
//	@Failure		422				{object}	entities.RetFailed
//	@Failure		500				{object}	entities.RetFailed
//	@Router			/login [post]
func (u *Users) Login(w http.ResponseWriter, r *http.Request) error {
	slog.Debug("Login")

	// varify user credentials
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return entities.NewRetFailed(ErrorAuthorizationHeaderEmpty, http.StatusPreconditionFailed).WriteJSON(w)
	}

	inUser, err := readCredentials(authHeader)
	if err != nil {
		return entities.NewRetFailed(err, http.StatusUnprocessableEntity).WriteJSON(w)
	}

	user, err := u.repo.Get(r.Context())
	if err != nil {
		return entities.NewRetFailed(err, http.StatusInternalServerError).WriteJSON(w)
	}

	authorized := verifyUser(*inUser, *user)
	if !authorized {
		return entities.NewRetFailed(ErrorAuthorizationFailed, http.StatusBadRequest).WriteJSON(w)
	}

	// generate jwt token
	newToken, err := u.jwt.GenJWT(inUser.Name)
	if err != nil {
		return entities.NewRetFailed(err, http.StatusInternalServerError).WriteJSON(w)
	}

	// update stored jwt token
	if err := u.repo.UpdateJWT(r.Context(), newToken); err != nil {

		if sqliteErr, ok := getSQLiteError(err); ok {
			slog.Error("got sqlite error", "error code", sqliteErr.Code, "extended error code", sqliteErr.ExtendedCode)
			return entities.NewRetFailedCustom(err, int(sqliteErr.ExtendedCode), http.StatusInternalServerError).WriteJSON(w)
		}

		return entities.NewRetFailed(err, http.StatusInternalServerError).WriteJSON(w)
	}

	return entities.NewRetSuccess(*entities.NewJWT(newToken)).WriteJSON(w)
}

// Logout
//
//	@Summary		Logout
//	@Description	logout, deletes jwt token, needs to have valid token in the first place
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Param			Authorization	header		string	true	"jwt token"
//	@Success		200				{object}	entities.RetSuccess[string]
//	@Failure		412				{object}	entities.RetFailed
//	@Failure		403				{object}	entities.RetFailed
//	@Failure		500				{object}	entities.RetFailed
//	@Router			/logout [post]
func (u *Users) Logout(w http.ResponseWriter, r *http.Request) error {
	slog.Debug("Logout")

	// varifiy jwt token
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return entities.NewRetFailed(ErrorAuthorizationHeaderEmpty, http.StatusPreconditionFailed).WriteJSON(w)
	}

	token := readToken(authHeader)
	err := u.jwt.VerifyJWT(token)
	if err != nil {
		return entities.NewRetFailed(err, http.StatusForbidden).WriteJSON(w)
	}

	// delete jwt from user
	if err := u.repo.ClearJWT(r.Context()); err != nil {

		if sqliteErr, ok := getSQLiteError(err); ok {
			slog.Error("got sqlite error", "error code", sqliteErr.Code, "extended error code", sqliteErr.ExtendedCode)
			return entities.NewRetFailedCustom(err, int(sqliteErr.ExtendedCode), http.StatusInternalServerError).WriteJSON(w)
		}

		return entities.NewRetFailed(err, http.StatusInternalServerError).WriteJSON(w)
	}

	return entities.NewRetSuccess("logout success").WriteJSON(w)
}

// AuthorizeCheck
//
//	@Summary		AuthorizeCheck
//	@Description	Checks if jwt is valid
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Param			Authorization	header		string	true	"jwt token"
//	@Success		200				{object}	entities.RetSuccess[string]
//	@Failure		403				{object}	entities.RetFailed
//	@Router			/auth-check [post]
func (u *Users) AuthorizeCheck(w http.ResponseWriter, r *http.Request) error {
	slog.Debug("AuthorizeCheck")

	// authorization
	authorized, err := u.auth.Verify(r)
	if err != nil || !authorized {
		slog.Warn("AuthorizeCheck: authorization failed", "error", err.Error())
		return entities.NewRetFailed(err, http.StatusForbidden).WriteJSON(w)
	}

	return entities.NewRetSuccess("pass").WriteJSON(w)
}

func readCredentials(authHeader string) (*entities.InUser, error) {
	// Authorization: Basic <base64 encoded stuff>
	encodedData := authHeader[6:]
	data, err := base64.StdEncoding.DecodeString(encodedData)
	if err != nil {
		return &entities.InUser{}, fmt.Errorf("readAuthorization: base64 decode failed: %w", err)
	}
	ret := strings.Split(string(data), ":")
	return entities.NewInUser(ret[0], ret[1]), nil
}

func readToken(authHeader string) string {
	// Authorization: Bearer <jwt token>
	return authHeader[7:]
}
