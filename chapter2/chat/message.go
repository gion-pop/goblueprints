package main

import (
	"time"
)

// message は 1 つのメッセージを表す
type message struct {
	Name    string
	Message string
	When    time.Time
}
