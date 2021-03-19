package main

import (
	"fmt"
	"github.com/Shopify/go-lua"
	"github.com/docker/cli/cli/command"
	cliconfig "github.com/docker/cli/cli/config"
	cliflags "github.com/docker/cli/cli/flags"
	"github.com/docker/compose-cli/api/backend"
	"github.com/docker/compose-cli/api/context/store"
	"github.com/docker/compose-cli/local"
	"jolicode.com/docker-starter/functions"
)

func main() {
	configDir := ".docker"

	s, _ := store.New(configDir)
	store.WithContextStore(s)

	service, err := createBackend(configDir)

	if err != nil {
		panic(err)
	}

	backend.WithBackend(service)

	l := lua.NewState()

	lua.OpenLibraries(l)
	functions.LoadLibrary(l)

	var name string
	var commandDef interface{}
	var commandCall interface{}

	l.Register("create_task", func(l *lua.State) int {

		name, _ = l.ToString(1)
		commandDef = l.ToValue(2)
		commandCall = l.ToValue(3)

		return 0
	})

	if err := lua.DoFile(l, "test.lua"); err != nil {
		panic(err)
	}

	fmt.Println(name)

	if commandDef != nil {
		l.PushLightUserData(commandDef)
		l.Call(0, 0)
	}

	if commandCall != nil {
		l.PushLightUserData(commandCall)
		l.Call(0, 0)
	}
}

func createBackend(configDir string) (backend.Service, error) {
	configFile, err := cliconfig.Load(configDir)

	if err != nil {
		return nil, err
	}

	options := cliflags.CommonOptions{}
	apiClient, err := command.NewAPIClientFromFlags(&options, configFile)

	if err != nil {
		return nil, err
	}

	return local.NewService(apiClient), nil
}