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

// DatabaseListFeedsRowToFeed maps a row of the paginated listing, which carries
// the pagination total alongside the columns and so comes back as its own
// generated type.
func DatabaseListFeedsRowToFeed(row database.ListFeedsRow) Feed {
	return DatabaseFeedToFeed(database.Feed{
		ID:            row.ID,
		CreatedAt:     row.CreatedAt,
		UpdatedAt:     row.UpdatedAt,
		Url:           row.Url,
		Name:          row.Name,
		LastFetchedAt: row.LastFetchedAt,
	})
}
