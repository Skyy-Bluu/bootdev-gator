package handlers

import (
	"context"
	"fmt"
)

func HandlerReset(s *State, cmd Command) error {
	err := s.DB.DeleteUsers(context.Background())

	if err != nil {
		return err
	}

	fmt.Println("Users table cleared!")

	return nil
}
