package functions

import (
	"github.com/Shopify/go-lua"
	"github.com/spf13/cobra"
)

type Task struct {
	Name              string
	definitionClosure interface{}
	callClosure       interface{}
	Command           *cobra.Command
}

type Argument struct {
	name      	 string
	optional 	 bool
	defaultValue string
}

var TaskRegistry []*Task
var LState       *lua.State
var currentTask  *Task

func CreateTask(l *lua.State) int {
	if !l.IsString(1) {
		l.PushString("first argument must be a string")
		l.Error()

		return 0
	}

	if !l.IsFunction(2) {
		l.PushString("second argument must be a function")
		l.Error()

		return 0
	}

	if !l.IsFunction(3) {
		l.PushString("third argument must be a function")
		l.Error()

		return 0
	}

	name, _ := l.ToString(1)
	commandDef := l.ToValue(2)
	commandCall := l.ToValue(3)

	task := &Task{
		Name:              name,
		definitionClosure: commandDef,
		callClosure:       commandCall,
		Command:           &cobra.Command{
			Use: name,
		},
	}

	currentTask = task
	l.PushValue(2)
	l.Call(0, 0)
	currentTask = nil

	task.Command.Run = task.Run
	TaskRegistry = append(TaskRegistry, task)

	return 0
}

func SetShortDescription(l *lua.State) int {
	if currentTask == nil {
		l.PushString("not in context of command definition")
		l.Error()

		return 0
	}

	if !l.IsString(1) {
		l.PushString("first argument must be a string")
		l.Error()

		return 0
	}

	desc, _ := l.ToString(1)
	currentTask.Command.Short = desc

	return 0
}

func SetLongDescription(l *lua.State) int {
	if currentTask == nil {
		l.PushString("not in context of command definition")
		l.Error()

		return 0
	}

	if !l.IsString(1) {
		l.PushString("wrong argument")
		l.Error()

		return 0
	}

	desc, _ := l.ToString(1)
	currentTask.Command.Long = desc

	return 0
}

var taskLibrary = []lua.RegistryFunction{
	{"create", CreateTask},
	{"set_short_description", SetShortDescription},
	{"set_long_description", SetLongDescription},
}

func AddTaskLibrary(l *lua.State) {
	lua.SubTable(l, lua.RegistryIndex, "_PRELOAD")

	l.PushGoFunction(TaskLibraryOpen)
	l.SetField(-2, "task")
	l.Pop(1)
}

func TaskLibraryOpen(l *lua.State) int {
	lua.NewLibrary(l, taskLibrary)
	return 1
}

func (t *Task) Call(l *lua.State) {

}

func (t *Task) Run(cmd *cobra.Command, args []string) {
	// @TODO Parse args

	// @TODO IF args error show help of command

	// Run command
	if LState != nil {
		LState.PushLightUserData(t.callClosure)

		// @TODO Pass args to function
		LState.Call(0, 0)
	}
}