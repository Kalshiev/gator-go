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

	ctx := context.Background()

	getUser, err := s.db.GetUser(ctx, cmd.argv[0])
	if err != nil {
		return err
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

	ctx := context.Background()

	insertUser, err := s.db.CreateUser(ctx, database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      cmd.argv[0],
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

func handlerAgg(s *state, cmd command) error {
	url := "https://www.wagslane.dev/index.xml"

	ctx := context.Background()

	feed, err := fetchFeed(ctx, url)
	if err != nil {
		return err
	}

	fmt.Print(feed)
	return nil
}

func handlerAddFeed(s *state, cmd command) error {
	if len(cmd.argv) < 2 {
		return errors.New("No name and url provided")
	}

	ctx := context.Background()

	fName := cmd.argv[0]
	fUrl := cmd.argv[1]

	currentUser, err := s.db.GetUser(ctx, s.config.Current_user_name)
	if err != nil {
		return err
	}

	feed, err := s.db.CreateFeed(ctx, database.CreateFeedParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      fName,
		Url:       fUrl,
		UserID:    currentUser.ID,
	})
	if err != nil {
		return err
	}

	fmt.Printf("Added %s to user %s's feed\n", feed.Name, currentUser.Name)
	fmt.Print(feed)

	return nil

}

func handlerFeeds(s *state, cmd command) error {
	ctx := context.Background()

	feeds, err := s.db.ListFeeds(ctx)
	if err != nil {
		return nil
	}

	for i := 0; i < len(feeds); i++ {
		userName, err := s.db.GetUserByID(ctx, feeds[i].UserID)
		if err != nil {
			return err
		}

		fmt.Printf("Feed: %d\n", i)
		fmt.Printf("name: %s\n", feeds[i].Name)
		fmt.Printf("URL: %s\n", feeds[i].Url)
		fmt.Printf("User: %s\n\n", userName.Name)
	}

	return nil
}
