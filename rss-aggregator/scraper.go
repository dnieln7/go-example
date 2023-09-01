package main

import (
	"context"
	"database/sql"
	"log"
	"strings"
	"sync"
	"time"

	"github.com/dnieln7/go-examples/rss-aggregator/internal/database"
	"github.com/google/uuid"
)

func startScraping(db *database.Queries, concurrency int, timeBetweenRequests time.Duration) {
	log.Printf("Scraping on %d goroutines every %v duration\n", concurrency, timeBetweenRequests)

	ticker := time.NewTicker(timeBetweenRequests)

	// wait timeBetweenRequests then run
	// for c := range ticker.C {

	// run immediatetly
	for ; ; <-ticker.C {
		tbFeeds, err := db.GetNextFeedsToFech(context.Background(), int32(concurrency))

		if err != nil {
			log.Println("Error feching feeds from database: ", err)
			continue
		}

		waitGroup := sync.WaitGroup{}

		for _, tbFeed := range tbFeeds {
			waitGroup.Add(1) // Increases counter by 1

			go scrapeFeed(&waitGroup, db, tbFeed)
		}

		waitGroup.Wait() // Waits until counter is 0
	}
}

func scrapeFeed(waitGroup *sync.WaitGroup, db *database.Queries, tbFeed database.TbFeed) {
	defer waitGroup.Done() // Decreases counter by 1

	_, err := db.MarkFeedAsFetched(context.Background(), tbFeed.ID)

	if err != nil {
		log.Println("MarkFeedAsFetched error: ", err)
		return
	}

	rssFeed, err := urlToFeed(tbFeed.Url)

	if err != nil {
		log.Println("urlToFeed error: ", err)
		return
	}

	for _, rssItem := range rssFeed.Channel.Items {
		description := sql.NullString{}

		if rssItem.Description != "" {
			description.String = rssItem.Description
			description.Valid = true
		}

		publishedAt, err := time.Parse(time.RFC1123Z, rssItem.PubDate)

		if err != nil {
			log.Println("Could not parse date: ", err)
			continue
		}

		_, err = db.CreatePost(context.Background(), database.CreatePostParams{
			ID:          uuid.New(),
			Title:       rssItem.Title,
			Description: description,
			Url:         rssItem.Link,
			PublishedAt: publishedAt,
			FeedID:      tbFeed.ID,
			CreatedAt:   time.Now().UTC(),
			UpdatedAt:   time.Now().UTC(),
		})

		if err != nil {
			if strings.Contains(err.Error(), "duplicate key") {
				continue
			}

			log.Println("Failed to create post: ", err)
		}
	}

	log.Printf("Feed %s collected, %d posts found\n", tbFeed.Name.String, len(rssFeed.Channel.Items))
}
