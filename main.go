package main

import (
	"fmt"
	"os"

	"github.com/kalshiev/gator-go/internal/config"
)

type state struct {
	config *config.Config
}

func main() {

	var currentState state

	cfg, err := config.Read()
	if err != nil {
		os.Exit(1)
	}

	currentState.config = &cfg

	cmds := commands{
		cmdMap: make(map[string]func(*state, command) error),
	}
	cmds.register("login", handlerLogin)

	if len(os.Args) < 2 {
		fmt.Println("Not enough arguments provided")
		os.Exit(1)
	}

	var currentCommand command

	currentCommand.name = os.Args[1]
	currentCommand.argv = os.Args[2:]

	err = cmds.run(&currentState, currentCommand)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
