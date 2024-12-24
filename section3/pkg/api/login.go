package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type LoginRequest struct {
	Password string `json:"password"`
}

type LoginResponse struct {
	Token string `json:"token"`
}

func doLoginRequest(client http.Client, loginURL, password string) (string, error) {
	loginRequest := LoginRequest{
		Password: password,
	}

	requestBody, err := json.Marshal(loginRequest)
	if err != nil {
		return "", fmt.Errorf("LoginRequest Marshall error: %s", err)
	}

	response, err := client.Post(loginURL, "application/json", bytes.NewBuffer(requestBody))
	if err != nil {
		return "", fmt.Errorf("http Post error: %s", err)
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return "", fmt.Errorf("ReadAll error: %s", err)
	}

	if response.StatusCode != 200 {
		return "", fmt.Errorf("Invalid output (HTTP Code %d): %s", response.StatusCode, string(requestBody))
	}
	if !json.Valid(body) {
		return "", RequestError{
			HTTPCode: response.StatusCode,
			Body:     string(body),
			Err:      fmt.Sprintf("No valid json returned"),
		}
	}

	var loginResponse LoginResponse

	err = json.Unmarshal(body, &loginResponse)

	if err != nil {
		return "", RequestError{
			HTTPCode: response.StatusCode,
			Body:     string(body),
			Err:      fmt.Sprintf("LoginResponse unmarshall error: %s", err),
		}
	}

	return loginResponse.Token, nil
}
