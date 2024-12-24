package main

import (
	"flag"
	"fmt"
	"net/url"
	"os"

	"github.com/cpmachado/golang-for-devops-and-cloud-engineers/section3/shgo/pkg/api"
)

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
	loginURL := fmt.Sprintf("%s://%s/login", parsedURL.Scheme, parsedURL.Host)

	apiInstance := api.New(api.Options{
		Password: password,
		LoginURL: loginURL,
	})

	response, err := apiInstance.DoRequest(parsedURL.String())
	if err != nil {
		if requestError, ok := err.(api.RequestError); ok {
			fmt.Printf(
				"Error: %s, Status code: %d, body: %s\n",
				requestError.Err, requestError.HTTPCode, requestError.Body)
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
