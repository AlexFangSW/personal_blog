package main

import (
	"blog/entities"
	"bufio"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"golang.org/x/term"
)

func getJWT(baseURL, username, password string, client *http.Client) (oStr string, oErr error) {

	// setup api url
	url := fmt.Sprintf("%s/login", baseURL)

	// build request
	req, err := http.NewRequest(http.MethodPost, url, nil)
	if err != nil {
		return "", fmt.Errorf("getJWT: create new request failed: %w", err)
	}
	encodedCredential := base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%s", username, password)))
	req.Header.Set("Authorization", "Basic "+encodedCredential)
	req.Header.Set("Content-type", "application/json")

	// send request
	res, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("getJWT: request failed: %w", err)
	}

	// cleanup
	defer func() {
		// drain request body and close
		r := io.LimitReader(res.Body, LimitReaderSize)
		if n, err := io.Copy(io.Discard, r); err != nil {
			oErr = errors.Join(oErr, fmt.Errorf("getJWT: failed reading response body after %d bytes: %w", n, err))
		}
		oErr = errors.Join(oErr, res.Body.Close())
	}()

	// read request
	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		return "", fmt.Errorf("getJWT: read body failed: %w", err)
	}

	// fail
	if res.StatusCode >= 400 {
		return "", errors.New(string(resBody))
	}

	// success
	msg := entities.RetSuccess[entities.JWT]{}
	if err := json.Unmarshal(resBody, &msg); err != nil {
		return "", fmt.Errorf("getJWT: decode success response failed: %w", err)
	}
	return msg.Msg.JWT, nil
}

// reads username and password and get jwt token
func login(baseURL string, client *http.Client) (str string, err error) {
	// read username
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Username: ")
	username, err := reader.ReadString('\n')
	if err != nil {
		return "", fmt.Errorf("Read username error: %w", err)
	}
	username = strings.TrimSuffix(username, "\n")

	// read password
	fmt.Print("Password: ")
	bytepw, err := term.ReadPassword(int(os.Stdin.Fd()))
	if err != nil {
		return "", fmt.Errorf("Read password error: %w", err)
	}
	password := string(bytepw)

	// get jwt
	jwt, err := getJWT(baseURL, username, password, client)
	if err != nil {
		return "", fmt.Errorf("Get jwt token error: %w", err)
	}

	return jwt, nil
}
