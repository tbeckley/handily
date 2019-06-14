package main

import (
	"fmt"
	"time"
)

type messageHandlerType = func(Event) bool

func setupHandlers() {
	handlers = append(handlers, debugHandler, secondHandler)
}

// Put handlers here

//Basic debug handler
var debugHandler = func(e Event) bool {
	msgTime := time.Unix(int64(e.Event.ServerTS)/1000, 0).String()
	toPrint := fmt.Sprintf("%s said %s at %s", e.Event.Sender, e.Event.Content["body"], msgTime)

	sendMessage(e.Room, toPrint)

	return false
}

// This should get called
var secondHandler = func(e Event) bool {
	return false
}
