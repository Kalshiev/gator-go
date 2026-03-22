package main

import (
	"context"

	"github.com/kalshiev/gator-go/internal/database"
)

func middlewareLoggedIn(handler func(s *state, cmd command, user database.User) error) func(*state, command) error {
	return func(s *state, cmd command) error {
		ctx := context.Background()

		user, err := s.db.GetUser(ctx, s.config.Current_user_name)
		if err != nil {
			return err
		}

		err = handler(s, cmd, user)
		if err != nil {
			return err
		}
		return nil
	}
}
