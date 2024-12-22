package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
)

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

	fmt.Printf("HTTP Status Code: %d\n", res.StatusCode)
	fmt.Printf("Body: %s\n", body)
}
