package models

import (
	"database/sql"
	"testing"

	"rss-aggregator/internal/database"

	"github.com/google/uuid"
)

func TestDatabaseChannelToChannel(t *testing.T) {
	feedID := uuid.New()

	channel := DatabaseChannelToChannel(database.Channel{
		ID:          uuid.New(),
		Title:       "Nyaa - Home - Torrent File RSS",
		Description: sql.NullString{String: "RSS Feed for Home", Valid: true},
		AtomLink:    sql.NullString{String: "https://nyaa.si/?page=rss", Valid: true},
		FeedID:      feedID,
	})

	if channel.Title == nil || *channel.Title != "Nyaa - Home - Torrent File RSS" {
		t.Errorf("title = %v, want the channel title", channel.Title)
	}

	// link was NULL in the row and must not be reported as an empty string.
	if channel.Link != nil {
		t.Errorf("link = %q, want nil", *channel.Link)
	}

	if channel.AtomLink == nil || *channel.AtomLink != "https://nyaa.si/?page=rss" {
		t.Errorf("atom link = %v, want the feed url", channel.AtomLink)
	}

	if channel.FeedID == nil || *channel.FeedID != feedID {
		t.Errorf("feed id = %v, want %v", channel.FeedID, feedID)
	}
}

func TestDatabaseChannelsToChannelsEmpty(t *testing.T) {
	channels := DatabaseChannelsToChannels(nil)

	if channels == nil {
		t.Fatal("DatabaseChannelsToChannels(nil) = nil, want an empty slice so the API serialises [] instead of null")
	}
}
