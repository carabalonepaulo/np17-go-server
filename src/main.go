package main

import (
	"io"
	"log"
	"os"
	"server/src/listener"
	"server/src/script"
	"server/src/script/libs"
	"time"
)

const Host = "0.0.0.0"
const Port = 5000
const MaxClients = 1024
const MaxWorkers = 16

func main() {
	f, err := os.OpenFile("data/log.txt", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("Error opening file: %v", err)
	}
	defer f.Close()

	logOutput := io.MultiWriter(os.Stdout, f)
	log.SetOutput(logOutput)

	workers := NewWorkerPool(MaxWorkers)
	nativeHandlers := NativeHandlers()

	script := script.NewEngine()
	log.Println("Script engine initialized!")

	listener, err := listener.NewListener(Host, Port, MaxClients)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Listener initialized!")

	state := &State{listener: listener, script: script, workers: workers}

	listener.OnClientConnected = script.ClientConnected
	listener.OnClientDisconnected = script.ClientDisconnected
	listener.OnMessageReceived = func(sender int, rawMessage string) {
		message := ParseRawMessage(rawMessage)
		handlerFn := FindHandlerFn(nativeHandlers, message.Name)

		if handlerFn != nil {
			handlerFn(state, sender, message)
		} else {
			script.MessageReceived(sender, message.Name, message.Content)
		}
	}

	libs.RegisterListenerLib(script.L, listener)
	script.Init()

	for listener.Running() {
		time.Sleep(time.Millisecond)
		listener.DispatchEvents()
		workers.Poll()
	}

	script.Deinit()
	listener.Close()
}
