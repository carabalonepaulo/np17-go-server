package main

import (
	"strings"
)

type Message struct {
	Name    string
	Content string
}

type InvalidMessage struct{}

func (e InvalidMessage) Error() string {
	return "Invalid message."
}

func ParseRawMessage(raw string) *Message {
	idx := strings.Index(raw, ">")
	if idx == -1 {
		return nil
	}

	name := raw[1:idx]

	lastIdx := strings.LastIndex(raw, "<")
	if lastIdx == -1 || lastIdx <= idx {
		return nil
	}

	return &Message{Name: name, Content: raw[idx+1 : lastIdx]}
}
