package main

import (
	"encoding/json"
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
	args := os.Args
	progName = args[0]

	if len(args) != 2 {
		help()
	}

	res, err := doRequest(args[1])
	if err != nil {
		if reqError, ok := err.(RequestError); ok {
			fmt.Printf("Error: %s, Status code: %d, body: %s\n", reqError.Err, reqError.HTTPCode, reqError.Body)
		} else {
			fmt.Printf("Error: %s\n", err)
		}
		os.Exit(1)
	}
	if res == nil {
		fmt.Printf("No response\n")
		os.Exit(1)
	}
	fmt.Printf("Response:\n%s\n", res.GetResponse())
}

// help prints help and exits
func help() {
	fmt.Printf("Usage: %v <url>\n", progName)
	os.Exit(1)
}

func doRequest(target string) (Response, error) {
	if _, err := url.ParseRequestURI(target); err != nil {
		return nil, fmt.Errorf("Invalid URL: %v\n", target)
	}

	res, err := http.Get(target)
	if err != nil {
		return nil, fmt.Errorf("Get error: %s\n", err)
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)

	if err != nil {
		return nil, fmt.Errorf("ReadAll(Body) error: %s\n", err)
	}

	if res.StatusCode != 200 {
		return nil, fmt.Errorf("Invalid output(HTTP code %d): %s\n", res.StatusCode, body)
	}

	var page Page
	err = json.Unmarshal(body, &page)
	if err != nil {
		return nil, RequestError{
			HTTPCode: res.StatusCode,
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
				HTTPCode: res.StatusCode,
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
				HTTPCode: res.StatusCode,
				Body:     string(body),
				Err:      fmt.Sprintf("Unmarshal(Occurrence) error: %s\n", err),
			}
		}
		return &occurence, nil
	}

	return nil, nil
}
