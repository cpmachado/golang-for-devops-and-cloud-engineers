package main

// RequestError is an error to manage after retrieving response
type RequestError struct {
	HTTPCode int
	Body     string
	Err      string
}

func (r RequestError) Error() string {
	return r.Err
}
