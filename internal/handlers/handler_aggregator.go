package handlers

import (
	"fmt"
	"time"
)

func HandlerAggregator(s *State, cmd Command) error {

	err := checkIfArgumentPresent(cmd, 1)

	if err != nil {
		return err
	}

	time_between_reqs := cmd.Argurments[0]

	fmt.Printf("Collecting feeds every %s \n", time_between_reqs)

	timeBetweenRequests, err := time.ParseDuration(time_between_reqs)

	if err != nil {
		return err
	}

	ticker := time.NewTicker(timeBetweenRequests)

	for ; ; <-ticker.C {
		fmt.Println("Scraping....")
		err := scrapeFeeds(s)
		if err != nil {
			return err
		}
	}
}
