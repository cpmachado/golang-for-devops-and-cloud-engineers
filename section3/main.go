package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
)

// Reponse is an interface to manage the parsing of the API Payload
type Response interface {
	// GetResponse retrieves a serialized version of the content
	GetResponse() string
}

// Page represents the type of page being parsed
type Page struct {
	Name string `json:"page"`
}

// Words represents the type associated with the /words endpoint
type Words struct {
	Input string   `json:"input"`
	Words []string `json:"words"`
}

func (w *Words) GetResponse() string {
	return fmt.Sprintf("- Input: %s\n- Words: %s\n", w.Input, strings.Join(w.Words, ", "))
}

// Occurrence represents the type associated with the /occurrence endpoint
type Occurence struct {
	Words map[string]int `json:"words"`
}

func (o *Occurence) GetResponse() string {
	out := []string{}

	for word, count := range o.Words {
		out = append(out, fmt.Sprintf("%s (%d)", word, count))
	}
	return fmt.Sprintf("- Words: %s\n", strings.Join(out, ", "))
}

// Global variables

var progName string = ""

func main() {
	var (
		requestURL string
		password   string
		parsedURL  *url.URL
		err        error
	)

	flag.StringVar(&requestURL, "url", "", "url to access")
	flag.StringVar(&password, "password", "", "use a password to access our api")

	flag.Parse()

	if parsedURL, err = url.ParseRequestURI(requestURL); err != nil {
		fmt.Printf("Invalid URL: %v\n", err)
		flag.Usage()
		os.Exit(1)
	}

	response, err := doRequest(parsedURL.String())
	if err != nil {
		if reqError, ok := err.(RequestError); ok {
			fmt.Printf("Error: %s, Status code: %d, body: %s\n", reqError.Err, reqError.HTTPCode, reqError.Body)
		} else {
			fmt.Printf("Error: %s\n", err)
		}
		os.Exit(1)
	}
	if response == nil {
		fmt.Printf("No response\n")
		os.Exit(1)
	}
	fmt.Printf("Response:\n%s", response.GetResponse())
}

func doRequest(target string) (Response, error) {
	if _, err := url.ParseRequestURI(target); err != nil {
		return nil, fmt.Errorf("Invalid URL: %v\n", target)
	}

	response, err := http.Get(target)
	if err != nil {
		return nil, fmt.Errorf("Get error: %s\n", err)
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)

	if err != nil {
		return nil, fmt.Errorf("ReadAll(Body) error: %s\n", err)
	}

	if response.StatusCode != 200 {
		return nil, fmt.Errorf("Invalid output(HTTP code %d): %s\n", response.StatusCode, body)
	}

	var page Page
	err = json.Unmarshal(body, &page)
	if err != nil {
		return nil, RequestError{
			HTTPCode: response.StatusCode,
			Body:     string(body),
			Err:      fmt.Sprintf("Unmarshal(Page) error: %s\n", err),
		}
	}

	switch page.Name {
	case "words":
		var words Words
		err = json.Unmarshal(body, &words)
		if err != nil {
			return nil, RequestError{
				HTTPCode: response.StatusCode,
				Body:     string(body),
				Err:      fmt.Sprintf("Unmarshal(Words) error: %s\n", err),
			}
		}
		return &words, nil
	case "occurrence":
		var occurence Occurence
		err = json.Unmarshal(body, &occurence)
		if err != nil {
			return nil, RequestError{
				HTTPCode: response.StatusCode,
				Body:     string(body),
				Err:      fmt.Sprintf("Unmarshal(Occurrence) error: %s\n", err),
			}
		}
		return &occurence, nil
	}

	return nil, nil
}
