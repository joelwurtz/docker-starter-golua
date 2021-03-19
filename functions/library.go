package functions

import "github.com/Shopify/go-lua"

func LoadLibrary(l *lua.State) {
	lua.SubTable(l, lua.RegistryIndex, "_PRELOAD")
	l.PushGoFunction(LComposeRequire)
	l.SetField(-2, "compose")
	l.Pop(1)
}