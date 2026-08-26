package rss

import (
	"context"
	"database/sql"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"

	database "rss-aggregator/internal/database"

	"github.com/google/uuid"
)

// nyaaBaseURL is the RSS entry point of nyaa.si.
const nyaaBaseURL = "https://nyaa.si/"

// defaultTimeout bounds a single feed fetch.
const defaultTimeout = 30 * time.Second

type RSS struct {
	XMLName xml.Name `xml:"rss"`
	Version string   `xml:"version,attr"`
	Channel Channel  `xml:"channel"`
}

// Channel mirrors the <channel> element. AtomLink is declared before Link so
// that <atom:link> is matched by the namespaced field first: a field without a
// namespace matches any local name, Link included.
type Channel struct {
	Title       string `xml:"title"`
	Description string `xml:"description"`
	AtomLink    struct {
		Href string `xml:"href,attr"`
	} `xml:"http://www.w3.org/2005/Atom link"`
	Link  string `xml:"link"`
	Items []Item `xml:"item"`
}

// Item mirrors an <item> element. The nyaa:* fields are matched on their local
// name only, which keeps the mapping working whatever prefix or namespace URI
// nyaa.si serves them under.
type Item struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	GUID        string `xml:"guid"`
	PubDate     string `xml:"pubDate"`
	Seeders     int    `xml:"seeders"`
	Leechers    int    `xml:"leechers"`
	Downloads   int    `xml:"downloads"`
	InfoHash    string `xml:"infoHash"`
	CategoryID  string `xml:"categoryId"`
	Category    string `xml:"category"`
	Size        string `xml:"size"`
	Comments    int    `xml:"comments"`
	Trusted     string `xml:"trusted"`
	Remake      string `xml:"remake"`
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

// Languages and Resolutions list the accepted values, so the API can validate a
// request and report what it would have accepted.
var Languages = []Language{ENG, PORBR, SPALA, SPA, ARA, FRE, GER, ITA, RUS}

var Resolutions = []Resolution{Resolution480p, Resolution720p, Resolution1080p, Resolution4K}

// ParseLanguage matches a language, case-insensitively. An empty value yields
// the default, ENG.
func ParseLanguage(value string) (Language, error) {
	if value == "" {
		return ENG, nil
	}

	for _, language := range Languages {
		if strings.EqualFold(value, string(language)) {
			return language, nil
		}
	}

	return "", fmt.Errorf("unknown language %q, expected one of: %s", value, join(Languages))
}

// ParseResolution matches a resolution, case-insensitively. An empty value
// yields the default, 1080p.
func ParseResolution(value string) (Resolution, error) {
	if value == "" {
		return Resolution1080p, nil
	}

	for _, resolution := range Resolutions {
		if strings.EqualFold(value, string(resolution)) {
			return resolution, nil
		}
	}

	return "", fmt.Errorf("unknown resolution %q, expected one of: %s", value, join(Resolutions))
}

func join[T ~string](values []T) string {
	parts := make([]string, 0, len(values))

	for _, value := range values {
		parts = append(parts, string(value))
	}

	return strings.Join(parts, ", ")
}

type FetchAndParseRSSRequest struct {
	Language   Language
	Resolution Resolution
}

// NyaaFeedURL builds the nyaa.si RSS URL for the requested language and resolution.
func NyaaFeedURL(req *FetchAndParseRSSRequest) string {
	query := url.Values{
		"page": {"rss"},
		"c":    {"1_0"},
		"f":    {"0"},
		"q":    {string(req.Language) + " " + string(req.Resolution)},
	}

	return nyaaBaseURL + "?" + query.Encode()
}

// ParseRSS decodes an RSS document.
func ParseRSS(r io.Reader) (*RSS, error) {
	rss := &RSS{}

	if err := xml.NewDecoder(r).Decode(rss); err != nil {
		return nil, fmt.Errorf("decoding rss feed: %w", err)
	}

	return rss, nil
}

