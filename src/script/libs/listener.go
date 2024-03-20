package libs

import (
	"fmt"
	"server/src/listener"
	"server/src/state"

	"github.com/Shopify/go-lua"
)

func DumpStack(L *lua.State) {
	size := L.Top()
	fmt.Println("-----------------------------------")
	fmt.Printf("- Stack: %d\n", size)
	fmt.Println("-----------------------------------")
	for i := 1; i <= size; i++ {
		fmt.Printf("> [%d / -%d] ", i, size-i+1)
		switch {
		case L.IsNumber(i):
			value, _ := L.ToNumber(i)
			fmt.Printf("%.2f\n", value)
		case L.IsBoolean(i):
			value := L.ToBoolean(i)
			fmt.Printf("%t\n", value)
		case L.IsFunction(i):
			fmt.Println("function")
		case L.IsTable(i):
			fmt.Println("table")
		case L.IsGoFunction(i):
			fmt.Println("native function")
		case L.IsUserData(i):
			fmt.Println("userdata")
		case L.IsLightUserData(i):
			fmt.Println("lighuserdata")
		case L.IsThread(i):
			fmt.Println("coroutine")
		case L.IsNil(i):
			fmt.Println("nil")
		case L.IsString(i):
			value, _ := L.ToString(i)
			fmt.Println(value)
		}
	}
	fmt.Println("-----------------------------------")
}

func RegisterListenerLib(state *state.State) {
	l := state.Script.L
	l.PushUserData(state.Listener)

	if lua.NewMetaTable(l, "Listener") {
		// local t = {}
		l.CreateTable(0, 0)

		// t:get_total_sent(id)
		l.PushGoFunction(totalSent)
		l.SetField(-2, "get_tottal_sent")

		// t:get_total_received(id)
		l.PushGoFunction(totalReceived)
		l.SetField(-2, "get_total_received")

		// t:send_to(id, message)
		l.PushGoFunction(sendTo)
		l.SetField(-2, "send_to")

		// t:send_to_many(message, filter)
		l.PushGoFunction(sendToMany)
		l.SetField(-2, "send_to_many")

		// t:send_to_all(message)
		l.PushGoFunction(sendToAll)
		l.SetField(-2, "send_to_all")

		// t:kick(id)
		l.PushGoFunction(kick)
		l.SetField(-2, "kick")

		// t:kick_all()
		l.PushGoFunction(kickAll)
		l.SetField(-2, "kick_all")

		// mt.__index = t
		l.SetField(-2, "__index")
	}
	l.SetMetaTable(-2)

	l.Global("package")
	l.Field(-1, "loaded")
	l.PushValue(-3)
	l.SetField(-2, "listener")
	l.Pop(3)

	// l.PushUserData(state.Listener)
	// l.CreateTable(0, 0)
	// l.PushGoFunction(sendTo)
	// l.SetField(-2, "send_to")
	// l.SetMetaTable(-2)
	// l.SetGlobal("listener")
	// state.Script.DumpStack()
}

// listener:get_total_sent(id)
func totalSent(l *lua.State) int {
	lua.ArgumentCheck(l, l.IsUserData(1), 1, "expected userdata")
	lua.ArgumentCheck(l, l.IsNumber(2), 2, "expected number")

	id, _ := l.ToNumber(2)
	listener, ok := l.ToUserData(1).(*listener.Listener)
	if !ok {
		l.PushString("Failed to cast userdata!")
		l.Error()
	}

	l.PushInteger(listener.TotalSent(int(id)))

	return 1
}

// listener:get_total_received(id)
func totalReceived(l *lua.State) int {
	lua.ArgumentCheck(l, l.IsUserData(1), 1, "expected userdata")
	lua.ArgumentCheck(l, l.IsNumber(2), 2, "expected number")

	id, _ := l.ToNumber(2)
	listener, ok := l.ToUserData(1).(*listener.Listener)
	if !ok {
		l.PushString("Failed to cast userdata!")
		l.Error()
	}

	l.PushInteger(listener.TotalReceived(int(id)))

	return 1
}

// listener:send_to(id, message)
func sendTo(l *lua.State) int {
	lua.ArgumentCheck(l, l.IsUserData(1), 1, "expected userdata")
	lua.ArgumentCheck(l, l.IsNumber(2), 2, "expected number")
	lua.ArgumentCheck(l, l.IsString(3), 3, "expected string")

	listener, ok := l.ToUserData(1).(*listener.Listener)
	if !ok {
		l.PushString("Failed to cast userdata!")
		l.Error()
	}

	id, _ := l.ToNumber(2)
	message, _ := l.ToString(3)
	if !listener.SendTo(int(id), message) {
		l.PushFString("Failed to send message to client `%d`!", int(id))
		l.Error()
	}

	return 0
}

// listener:send_to_many(message, filter)
func sendToMany(l *lua.State) int {
	lua.ArgumentCheck(l, l.IsUserData(1), 1, "expected userdata")
	lua.ArgumentCheck(l, l.IsString(2), 2, "expected string")
	lua.ArgumentCheck(l, l.IsFunction(3), 3, "expected function")

	listener, ok := l.ToUserData(1).(*listener.Listener)
	if !ok {
		l.PushString("Failed to cast userdata!")
		l.Error()
	}

	message, _ := l.ToString(2)

	listener.SendToMany(message, func(id int) bool {
		l.PushValue(3)
		l.PushInteger(id)
		l.Call(1, 1)
		return l.ToBoolean(-1)
	})

	return 0
}

// listener:send_to_all(message)
func sendToAll(l *lua.State) int {
	lua.ArgumentCheck(l, l.IsUserData(1), 1, "expected userdata")
	lua.ArgumentCheck(l, l.IsString(2), 2, "expected string")

	listener, ok := l.ToUserData(1).(*listener.Listener)
	if !ok {
		l.PushString("Failed to cast userdata!")
		l.Error()
	}

	message, _ := l.ToString(2)
	listener.SendToAll(message)

	return 0
}

// listener:kick(id)
func kick(l *lua.State) int {
	lua.ArgumentCheck(l, l.IsUserData(1), 1, "expected userdata")
	lua.ArgumentCheck(l, l.IsNumber(2), 2, "expected number")

	listener, ok := l.ToUserData(1).(*listener.Listener)
	if !ok {
		l.PushString("Failed to cast userdata!")
		l.Error()
	}

	id, _ := l.ToNumber(2)
	listener.Kick(int(id))

	return 0
}

// listener:kick_all()
func kickAll(l *lua.State) int {
	lua.ArgumentCheck(l, l.IsUserData(1), 1, "expected userdata")

	listener, ok := l.ToUserData(1).(*listener.Listener)
	if !ok {
		l.PushString("Failed to cast userdata!")
		l.Error()
	}

	listener.KickAll()

	return 0
}
