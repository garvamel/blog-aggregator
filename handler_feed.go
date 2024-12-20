package main

import (
	"context"
	"database/sql"
	"encoding/xml"
	"fmt"
	"html"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/giapoldo/blog-aggregator/internal/database"
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

	if l := len(cmd.args); l < 1 {
		fmt.Println("This command requires a time interval")
	}

	time_between_reqs := cmd.args[0]

	req_interval, err := time.ParseDuration(time_between_reqs)
	if err != nil {
		fmt.Println("Not an time interval string (1s, 1m, 1h)")
		return err
	}

	fmt.Printf("Collecting feeds every: %s\n", time_between_reqs)

	ticker := time.NewTicker(req_interval)
	for ; ; <-ticker.C {
		err := scrapeFeeds(s)
		if err != nil {
			fmt.Println(err)
			return err
		}
	}
	// feed, err := fetchFeed(context.Background(), "https://www.wagslane.dev/index.xml")

	// fmt.Println(feed)
	return nil
}

func scrapeFeeds(s *state) error {

	feed, err := s.db.GetNextFeedToFetch(context.Background())

	if err != nil {
		fmt.Println(err)
		return err
	}

	updated_fetched_feed, err := s.db.MarkFetchedFeed(context.Background(), database.MarkFetchedFeedParams{
		ID: feed.ID,
		LastFetchedAt: sql.NullTime{
			Time:  time.Now(),
			Valid: true,
		},
		UpdatedAt: time.Now(),
	})

	if err != nil {
		fmt.Println(err)
		return err
	}

	rss_feed, err := fetchFeed(context.Background(), updated_fetched_feed.Url)
	if err != nil {
		fmt.Println(err)
		return err
	}
	fmt.Println(rss_feed.Channel.Item[0].PubDate)
	for _, fd := range rss_feed.Channel.Item {
		pubTime, err := time.Parse(time.RFC1123Z, fd.PubDate)
		if err != nil {
			fmt.Println(err)
			return err
		}
		_, err = s.db.CreatePost(context.Background(), database.CreatePostParams{
			ID:          uuid.New(),
			CreatedAt:   updated_fetched_feed.CreatedAt,
			UpdatedAt:   updated_fetched_feed.UpdatedAt,
			Title:       sql.NullString{String: fd.Title, Valid: true},
			Description: sql.NullString{String: fd.Description, Valid: true},
			PublishedAt: sql.NullTime{Time: pubTime, Valid: true},
			FeedID:      feed.ID,
		})
		if err == sql.ErrTxDone {
		} else if err != nil {
			fmt.Println(err)
			return nil
		}
	}

	// fmt.Println(rss_feed.Channel.Title)
	// for _, feeds := range rss_feed.Channel.Item {
	// 	fmt.Printf("  * %s\n", feeds.Title)
	// }

	return nil
}

func handlerBrowse(s *state, cmd command) error {

	var limit int
	var err error

	if l := len(cmd.args); l < 1 {
		limit = 2
	} else {
		limit, err = strconv.Atoi(cmd.args[0])
		if err != nil {
			fmt.Println("Conversion")

			fmt.Println(err)
			return err
		}
	}

	posts, err := s.db.GetPosts(context.Background(), int32(limit))
	if err != nil {
		fmt.Println("GetPost")

		fmt.Println(err)
		return err
	}

	for _, post := range posts {
		fmt.Println(post)
	}

	return nil

}

func handlerAddFeed(s *state, cmd command, user database.User) error {

	if l := len(cmd.args); l < 2 {
		return fmt.Errorf("Enter a name and url for the feed")
	}

	name := cmd.args[0]
	url := cmd.args[1]

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

	_, err = s.db.CreateFeedFollows(context.Background(), database.CreateFeedFollowsParams{
		ID:        uuid.New(),
		CreatedAt: feed.CreatedAt,
		UpdatedAt: feed.UpdatedAt,
		UserID:    feed.UserID,
		FeedID:    feed.ID,
	})
	if err != nil {
		fmt.Println("Feed_follow entry creation failed")
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
