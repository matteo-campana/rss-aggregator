// Package scraper keeps the stored feeds up to date in the background, so the
// API serves items that were collected ahead of the request instead of fetching
// them while a client waits.
package scraper

import (
	"context"
	"log"
	"net/http"
	"os"
	"strconv"
	"sync"
	"sync/atomic"
	"time"

	"rss-aggregator/internal/cache"
	"rss-aggregator/internal/database"
	"rss-aggregator/pkg/rss"

	"github.com/google/uuid"
)

const (
	defaultConcurrency = 3
	defaultInterval    = 10 * time.Minute
	fetchTimeout       = 30 * time.Second

	// lockTTL outlives a fetch, so a worker that dies mid-run does not hold a
	// feed hostage, and a healthy one is never cut off half way.
	lockTTL = 2 * fetchTimeout
)

// Config holds the scraper settings, read from the environment.
type Config struct {
	Enabled     bool
	Concurrency int
	Interval    time.Duration
}

// ConfigFromEnv reads SCRAPER_ENABLED, SCRAPER_CONCURRENCY and
// SCRAPER_INTERVAL, falling back to the defaults on missing or invalid values.
func ConfigFromEnv() Config {
	config := Config{
		Enabled:     os.Getenv("SCRAPER_ENABLED") == "true",
		Concurrency: defaultConcurrency,
		Interval:    defaultInterval,
	}

	if raw := os.Getenv("SCRAPER_CONCURRENCY"); raw != "" {
		concurrency, err := strconv.Atoi(raw)
		if err != nil || concurrency < 1 {
			log.Printf("scraper: ignoring SCRAPER_CONCURRENCY=%q, using %d", raw, defaultConcurrency)
		} else {
			config.Concurrency = concurrency
		}
	}

	if raw := os.Getenv("SCRAPER_INTERVAL"); raw != "" {
		interval, err := time.ParseDuration(raw)
		if err != nil || interval <= 0 {
			log.Printf("scraper: ignoring SCRAPER_INTERVAL=%q, using %s", raw, defaultInterval)
		} else {
			config.Interval = interval
		}
	}

	return config
}

// FeedStore is the slice of the database layer the scraper depends on.
type FeedStore interface {
	GetNextFeedsToFetch(ctx context.Context, limit int32) ([]database.Feed, error)
	MarkFeedAsFetched(ctx context.Context, id uuid.UUID) (database.Feed, error)
}

// Scraper periodically refreshes the feeds that were fetched least recently.
type Scraper struct {
	store       FeedStore
	cache       *cache.Client
	syncFeed    func(ctx context.Context, feed database.Feed) error
	concurrency int
	interval    time.Duration
}

// New builds a scraper backed by the generated queries. The cache client may be
// a disabled one: locking and invalidation then simply do nothing.
func New(db *database.Queries, config Config, cacheClient *cache.Client) *Scraper {
	client := &http.Client{Timeout: fetchTimeout}

	if cacheClient == nil {
		cacheClient = cache.Disabled()
	}

	return &Scraper{
		store:       db,
		cache:       cacheClient,
		concurrency: config.Concurrency,
		interval:    config.Interval,
		syncFeed: func(ctx context.Context, feed database.Feed) error {
			parsed, err := rss.FetchRSS(ctx, client, feed.Url)
			if err != nil {
				return err
			}

			_, items, err := rss.SyncFeed(ctx, db, feed.Name, feed.Url, parsed)
			if err != nil {
				return err
			}

			log.Printf("scraper: %s: stored %d items", feed.Name, len(items))

			return nil
		},
	}
}

// Run scrapes once immediately and then on every tick, until ctx is cancelled.
func (s *Scraper) Run(ctx context.Context) {
	log.Printf("scraper: starting with %d workers every %s", s.concurrency, s.interval)

	ticker := time.NewTicker(s.interval)
	defer ticker.Stop()

	for {
		s.ScrapeOnce(ctx)

		select {
		case <-ctx.Done():
			log.Print("scraper: stopped")
			return
		case <-ticker.C:
		}
	}
}

// ScrapeOnce refreshes one batch of feeds.
func (s *Scraper) ScrapeOnce(ctx context.Context) {
	feeds, err := s.store.GetNextFeedsToFetch(ctx, int32(s.concurrency))
	if err != nil {
		log.Printf("scraper: fetching the feed list: %v", err)
		return
	}

	if len(feeds) == 0 {
		return
	}

	waitGroup := &sync.WaitGroup{}
	refreshed := &atomic.Bool{}

	for _, feed := range feeds {
		waitGroup.Add(1)

		go func(feed database.Feed) {
			defer waitGroup.Done()

			if s.scrape(ctx, feed) {
				refreshed.Store(true)
			}
		}(feed)
	}

	waitGroup.Wait()

	// The stored items changed, so anything cached from them is stale.
	if refreshed.Load() {
		s.cache.BumpVersion(ctx, "items")
	}
}

// scrape refreshes a single feed. A failure is logged and skipped: one broken
// feed must not stop the others or the loop.
func (s *Scraper) scrape(ctx context.Context, feed database.Feed) bool {
	// With several instances running, only one of them should work a given
	// feed. Without Redis the lock is always granted, which is the behaviour a
	// single instance had all along.
	release, acquired := s.cache.AcquireLock(ctx, "lock:feed:"+feed.ID.String(), lockTTL)

	if !acquired {
		log.Printf("scraper: feed %s is being refreshed elsewhere, skipping", feed.Name)
		return false
	}

	defer release()

	if _, err := s.store.MarkFeedAsFetched(ctx, feed.ID); err != nil {
		log.Printf("scraper: marking feed %s as fetched: %v", feed.Name, err)
		return false
	}

	if err := s.syncFeed(ctx, feed); err != nil {
		log.Printf("scraper: syncing feed %s: %v", feed.Name, err)
		return false
	}

	return true
}
