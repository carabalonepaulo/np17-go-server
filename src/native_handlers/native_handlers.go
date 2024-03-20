package native_handlers

import (
	"log"
	"server/src/message"
	"server/src/state"
)

type HandlerFn func(
	state *state.State,
	sender int,
	message *message.Message,
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

func HandleAuth(state *state.State, sender int, message *message.Message) {
	log.Println("Auth message received!")
	state.Listener.Kick(sender)
}

func HandleSignIn(state *state.State, sender int, message *message.Message) {
	log.Println("SignIn message received!")
}

func HandleSignUp(state *state.State, sender int, message *message.Message) {
	log.Println("SignUp message received!")
}
