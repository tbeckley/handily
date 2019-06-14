package main

import (
	"bytes"
	"encoding/json"
	"net/http"
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
