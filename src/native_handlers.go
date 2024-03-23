package main

import (
	"log"

	"github.com/carabalonepaulo/weave"
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
		{"0", HandleAuth},
		{"sign_in", HandleSignIn},
		{"sign_up", HandleSignUp},
	}
}

func HandleAuth(state *State, sender int, message *Message) {
	type Context struct {
		matches bool
	}
	s := &Context{matches: false}

	state.workers.Dispatch(weave.NewChain(s, 2).Add(weave.Main, func(ctx *Context) {
		log.Println("executing on main thread...")
	}).Add(weave.Background, func(ctx *Context) {
		log.Println("executing on background thread...")
	}).Add(weave.Main, func(ctx *Context) {
		log.Println("back to main thread")
	}))

	log.Println("Auth message received!")
}

func HandleSignIn(state *State, sender int, message *Message) {

	// parts := strings.Split(message.Content, ":")

	// state.workers.Dispatch(weave.NewChain(s, 2).Add(weave.Background, func(ctx *Context) {
	// 	query := "select 1 from accounts where email=? and password = ? limit 1"
	// 	_, err := state.db.Query(query, ctx.email, ctx.password)
	// 	ctx.matches = err == nil
	// }).Add(weave.Main, func(ctx *Context) {
	// 	state.listener.SendTo(sender, fmt.Sprintf("<0 %d>'e' n=%t</0>", sender, ctx.matches))
	// }))

	log.Println("SignIn message received!")
}

func HandleSignUp(state *State, sender int, message *Message) {
	log.Println("SignUp message received!")
}
