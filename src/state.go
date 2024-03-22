package main

import (
	"server/src/listener"
	"server/src/script"
)

type State struct {
	listener *listener.Listener
	script   *script.Engine
	workers  *WorkerPool
}
