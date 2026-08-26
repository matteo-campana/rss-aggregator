package scraper

import (
	"context"
	"errors"
	"sync"
	"testing"
	"time"

	"rss-aggregator/internal/database"

	"github.com/google/uuid"
)

type fakeStore struct {
	mutex  sync.Mutex
	feeds  []database.Feed
	marked []uuid.UUID
	err    error
}

func (f *fakeStore) GetNextFeedsToFetch(ctx context.Context, limit int32) ([]database.Feed, error) {
	if f.err != nil {
		return nil, f.err
	}

	if int(limit) < len(f.feeds) {
		return f.feeds[:limit], nil
	}

	return f.feeds, nil
}

func (f *fakeStore) MarkFeedAsFetched(ctx context.Context, id uuid.UUID) (database.Feed, error) {
	f.mutex.Lock()
	defer f.mutex.Unlock()

	f.marked = append(f.marked, id)

	return database.Feed{ID: id}, nil
}

func newTestScraper(store FeedStore, sync func(ctx context.Context, feed database.Feed) error) *Scraper {
	return &Scraper{
		store:       store,
		syncFeed:    sync,
		concurrency: 2,
		interval:    time.Hour,
	}
}

func TestScrapeOnceMarksAndSyncsEveryFeed(t *testing.T) {
	store := &fakeStore{feeds: []database.Feed{
		{ID: uuid.New(), Name: "Nyaa", Url: "https://nyaa.si/?page=rss"},
		{ID: uuid.New(), Name: "Other", Url: "https://example.com/rss"},
	}}

	var mutex sync.Mutex
	synced := []string{}

	scraper := newTestScraper(store, func(ctx context.Context, feed database.Feed) error {
		mutex.Lock()
		defer mutex.Unlock()

		synced = append(synced, feed.Name)

		return nil
	})

	scraper.ScrapeOnce(context.Background())

	if got, want := len(store.marked), 2; got != want {
		t.Errorf("marked feeds = %d, want %d", got, want)
	}

	if got, want := len(synced), 2; got != want {
		t.Errorf("synced feeds = %d, want %d", got, want)
	}
}

// A failing feed must not stop the other feeds in the same batch.
func TestScrapeOnceKeepsGoingAfterAFailure(t *testing.T) {
	store := &fakeStore{feeds: []database.Feed{
		{ID: uuid.New(), Name: "broken"},
		{ID: uuid.New(), Name: "working"},
	}}

	var mutex sync.Mutex
	succeeded := 0

	scraper := newTestScraper(store, func(ctx context.Context, feed database.Feed) error {
		if feed.Name == "broken" {
			return errors.New("feed is unreachable")
		}

		mutex.Lock()
		defer mutex.Unlock()

		succeeded++

		return nil
	})

	scraper.ScrapeOnce(context.Background())

	if got, want := succeeded, 1; got != want {
		t.Errorf("successful syncs = %d, want %d", got, want)
	}
}

func TestScrapeOnceWithoutFeeds(t *testing.T) {
	scraper := newTestScraper(&fakeStore{}, func(ctx context.Context, feed database.Feed) error {
		t.Error("syncFeed was called even though there are no feeds")
		return nil
	})

	scraper.ScrapeOnce(context.Background())
}

func TestRunStopsOnContextCancellation(t *testing.T) {
	scraper := newTestScraper(&fakeStore{}, func(ctx context.Context, feed database.Feed) error {
		return nil
	})

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	done := make(chan struct{})

	go func() {
		scraper.Run(ctx)
		close(done)
	}()

	select {
	case <-done:
	case <-time.After(2 * time.Second):
		t.Fatal("Run did not return after the context was cancelled")
	}
}

func TestConfigFromEnv(t *testing.T) {
	t.Setenv("SCRAPER_ENABLED", "true")
	t.Setenv("SCRAPER_CONCURRENCY", "7")
	t.Setenv("SCRAPER_INTERVAL", "45s")

	config := ConfigFromEnv()

	if !config.Enabled {
		t.Error("enabled = false, want true")
	}

	if got, want := config.Concurrency, 7; got != want {
		t.Errorf("concurrency = %d, want %d", got, want)
	}

	if got, want := config.Interval, 45*time.Second; got != want {
		t.Errorf("interval = %s, want %s", got, want)
	}
}

func TestConfigFromEnvFallsBackOnInvalidValues(t *testing.T) {
	t.Setenv("SCRAPER_ENABLED", "yes")
	t.Setenv("SCRAPER_CONCURRENCY", "0")
	t.Setenv("SCRAPER_INTERVAL", "soon")

	config := ConfigFromEnv()

	if config.Enabled {
		t.Error(`enabled = true, want false: only "true" enables the scraper`)
	}

	if got, want := config.Concurrency, defaultConcurrency; got != want {
		t.Errorf("concurrency = %d, want the default %d", got, want)
	}

	if got, want := config.Interval, defaultInterval; got != want {
		t.Errorf("interval = %s, want the default %s", got, want)
	}
}
