package main

import (
	"io"
	"log"
	"os"
	"server/src/listener"
	"server/src/message"
	"server/src/native_handlers"
	"server/src/script"
	"server/src/script/libs"
	"server/src/state"
	"time"
)

const host string = "0.0.0.0"
const port int = 50000
const maxClients int = 1024

func main() {
	f, err := os.OpenFile("data/log.txt", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("Error opening file: %v", err)
	}
	defer f.Close()

	logOutput := io.MultiWriter(os.Stdout, f)
	log.SetOutput(logOutput)

	nativeHandlers := native_handlers.NativeHandlers()

	script := script.NewEngine()
	log.Println("Script engine initialized!")

	listener, err := listener.NewListener(host, port, maxClients)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Listener initialized!")

	state := &state.State{Listener: listener, Script: script}

	listener.OnClientConnected = script.ClientConnected
	listener.OnClientDisconnected = script.ClientDisconnected
	listener.OnMessageReceived = func(sender int, rawMessage string) {
		message := message.ParseRawMessage(rawMessage)
		handlerFn := native_handlers.FindHandlerFn(nativeHandlers, message.Name)

		if handlerFn != nil {
			handlerFn(state, sender, message)
		} else {
			script.MessageReceived(sender, message.Name, message.Content)
		}
	}

	libs.RegisterListenerLib(state)
	script.Init()

	for listener.Running() {
		time.Sleep(time.Millisecond)
		listener.DispatchEvents()
	}

	script.Deinit()
	listener.Close()
}
