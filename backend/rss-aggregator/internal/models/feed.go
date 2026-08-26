package models

import (
	"rss-aggregator/internal/database"
	"time"

	"github.com/google/uuid"
)

type Feed struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Url       string    `json:"url"`
	Name      string    `json:"name"`
	// LastFetchedAt is nil until the scraper refreshes the feed for the first time.
	LastFetchedAt *time.Time `json:"last_fetched_at"`
}

func DatabaseFeedToFeed(dbFeed database.Feed) Feed {
	var lastFetchedAt *time.Time

	if dbFeed.LastFetchedAt.Valid {
		lastFetchedAt = &dbFeed.LastFetchedAt.Time
	}

	return Feed{
		ID:        dbFeed.ID,
		CreatedAt: dbFeed.CreatedAt,
		UpdatedAt: dbFeed.UpdatedAt,
		Url:       dbFeed.Url,
		Name:      dbFeed.Name,

		LastFetchedAt: lastFetchedAt,
	}
}

func DatabaseFeedsToFeeds(dbFeeds []database.Feed) []Feed {
	feeds := []Feed{}

	for _, dbFeed := range dbFeeds {
		feeds = append(feeds, DatabaseFeedToFeed(dbFeed))
	}

	return feeds
}
