package rss

import (
	"context"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"strings"
	"testing"
	"time"
)

func loadSample(t *testing.T) *RSS {
	t.Helper()

	file, err := os.Open("testdata/nyaa_sample.xml")
	if err != nil {
		t.Fatalf("opening fixture: %v", err)
	}

	defer file.Close()

	rss, err := ParseRSS(file)
	if err != nil {
		t.Fatalf("ParseRSS: %v", err)
	}

	return rss
}

func TestParseRSSChannel(t *testing.T) {
	rss := loadSample(t)

	if got, want := rss.Channel.Title, "Nyaa - Home - Torrent File RSS"; got != want {
		t.Errorf("channel title = %q, want %q", got, want)
	}

	if got, want := rss.Channel.Link, "https://nyaa.si/"; got != want {
		t.Errorf("channel link = %q, want %q", got, want)
	}

	// <atom:link> shares its local name with <link>: the namespaced field has to
	// win, and it must not swallow the plain <link> value.
	if got, want := rss.Channel.AtomLink.Href, "https://nyaa.si/?page=rss&c=1_0&f=0&q=ENG+1080p"; got != want {
		t.Errorf("channel atom link = %q, want %q", got, want)
	}

	if got, want := len(rss.Channel.Items), 2; got != want {
		t.Fatalf("item count = %d, want %d", got, want)
	}
}

// TestParseRSSNyaaNamespacedFields covers the nyaa:* elements, which used to be
// dropped silently and stored as zero values.
func TestParseRSSNyaaNamespacedFields(t *testing.T) {
	rss := loadSample(t)

	item := rss.Channel.Items[0]

	tests := []struct {
		name string
		got  any
		want any
	}{
		{"title", item.Title, "[SubsPlease] Example Anime - 01 (1080p) [ABCD1234].mkv"},
		{"link", item.Link, "https://nyaa.si/download/1800001.torrent"},
		{"guid", item.GUID, "https://nyaa.si/view/1800001"},
		{"pubDate", item.PubDate, "Tue, 25 Jun 2024 10:15:00 -0000"},
		{"seeders", item.Seeders, 1234},
		{"leechers", item.Leechers, 56},
		{"downloads", item.Downloads, 7890},
		{"infoHash", item.InfoHash, "1f1b2c3d4e5f60718293a4b5c6d7e8f901234567"},
		{"categoryId", item.CategoryID, "1_2"},
		{"category", item.Category, "Anime - English-translated"},
		{"size", item.Size, "1.4 GiB"},
		{"comments", item.Comments, 3},
		{"trusted", item.Trusted, "Yes"},
		{"remake", item.Remake, "No"},
	}

	for _, tt := range tests {
		if tt.got != tt.want {
			t.Errorf("%s = %v, want %v", tt.name, tt.got, tt.want)
		}
	}

	if !strings.Contains(item.Description, "#1800001") {
		t.Errorf("description = %q, want it to contain %q", item.Description, "#1800001")
	}

	if got, want := rss.Channel.Items[1].Seeders, 42; got != want {
		t.Errorf("second item seeders = %d, want %d", got, want)
	}
}

func TestParseRSSInvalidDocument(t *testing.T) {
	if _, err := ParseRSS(strings.NewReader("<rss><channel>")); err == nil {
		t.Fatal("ParseRSS on a truncated document: got nil error, want an error")
	}
}

func TestNyaaFeedURL(t *testing.T) {
	feedURL := NyaaFeedURL(&FetchAndParseRSSRequest{Language: ENG, Resolution: Resolution1080p})

	parsed, err := url.Parse(feedURL)
	if err != nil {
		t.Fatalf("parsing built url %q: %v", feedURL, err)
	}

	query := parsed.Query()

	if got, want := query.Get("q"), "ENG 1080p"; got != want {
		t.Errorf("q = %q, want %q", got, want)
	}

	if got, want := query.Get("page"), "rss"; got != want {
		t.Errorf("page = %q, want %q", got, want)
	}
}

func TestFetchRSS(t *testing.T) {
	fixture, err := os.ReadFile("testdata/nyaa_sample.xml")
	if err != nil {
		t.Fatalf("reading fixture: %v", err)
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/rss+xml")
		w.Write(fixture)
	}))

	defer server.Close()

	rss, err := FetchRSS(context.Background(), server.Client(), server.URL)
	if err != nil {
		t.Fatalf("FetchRSS: %v", err)
	}

	if got, want := len(rss.Channel.Items), 2; got != want {
		t.Fatalf("item count = %d, want %d", got, want)
	}

	if got, want := rss.Channel.Items[0].Seeders, 1234; got != want {
		t.Errorf("seeders = %d, want %d", got, want)
	}
}

func TestFetchRSSNonOKStatus(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "nope", http.StatusServiceUnavailable)
	}))

	defer server.Close()

	if _, err := FetchRSS(context.Background(), server.Client(), server.URL); err == nil {
		t.Fatal("FetchRSS on a 503 response: got nil error, want an error")
	}
}

func TestParsePubDate(t *testing.T) {
	tests := []struct {
		name  string
		value string
		want  string
		valid bool
	}{
		// The format nyaa.si actually serves: a numeric offset, which RFC1123
		// alone does not accept.
		{name: "nyaa", value: "Tue, 25 Jun 2024 10:15:00 -0000", want: "2024-06-25T10:15:00Z", valid: true},
		{name: "named zone", value: "Tue, 25 Jun 2024 10:15:00 UTC", want: "2024-06-25T10:15:00Z", valid: true},
		{name: "offset", value: "Tue, 25 Jun 2024 12:15:00 +0200", want: "2024-06-25T10:15:00Z", valid: true},
		{name: "empty", value: ""},
		{name: "garbage", value: "yesterday"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := parsePubDate(tt.value)

			if got.Valid != tt.valid {
				t.Fatalf("valid = %t, want %t", got.Valid, tt.valid)
			}

			if tt.valid && got.Time.UTC().Format(time.RFC3339) != tt.want {
				t.Errorf("time = %s, want %s", got.Time.UTC().Format(time.RFC3339), tt.want)
			}
		})
	}
}

func TestParseRSSItemDateIsParsable(t *testing.T) {
	rss := loadSample(t)

	if got := parsePubDate(rss.Channel.Items[0].PubDate); !got.Valid {
		t.Errorf("the fixture pubDate %q was not parsed", rss.Channel.Items[0].PubDate)
	}
}
