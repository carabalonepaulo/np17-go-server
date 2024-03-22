package main

import (
	"log"
)

type HandlerFn func(
	state *State,
	sender int,
	message *Message,
)

type Handler struct {
	string
	HandlerFn
}

func FindHandlerFn(handlers []Handler, messageName string) HandlerFn {
	for i := 0; i < len(handlers); i++ {
		if handlers[i].string == messageName {
			return handlers[i].HandlerFn
		}
	}
	return nil
}

func NativeHandlers() []Handler {
	return []Handler{
		// {"0", HandleAuth},
		{"sign_in", HandleSignIn},
		{"sign_up", HandleSignUp},
	}
}

func HandleAuth(state *State, sender int, message *Message) {
	log.Println("Auth message received!")
}

func HandleSignIn(state *State, sender int, message *Message) {
	log.Println("SignIn message received!")
}

func HandleSignUp(state *State, sender int, message *Message) {
	log.Println("SignUp message received!")
}
