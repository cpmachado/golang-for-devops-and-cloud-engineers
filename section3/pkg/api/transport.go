package api

import (
	"net/http"
)

type JotTransport struct {
	token     string
	transport http.RoundTripper
	password  string
	loginURL  string
}

func (j *JotTransport) RoundTrip(request *http.Request) (*http.Response, error) {
	if j.token == "" && j.password != "" {
		token, err := doLoginRequest(http.Client{}, j.loginURL, j.password)
		if err != nil {
			return nil, err
		}
		j.token = token
	}
	if j.token != "" {
		request.Header.Add("Authorization", "Bearer "+j.token)
	}

	return j.transport.RoundTrip(request)
}
