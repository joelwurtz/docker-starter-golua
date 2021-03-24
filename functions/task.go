package functions

import (
	"github.com/yuin/gopher-lua"
	"github.com/spf13/cobra"
)

type Task struct {
	Name              string
	definitionClosure *lua.LFunction
	callClosure       *lua.LFunction
	Command           *cobra.Command
}

type Argument struct {
	name      	 string
	optional 	 bool
	defaultValue string
}

var TaskRegistry []*Task
var LState       *lua.LState
var currentTask  *Task

func CreateTask(l *lua.LState) int {
	name := l.CheckString(1)
	commandDef := l.CheckFunction(2)
	commandCall := l.CheckFunction(3)

	task := &Task{
		Name:              name,
		definitionClosure: commandDef,
		callClosure:       commandCall,
		Command:           &cobra.Command{
			Use: name,
		},
	}

	currentTask = task
	l.Push(commandDef)
	l.Call(0, 0)
	currentTask = nil

	task.Command.Run = task.Run
	TaskRegistry = append(TaskRegistry, task)

	return 0
}

func SetShortDescription(l *lua.LState) int {
	if currentTask == nil {
		l.Error(lua.LString("not in context of command definition"), 0)

		return 0
	}

	desc := l.CheckString(1)
	currentTask.Command.Short = desc

	return 0
}

func SetLongDescription(l *lua.LState) int {
	if currentTask == nil {
		l.Error(lua.LString("not in context of command definition"), 0)

		return 0
	}

	desc := l.CheckString(1)
	currentTask.Command.Long = desc

	return 0
}

func AddArgument(l *lua.LState) int {
	if currentTask == nil {
		l.Error(lua.LString("not in context of command definition"), 0)

		return 0
	}

	return 0
}

func AddOption(l *lua.LState) int {
	if currentTask == nil {
		l.Error(lua.LString("not in context of command definition"), 0)

		return 0
	}

	return 0
}

var taskApi = map[string]lua.LGFunction{
	"create": CreateTask,
	"set_short_description": SetShortDescription,
	"set_long_description": SetLongDescription,
	"add_argument": AddArgument,
	"add_option": AddOption,
}

func LoadTask(l *lua.LState) int {
	api := l.NewTable()
	l.SetFuncs(api, taskApi)
	l.Push(api)

	return 1
}

func PreloadTask(l *lua.LState) {
	l.PreloadModule("task", LoadTask)
}

func (t *Task) Call(l *lua.LState) {

}

func (t *Task) Run(cmd *cobra.Command, args []string) {
	// @TODO Parse args

	// @TODO IF args error show help of command

	// Run command
	if LState != nil {
		LState.Push(t.callClosure)

		// @TODO Pass args to function
		LState.Call(0, 0)
	}
}