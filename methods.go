package main

import (
	"encoding/json"

	"github.com/google/uuid"
)

func login(username string, password string) LoginResponse {
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

	req := &LoginRequest{
		DisplayName: "Matrix Bot",
		Password:    password,
		AuthType:    "m.login.password",
		User:        username,
	}
	req.Identifier.UserType = "m.id.user"
	req.Identifier.User = username

	httpReq := assembleRequest("/r0/login", "POST", req)

	resp, _ := client.Do(httpReq)
	defer resp.Body.Close()

	var respObj LoginResponse
	json.NewDecoder(resp.Body).Decode(&respObj)

	return respObj
}

// Sync - Goes out to the server and fetches response since last update
func Sync() SyncResponse {
	endpoint := "/r0/sync"

	if since != "" {
		endpoint = endpoint + "?since=" + since
	}

	resp, err := client.Do(assembleRequest(endpoint, "GET", nil))

	if err != nil {
		panic("Sync Error: " + err.Error())
	}

	defer resp.Body.Close()

	var respObj SyncResponse
	json.NewDecoder(resp.Body).Decode(&respObj)

	since = respObj.NextBatch
	return respObj
}

func getNewEvents() []Event {
	syncResp := Sync()

	var events []Event

	for roomName, roomContents := range syncResp.Rooms.Join {
		for _, roomEvent := range roomContents.Timeline.Events {
			events = append(events, Event{roomEvent, roomName})
		}
	}

	return events
}

func sendMessage(roomID string, message string) {
	type MessageRequest struct {
		MsgType string `json:"msgtype"`
		Body    string `json:"body"`
	}

	uuid, _ := uuid.NewUUID()

	endpoint := "/r0/rooms/" + roomID + "/send/m.room.message/" + uuid.String()

	client.Do(assembleRequest(endpoint, "PUT", MessageRequest{"m.text", message}))
}
