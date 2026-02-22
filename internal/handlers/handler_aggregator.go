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

	timeBetweenReqs := cmd.Argurments[0]

	fmt.Printf("Collecting feeds every %s \n", timeBetweenReqs)

	timeBetweenRequests, err := time.ParseDuration(timeBetweenReqs)

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
