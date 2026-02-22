package handlers

import (
	"context"
	"database/sql"
	"encoding/xml"
	"html"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"

	"fmt"

	_ "github.com/lib/pq"

	config "github.com/Skyy-Bluu/bootdev-gator/internal/config"
	database "github.com/Skyy-Bluu/bootdev-gator/internal/database"
	rss "github.com/Skyy-Bluu/bootdev-gator/internal/rss"
)

type State struct {
	Config *config.Config
	DB     *database.Queries
}

type Command struct {
	Name       string
	Argurments []string
}

func checkIfArgumentPresent(cmd Command, numberOfArguments int) error {
	if len(cmd.Argurments) == 0 {
		return fmt.Errorf(" Expected %v argumen(s)): Arguments cannot be empty", numberOfArguments)
	} else if len(cmd.Argurments) != numberOfArguments {
		return fmt.Errorf("Expected %v argument(s)", numberOfArguments)
	}

	return nil
}

func fetchFeed(ctx context.Context, feedURL string) (*rss.RSSFeed, error) {
	request, err := http.NewRequestWithContext(ctx, http.MethodGet, feedURL, nil)

	if err != nil {
		return nil, err
	}

	client := http.Client{}

	request.Header.Set("User-Agent", "gator")

	response, err := client.Do(request)

	if err != nil {
		return nil, err
	}

	data, err := io.ReadAll(response.Body)

	if err != nil {
		return nil, err
	}

	var rssFeed rss.RSSFeed

	if err = xml.Unmarshal(data, &rssFeed); err != nil {
		return nil, err
	}

	rssFeed.Channel.Title = html.UnescapeString(rssFeed.Channel.Title)
	rssFeed.Channel.Description = html.UnescapeString(rssFeed.Channel.Description)
	for _, item := range rssFeed.Channel.Item {
		item.Title = html.UnescapeString(item.Title)
		item.Description = html.UnescapeString(item.Description)
	}

	return &rssFeed, nil
}

func scrapeFeeds(s *State) error {

	feed, err := s.DB.GetNextFeedToFetch(context.Background())

	if err != nil {
		return err
	}

	markFeedArgs := database.MarkFeedFetchedByIDParams{
		ID: feed.ID,
		LastFetchedAt: sql.NullTime{
			Time:  time.Now(),
			Valid: true,
		},
	}

	if err = s.DB.MarkFeedFetchedByID(context.Background(), markFeedArgs); err != nil {
		return err
	}

	rssFeed, err := fetchFeed(context.Background(), feed.Url)

	if err != nil {
		return err
	}

	for _, item := range rssFeed.Channel.Item {
		publishedDateTime, err := time.Parse("Mon, 02 Jan 2006 15:04:05 -0700", item.PubDate)

		if err != nil {
			return err
		}

		var itemDescription sql.NullString

		if item.Description == "" {
			itemDescription.String = ""
			itemDescription.Valid = false

		} else {
			itemDescription.String = item.Description
			itemDescription.Valid = true
		}

		addPostArgs := database.CreatePostParams{
			ID:          uuid.New(),
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
			Title:       item.Title,
			Url:         item.Link,
			Description: itemDescription,
			PublishedAt: publishedDateTime,
			FeedID:      feed.ID,
		}

		post, err := s.DB.CreatePost(context.Background(), addPostArgs)

		if err != nil && !strings.Contains(err.Error(), "duplicate key value violates unique constraint") {
			return err
		}

		if post.ID.String() != "00000000-0000-0000-0000-000000000000" {
			fmt.Println("Post created: ", post)
		}
	}

	return nil
}