// FetchRSS downloads and decodes the RSS document served at feedURL.
func FetchRSS(ctx context.Context, client *http.Client, feedURL string) (*RSS, error) {
	if client == nil {
		client = &http.Client{Timeout: defaultTimeout}
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, feedURL, nil)
	if err != nil {
		return nil, fmt.Errorf("building request for %s: %w", feedURL, err)
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("fetching %s: %w", feedURL, err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("fetching %s: unexpected status %s", feedURL, resp.Status)
	}

	return ParseRSS(resp.Body)
}

// SyncFeed stores the parsed feed: it upserts the feed, its channel and every
// item, keyed on the feed URL, the channel title and the item GUID.
func SyncFeed(ctx context.Context, db *database.Queries, feedName string, feedURL string, rss *RSS) (database.Channel, []database.Item, error) {

	feed, err := getOrCreateFeed(ctx, db, feedName, feedURL)
	if err != nil {
		return database.Channel{}, nil, err
	}

	channel, err := getOrCreateChannel(ctx, db, feed, rss.Channel)
	if err != nil {
		return database.Channel{}, nil, err
	}

	items := make([]database.Item, 0, len(rss.Channel.Items))

	for _, item := range rss.Channel.Items {
		stored, err := upsertItem(ctx, db, channel.ID, item)
		if err != nil {
			return database.Channel{}, nil, err
		}

		items = append(items, stored)
	}

	return channel, items, nil
}

func getOrCreateFeed(ctx context.Context, db *database.Queries, feedName string, feedURL string) (database.Feed, error) {
	feed, err := db.GetFeedByUrl(ctx, feedURL)

	if err == nil {
		return feed, nil
	}

	if !errors.Is(err, sql.ErrNoRows) {
		return database.Feed{}, fmt.Errorf("looking up feed %s: %w", feedURL, err)
	}

	now := time.Now().UTC()

	feed, err = db.CreateFeed(ctx, database.CreateFeedParams{
		ID:        uuid.New(),
		CreatedAt: now,
		UpdatedAt: now,
		Url:       feedURL,
		Name:      feedName,
	})

	if err != nil {
		return database.Feed{}, fmt.Errorf("creating feed %s: %w", feedURL, err)
	}

	return feed, nil
}

func getOrCreateChannel(ctx context.Context, db *database.Queries, feed database.Feed, parsed Channel) (database.Channel, error) {
	channel, err := db.GetChannelByTitle(ctx, parsed.Title)

	if err == nil {
		return channel, nil
	}

	if !errors.Is(err, sql.ErrNoRows) {
		return database.Channel{}, fmt.Errorf("looking up channel %q: %w", parsed.Title, err)
	}

	now := time.Now().UTC()

	channel, err = db.CreateChannel(ctx, database.CreateChannelParams{
		ID:          uuid.New(),
		CreatedAt:   now,
		UpdatedAt:   now,
		Title:       parsed.Title,
		Description: nullString(parsed.Description),
		Link:        nullString(parsed.Link),
		AtomLink:    nullString(parsed.AtomLink.Href),
		FeedID:      feed.ID,
	})

	if err != nil {
		return database.Channel{}, fmt.Errorf("creating channel %q: %w", parsed.Title, err)
	}

	return channel, nil
}

func upsertItem(ctx context.Context, db *database.Queries, channelID uuid.UUID, item Item) (database.Item, error) {
	// The GUID is only unique within its channel: the same torrent can legitimately
	// appear in two feeds, and each keeps its own row.
	current, err := db.GetItemByChannelIdAndGuid(ctx, database.GetItemByChannelIdAndGuidParams{
		ChannelID: channelID,
		Guid:      item.GUID,
	})

	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return database.Item{}, fmt.Errorf("looking up item %s: %w", item.GUID, err)
	}

	now := time.Now().UTC()
	publishedAt := parsePubDate(item.PubDate)

	if errors.Is(err, sql.ErrNoRows) {
		created, err := db.CreateItem(ctx, database.CreateItemParams{
			ID:          uuid.New(),
			Title:       nullString(item.Title),
			Link:        nullString(item.Link),
			Guid:        item.GUID,
			Pubdate:     nullString(item.PubDate),
			PublishedAt: publishedAt,
			Seeders:     nullInt32(item.Seeders),
			Leechers:    nullInt32(item.Leechers),
			Downloads:   nullInt32(item.Downloads),
			Infohash:    nullString(item.InfoHash),
			CategoryID:  nullString(item.CategoryID),
			Category:    nullString(item.Category),
			Size:        nullString(item.Size),
			Comments:    nullInt32(item.Comments),
			Trusted:     nullString(item.Trusted),
			Remake:      nullString(item.Remake),
			Description: nullString(item.Description),
			CreatedAt:   now,
			UpdatedAt:   now,
			ChannelID:   channelID,
		})

		if err != nil {
			return database.Item{}, fmt.Errorf("creating item %s: %w", item.GUID, err)
		}

		return created, nil
	}

	updated, err := db.UpdateItem(ctx, database.UpdateItemParams{
		ID:          current.ID,
		Title:       nullString(item.Title),
		Link:        nullString(item.Link),
		Guid:        item.GUID,
		Pubdate:     nullString(item.PubDate),
		PublishedAt: publishedAt,
		Seeders:     nullInt32(item.Seeders),
		Leechers:    nullInt32(item.Leechers),
		Downloads:   nullInt32(item.Downloads),
		Infohash:    nullString(item.InfoHash),
		CategoryID:  nullString(item.CategoryID),
		Category:    nullString(item.Category),
		Size:        nullString(item.Size),
		Comments:    nullInt32(item.Comments),
		Trusted:     nullString(item.Trusted),
		Remake:      nullString(item.Remake),
		Description: nullString(item.Description),
		UpdatedAt:   now,
		ChannelID:   channelID,
	})

	if err != nil {
		return database.Item{}, fmt.Errorf("updating item %s: %w", item.GUID, err)
	}

	return updated, nil
}

