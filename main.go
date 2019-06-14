package main

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"time"
)

// LoginResponse - Respose structure for logins
type LoginResponse struct {
	HomeServer  string `json:"home_server"`
	AccessToken string `json:"access_token"`
	DeviceID    string `json:"device_id"`
	UserID      string `json:"user_id"`
}

// RoomEvent - Event in a room
type RoomEvent struct {
	Sender    string            `json:"sender"`
	EventType string            `json:"type"`
	EventID   string            `json:"event_id"`
	ServerTS  uint64            `json:"origin_server_ts"`
	Content   map[string]string `json:"content"`
	Unsigned  map[string]string `json:"unsgined"`
}

// Event - Event to which we might have to respond to
type Event struct {
	Event RoomEvent
	Room  string
}

var client http.Client
var homeserverURL string

const userAgent = "MatrixBot/0.0 golang"

var since string

var loginToken string
var selfID string

const bufferSize = 1000

var handlers []messageHandlerType

func main() {
	homeserverURL = "https://matrix.test.c583.psiroom.net"
	client = http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
		Timeout: time.Second * 10,
	}

	botUserName := "***REMOVED***"
	botPassword := "***REMOVED***"

	loginRes := login(botUserName, botPassword)
	loginToken = loginRes.AccessToken
	selfID = loginRes.UserID

	sync()

	setupHandlers()

	eventChannel := make(chan Event, bufferSize)
	go vigilant(eventChannel)
	go rootHandler(eventChannel)

	// Blocks for input

	var input string
	fmt.Scanln(&input)
}

// Continually syncs and pipes new events to the channel
func vigilant(ch chan Event) {
	for {
		results := getNewEvents()

		for _, event := range results {
			ch <- event
		}

		time.Sleep(time.Second * 1)
	}
}

// Handles events from the channel in an async buffered way
func rootHandler(ch chan Event) {
	for {
		event := <-ch

		// Ignore self generated events
		if event.Event.Sender != selfID {

			// Try every handler
			for _, handler := range handlers {
				// Indicates if it should act as a final handler and
				// disallow future actions (EG kicking, banning) should not warrant a text response
				if handler(event) == true {
					break
				}
			}
		}
	}
}
