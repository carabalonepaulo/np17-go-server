package libs

import (
	"log"

	"github.com/Shopify/go-lua"
)

func RegisterLogLib(l *lua.State) {
	l.CreateTable(0, 0)

	l.PushGoFunction(print)
	l.SetField(-2, "print")

	l.PushGoFunction(fatal)
	l.SetField(-2, "fatal")

	CachePackage(l, "log")
}

func print(l *lua.State) int {
	lua.ArgumentCheck(l, l.IsString(1), 1, "expected string")

	line, _ := l.ToString(1)
	log.Println(line)

	return 0
}

func fatal(l *lua.State) int {
	lua.ArgumentCheck(l, l.IsString(1), 1, "expected string")

	line, _ := l.ToString(1)
	log.Fatalln(line)

	return 0
}
