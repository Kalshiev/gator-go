package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strconv"
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
	if len(cmd.argv) == 0 {
		return errors.New("Provide a time interval")
	}
	time_between_reqs := cmd.argv[0]

	if time_between_reqs == "" {
		return errors.New("No time interval provided")
	}

	interval, err := time.ParseDuration(time_between_reqs)
	if err != nil {
		return err
	}

	fmt.Printf("Collecting feeds every %v\n", interval)

	ticker := time.NewTicker(interval)
	for ; ; <-ticker.C {
		err := scrapeFeeds(s)
		if err != nil {
			return err
		}
	}

}

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

func scrapeFeeds(s *state) error {
	ctx := context.Background()

	nextFeed, err := s.db.GetNextFeedToFetch(ctx)
	if err != nil {
		return err
	}

	err = s.db.MarkFeedFetched(ctx, database.MarkFeedFetchedParams{
		ID:            nextFeed.ID,
		LastFetchedAt: sql.NullTime{Time: time.Now(), Valid: true},
	})
	if err != nil {
		return err
	}

	feed, err := fetchFeed(ctx, nextFeed.Url)
	if err != nil {
		return err
	}

	numPosts := len(feed.Channel.Item)
	fmt.Printf("Saving %d posts from: %s\n", numPosts, feed.Channel.Title)

	for _, item := range feed.Channel.Item {
		publishDate, err := time.Parse("2006-01-02", item.PubDate)
		if err != nil {
			publishDate = time.Now()
		}

		post, err := s.db.CreatePost(ctx, database.CreatePostParams{
			ID:          uuid.New(),
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
			Title:       item.Title,
			Url:         item.Link,
			Description: sql.NullString{String: item.Description, Valid: true},
			PublishedAt: publishDate,
			FeedID:      nextFeed.ID,
		})
		if err != nil {
			return err
		}

		fmt.Printf("Post: %s from blog %s created!\n \n", post.Title, feed.Channel.Title)
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
