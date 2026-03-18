package main

import (
	"errors"
	"fmt"
)

type command struct {
	name string
	argv []string
}

type commands struct {
	cmdMap map[string]func(*state, command) error
}

func (c *commands) run(s *state, cmd command) error {
	if c.cmdMap[cmd.name] == nil {
		return fmt.Errorf("No command registered with %s", cmd.name)
	}

	return c.cmdMap[cmd.name](s, cmd)
}

func (c *commands) register(name string, f func(*state, command) error) {
	c.cmdMap[name] = f
}

func handlerLogin(s *state, cmd command) error {
	if len(cmd.argv) == 0 {
		return errors.New("No username found")
	}

	err := s.config.SetUser(cmd.argv[0])
	if err != nil {
		return errors.New("Could not set user")
	}

	fmt.Printf("User %s set", cmd.argv[0])

	return nil
}
