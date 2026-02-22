package main

import (
	"context"
	"database/sql"

	"fmt"

	"log"

	"os"

	_ "github.com/lib/pq"

	config "github.com/Skyy-Bluu/bootdev-gator/internal/config"
	database "github.com/Skyy-Bluu/bootdev-gator/internal/database"
	handlers "github.com/Skyy-Bluu/bootdev-gator/internal/handlers"
)

type state = handlers.State

type command = handlers.Command

type commands struct {
	commandsHandler map[string]func(s *state, cmd command) error
}

func (c *commands) run(s *state, cmd command) error {
	runner, ok := c.commandsHandler[cmd.Name]
	if !ok {
		return fmt.Errorf("Command %s does not exist", cmd.Name)
	}

	err := runner(s, cmd)

	if err != nil {
		return err
	}

	return nil
}

func (c *commands) register(name string, f func(s *state, cmd command) error) {
	c.commandsHandler[name] = f
}

var args = os.Args
var c_state = state{}
var cmd command

func main() {

	config, err := config.Read()

	if err != nil {
		log.Fatalf("Error reading config file:  %v", err)
	}

	db, err := sql.Open("postgres", config.DB_URL)

	c_state.DB = database.New(db)

	c_state.Config = &config

	if len(args) < 2 {
		log.Fatalln("Not enough arguments. Exiting program")
		return
	}
	cmd = command{
		Name:       args[1],
		Argurments: args[2:],
	}

	commands := commands{
		commandsHandler: make(map[string]func(*state, command) error),
	}

	commands.register("login", handlers.HandlerLogin)
	commands.register("register", handlers.HandlerRegister)
	commands.register("reset", handlers.HandlerReset)
	commands.register("users", handlers.HandlerUsers)
	commands.register("feeds", handlers.HandlerGetFeeds)
	commands.register("agg", handlers.HandlerAggregator)
	commands.register("addfeed", middlewareLoggedIn(handlers.HandlerAddFeed))
	commands.register("follow", middlewareLoggedIn(handlers.HandlerFollow))
	commands.register("following", middlewareLoggedIn(handlers.HandlerFollowing))
	commands.register("unfollow", middlewareLoggedIn(handlers.HandlerUnfollow))
	commands.register("browse", middlewareLoggedIn(handlers.HandlerBrowse))

	if err = commands.run(&c_state, cmd); err != nil {
		log.Fatalf("[Error] %v", err)
	}
}

func middlewareLoggedIn(handler func(s *handlers.State, cmd handlers.Command, user database.User) error) func(*state, command) error {

	return func(s *state, c command) error {
		user, err := c_state.DB.GetUserByName(context.Background(), c_state.Config.CurrentUser)

		if err != nil {
			log.Fatalf("[Error] %v", err)
		}

		return handler(&c_state, cmd, user)
	}
}
