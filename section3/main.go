package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
)

type Words struct {
	Page  string   `json:"page"`
	Input string   `json:"input"`
	Words []string `json:"words"`
}

func main() {
	args := os.Args
	prog := args[0]

	if len(args) != 2 {
		fmt.Printf("Usage: %v <url>\n", prog)
		os.Exit(1)
	}

	target := args[1]

	if _, err := url.ParseRequestURI(target); err != nil {
		fmt.Errorf("Usage: %v <url>\n", prog)
		fmt.Errorf("Invalid URL: %v\n", target)
		os.Exit(1)
	}

	res, err := http.Get(target)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)

	if err != nil {
		log.Fatal(err)
	}

	if res.StatusCode != 200 {
		fmt.Errorf("Invalid output(HTTP code %d): %s\n", res.StatusCode, body)
		os.Exit(1)
	}

	var words Words

	err = json.Unmarshal(body, &words)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("JSON Parsed:\n")
	fmt.Printf("- page: %s\n", words.Page)
	fmt.Printf("- input: %v\n", words.Input)
	fmt.Printf("- words: %s\n", strings.Join(words.Words, ", "))
}
