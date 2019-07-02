package main

// These are all large "public" types used by the rest of the bot
// Other types should be kept in their respective files

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

// Config - Configuration file schema
type config struct {
	HomeserverURL   string `json:"homeserverURL"`
	BotUsername     string `json:"botUsername"`
	BotPassword     string `json:"botPassword"`
	CustomUserAgent string `json:"customUserAgent"`
}
