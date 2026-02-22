package handlers

import (
	"context"
	"fmt"
	"strconv"

	database "github.com/Skyy-Bluu/bootdev-gator/internal/database"
)

const defaultLimit = 2

func HandlerBrowse(s *State, cmd Command, user database.User) error {
	var limit int32

	err := checkIfArgumentPresent(cmd, 1)

	if err != nil {
		limit = defaultLimit
	} else {
		i, err := strconv.ParseInt(cmd.Argurments[0], 0, 32)

		if err != nil {
			return err
		}
		limit = int32(i)
	}

	args := database.GetPostsForUserByUserIDParams{
		UserID: user.ID,
		Limit:  limit,
	}

	fmt.Println()

	posts, err := s.DB.GetPostsForUserByUserID(context.Background(), args)

	if err != nil {
		return err
	}

	for _, post := range posts {
		fmt.Println("Post: ")
		fmt.Println("Title: ", post.Title)
		fmt.Println("Description: ", post.Description)
		fmt.Println("URL: ", post.Url)
	}

	return nil
}
