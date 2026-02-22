package handlers

import (
	"context"
	"fmt"
	"time"

	database "github.com/Skyy-Bluu/bootdev-gator/internal/database"
	"github.com/google/uuid"
)

func HandlerRegister(s *State, cmd Command) error {

	if err := checkIfArgumentPresent(cmd, 1); err != nil {
		return err
	}

	name := cmd.Argurments[0]

	args := database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      name,
	}

	user, err := s.DB.CreateUser(context.Background(), args)

	if err != nil {
		return err
	}

	err = s.Config.SetUser(name)

	if err != nil {
		return err
	}

	fmt.Println("User added: ", user)

	return nil
}
