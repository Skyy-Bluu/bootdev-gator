package handlers

import (
	"context"
	"fmt"
	"time"

	database "github.com/Skyy-Bluu/bootdev-gator/internal/database"
	"github.com/google/uuid"
)

func HandlerAddFeed(s *State, cmd Command, user database.User) error {

	if err := checkIfArgumentPresent(cmd, 2); err != nil {
		return err
	}

	feedID := uuid.New()

	argsCreateFeed := database.CreateFeedParams{
		ID:        feedID,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      cmd.Argurments[0],
		Url:       cmd.Argurments[1],
		UserID:    user.ID,
	}

	argsCreateFeedFollow := database.CreateFeedFollowsParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID:    user.ID,
		FeedID:    feedID,
	}

	feed, err := s.DB.CreateFeed(context.Background(), argsCreateFeed)

	if err != nil {
		return err
	}

	data, err := s.DB.CreateFeedFollows(context.Background(), argsCreateFeedFollow)

	if err != nil {
		return err
	}

	fmt.Println("Feed created: ", feed)
	fmt.Println("Follow feed created for user: ", data)

	return nil
}
