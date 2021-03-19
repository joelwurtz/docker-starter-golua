package functions

import (
	"github.com/Shopify/go-lua"
	"github.com/docker/compose-cli/api/context/store"
	"github.com/docker/compose-cli/cli/cmd/compose"
)

func LComposeRequire(l *lua.State) int {
	// Composer table
	l.NewTable()

	// Metatable
	l.CreateTable(0, 1)
	l.PushGoFunction(LComposeCommand)

	// Push function to call
	l.SetField(-2, "__index")

	// Set current table as metatable
	l.SetMetaTable(-2)

	return 1
}

func LComposeCommand(l *lua.State) int {
	name, _ := l.ToString(2)

	// Composer table
	l.NewTable()

	// Metatable
	l.CreateTable(0, 1)
	l.PushGoFunction(CreateLComposeCommand(name))

	// Push function to call
	l.SetField(-2, "__call")

	// Set current table as metatable
	l.SetMetaTable(-2)

	return 1
}

func CreateLComposeCommand(name string) func (l *lua.State) int {
	return func (l *lua.State) int {
		args := []string{name}
		ok := true
		index := 2

		for ok {
			arg, argok := l.ToString(index)

			if argok {
				args = append(args, arg)
			}

			index = index + 1
			ok = argok
		}

		cmd := compose.Command(store.LocalContextType)
		cmd.SetArgs(args)

		err := cmd.Execute()

		if err != nil {
			l.PushString(err.Error())
			l.Error()
		}

		return 0
	}
}

func AddComposeLibrary(l *lua.State) {
	lua.SubTable(l, lua.RegistryIndex, "_PRELOAD")
	l.PushGoFunction(LComposeRequire)
	l.SetField(-2, "compose")
	l.Pop(1)
}