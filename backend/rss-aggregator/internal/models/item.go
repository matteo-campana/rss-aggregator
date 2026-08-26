package models

import (
	"rss-aggregator/internal/database"
	"time"

	"github.com/google/uuid"
)

type Item struct {
	ID          uuid.UUID  `json:"id"`
	Title       *string    `json:"title"`
	Link        *string    `json:"link"`
	Guid        string     `json:"guid"`
	Pubdate     *string    `json:"pubdate"`
	PublishedAt *time.Time `json:"published_at"`
	Seeders     *int32     `json:"seeders"`
	Leechers    *int32     `json:"leechers"`
	Downloads   *int32     `json:"downloads"`
	Infohash    *string    `json:"infohash"`
	CategoryID  *string    `json:"category_id"`
	Category    *string    `json:"category"`
	Size        *string    `json:"size"`
	Comments    *int32     `json:"comments"`
	Trusted     *string    `json:"trusted"`
	Remake      *string    `json:"remake"`
	Description *string    `json:"description"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	ChannelID   uuid.UUID  `json:"channel_id"`
}

func DatabaseItemToItem(dbItem database.Item) Item {

	var title, link, pubdate, infohash, category_id, category, size, trusted, remake, description *string
	var seeders, leechers, downloads, comments *int32
	var publishedAt *time.Time

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

	if dbItem.Seeders.Valid {
		seeders = &dbItem.Seeders.Int32
	}

	if dbItem.Leechers.Valid {
		leechers = &dbItem.Leechers.Int32
	}

	if dbItem.Downloads.Valid {
		downloads = &dbItem.Downloads.Int32
	}

	if dbItem.Comments.Valid {
		comments = &dbItem.Comments.Int32
	}

	if dbItem.PublishedAt.Valid {
		publishedAt = &dbItem.PublishedAt.Time
	}

	return Item{
		ID:          dbItem.ID,
		Title:       title,
		Link:        link,
		Guid:        dbItem.Guid,
		Pubdate:     pubdate,
		PublishedAt: publishedAt,
		Seeders:     seeders,
		Leechers:    leechers,
		Downloads:   downloads,
		Infohash:    infohash,
		CategoryID:  category_id,
		Category:    category,
		Size:        size,
		Comments:    comments,
		Trusted:     trusted,
		Remake:      remake,
		Description: description,
		CreatedAt:   dbItem.CreatedAt,
		UpdatedAt:   dbItem.UpdatedAt,
		ChannelID:   dbItem.ChannelID,
	}
}

func DatabaseItemsToItems(dbItem []database.Item) []Item {
	items := []Item{}
	for _, i := range dbItem {
		items = append(items, DatabaseItemToItem(i))
	}
	return items
}

// DatabaseListItemsRowToItem maps a row of the paginated listing. The listing
// carries the pagination total alongside the columns, so it comes back as its
// own generated type; the fields themselves are the same, and the conversion
// keeps DatabaseItemToItem as the single definition of the JSON shape.
func DatabaseListItemsRowToItem(row database.ListItemsRow) Item {
	return DatabaseItemToItem(database.Item{
		ID:          row.ID,
		Title:       row.Title,
		Link:        row.Link,
		Guid:        row.Guid,
		Pubdate:     row.Pubdate,
		PublishedAt: row.PublishedAt,
		Seeders:     row.Seeders,
		Leechers:    row.Leechers,
		Downloads:   row.Downloads,
		Infohash:    row.Infohash,
		CategoryID:  row.CategoryID,
		Category:    row.Category,
		Size:        row.Size,
		Comments:    row.Comments,
		Trusted:     row.Trusted,
		Remake:      row.Remake,
		Description: row.Description,
		CreatedAt:   row.CreatedAt,
		UpdatedAt:   row.UpdatedAt,
		ChannelID:   row.ChannelID,
	})
}
