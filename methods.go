package main

import (
	"encoding/json"
	"fmt"
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

	_ = assembleRequest("/r0/login", "POST", req)

	//resp, _ := client.Do(httpReq)
	//defer resp.Body.Close()

	var respObj LoginResponse
	//json.NewDecoder(resp.Body).Decode(&respObj)

	return respObj
}

func sync() []Event {
	type RoomState struct {
		Timeline struct {
			Limited   bool        `json:"limited"`
			PrevBatch string      `json:"prev_batch"`
			Events    []RoomEvent `json:"events"`
		} `json:"timeline"`
		AccountData         map[string]interface{} `json:"account_data"`
		Ephemeral           map[string]interface{} `json:"ephemeral"`
		UnreadNotifications struct {
			HighlightCount    uint16 `json:"highlight_count"`
			NotificationCount uint16 `json:"notification_count"`
		} `json:"unread_notifications"`
	}

	type SyncResponse struct {
		NextBatch              string `json:"next_batch"`
		DeviceOneTimeKeysCount string `json:"device_one_time_keys_count"`
		AccountData            string `json:"account_data"`
		Presence               string `json:"presence"`
		Rooms                  struct {
			Leave  map[string]RoomState `json:"leave"`
			Join   map[string]RoomState `json:"join"`
			Invite map[string]RoomState `json:"invite"`
		}
	}

	endpoint := "/r0/sync"

	if since != "" {
		endpoint = endpoint + "?since=" + since
	}

	resp, err := client.Do(assembleRequest(endpoint, "GET", nil))

	if err != nil {
		fmt.Println("Error: ", err.Error())
	}

	defer resp.Body.Close()

	var respObj SyncResponse
	json.NewDecoder(resp.Body).Decode(&respObj)

	var events []Event

	for roomName, roomContents := range respObj.Rooms.Join {
		for _, roomEvent := range roomContents.Timeline.Events {
			events = append(events, Event{roomEvent, roomName})
		}
	}

	since = respObj.NextBatch
	return events
}
