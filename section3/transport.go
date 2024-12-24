package main

import "net/http"

type JotTransport struct {
	token     string
	transport http.RoundTripper
}

func (j JotTransport) RoundTrip(request *http.Request) (*http.Response, error) {
	if j.token != "" {
		request.Header.Add("Authorization", "Bearer "+j.token)
	}

	return j.transport.RoundTrip(request)
}
