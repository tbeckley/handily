package main

import (
	"fmt"
	"time"
)

type messageHandlerType = func(Event) bool

func setupHandlers() {
	handlers = append(handlers, debugHandler)
}

// Put handlers here

//Basic debug handler
var debugHandler = func(e Event) bool {
	if e.Event.EventType == "m.room.message" {
		msgTime := time.Unix(int64(e.Event.ServerTS)/1000, 0).String()
		toPrint := fmt.Sprintf("%s said %s at %s", e.Event.Sender, e.Event.Content["body"], msgTime)

		sendMessage(e.Room, toPrint)
	}
	return false
}
