package main

import (
	"blog/entities"
	"bufio"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
	"strings"

	"golang.org/x/term"
)

type Credentials struct {
	Username string
	Password string
}

func NewCredentials(username, password string) Credentials {
	return Credentials{
		Username: username,
		Password: password,
	}
}

func getJWT(baseURL, username, password string) (oStr string, oErr error) {
	slog.Debug("getJWT")

	// setup api url
	// build request
	req, err := http.NewRequest(http.MethodPost, baseURL+"/login", nil)
	if err != nil {
		return "", fmt.Errorf("getJWT: create new request failed: %w", err)
	}
	req.SetBasicAuth(username, password)
	req.Header.Set("Content-type", "application/json")

	// send request
	res, err := httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("getJWT: request failed: %w", err)
	}

	// cleanup
	defer func() {
		oErr = errors.Join(oErr, drainAndClose(res.Body))
	}()

	// read request
	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		return "", fmt.Errorf("getJWT: read body failed: %w", err)
	}

	// fail
	if res.StatusCode >= 400 {
		return "", fmt.Errorf("getJWT: status code %d, msg: %s", res.StatusCode, string(resBody))
	}

	// success
	msg := entities.RetSuccess[entities.JWT]{}
	if err := json.Unmarshal(resBody, &msg); err != nil {
		return "", fmt.Errorf("getJWT: decode success response failed: %w", err)
	}
	return msg.Msg.JWT, nil
}

// reads username and password and get jwt token
func login(ctx context.Context, done chan<- bool, baseURL, username, password string) (oStr string, oErr error) {
	slog.Info("login")

	defer func() {
		done <- true
	}()

	if username == "" && password == "" {
		// get current terminal state
		currFd := int(os.Stdin.Fd())
		currState, err := term.GetState(currFd)
		if err != nil {
			return "", fmt.Errorf("Get treminal current state error: %w", err)
		}

		defer func() {
			slog.Debug("restore terminal state")
			oErr = errors.Join(oErr, term.Restore(currFd, currState))
		}()
	}

	processErr := make(chan error, 1)
	result := make(chan string, 1)

	go func() {
		if username == "" && password == "" {
			cred, err := stdinCredentials()
			if err != nil {
				processErr <- fmt.Errorf("stdinCredentials error: %w", err)
			}
			username = cred.Username
			password = cred.Password
		}

		// get jwt
		jwt, err := getJWT(baseURL, username, password)
		if err != nil {
			processErr <- fmt.Errorf("Get jwt token error: %w", err)
			return
		}
		result <- jwt
	}()

	select {
	case <-ctx.Done():
		slog.Warn("got done")
		return "", errors.New("login: canceled")
	case err := <-processErr:
		return "", err
	case jwt := <-result:
		return jwt, nil
	}
}

func stdinCredentials() (Credentials, error) {
	// read username
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Username: ")
	username, err := reader.ReadString('\n')
	if err != nil {
		return Credentials{}, fmt.Errorf("Read username error: %w", err)
	}
	username = strings.TrimSuffix(username, "\n")

	// read password
	fmt.Print("Password: ")
	bytepw, err := term.ReadPassword(int(os.Stdin.Fd()))
	if err != nil {
		return Credentials{}, fmt.Errorf("Read password error: %w", err)
	}
	password := string(bytepw)

	return NewCredentials(username, password), nil
}
