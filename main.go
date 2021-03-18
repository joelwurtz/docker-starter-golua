package main

import (
	"fmt"
	"github.com/Shopify/go-lua"
	"github.com/docker/compose-cli/cli/cmd/compose"
	"github.com/docker/compose-cli/api/context/store"
)

func main() {
	s, _ := store.New(".docker")
	s.Create("", "local", "", nil)

	store.WithContextStore(s)

	l := lua.NewState()
	lua.BaseOpen(l)

	lua.SetFunctions(l, []lua.RegistryFunction{{"compose", func(l *lua.State) int {
		ok := true
		args := []string{}
		index := 1

		for ok {
			arg, argok := l.ToString(index)

			if argok {
				args = append(args, arg)
			}

			index = index + 1
			ok = argok
		}

		execCompose(args...)

		return 0
	}}}, 0)

	lua.OpenLibraries(l)

	if err := lua.DoFile(l, "test.lua"); err != nil {
		panic(err)
	}
}

func execCompose(args ...string) {
	cmd := compose.Command("local")
	cmd.SetArgs(args)

	err := cmd.Execute()

	if err != nil {
		fmt.Errorf("error executing compose: %w", err)
	}
}