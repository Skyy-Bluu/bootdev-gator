package handlers

import (
	"context"
	"fmt"
)

func HandlerGetFeeds(s *State, cmd Command) error {
	feeds, err := s.DB.GetFeeds(context.Background())

	if err != nil {
		return err
	}

	for _, feed := range feeds {
		user, err := s.DB.GetUserByID(context.Background(), feed.UserID)

		if err != nil {
			return err
		}

		fmt.Println("Feed:- ")
		fmt.Printf("Name: %s \n", feed.Name)
		fmt.Printf("URL: %s \n", feed.Url)
		fmt.Printf("Name: %s \n", user.Name)
	}

	return nil
}
