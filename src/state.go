package main

import (
	"server/src/listener"
	"server/src/script"
)

type State struct {
	Listener *listener.Listener
	Script   *script.Engine
}
