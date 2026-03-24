package main

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/kalshiev/gator-go/internal/database"
)

type command struct {
	name string
	argv []string
	help map[string]usage
}

type commands struct {
	cmdMap   map[string]func(*state, command) error
	cmdUsage map[string]usage
}

type usage struct {
	usage       string
	description string
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

func (c *commands) registerUsage(n string, u string, d string) {
	c.cmdUsage[n] = usage{
		usage:       u,
		description: d,
	}
}

func handlerLogin(s *state, cmd command) error {
	if len(cmd.argv) == 0 {
		return errors.New("No username found")
	}

	var passwd string
	fmt.Printf("Please enter your password: ")
	fmt.Scanln(&passwd)

	ctx := context.Background()

	getUser, err := s.db.GetUser(ctx, cmd.argv[0])
	if err != nil {
		return err
	}

	match := VerifyPassword(passwd, getUser.PasswordHash)
	if !match {
		return fmt.Errorf("Incorrect Password!")
	}

	err = s.config.SetUser(getUser.Name)
	if err != nil {
		return errors.New("Could not set user")
	}

	fmt.Printf("User %s set \n", cmd.argv[0])

	return nil
}

func handlerRegister(s *state, cmd command) error {
	if len(cmd.argv) == 0 {
		return errors.New("No username submitted")
	}

	var passwd string
	fmt.Printf("Please Enter your password: ")
	fmt.Scanln(&passwd)

	hash, _ := HashPassword(passwd)

	ctx := context.Background()

	insertUser, err := s.db.CreateUser(ctx, database.CreateUserParams{
		ID:           uuid.New(),
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
		Name:         cmd.argv[0],
		PasswordHash: hash,
	})
	if err != nil {
		return err
	}

	err = s.config.SetUser(insertUser.Name)
	if err != nil {
		return err
	}

	fmt.Println("User created succesfully!")

	return nil
}

func handlerReset(s *state, cmd command) error {
	ctx := context.Background()

	err := s.db.ResetUsers(ctx)
	if err != nil {
		return err
	}

	fmt.Println("All users deleted!")

	return nil
}

func handlerGetUsers(s *state, cmd command) error {
	ctx := context.Background()

	users, err := s.db.GetUsers(ctx)
	if err != nil {
		return err
	}

	for i := 0; i < len(users); i++ {
		if users[i].Name == s.config.Current_user_name {
			fmt.Printf("%s (current) \n", users[i].Name)
			continue
		}
		fmt.Println(users[i].Name)
	}

	return nil
}
