package functions

import (
	"github.com/yuin/gopher-lua"
	"os"
)

func GetCwd(l *lua.LState) int {
	path, err := os.Getwd()

	if err != nil {
		l.Push(lua.LNil)

		return 1
	}

	l.Push(lua.LString(path))

	return 1
}

func AddToOsLibrary(l* lua.LState) {

}