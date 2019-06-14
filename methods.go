package main

import (
	"encoding/json"

	"github.com/google/uuid"
)

// RoomState - The state of messages from a room
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

// SyncResponse - Sync response schema object
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

	respObj = LoginResponse{
		HomeServer:  "matrix.test.c583.psiroom.net",
		AccessToken: "MDAyYWxvY2F0aW9uIG1hdHJpeC50ZXN0LmM1ODMucHNpcm9vbS5uZXQKMDAxM2lkZW50aWZpZXIga2V5CjAwMTBjaWQgZ2VuID0gMQowMDNhY2lkIHVzZXJfaWQgPSBAbWF0cml4Ym90Om1hdHJpeC50ZXN0LmM1ODMucHNpcm9vbS5uZXQKMDAxNmNpZCB0eXBlID0gYWNjZXNzCjAwMjFjaWQgbm9uY2UgPSAyand-amVaQGN6cmlWS0t4CjAwMmZzaWduYXR1cmUgnUX4_W6J6rsiPTQC8x8jqRZpFVfaZCwhjYypY06vQrYK",
		DeviceID:    "PUHNYCEVPT",
		UserID:      "@***REMOVED***:matrix.test.c583.psiroom.net",
	}

	return respObj
}

func sync() SyncResponse {
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
	syncResp := sync()

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
