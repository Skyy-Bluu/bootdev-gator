package handlers

import (
	"context"
	"fmt"

	database "github.com/Skyy-Bluu/bootdev-gator/internal/database"
)

func HandlerFollowing(s *State, cmd Command, user database.User) error {
	feedsFollowed, err := s.DB.GetFeedFollowsByUser(context.Background(), user.ID)

	if err != nil {
		return err
	}

	for _, userAndFeed := range feedsFollowed {
		fmt.Printf("User %s is following: \n", userAndFeed.UserName)
		fmt.Printf("- %s \n", userAndFeed.FeedName)
	}

	return nil
}
