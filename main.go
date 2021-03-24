package main

import (
	"github.com/yuin/gopher-lua"
	"github.com/docker/cli/cli/command"
	cliconfig "github.com/docker/cli/cli/config"
	cliflags "github.com/docker/cli/cli/flags"
	"github.com/docker/compose-cli/api/backend"
	"github.com/docker/compose-cli/api/context"
	"github.com/docker/compose-cli/api/context/store"
	"github.com/docker/compose-cli/local"
	"github.com/spf13/cobra"
	"jolicode.com/docker-starter/functions"
)

func main() {
	rootCmd := &cobra.Command{
		Use:	"docker-starter",
		Short:	"Docker start allow to manage project / boostrap",
	}

	functions.LState = loadLua()

	for _, task := range functions.TaskRegistry {
		rootCmd.AddCommand(task.Command)
	}

	rootCmd.Execute()
}

func loadLua() *lua.LState {
	configDir := ".docker"

	s, _ := store.New(configDir)
	store.WithContextStore(s)
	context.WithCurrentContext(store.LocalContextType)
	s.Create(store.LocalContextType, store.LocalContextType, "local context", nil)

	service, err := createBackend(configDir)

	if err != nil {
		panic(err)
	}

	backend.WithBackend(service)

	l := lua.NewState()

	//lua.OpenLibraries(l)

	functions.PreloadCompose(l)
	functions.PreloadTask(l)

	if err := l.DoFile("bootstrap.lua"); err != nil {
		panic(err)
	}

	return l
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