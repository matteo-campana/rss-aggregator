package rss

import (
	"database/sql"
	"encoding/xml"
	"net/http"
	database "rss-aggregator/internal/database"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type RSS struct {
	XMLName xml.Name `xml:"rss"`
	Version string   `xml:"version,attr"`
	Channel Channel  `xml:"channel"`
}

type Channel struct {
	Title       string `xml:"title"`
	Description string `xml:"description"`
	Link        string `xml:"link"`
	AtomLink    struct {
		Href string `xml:"href,attr"`
	} `xml:"atom\\:link"`
	Items []Item `xml:"item"`
}

type Item struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	GUID        string `xml:"guid"`
	PubDate     string `xml:"pubDate"`
	Seeders     int    `xml:"nyaa\\:seeders"`
	Leechers    int    `xml:"nyaa\\:leechers"`
	Downloads   int    `xml:"nyaa\\:downloads"`
	InfoHash    string `xml:"nyaa\\:infoHash"`
	CategoryID  string `xml:"nyaa\\:categoryId"`
	Category    string `xml:"nyaa\\:category"`
	Size        string `xml:"nyaa\\:size"`
	Comments    int    `xml:"nyaa\\:comments"`
	Trusted     string `xml:"nyaa\\:trusted"`
	Remake      string `xml:"nyaa\\:remake"`
	Description string `xml:"description"`
}

type Resolution string

const (
	Resolution480p  Resolution = "480p"
	Resolution720p  Resolution = "720p"
	Resolution1080p Resolution = "1080p"
	Resolution4K    Resolution = "4K"
)

type Language string

const (
	ENG   Language = "ENG"
	PORBR Language = "POR-BR"
	SPALA Language = "SPA-LA"
	SPA   Language = "SPA"
	ARA   Language = "ARA"
	FRE   Language = "FRE"
	GER   Language = "GER"
	ITA   Language = "ITA"
	RUS   Language = "RUS"
)

type FetchAndParseRSSRequest struct {
	Language   Language
	Resolution Resolution
}

// ParseRSS parses the RSS feed from the given byte slice.
func FetchAndParseRSS(c *gin.Context,
	fetchAndParseRSSRequest *FetchAndParseRSSRequest,
	db *database.Queries) (*database.Channel, *[]database.Item, error) {

	// language Language, resolution Resolution

	url := "https://nyaa.si/?page=rss&c=1_0&f=0&q=" + string(fetchAndParseRSSRequest.Language) +
		"+" + string(fetchAndParseRSSRequest.Resolution)

	rss := &RSS{}
	httpClient := &http.Client{
		Timeout: time.Second * 30,
	}

	resp, err := httpClient.Get(url)

	if err != nil {
		return nil, nil, err
	}

	defer resp.Body.Close()

	err = xml.NewDecoder(resp.Body).Decode(rss)

	if err != nil {
		return nil, nil, err
	}

	// check if the feed already exists

	nyaa_feed, err := db.GetFeedByUrl(c, url)

	if err != nil {
		// create the feed

		nyaa_feed, err = db.CreateFeed(c, database.CreateFeedParams{
			Url:  url,
			Name: "Nyaa",
		})

		if err != nil {
			return nil, nil, err
		}
	}

	// check if the channel already exists

	channel, err := db.GetChannelByTitle(c, rss.Channel.Title)

	if err != nil {
		// create the channel

		channel, err = db.CreateChannel(c, database.CreateChannelParams{
			ID:          uuid.New(),
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
			Title:       rss.Channel.Title,
			Description: sql.NullString{String: rss.Channel.Description, Valid: true},
			Link:        sql.NullString{String: rss.Channel.Link, Valid: true},
			AtomLink:    sql.NullString{String: rss.Channel.AtomLink.Href, Valid: true},
			FeedID:      nyaa_feed.ID,
		})

		if err != nil {
			return nil, nil, err
		}

	}

	// check if the items already exist

	databaseItems := []database.Item{}

	for _, item := range rss.Channel.Items {

		current_item, err := db.GetItemByGuid(c, item.GUID)

		if err != nil {
			created_item, err := db.CreateItem(c, database.CreateItemParams{
				ID:          uuid.New(),
				Title:       sql.NullString{String: item.Title, Valid: true},
				Link:        sql.NullString{String: item.Link, Valid: true},
				Guid:        item.GUID,
				Pubdate:     sql.NullString{String: item.PubDate, Valid: true},
				Seeders:     sql.NullInt32{Int32: int32(item.Seeders), Valid: true},
				Leechers:    sql.NullInt32{Int32: int32(item.Leechers), Valid: true},
				Downloads:   sql.NullInt32{Int32: int32(item.Downloads), Valid: true},
				Infohash:    sql.NullString{String: item.InfoHash, Valid: true},
				CategoryID:  sql.NullString{String: item.CategoryID, Valid: true},
				Category:    sql.NullString{String: item.Category, Valid: true},
				Size:        sql.NullString{String: item.Size, Valid: true},
				Comments:    sql.NullInt32{Int32: int32(item.Comments), Valid: true},
				Trusted:     sql.NullString{String: item.Trusted, Valid: true},
				Remake:      sql.NullString{String: item.Remake, Valid: true},
				Description: sql.NullString{String: item.Description, Valid: true},
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
				ChannelID:   channel.ID,
			})

			if err != nil {
				return nil, nil, err
			}

			databaseItems = append(databaseItems, created_item)

		} else {

			updated_item, err := db.UpdateItem(c, database.UpdateItemParams{
				ID:          current_item.ID,
				Title:       sql.NullString{String: item.Title, Valid: true},
				Link:        sql.NullString{String: item.Link, Valid: true},
				Guid:        item.GUID,
				Pubdate:     sql.NullString{String: item.PubDate, Valid: true},
				Seeders:     sql.NullInt32{Int32: int32(item.Seeders), Valid: true},
				Leechers:    sql.NullInt32{Int32: int32(item.Leechers), Valid: true},
				Downloads:   sql.NullInt32{Int32: int32(item.Downloads), Valid: true},
				Infohash:    sql.NullString{String: item.InfoHash, Valid: true},
				CategoryID:  sql.NullString{String: item.CategoryID, Valid: true},
				Category:    sql.NullString{String: item.Category, Valid: true},
				Size:        sql.NullString{String: item.Size, Valid: true},
				Comments:    sql.NullInt32{Int32: int32(item.Comments), Valid: true},
				Trusted:     sql.NullString{String: item.Trusted, Valid: true},
				Remake:      sql.NullString{String: item.Remake, Valid: true},
				Description: sql.NullString{String: item.Description, Valid: true},
				UpdatedAt:   time.Now(),
				ChannelID:   channel.ID,
			})

			if err != nil {
				return nil, nil, err
			}

			databaseItems = append(databaseItems, updated_item)
		}

	}

	return &channel, &databaseItems, nil
}
