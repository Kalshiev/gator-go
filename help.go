package main

import "fmt"

func handlerHelp(s *state, cmd command) error {
	if len(cmd.argv) == 0 {
		for name := range cmd.help {
			fmt.Printf("command: %s\n", name)
		}
		return nil
	} else {
		if _, exists := cmd.help[cmd.argv[0]]; !exists {
			return fmt.Errorf("Command does not exist")
		}

		fmt.Printf("command: %s\n", cmd.argv[0])
		fmt.Println(cmd.help[cmd.argv[0]].usage)
		fmt.Println(cmd.help[cmd.argv[0]].description)
		return nil
	}
}
