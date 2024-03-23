package script

import (
	"log"
	"server/src/script/libs"

	"github.com/Shopify/go-lua"
)

const (
	Init = iota + 1
	Deinit
	Update
	OnClientConnected
	OnClientDisconnected
	OnMessageReceived
)

type Engine struct {
	L          *lua.State
	shouldQuit bool
	commands   chan Command
}

type Command interface {
	Handle(l *lua.State) error
}

func NewEngine() *Engine {
	l := lua.NewState()

	s := &Engine{L: l, shouldQuit: false, commands: make(chan Command)}

	lua.OpenLibraries(l)
	lua.DoString(l, "package.path = package.path .. ';' .. './?/init.lua;./scripts/?/init.lua;./scripts/?.lua'")

	return s
}

func (s *Engine) DumpStack() {
	libs.DumpStack(s.L)
}

func (s *Engine) Init() {
	if err := lua.DoFile(s.L, "./scripts/init.lua"); err != nil {
		log.Fatal(err)
	}

	s.L.Field(1, "init")
	s.L.Field(1, "deinit")
	s.L.Field(1, "update")
	s.L.Field(1, "on_client_connected")
	s.L.Field(1, "on_client_disconnected")
	s.L.Field(1, "on_data_received")

	s.L.Remove(1)

	s.L.PushValue(Init)
	s.L.Call(0, 0)
}

func (s *Engine) Deinit() {
	s.L.PushValue(Deinit)
	s.L.Call(0, 0)
}

func (s *Engine) ClientConnected(clientId int) {
	s.L.PushValue(OnClientConnected)
	s.L.PushInteger(clientId)
	s.L.Call(1, 0)
}

func (s *Engine) ClientDisconnected(clientId int) {
	s.L.PushValue(OnClientDisconnected)
	s.L.PushInteger(clientId)
	s.L.Call(1, 0)
}

func (s *Engine) MessageReceived(clientId int, messageName, messageContent string) {
	s.L.PushValue(OnMessageReceived)
	s.L.PushInteger(clientId)
	s.L.PushString(messageName)
	s.L.PushString(messageContent)
	s.L.Call(3, 0)
}

func (s *Engine) Update() {
	s.L.PushValue(Update)
	s.L.Call(0, 0)
}
