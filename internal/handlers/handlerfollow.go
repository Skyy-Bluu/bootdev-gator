package handlers

import (
	"context"
	"fmt"
	"time"

	database "github.com/Skyy-Bluu/bootdev-gator/internal/database"
	"github.com/google/uuid"
)

func HandlerFollow(s *State, cmd Command, user database.User) error {
	if err := checkIfArgumentPresent(cmd, 1); err != nil {
		return err
	}

	url := cmd.Argurments[0]

	feed, err := s.DB.GetFeedByURL(context.Background(), url)

	if err != nil {
		return err
	}

	arg := database.CreateFeedFollowsParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID:    user.ID,
		FeedID:    feed.ID,
	}

	data, err := s.DB.CreateFeedFollows(context.Background(), arg)

	if err != nil {
		return err
	}

	fmt.Printf("User is: %s, Feed is: %s \n", data.UserName, data.FeedName)
	return nil
}
