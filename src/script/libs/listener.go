package libs

import (
	"server/src/listener"

	"github.com/Shopify/go-lua"
)

func RegisterListenerLib(l *lua.State, listener *listener.Listener) {
	l.PushUserData(listener)

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

	CachePackage(l, "listener")
}

func listenerFromStack(l *lua.State) *listener.Listener {
	listener, ok := l.ToUserData(1).(*listener.Listener)
	if !ok {
		l.PushString("Failed to cast userdata!")
		l.Error()
	}
	return listener
}

// listener:get_total_sent(id)
func totalSent(l *lua.State) int {
	lua.ArgumentCheck(l, l.IsUserData(1), 1, "expected userdata")
	lua.ArgumentCheck(l, l.IsNumber(2), 2, "expected number")

	id, _ := l.ToNumber(2)
	listener := listenerFromStack(l)

	l.PushInteger(listener.TotalSent(int(id)))

	return 1
}

// listener:get_total_received(id)
func totalReceived(l *lua.State) int {
	lua.ArgumentCheck(l, l.IsUserData(1), 1, "expected userdata")
	lua.ArgumentCheck(l, l.IsNumber(2), 2, "expected number")

	id, _ := l.ToNumber(2)
	listener := listenerFromStack(l)

	l.PushInteger(listener.TotalReceived(int(id)))

	return 1
}

// listener:send_to(id, message)
func sendTo(l *lua.State) int {
	lua.ArgumentCheck(l, l.IsUserData(1), 1, "expected userdata")
	lua.ArgumentCheck(l, l.IsNumber(2), 2, "expected number")
	lua.ArgumentCheck(l, l.IsString(3), 3, "expected string")

	listener := listenerFromStack(l)
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

	listener := listenerFromStack(l)
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

	listener := listenerFromStack(l)
	message, _ := l.ToString(2)
	listener.SendToAll(message)

	return 0
}

// listener:kick(id)
func kick(l *lua.State) int {
	lua.ArgumentCheck(l, l.IsUserData(1), 1, "expected userdata")
	lua.ArgumentCheck(l, l.IsNumber(2), 2, "expected number")

	listener := listenerFromStack(l)
	id, _ := l.ToNumber(2)
	listener.Kick(int(id))

	return 0
}

// listener:kick_all()
func kickAll(l *lua.State) int {
	lua.ArgumentCheck(l, l.IsUserData(1), 1, "expected userdata")

	listener := listenerFromStack(l)
	listener.KickAll()

	return 0
}
