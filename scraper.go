package main

import (
	"context"
	"log"
	"sync"
	"time"

	"github.com/Yom3n/rssagg/internal/database"
)

func startScrapping(db *database.Queries,
	concurrency int,
	timeBetweenRequests time.Duration,
) {
	log.Printf("Scrapping on %v gouroutines every %s duration", concurrency, timeBetweenRequests)
	ticker := time.NewTicker(timeBetweenRequests)

	for ; ; <-ticker.C {
		feeds, err := db.GetNextFeedsToFetch(
			context.Background(),
			int32(concurrency),
		)
		if err != nil {
			log.Println("Error fetching feeds: ", err)
			continue
		}
		wg := sync.WaitGroup{}
		for _, feed := range feeds {
			wg.Add(1)
			go scrapeFeed(db, &wg, feed)
		}
		wg.Wait()
	}
}

func scrapeFeed(db *database.Queries, wg *sync.WaitGroup, feed database.Feed) {
	defer wg.Done()
	_, err := db.MarkFeedAsFetched(context.Background(), feed.ID)
	if err != nil {
		log.Print("Error marking feed as fetched: ", err)
		return
	}

	rssFeed, err := urlToFeed(feed.Url)
	if err != nil {
		log.Println("Error fetching rssFeed: ", err)
		return
	}

	for _, item := range rssFeed.Channel.Item {
		log.Println("Found post ", item.Title, "on feed ", feed.Name)
	}
	log.Printf("Feed %s collected, %v posts found", feed.Name, len(rssFeed.Channel.Item))

}
