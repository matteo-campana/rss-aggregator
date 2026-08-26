package models

import (
	"database/sql"
	"testing"
	"time"

	"rss-aggregator/internal/database"

	"github.com/google/uuid"
)

func TestDatabaseItemToItemNullFields(t *testing.T) {
	item := DatabaseItemToItem(database.Item{
		ID:        uuid.New(),
		Guid:      "https://nyaa.si/view/1",
		ChannelID: uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	})

	if item.Title != nil {
		t.Errorf("title = %v, want nil", *item.Title)
	}

	// A NULL count must stay null instead of being reported as a real zero.
	if item.Seeders != nil {
		t.Errorf("seeders = %v, want nil", *item.Seeders)
	}

	if item.Leechers != nil {
		t.Errorf("leechers = %v, want nil", *item.Leechers)
	}

	if item.Downloads != nil {
		t.Errorf("downloads = %v, want nil", *item.Downloads)
	}

	if item.Comments != nil {
		t.Errorf("comments = %v, want nil", *item.Comments)
	}
}

func TestDatabaseItemToItemValidFields(t *testing.T) {
	item := DatabaseItemToItem(database.Item{
		ID:        uuid.New(),
		Guid:      "https://nyaa.si/view/1",
		Title:     sql.NullString{String: "an episode", Valid: true},
		Size:      sql.NullString{String: "1.4 GiB", Valid: true},
		Seeders:   sql.NullInt32{Int32: 1234, Valid: true},
		Leechers:  sql.NullInt32{Int32: 56, Valid: true},
		Downloads: sql.NullInt32{Int32: 7890, Valid: true},
		Comments:  sql.NullInt32{Int32: 0, Valid: true},
		ChannelID: uuid.New(),
	})

	if item.Title == nil || *item.Title != "an episode" {
		t.Errorf("title = %v, want %q", item.Title, "an episode")
	}

	if item.Size == nil || *item.Size != "1.4 GiB" {
		t.Errorf("size = %v, want %q", item.Size, "1.4 GiB")
	}

	if item.Seeders == nil || *item.Seeders != 1234 {
		t.Errorf("seeders = %v, want 1234", item.Seeders)
	}

	// Zero is a legitimate value and must survive as a non-nil pointer.
	if item.Comments == nil || *item.Comments != 0 {
		t.Errorf("comments = %v, want 0", item.Comments)
	}
}

func TestDatabaseItemsToItemsEmpty(t *testing.T) {
	items := DatabaseItemsToItems(nil)

	if items == nil {
		t.Fatal("DatabaseItemsToItems(nil) = nil, want an empty slice so the API serialises [] instead of null")
	}

	if len(items) != 0 {
		t.Errorf("len = %d, want 0", len(items))
	}
}
