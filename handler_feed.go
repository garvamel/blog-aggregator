package main

import (
	"context"
	"database/sql"
	"encoding/xml"
	"fmt"
	"html"
	"io"
	"net/http"
	"time"

	"github.com/garvamel/blog-aggregator/internal/database"
	"github.com/google/uuid"
)

type RSSFeed struct {
	Channel struct {
		Title       string    `xml:"title"`
		Link        string    `xml:"link"`
		Description string    `xml:"description"`
		Item        []RSSItem `xml:"item"`
	} `xml:"channel"`
}

type RSSItem struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description string `xml:"description"`
	PubDate     string `xml:"pubDate"`
}

func fetchFeed(ctx context.Context, feedURL string) (*RSSFeed, error) {

	feed := RSSFeed{}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, feedURL, nil)

	if err != nil {
		fmt.Println(err)
		return &feed, err
	}

	req.Header.Set("User-Agent", "gator")

	client := http.Client{}
	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return &feed, err
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return &feed, err
	}

	err = xml.Unmarshal(body, &feed)
	if err != nil {
		fmt.Println(err)
		return &feed, err
	}

	feed.Channel.Title = html.UnescapeString(feed.Channel.Title)
	feed.Channel.Description = html.UnescapeString(feed.Channel.Description)

	for i, rssitem := range feed.Channel.Item {
		feed.Channel.Item[i].Title = html.UnescapeString(rssitem.Title)
		feed.Channel.Item[i].Description = html.UnescapeString(rssitem.Description)

	}

	return &feed, nil

}

func handlerAgg(s *state, cmd command) error {

	feed, err := fetchFeed(context.Background(), "https://www.wagslane.dev/index.xml")

	if err != nil {
		fmt.Println(err)
		return err
	}

	fmt.Println(feed)
	return nil
}

func handlerAddFeed(s *state, cmd command) error {

	if l := len(cmd.args); l < 2 {
		return fmt.Errorf("Enter a name and url for the feed")
	}

	username := s.cfg.CurrentUserName
	name := cmd.args[0]
	url := cmd.args[1]

	user, err := s.db.GetUser(context.Background(), username)
	if err == sql.ErrNoRows {
		// Current user should exist, create it if not
		fmt.Println("Current user not set, add one first")
		return err
	} else if err != nil {
		fmt.Println(err)
		return err
	}

	if user.Name != username {
		return fmt.Errorf("Database returned a different user")
	}

	feed, err := s.db.CreateFeed(context.Background(), database.CreateFeedParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      name,
		Url:       url,
		UserID:    user.ID,
	})
	if err != nil {
		fmt.Println("Feed entry creation failed")
		return err
	}

	fmt.Println(feed)

	return nil
}

func handlerFeeds(s *state, cmd command) error {

	if l := len(cmd.args); l > 0 {
		fmt.Println("This command doesn't need/use arguments")
	}

	feeds, err := s.db.GetFeeds(context.Background())
	if err == sql.ErrNoRows {
		fmt.Println("Database is emtpy")
		return err

	} else if err != nil {
		fmt.Println(err)
		return err
	}

	for _, feed := range feeds {
		// if feed exists then user exists, no need to check errors.
		user, _ := s.db.GetUserByID(context.Background(), feed.UserID)
		fmt.Println(feed.Name)
		fmt.Println(feed.Url)
		fmt.Println(user.Name)
	}
	return nil
}
