package handlers

import (
	"blog/config"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var (
	defaultJWTAud = "coding.notes.backend"
)

type jwtHelper interface {
	GenJWT(name string) (string, error)
	VerifyJWT(token string) error
}

type JWTHelper struct {
	config config.JWTSetting
}

func NewJWTHelper(config config.JWTSetting) *JWTHelper {
	return &JWTHelper{
		config: config,
	}
}

func (j *JWTHelper) GenJWT(name string) (string, error) {
	token := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		jwt.MapClaims{
			"iss": j.config.Issuer,
			"sub": name,
			"exp": time.Now().UTC().Add(time.Duration(j.config.Expire) * time.Hour).Unix(),
			"nbf": time.Now().UTC().Unix(),
			"iat": time.Now().UTC().Unix(),
			"aud": defaultJWTAud,
		},
	)
	if j.config.Secret == "" {
		return "", fmt.Errorf("genJWT: must have secret")
	}
	signedToken, err := token.SignedString([]byte(j.config.Secret))
	if err != nil {
		return "", fmt.Errorf("genJWT: sign jwt failed: %w", err)
	}
	return signedToken, nil
}

// returns nil on success
func (j *JWTHelper) VerifyJWT(token string) error {
	parseFunc := func(token *jwt.Token) (interface{}, error) {
		return []byte(j.config.Secret), nil
	}
	parsedToken, err := jwt.Parse(
		token,
		parseFunc,
		jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Name}),
		jwt.WithIssuer(j.config.Issuer),
		jwt.WithAudience(defaultJWTAud),
	)

	if err != nil {
		return fmt.Errorf("verifyJWT: validate failed: %w", err)
	}

	if parsedToken.Valid {
		return nil
	}
	return fmt.Errorf("verifyJWT: unexpedted error")
}
