package main

import "fmt"

func handlerHelp(s *state, cmd command) error {
	for name, usage := range cmd.help {
		fmt.Printf("command: %s\n", name)
		fmt.Println(usage.usage)
		fmt.Println(usage.description)
		fmt.Println()
	}
	return nil
}
