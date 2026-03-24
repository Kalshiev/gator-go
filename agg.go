package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/kalshiev/gator-go/internal/database"
)

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
		publishDate, err := time.Parse(time.RFC1123, item.PubDate)
		if err != nil {
			fmt.Println("Could not parse date")
			continue
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
			fmt.Println(err)
			continue
		}

		fmt.Printf("Post: %s from blog %s created!\n \n", post.Title, feed.Channel.Title)
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
