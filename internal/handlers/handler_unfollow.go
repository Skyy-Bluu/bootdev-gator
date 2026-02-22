package handlers

import (
	"context"
	"fmt"

	database "github.com/Skyy-Bluu/bootdev-gator/internal/database"
)

func HandlerUnfollow(s *State, cmd Command, user database.User) error {
	err := checkIfArgumentPresent(cmd, 1)

	if err != nil {
		return nil
	}

	url := cmd.Argurments[0]

	feed, err := s.DB.GetFeedByURL(context.Background(), url)

	if err != nil {
		return err
	}

	args := database.DeleteFeedFollowEntryByUserIDAndFeedIDParams{
		UserID: user.ID,
		FeedID: feed.ID,
	}

	if err = s.DB.DeleteFeedFollowEntryByUserIDAndFeedID(context.Background(), args); err != nil {
		return err
	}

	fmt.Println("Unfollowed feed!")

	return nil
}
