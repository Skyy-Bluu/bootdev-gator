package main

import (
	"context"
	"database/sql"
	"encoding/xml"
	"fmt"
	"html"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	_ "github.com/lib/pq"

	config "github.com/Skyy-Bluu/bootdev-gator/internal/config"
	database "github.com/Skyy-Bluu/bootdev-gator/internal/database"
	rss "github.com/Skyy-Bluu/bootdev-gator/internal/rss"
)

type state struct {
	config *config.Config
	db     *database.Queries
}

type command struct {
	name       string
	argurments []string
}

type commands struct {
	commandsHandler map[string]func(s *state, cmd command) error
}

const defaultLimit = 2

func (c *commands) run(s *state, cmd command) error {
	runner, ok := c.commandsHandler[cmd.name]
	if !ok {
		return fmt.Errorf("Command %s does not exist", cmd.name)
	}

	err := runner(s, cmd)

	if err != nil {
		return err
	}

	return nil
}

func (c *commands) register(name string, f func(s *state, cmd command) error) {
	c.commandsHandler[name] = f
	//fmt.Printf("Command %s registered", name)
}

func fetchFeed(ctx context.Context, feedURL string) (*rss.RSSFeed, error) {
	request, err := http.NewRequestWithContext(ctx, http.MethodGet, feedURL, nil)

	if err != nil {
		return nil, err
	}

	client := http.Client{}

	request.Header.Set("User-Agent", "gator")

	response, err := client.Do(request)

	if err != nil {
		return nil, err
	}

	data, err := io.ReadAll(response.Body)

	if err != nil {
		return nil, err
	}

	var rssFeed rss.RSSFeed

	if err = xml.Unmarshal(data, &rssFeed); err != nil {
		return nil, err
	}

	rssFeed.Channel.Title = html.UnescapeString(rssFeed.Channel.Title)
	rssFeed.Channel.Description = html.UnescapeString(rssFeed.Channel.Description)
	for _, item := range rssFeed.Channel.Item {
		item.Title = html.UnescapeString(item.Title)
		item.Description = html.UnescapeString(item.Description)
	}

	return &rssFeed, nil
}

func scrapeFeeds(s *state) error {

	feed, err := s.db.GetNextFeedToFetch(context.Background())

	fmt.Println("Found feed: ", feed)

	if err != nil {
		return err
	}

	markFeedArgs := database.MarkFeedFetchedByIDParams{
		ID: feed.ID,
		LastFetchedAt: sql.NullTime{
			Time:  time.Now(),
			Valid: true,
		},
	}

	if err = s.db.MarkFeedFetchedByID(context.Background(), markFeedArgs); err != nil {
		return err
	}

	rssFeed, err := fetchFeed(context.Background(), feed.Url)

	if err != nil {
		return err
	}

	for _, item := range rssFeed.Channel.Item {
		//fmt.Println("Got item: ", item)
		publishedDateTime, err := time.Parse("Mon, 02 Jan 2006 15:04:05 -0700", item.PubDate)

		if err != nil {
			return err
		}

		var itemDescription sql.NullString

		if item.Description == "" {
			itemDescription.String = ""
			itemDescription.Valid = false

		} else {
			itemDescription.String = item.Description
			itemDescription.Valid = true
		}

		addPostArgs := database.CreatePostParams{
			ID:          uuid.New(),
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
			Title:       item.Title,
			Url:         item.Link,
			Description: itemDescription,
			PublishedAt: publishedDateTime,
			FeedID:      feed.ID,
		}

		//fmt.Println("Creating post with the args: ", addPostArgs)

		post, err := s.db.CreatePost(context.Background(), addPostArgs)

		if err != nil && !strings.Contains(err.Error(), "duplicate key value violates unique constraint") {
			return err
		}

		fmt.Println("Post created: ", post)
	}

	return nil
}

