package models

import (
	"rss-aggregator/internal/database"
	"time"

	"github.com/google/uuid"
)

type Item struct {
	ID          uuid.UUID `json:"id"`
	Title       *string   `json:"title"`
	Link        *string   `json:"link"`
	Guid        string    `json:"guid"`
	Pubdate     *string   `json:"pubdate"`
	Seeders     *int32    `json:"seeders"`
	Leechers    *int32    `json:"leechers"`
	Downloads   *int32    `json:"downloads"`
	Infohash    *string   `json:"infohash"`
	CategoryID  *string   `json:"category_id"`
	Category    *string   `json:"category"`
	Size        *string   `json:"size"`
	Comments    *int32    `json:"comments"`
	Trusted     *string   `json:"trusted"`
	Remake      *string   `json:"remake"`
	Description *string   `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	ChannelID   int32     `json:"channel_id"`
}

func DatabaseItemToItem(dbItem database.Item) Item {

	var title, link, pubdate, infohash, category_id, category, size, trusted, remake, description *string

	if dbItem.Title.Valid {
		title = &dbItem.Title.String
	}

	if dbItem.Link.Valid {
		link = &dbItem.Link.String
	}

	if dbItem.Pubdate.Valid {
		pubdate = &dbItem.Pubdate.String
	}

	if dbItem.Infohash.Valid {
		infohash = &dbItem.Infohash.String
	}

	if dbItem.CategoryID.Valid {
		category_id = &dbItem.CategoryID.String
	}

	if dbItem.Category.Valid {
		category = &dbItem.Category.String
	}

	if dbItem.Size.Valid {
		size = &dbItem.Size.String
	}

	if dbItem.Trusted.Valid {
		trusted = &dbItem.Trusted.String
	}

	if dbItem.Remake.Valid {
		remake = &dbItem.Remake.String
	}

	if dbItem.Description.Valid {
		description = &dbItem.Description.String
	}

	return Item{
		ID:          dbItem.ID,
		Title:       title,
		Link:        link,
		Guid:        dbItem.Guid,
		Pubdate:     pubdate,
		Seeders:     &dbItem.Seeders.Int32,
		Leechers:    &dbItem.Leechers.Int32,
		Downloads:   &dbItem.Downloads.Int32,
		Infohash:    infohash,
		CategoryID:  category_id,
		Category:    category,
		Size:        size,
		Comments:    &dbItem.Comments.Int32,
		Trusted:     trusted,
		Remake:      remake,
		Description: description,
		CreatedAt:   dbItem.CreatedAt,
		UpdatedAt:   dbItem.UpdatedAt,
		ChannelID:   dbItem.ChannelID,
	}
}

func DatabaseItemsToItems(dbItem []database.Item) []Item {
	var items []Item
	for _, i := range dbItem {
		items = append(items, DatabaseItemToItem(i))
	}
	return items
}
