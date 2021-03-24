package functions

import (
	"github.com/yuin/gopher-lua"
	"github.com/docker/compose-cli/api/context/store"
	"github.com/docker/compose-cli/cli/cmd/compose"
)

var ComposeDefaultArgs []string

func LoadCompose(l *lua.LState) int {
	// Composer table
	composerTable := l.CreateTable(0, 1)

	// Push set default args command
	l.SetField(composerTable, "set_default_args", l.NewFunction(SetDefaultArgs))

	// Metatable
	metatable := l.CreateTable(0, 1)
	l.SetField(metatable, "__index", l.NewFunction(LComposeCommand))
	l.SetMetatable(composerTable, metatable)
	l.Push(composerTable)

	return 1
}

func LComposeCommand(l *lua.LState) int {
	name := l.ToString(2)

	// Composer table
	composerTable := l.NewTable()

	// Metatable
	metatable := l.CreateTable(0, 1)
	l.SetField(metatable, "__call", l.NewFunction(CreateLComposeCommand(name)))
	l.SetMetatable(composerTable, metatable)

	// Set current table as metatable
	l.Push(composerTable)

	return 1
}

func SetDefaultArgs(l *lua.LState) int {
	var args []string
	ok := true
	index := 1

	for ok {
		arg := l.ToString(index)

		if arg != "" {
			args = append(args, arg)
		} else {
			ok = false
		}

		index = index + 1
	}

	ComposeDefaultArgs = args

	return 0
}

func CreateLComposeCommand(name string) func (l *lua.LState) int {
	return func (l *lua.LState) int {
		args := ComposeDefaultArgs
		args = append(args, name)

		ok := true
		index := 2

		for ok {
			arg := l.ToString(index)

			if arg != "" {
				args = append(args, arg)
			} else {
				ok = false
			}

			index = index + 1
		}

		cmd := compose.Command(store.LocalContextType)
		cmd.SetArgs(args)

		err := cmd.Execute()

		if err != nil {
			l.Error(lua.LString(err.Error()), 0)
		}

		return 0
	}
}

func PreloadCompose(l *lua.LState) {
	l.PreloadModule("compose", LoadCompose)
}