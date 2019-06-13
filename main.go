package main

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"
)

// LoginRequest - Request structure for logins
type LoginRequest struct {
	Identifier struct {
		UserType string `json:"type"`
		User     string `json:"user"`
	}
	DisplayName string `json:"initial_device_display_name"`
	Password    string `json:"password"`
	AuthType    string `json:"type"`
	User        string `json:"user"`
}

// LoginResponse - Respose structure for logins
type LoginResponse struct {
	HomeServer  string `json:"home_server"`
	AccessToken string `json:"access_token"`
	DeviceID    string `json:"device_id"`
	UserID      string `json:"user_id"`
}

var client http.Client
var homeserverURL string

const userAgent = "MatrixBot/0.0 golang"

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

	return httpReq
}

func login(username string, password string) LoginResponse {
	req := &LoginRequest{
		DisplayName: "Matrix Bot",
		Password:    password,
		AuthType:    "m.login.password",
		User:        username,
	}
	req.Identifier.UserType = "m.id.user"
	req.Identifier.User = username

	_ = assembleRequest("/r0/login", "POST", req)

	//resp, _ := client.Do(httpReq)
	//defer resp.Body.Close()

	var respObj LoginResponse
	//json.NewDecoder(resp.Body).Decode(&respObj)

	return respObj
}

func main() {
	proxyURL, _ := url.Parse("http://localhost:5555")

	homeserverURL = "https://matrix.test.c583.psiroom.net"
	client = http.Client{
		Transport: &http.Transport{
			Proxy:           http.ProxyURL(proxyURL),
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
		Timeout: time.Second * 10,
	}

	loginResponse := login(botUserName, botPassword)

	fmt.Printf("%+v\n", loginResponse)
}
