package script

import (
	"fmt"
	"log"

	"github.com/Shopify/go-lua"
)

const CALLBACKS = 1
const UPDATE = 2

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
	size := s.L.Top()
	fmt.Println("-----------------------------------")
	fmt.Printf("- Stack: %d\n", size)
	fmt.Println("-----------------------------------")
	for i := 1; i <= size; i++ {
		fmt.Printf("> [%d / -%d] ", i, size-i+1)
		switch {
		case s.L.IsNumber(i):
			value, _ := s.L.ToNumber(i)
			fmt.Printf("%.2f\n", value)
		case s.L.IsBoolean(i):
			value := s.L.ToBoolean(i)
			fmt.Printf("%t\n", value)
		case s.L.IsFunction(i):
			fmt.Println("function")
		case s.L.IsTable(i):
			fmt.Println("table")
		case s.L.IsGoFunction(i):
			fmt.Println("native function")
		case s.L.IsUserData(i):
			fmt.Println("userdata")
		case s.L.IsLightUserData(i):
			fmt.Println("lighuserdata")
		case s.L.IsThread(i):
			fmt.Println("coroutine")
		case s.L.IsNil(i):
			fmt.Println("nil")
		}
	}
	fmt.Println("-----------------------------------")
}

func (s *Engine) Init() {
	if err := lua.DoFile(s.L, "./scripts/init.lua"); err != nil {
		log.Fatal(err)
	}
	s.L.Field(CALLBACKS, "update")

	s.L.Field(CALLBACKS, "init")
	s.L.Call(0, 0)
}

func (s *Engine) Deinit() {
	s.L.Field(CALLBACKS, "deinit")
	s.L.Call(0, 0)
}

func (s *Engine) ClientConnected(clientId int) {
	s.L.Field(CALLBACKS, "on_client_connected")
	s.L.PushInteger(clientId)
	s.L.Call(1, 0)
}

func (s *Engine) ClientDisconnected(clientId int) {
	s.L.Field(CALLBACKS, "on_client_disconnected")
	s.L.PushInteger(clientId)
	s.L.Call(1, 0)
}

func (s *Engine) MessageReceived(clientId int, messageName, messageContent string) {
	s.L.Field(CALLBACKS, "on_data_received")
	s.L.PushInteger(clientId)
	s.L.PushString(messageName)
	s.L.PushString(messageContent)
	s.L.Call(3, 0)
}

func (s *Engine) Update() {
	s.L.PushValue(UPDATE)
	s.L.Call(0, 0)
}
