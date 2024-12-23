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

type Page struct {
	Name string `json:"page"`
}

type Words struct {
	Input string   `json:"input"`
	Words []string `json:"words"`
}

type Occurence struct {
	Words map[string]int `json:"words"`
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
		fmt.Printf("Usage: %v <url>\n", prog)
		fmt.Printf("Invalid URL: %v\n", target)
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
		fmt.Printf("Invalid output(HTTP code %d): %s\n", res.StatusCode, body)
		os.Exit(1)
	}

	var page Page
	err = json.Unmarshal(body, &page)
	if err != nil {
		log.Fatal(err)
	}

	switch page.Name {
	case "words":
		var words Words
		err = json.Unmarshal(body, &words)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("JSON Parsed:\n")
		fmt.Printf("- page: %s\n", page.Name)
		fmt.Printf("- input: %v\n", words.Input)
		fmt.Printf("- words: %s\n", strings.Join(words.Words, ", "))
	case "occurrence":
		var occurence Occurence
		err = json.Unmarshal(body, &occurence)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("JSON Parsed:\n")
		fmt.Printf("- page: %s\n", page.Name)
		fmt.Printf("- words: %v\n", occurence.Words)
	default:
		fmt.Printf("Page is Unknown\n")
	}
}
