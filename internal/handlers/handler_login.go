package handlers

import (
	"context"
	"fmt"
)

func HandlerLogin(s *State, cmd Command) error {

	if err := checkIfArgumentPresent(cmd, 1); err != nil {
		return err
	}

	name := cmd.Argurments[0]

	_, err := s.DB.GetUserByName(context.Background(), name)

	if err != nil {
		return err
	}

	if err := s.Config.SetUser(name); err != nil {
		return err
	}

	fmt.Println("User has been set!")

	return nil
}
