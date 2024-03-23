package libs

import (
	"fmt"

	"github.com/Shopify/go-lua"
)

// package.loaded[name] = <stack top>
func CachePackage(l *lua.State, name string) {
	l.Global("package")
	l.Field(-1, "loaded")
	l.PushValue(-3)
	l.SetField(-2, name)
	l.Pop(3)
}

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
