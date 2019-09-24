package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	auth "github.com/maddiesch/pco-auth"
)

func main() {
	logger := log.New(os.Stderr, "pco-auth --> ", 0)

	input := auth.PerformInput{}

	var scope string

	flag.IntVar(&input.Port, "port", 8080, "the callback port to listen on")
	flag.StringVar(&input.ClientID, "client_id", "", "the oauth client id")
	flag.StringVar(&input.ClientSecret, "client_secret", "", "the oauth client secret")
	flag.StringVar(&scope, "scope", "people", "the Oauth scopes to authenticate with")

	flag.Parse()

	input.Scopes = strings.Split(scope, ",")

	output, err := auth.Perform(&input)
	if err != nil {
		logger.Fatal(err)
	}

	dump, err := json.MarshalIndent(output, "", " ")
	if err != nil {
		logger.Fatal(err)
	}

	fmt.Fprintln(os.Stdout, string(dump))
}
