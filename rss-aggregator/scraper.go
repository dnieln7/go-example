package main

import (
	"context"
	"log"
	"sync"
	"time"

	"github.com/dnieln7/go-examples/rss-aggregator/internal/database"
)

func startScraping(database *database.Queries, concurrency int, timeBetweenRequests time.Duration) {
	log.Printf("Scraping on %d goroutines every %v duration\n", concurrency, timeBetweenRequests)

	ticker := time.NewTicker(timeBetweenRequests)

	// wait timeBetweenRequests then run
	// for c := range ticker.C {

	// run immediatetly
	for ; ; <-ticker.C {
		tbFeeds, err := database.GetNextFeedsToFech(context.Background(), int32(concurrency))

		if err != nil {
			log.Println("Error feching feeds from database: ", err)
			continue
		}

		waitGroup := sync.WaitGroup{}

		for _, tbFeed := range tbFeeds {
			waitGroup.Add(1) // Increases counter by 1

			go scrapeFeed(&waitGroup, database, tbFeed)
		}

		waitGroup.Wait() // Waits until counter is 0
	}
}

func scrapeFeed(waitGroup *sync.WaitGroup, database *database.Queries, tbFeed database.TbFeed) {
	defer waitGroup.Done() // Decreases counter by 1

	_, err := database.MarkFeedAsFetched(context.Background(), tbFeed.ID)

	if err != nil {
		log.Println("MarkFeedAsFetched error: ", err)
		return
	}

	rssFeed, err := urlToFeed(tbFeed.Url)

	if err != nil {
		log.Println("urlToFeed error: ", err)
		return
	}

	// for _, rssItem := range rssFeed.Channel.Items {
		// log.Println("Found post: ", rssItem.Title, " of feed: ", tbFeed.Name)
	// }

	log.Printf("Feed %s collected, %d posts found\n", tbFeed.Name.String, len(rssFeed.Channel.Items))
}