// FetchAndParseRSS fetches the nyaa.si feed for the given language and
// resolution and stores its content.
func FetchAndParseRSS(ctx context.Context,
	fetchAndParseRSSRequest *FetchAndParseRSSRequest,
	db *database.Queries) (database.Channel, []database.Item, error) {

	feedURL := NyaaFeedURL(fetchAndParseRSSRequest)

	rss, err := FetchRSS(ctx, nil, feedURL)
	if err != nil {
		return database.Channel{}, nil, err
	}

	return SyncFeed(ctx, db, "Nyaa", feedURL, rss)
}

// pubDateLayouts are the layouts accepted for an RSS <pubDate>. nyaa.si serves a
// numeric offset ("-0000"), which RFC1123 does not accept: it wants a zone name.
var pubDateLayouts = []string{
	time.RFC1123Z,
	time.RFC1123,
	time.RFC822Z,
	time.RFC822,
	time.RFC3339,
}

// parsePubDate turns the raw RSS date into a sortable timestamp. An unparsable
// date leaves the column NULL rather than failing the whole item: the original
// string is kept in pubdate either way.
func parsePubDate(value string) sql.NullTime {
	if value == "" {
		return sql.NullTime{}
	}

	for _, layout := range pubDateLayouts {
		if parsed, err := time.Parse(layout, value); err == nil {
			return sql.NullTime{Time: parsed.UTC(), Valid: true}
		}
	}

	log.Printf("rss: cannot parse pubDate %q, leaving published_at empty", value)

	return sql.NullTime{}
}

// nullString maps an empty string to a NULL column value.
func nullString(value string) sql.NullString {
	if value == "" {
		return sql.NullString{}
	}

	return sql.NullString{String: value, Valid: true}
}

func nullInt32(value int) sql.NullInt32 {
	return sql.NullInt32{Int32: int32(value), Valid: true}
}