func handlerBrowse(s *state, cmd command, user database.User) error {
	var limit int32

	err := checkIfArgumentPresent(cmd, 1)

	if err != nil {
		limit = defaultLimit
	} else {
		i, err := strconv.ParseInt(cmd.argurments[0], 0, 32)

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

	posts, err := s.db.GetPostsForUserByUserID(context.Background(), args)

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

func handlerFollowing(s *state, cmd command, user database.User) error {
	feedsFollowed, err := s.db.GetFeedFollowsByUser(context.Background(), user.ID)

	if err != nil {
		return err
	}

	for _, userAndFeed := range feedsFollowed {
		fmt.Printf("User %s is following: \n", userAndFeed.UserName)
		fmt.Printf("- %s \n", userAndFeed.FeedName)
	}

	return nil
}

func handlerFollow(s *state, cmd command, user database.User) error {
	if err := checkIfArgumentPresent(cmd, 1); err != nil {
		return err
	}

	url := cmd.argurments[0]

	feed, err := s.db.GetFeedByURL(context.Background(), url)

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

	data, err := s.db.CreateFeedFollows(context.Background(), arg)

	if err != nil {
		return err
	}

	fmt.Printf("User is: %s, Feed is: %s \n", data.UserName, data.FeedName)
	return nil
}

func handlerUnfollow(s *state, cmd command, user database.User) error {
	err := checkIfArgumentPresent(cmd, 1)

	if err != nil {
		return nil
	}

	url := cmd.argurments[0]

	feed, err := s.db.GetFeedByURL(context.Background(), url)

	if err != nil {
		return err
	}

	args := database.DeleteFeedFollowEntryByUserIDAndFeedIDParams{
		UserID: user.ID,
		FeedID: feed.ID,
	}

	if err = s.db.DeleteFeedFollowEntryByUserIDAndFeedID(context.Background(), args); err != nil {
		return err
	}

	fmt.Println("Unfollowed feed!")

	return nil
}

func handlerGetFeeds(s *state, cmd command) error {
	feeds, err := s.db.GetFeeds(context.Background())

	if err != nil {
		return err
	}

	for _, feed := range feeds {
		user, err := s.db.GetUserByID(context.Background(), feed.UserID)

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

func handlerAddFeed(s *state, cmd command, user database.User) error {

	if err := checkIfArgumentPresent(cmd, 2); err != nil {
		return err
	}

	feedID := uuid.New()

	argsCreateFeed := database.CreateFeedParams{
		ID:        feedID,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      cmd.argurments[0],
		Url:       cmd.argurments[1],
		UserID:    user.ID,
	}

	argsCreateFeedFollow := database.CreateFeedFollowsParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID:    user.ID,
		FeedID:    feedID,
	}

	feed, err := s.db.CreateFeed(context.Background(), argsCreateFeed)

	if err != nil {
		return err
	}

	data, err := s.db.CreateFeedFollows(context.Background(), argsCreateFeedFollow)

	if err != nil {
		return err
	}

	fmt.Println("Feed created: ", feed)
	fmt.Println("Follow feed created for user: ", data)

	return nil
}

func handlerAggregator(time_between_reqs string) error {

	fmt.Printf("Collecting feeds every %s \n", time_between_reqs)

	timeBetweenRequests, err := time.ParseDuration(time_between_reqs)

	if err != nil {
		return err
	}

	ticker := time.NewTicker(timeBetweenRequests)

	for ; ; <-ticker.C {
		fmt.Println("Scraping....")
		err := scrapeFeeds(&c_state)
		if err != nil {
			return err
		}
	}
}

func handlerUsers(s *state, cmd command) error {
	users, err := s.db.GetUsers(context.Background())

	if err != nil {
		return err
	}

	for _, user := range users {
		if user == s.config.CurrentUser {
			fmt.Printf("* %s (current)\n", user)
		} else {
			fmt.Printf("* %s\n", user)
		}
	}

	return nil
}

func handlerReset(s *state, cmd command) error {
	err := s.db.DeleteUsers(context.Background())

	if err != nil {
		return err
	}

	fmt.Println("Users table cleared!")

	return nil
}

func handlerRegister(s *state, cmd command) error {

	if err := checkIfArgumentPresent(cmd, 1); err != nil {
		return err
	}

	name := cmd.argurments[0]

	args := database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      name,
	}

	user, err := s.db.CreateUser(context.Background(), args)

	if err != nil {
		return err
	}

	err = s.config.SetUser(name)

	if err != nil {
		return err
	}

	fmt.Println("User added: ", user)

	return nil
}

func handelerLogin(s *state, cmd command) error {

	if err := checkIfArgumentPresent(cmd, 1); err != nil {
		return err
	}

	name := cmd.argurments[0]

	_, err := s.db.GetUserByName(context.Background(), name)

	if err != nil {
		return err
	}

	if err := s.config.SetUser(name); err != nil {
		return err
	}

	fmt.Println("User has been set!")

	return nil
}

var args = os.Args
var c_state = state{}
var cmd command

const g_time_between_requests = "10s"

func main() {

	config, err := config.Read()

	if err != nil {
		log.Fatalf("Error reading config file:  %v", err)
	}

	db, err := sql.Open("postgres", config.DB_URL)

	c_state.db = database.New(db)

	c_state.config = &config

	if len(args) < 2 {
		log.Fatalln("Not enough arguments. Exiting program")
		return
	}
	cmd = command{
		name:       args[1],
		argurments: args[2:],
	}

	commands := commands{
		commandsHandler: make(map[string]func(*state, command) error),
	}

	commands.register("login", handelerLogin)
	commands.register("register", handlerRegister)
	commands.register("reset", handlerReset)
	commands.register("users", handlerUsers)
	commands.register("feeds", handlerGetFeeds)
	commands.register("agg", middlewarRSSFeed(handlerAggregator))
	commands.register("addfeed", middlewareLoggedIn(handlerAddFeed))
	commands.register("follow", middlewareLoggedIn(handlerFollow))
	commands.register("following", middlewareLoggedIn(handlerFollowing))
	commands.register("unfollow", middlewareLoggedIn(handlerUnfollow))
	commands.register("browse", middlewareLoggedIn(handlerBrowse))

	if err = commands.run(&c_state, cmd); err != nil {
		log.Fatalf("[Error] %v", err)
	}

	//fmt.Println(config)
}

func checkIfArgumentPresent(cmd command, numberOfArguments int) error {
	if len(cmd.argurments) == 0 {
		return fmt.Errorf(" Expected %v argumen(s)): Arguments cannot be empty", numberOfArguments)
	} else if len(cmd.argurments) != numberOfArguments {
		return fmt.Errorf("Expected %v argument(s)", numberOfArguments)
	}

	return nil
}

// func registerCommands(commandsHandlerMap map[string]func(*state, command) error, cmd *commands) {

// 	for command, handler := range commandsHandlerMap {
// 		cmd.register(command, handler)
// 	}
// }

func middlewareLoggedIn(handler func(s *state, cmd command, user database.User) error) func(*state, command) error {

	return func(s *state, c command) error {
		user, err := c_state.db.GetUserByName(context.Background(), c_state.config.CurrentUser)

		if err != nil {
			log.Fatalf("[Error] %v", err)
		}

		return handler(&c_state, cmd, user)
	}
}

func middlewarRSSFeed(handler func(time_between_requests string) error) func(*state, command) error {
	return func(s *state, c command) error {
		return handler(g_time_between_requests)
	}
}
