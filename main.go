package main

import (
	"fmt"
	"github.com/Shopify/go-lua"
	"github.com/docker/cli/cli/command"
	cliconfig "github.com/docker/cli/cli/config"
	cliflags "github.com/docker/cli/cli/flags"
	"github.com/docker/compose-cli/cli/cmd/compose"
	"github.com/docker/compose-cli/api/backend"
	"github.com/docker/compose-cli/api/context/store"
	"github.com/docker/compose-cli/api/errdefs"
	apicontext "github.com/docker/compose-cli/api/context"
	cliopts "github.com/docker/compose-cli/cli/options"
	"github.com/docker/compose-cli/local"
)

func main() {
	var opts cliopts.GlobalOpts

	configDir := ".docker"
	currentContext := "default"

	apicontext.WithCurrentContext(currentContext)

	s, _ := store.New(configDir)
	store.WithContextStore(s)

	ctype := store.LocalContextType
	cc, _ := s.Get(currentContext)

	if cc != nil {
		ctype = cc.Type()
	}

	service, err := getBackend(ctype, configDir, opts)

	if err != nil {
		panic(err)
	}

	backend.WithBackend(service)

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
		//panic(err)
	}
}

func getBackend(ctype string, configDir string, opts cliopts.GlobalOpts) (backend.Service, error) {
	switch ctype {
	case store.DefaultContextType, store.LocalContextType:
		configFile, err := cliconfig.Load(configDir)
		if err != nil {
			return nil, err
		}
		options := cliflags.CommonOptions{
			Context:  opts.Context,
			Debug:    opts.Debug,
			Hosts:    opts.Hosts,
			LogLevel: opts.LogLevel,
		}

		if opts.TLSVerify {
			options.TLS = opts.TLS
			options.TLSVerify = opts.TLSVerify
			options.TLSOptions = opts.TLSOptions
		}
		apiClient, err := command.NewAPIClientFromFlags(&options, configFile)
		if err != nil {
			return nil, err
		}
		return local.NewService(apiClient), nil
	}
	service, err := backend.Get(ctype)
	if errdefs.IsNotFoundError(err) {
		return service, nil
	}
	return service, err
}

func execCompose(args ...string) {
	cmd := compose.Command("local")
	cmd.SetArgs(args)

	err := cmd.Execute()

	if err != nil {
		fmt.Println(err)
//		panic(err)
	}
}