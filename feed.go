package main

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/kalshiev/gator-go/internal/database"
)

func handlerAddFeed(s *state, cmd command, user database.User) error {
	if len(cmd.argv) < 2 {
		return errors.New("No name and url provided")
	}

	ctx := context.Background()

	fName := cmd.argv[0]
	fUrl := cmd.argv[1]

	feed, err := s.db.CreateFeed(ctx, database.CreateFeedParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      fName,
		Url:       fUrl,
		UserID:    user.ID,
	})
	if err != nil {
		return err
	}

	feedFollow, err := s.db.CreateFeedFollow(ctx, database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		FeedID:    feed.ID,
		UserID:    feed.UserID,
	})

	fmt.Printf("Added %s to user %s's feed\n", feed.Name, user.Name)
	fmt.Println(feed)
	fmt.Println(feedFollow)

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

func handlerFollow(s *state, cmd command, user database.User) error {
	if len(cmd.argv) < 1 {
		return errors.New("No url provided")
	}

	ctx := context.Background()

	feed, err := s.db.GetFeedByURL(ctx, cmd.argv[0])
	if err != nil {
		return err
	}

	follow, err := s.db.CreateFeedFollow(ctx, database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		FeedID:    feed.ID,
		UserID:    user.ID,
	})
	if err != nil {
		return err
	}

	fmt.Printf("Feed %s added to %s's feed!\n", feed.Url, user.Name)
	fmt.Print(follow)

	return nil
}

func handlerFollowing(s *state, cmd command, user database.User) error {
	ctx := context.Background()

	following, err := s.db.GetFeedFollowForUser(ctx, user.ID)
	if err != nil {
		return err
	}

	fmt.Printf("%s's feeds: \n", user.Name)
	for i := 0; i < len(following); i++ {
		fmt.Println(following[i].FeedName)
	}

	return nil
}

func handlerUnfollow(s *state, cmd command, user database.User) error {
	ctx := context.Background()

	feed, err := s.db.GetFeedByURL(ctx, cmd.argv[0])
	if err != nil {
		return err
	}

	err = s.db.DeleteFeedFollow(ctx, database.DeleteFeedFollowParams{
		FeedID: feed.ID,
		UserID: user.ID,
	})
	if err != nil {
		return err
	}

	return nil
}

func handlerBrowse(s *state, cmd command, user database.User) error {
	limit := int32(2)

	if len(cmd.argv) > 0 {
		parsedLimit, err := strconv.ParseInt(cmd.argv[0], 10, 32)
		if err != nil {
			return errors.New("Invalid limit")
		}
		limit = int32(parsedLimit)
	}

	ctx := context.Background()

	posts, err := s.db.GetPostsForUser(ctx, database.GetPostsForUserParams{
		UserID: user.ID,
		Limit:  limit,
	})
	if err != nil {
		return err
	}

	fmt.Printf("Posts for user %s:\n", user.Name)
	for _, post := range posts {
		fmt.Println(post.Title)
		fmt.Println(post.PublishedAt)
		fmt.Println(post.Url)
		fmt.Println()
	}

	return nil
}
