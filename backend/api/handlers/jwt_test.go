package handlers_test

import (
	"blog/api/handlers"
	"blog/config"
	"testing"
)

func TestJWTHelper(t *testing.T) {
	// pass
	jwtHelper := handlers.NewJWTHelper(
		config.JWTSetting{
			Issuer: "alexfangsw",
			Expire: 1,
			Secret: "123123",
		},
	)
	newToken, err := jwtHelper.GenJWT("alex")
	if err != nil {
		t.Fatalf("TestJWTHelper: gen jwt failed: %s", err)
	}
	if err := jwtHelper.VerifyJWT(newToken); err != nil {
		t.Fatalf("TestJWTHelper: verify jwt failed: %s", err)
	}

	// fail
	jwtHelper2 := handlers.NewJWTHelper(
		config.JWTSetting{
			Issuer: "alexfangsw",
			Expire: -1,
			Secret: "123123",
		},
	)
	newToken2, err2 := jwtHelper2.GenJWT("alex")
	if err2 != nil {
		t.Fatalf("TestJWTHelper: gen jwt failed: %s", err2)
	}
	if err := jwtHelper.VerifyJWT(newToken2); err == nil {
		t.Fatalf("TestJWTHelper: verify should have failed")
	}
}
