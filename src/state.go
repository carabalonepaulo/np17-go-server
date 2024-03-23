package main

import (
	"database/sql"
	"server/src/listener"
	"server/src/script"

	"github.com/carabalonepaulo/weave"
)

type State struct {
	listener *listener.Listener
	script   *script.Engine
	workers  *weave.WorkerPool
	db       *sql.DB
}
