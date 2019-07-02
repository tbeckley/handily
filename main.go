package main

import (
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

var client = http.Client{
	Timeout: time.Second * 10,
}

var since string

var loginToken string
var homeserverURL string
var userAgent string
var selfID string

// Make this relatively large to handle the most notifications you could expect to recieve in a second or so
const bufferSize = 1000

var handlers []messageHandlerType

func main() {
	config := parseConfig("config.json")
	homeserverURL = config.HomeserverURL

	loginResult := login(config.BotUsername, config.BotPassword)
	loginToken = loginResult.AccessToken // Sensitive!
	selfID = loginResult.UserID

	// Initial Sync to establish a baseline
	Sync()

	// Add the handlers to process the events
	setupHandlers()

	// Setup the listeners and handler goroutines
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
