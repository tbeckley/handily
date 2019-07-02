package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"os"
)

func assembleRequest(endpoint string, requestType string, body interface{}) *http.Request {
	buf := new(bytes.Buffer)
	json.NewEncoder(buf).Encode(body)

	fullURL := homeserverURL + "/_matrix/client" + endpoint
	httpReq, err := http.NewRequest(requestType, fullURL, buf)

	if err != nil {
		panic("Error assembling request:" + err.Error())
	}

	httpReq.Header.Add("Content-Type", "application/json")
	httpReq.Header.Add("User-Agent", userAgent)

	if loginToken != "" {
		httpReq.Header.Add("Authorization", "Bearer "+loginToken)
	}

	return httpReq
}

const defaultHomeserver = "https://matrix.org"
const defaultUserAgent = "Golang/0.1 MatrixBot"

func parseConfig(fileName string) config {
	file, _ := os.Open(fileName)
	defer file.Close()

	configVals := config{}
	json.NewDecoder(file).Decode(&configVals)

	// Default values
	if configVals.HomeserverURL == "" {
		configVals.HomeserverURL = defaultHomeserver
	}

	if configVals.CustomUserAgent == "" {
		configVals.CustomUserAgent = defaultUserAgent
	}

	return configVals
}
