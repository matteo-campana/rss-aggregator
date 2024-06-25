package models

import (
	"rss-aggregator/internal/database"
	"time"

	"github.com/google/uuid"
)

type Channel struct {
	ID          uuid.UUID  `json:"id"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	Title       *string    `json:"title"`
	Description *string    `json:"description"`
	Link        *string    `json:"link"`
	AtomLink    *string    `json:"atom_link"`
	FeedID      *uuid.UUID `json:"feed_id"`
}

func DatabaseChannelToChannel(dbChannel database.Channel) Channel {

	var description, link, atomLink *string

	if dbChannel.Description.Valid {
		description = &dbChannel.Description.String
	}

	if dbChannel.Link.Valid {
		link = &dbChannel.Link.String
	}

	if dbChannel.AtomLink.Valid {
		atomLink = &dbChannel.AtomLink.String
	}

	return Channel{
		ID:          dbChannel.ID,
		CreatedAt:   dbChannel.CreatedAt,
		UpdatedAt:   dbChannel.UpdatedAt,
		Title:       &dbChannel.Title,
		Description: description,
		Link:        link,
		AtomLink:    atomLink,
		FeedID:      &dbChannel.FeedID,
	}
}

func DatabaseChannelsToChannels(dbChannels []database.Channel) []Channel {
	var channels []Channel
	for _, dbChannel := range dbChannels {
		channels = append(channels, DatabaseChannelToChannel(dbChannel))
	}
	return channels
}
