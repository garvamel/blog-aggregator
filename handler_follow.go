package main

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/giapoldo/blog-aggregator/internal/database"
	"github.com/google/uuid"
)

func handlerFollow(s *state, cmd command, user database.User) error {

	if l := len(cmd.args); l < 1 {
		return fmt.Errorf("This command requieres a url")
	}

	url := cmd.args[0]

	feed, err := s.db.GetSingleFeedByURL(context.Background(), url)
	if err != nil {
		fmt.Println("Feed not found")
		return err
	}

	feed_follows, err := s.db.CreateFeedFollows(context.Background(), database.CreateFeedFollowsParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID:    user.ID,
		FeedID:    feed.ID,
	})
	if err != nil {
		fmt.Println("Feed_follows entry creation failed")
		return err
	}
	userID, err := s.db.GetUserByID(context.Background(), feed_follows.UserID)
	if err != nil {
		fmt.Println("User from feed_follow not found")
		return err
	}

	fmt.Println(userID, feed.Name)
	return nil
}

func handlerFollowing(s *state, cmd command, user database.User) error {

	if l := len(cmd.args); l > 0 {
		return fmt.Errorf("This command takes no arguments")
	}

	feed_follows, err := s.db.GetFeedFollowsForUserID(context.Background(), user.ID)
	if err == sql.ErrNoRows {
		fmt.Println("User doesn't follow any feeds")
		return err
	} else if err != nil {
		fmt.Println(err)
		return err
	}

	if len(feed_follows) == 0 {
		return nil
	}

	fmt.Printf("User %s is following:\n", feed_follows[0].UserName)
	for _, row := range feed_follows {
		fmt.Println(row.FeedName)
	}
	return nil
}

func handlerUnfollow(s *state, cmd command, user database.User) error {

	if l := len(cmd.args); l > 1 {
		return fmt.Errorf("This command takes 1 argument")
	}

	url := cmd.args[0]

	feed, err := s.db.GetSingleFeedByURL(context.Background(), url)
	if err != nil {
		fmt.Println(err)
		return err
	}

	err = s.db.DeleteFeedFollowRecord(context.Background(), database.DeleteFeedFollowRecordParams{
		UserID: user.ID,
		FeedID: feed.ID,
	})

	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil

}
