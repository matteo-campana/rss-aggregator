// Package cache holds the Redis-backed extras: a cache for the read endpoints,
// a per-key rate limiter, and a lock that stops two scrapers working the same
// feed.
//
// Redis is always optional. When it is disabled or unreachable every operation
// degrades to the behaviour the application had without it — the cache misses,
// the limiter allows, the lock is granted — and the error is logged, never
// returned to the caller. A broken cache must not take the API down with it.
package cache

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

const (
	defaultAddr              = "localhost:6379"
	defaultTTL               = time.Minute
	defaultRateLimitRequests = 120
	defaultRateLimitWindow   = time.Minute

	// operationTimeout keeps a slow or hung Redis from holding up a request.
	operationTimeout = 200 * time.Millisecond

	// errorCooldown throttles the degradation logging. A Redis that is down
	// fails on every single operation, and one line each would bury everything
	// else in the log.
	errorCooldown = 30 * time.Second
)

// Config holds the Redis settings, read from the environment.
type Config struct {
	Enabled           bool
	Addr              string
	Password          string
	DB                int
	TTL               time.Duration
	RateLimitRequests int
	RateLimitWindow   time.Duration
}

// ConfigFromEnv reads REDIS_* and RATE_LIMIT_*, falling back to the defaults on
// missing or invalid values.
func ConfigFromEnv() Config {
	config := Config{
		Enabled:           os.Getenv("REDIS_ENABLED") == "true",
		Addr:              defaultAddr,
		Password:          os.Getenv("REDIS_PASSWORD"),
		TTL:               defaultTTL,
		RateLimitRequests: defaultRateLimitRequests,
		RateLimitWindow:   defaultRateLimitWindow,
	}

	if addr := os.Getenv("REDIS_ADDR"); addr != "" {
		config.Addr = addr
	}

	if raw := os.Getenv("REDIS_DB"); raw != "" {
		if db, err := strconv.Atoi(raw); err == nil && db >= 0 {
			config.DB = db
		} else {
			log.Printf("cache: ignoring REDIS_DB=%q", raw)
		}
	}

	config.TTL = durationFromEnv("CACHE_TTL", defaultTTL)
	config.RateLimitWindow = durationFromEnv("RATE_LIMIT_WINDOW", defaultRateLimitWindow)

	if raw := os.Getenv("RATE_LIMIT_REQUESTS"); raw != "" {
		if requests, err := strconv.Atoi(raw); err == nil && requests > 0 {
			config.RateLimitRequests = requests
		} else {
			log.Printf("cache: ignoring RATE_LIMIT_REQUESTS=%q, using %d", raw, defaultRateLimitRequests)
		}
	}

	return config
}

func durationFromEnv(name string, fallback time.Duration) time.Duration {
	raw := os.Getenv(name)

	if raw == "" {
		return fallback
	}

	value, err := time.ParseDuration(raw)

	if err != nil || value <= 0 {
		log.Printf("cache: ignoring %s=%q, using %s", name, raw, fallback)
		return fallback
	}

	return value
}

// Client is the entry point for every Redis-backed feature. A zero Client, and
// one built from a disabled Config, are both usable: they simply do nothing.
type Client struct {
	redis  *redis.Client
	ttl    time.Duration
	limit  int
	window time.Duration

	errorsMutex       sync.Mutex
	lastErrorLoggedAt time.Time
	suppressedErrors  int
}

// report logs a degradation, at most once per cooldown. When a Redis outage
// silences the rest, the next line that gets through says how many were held
// back, so the volume is visible without being printed.
func (c *Client) report(format string, args ...any) {
	if c == nil {
		return
	}

	c.errorsMutex.Lock()
	defer c.errorsMutex.Unlock()

	now := time.Now()

	if !c.lastErrorLoggedAt.IsZero() && now.Sub(c.lastErrorLoggedAt) < errorCooldown {
		c.suppressedErrors++
		return
	}

	if c.suppressedErrors > 0 {
		log.Printf("cache: %d further errors suppressed in the last %s",
			c.suppressedErrors, now.Sub(c.lastErrorLoggedAt).Round(time.Second))

		c.suppressedErrors = 0
	}

	c.lastErrorLoggedAt = now

	log.Printf(format, args...)
}

// New builds a client. A Redis that cannot be reached at startup is not fatal:
// it may come back, and until it does every operation degrades.
func New(config Config) *Client {
	if !config.Enabled {
		log.Print("cache: disabled, set REDIS_ENABLED=true to enable it")
		return &Client{ttl: config.TTL, limit: config.RateLimitRequests, window: config.RateLimitWindow}
	}

	client := redis.NewClient(&redis.Options{
		Addr:     config.Addr,
		Password: config.Password,
		DB:       config.DB,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		log.Printf("cache: %s is not reachable yet (%v); continuing without it", config.Addr, err)
	} else {
		log.Printf("cache: using %s, ttl %s, rate limit %d per %s",
			config.Addr, config.TTL, config.RateLimitRequests, config.RateLimitWindow)
	}

	return &Client{
		redis:  client,
		ttl:    config.TTL,
		limit:  config.RateLimitRequests,
		window: config.RateLimitWindow,
	}
}

// Disabled returns a client that does nothing, for tests and for running
// without Redis.
func Disabled() *Client {
	return &Client{}
}

// Enabled reports whether a Redis connection was configured at all.
func (c *Client) Enabled() bool {
	return c != nil && c.redis != nil
}

// RateLimit reports the configured allowance.
func (c *Client) RateLimit() (int, time.Duration) {
	if c == nil {
		return defaultRateLimitRequests, defaultRateLimitWindow
	}

	return c.limit, c.window
}

func (c *Client) Close() error {
	if !c.Enabled() {
		return nil
	}

	return c.redis.Close()
}

// call runs an operation under the shared timeout.
func (c *Client) call(ctx context.Context, operation func(context.Context) error) error {
	if !c.Enabled() {
		return redis.ErrClosed
	}

	ctx, cancel := context.WithTimeout(ctx, operationTimeout)
	defer cancel()

	return operation(ctx)
}

// --- cached responses -------------------------------------------------------

// Key builds the cache key for an endpoint and its query string. The current
// namespace version is part of it, so bumping the version retires every key at
// once without scanning Redis.
func (c *Client) Key(ctx context.Context, namespace, endpoint, query string) string {
	return fmt.Sprintf("%s:v%d:%s:%s", namespace, c.Version(ctx, namespace), endpoint, query)
}

// Version reads the namespace counter. An unreachable Redis reports 0, which
// only means every key of this request shares one version.
func (c *Client) Version(ctx context.Context, namespace string) int64 {
	var version int64

	err := c.call(ctx, func(ctx context.Context) error {
		value, err := c.redis.Get(ctx, versionKey(namespace)).Int64()

		if errors.Is(err, redis.Nil) {
			return nil
		}

		if err != nil {
			return err
		}

		version = value

		return nil
	})

	if err != nil && !errors.Is(err, redis.ErrClosed) {
		c.report("cache: reading the %s version: %v", namespace, err)
	}

	return version
}

// BumpVersion retires everything cached under a namespace.
func (c *Client) BumpVersion(ctx context.Context, namespace string) {
	err := c.call(ctx, func(ctx context.Context) error {
		return c.redis.Incr(ctx, versionKey(namespace)).Err()
	})

	if err != nil && !errors.Is(err, redis.ErrClosed) {
		c.report("cache: bumping the %s version: %v", namespace, err)
	}
}

func versionKey(namespace string) string {
	return namespace + ":version"
}

// GetRaw returns a cached payload, if there is one.
func (c *Client) GetRaw(ctx context.Context, key string) ([]byte, bool) {
	var payload []byte

	err := c.call(ctx, func(ctx context.Context) error {
		value, err := c.redis.Get(ctx, key).Bytes()

		if errors.Is(err, redis.Nil) {
			return nil
		}

		if err != nil {
			return err
		}

		payload = value

		return nil
	})

	if err != nil && !errors.Is(err, redis.ErrClosed) {
		c.report("cache: reading %s: %v", key, err)
	}

	return payload, payload != nil
}

// SetRaw stores a payload for the configured TTL.
func (c *Client) SetRaw(ctx context.Context, key string, payload []byte) {
	err := c.call(ctx, func(ctx context.Context) error {
		return c.redis.Set(ctx, key, payload, c.ttl).Err()
	})

	if err != nil && !errors.Is(err, redis.ErrClosed) {
		c.report("cache: writing %s: %v", key, err)
	}
}

// --- rate limiting ----------------------------------------------------------

// Decision is the outcome of a rate limit check.
type Decision struct {
	Allowed    bool
	Limit      int
	Remaining  int
	RetryAfter time.Duration
}

// Allow counts one request in the current fixed window. Without Redis it always
// allows: a limiter that cannot count must not lock everybody out.
func (c *Client) Allow(ctx context.Context, identity string) Decision {
	limit, window := c.RateLimit()

	decision := Decision{Allowed: true, Limit: limit, Remaining: limit}

	if !c.Enabled() {
		return decision
	}

	key := "ratelimit:" + hashIdentity(identity)

	var count int64

	err := c.call(ctx, func(ctx context.Context) error {
		pipe := c.redis.TxPipeline()
		incr := pipe.Incr(ctx, key)
		// Only the first request of a window sets the expiry, so the window does
		// not slide forward with every call.
		pipe.ExpireNX(ctx, key, window)

		if _, err := pipe.Exec(ctx); err != nil {
			return err
		}

		count = incr.Val()

		return nil
	})

	if err != nil {
		c.report("cache: rate limiting: %v", err)
		return decision
	}

	decision.Remaining = limit - int(count)

	if decision.Remaining < 0 {
		decision.Remaining = 0
	}

	if int(count) > limit {
		decision.Allowed = false
		decision.RetryAfter = c.ttlOf(ctx, key, window)
	}

	return decision
}

func (c *Client) ttlOf(ctx context.Context, key string, fallback time.Duration) time.Duration {
	remaining := fallback

	err := c.call(ctx, func(ctx context.Context) error {
		value, err := c.redis.TTL(ctx, key).Result()
		if err != nil {
			return err
		}

		if value > 0 {
			remaining = value
		}

		return nil
	})

	if err != nil && !errors.Is(err, redis.ErrClosed) {
		c.report("cache: reading the ttl of %s: %v", key, err)
	}

	return remaining
}

// hashIdentity keeps API keys out of Redis: only their digest is stored.
func hashIdentity(identity string) string {
	sum := sha256.Sum256([]byte(identity))

	return hex.EncodeToString(sum[:])
}

// --- locking ----------------------------------------------------------------

// releaseScript deletes the lock only if it still holds our token, so a lock
// that already expired and was taken by somebody else is never released here.
var releaseScript = redis.NewScript(`
if redis.call("get", KEYS[1]) == ARGV[1] then
	return redis.call("del", KEYS[1])
end
return 0
`)

// AcquireLock takes a lock for at most ttl and returns a release function.
// Without Redis the lock is always granted: a single instance needs no
// coordination, which is exactly the behaviour before Redis existed.
func (c *Client) AcquireLock(ctx context.Context, key string, ttl time.Duration) (func(), bool) {
	if !c.Enabled() {
		return func() {}, true
	}

	token := uuid.NewString()
	acquired := false

	err := c.call(ctx, func(ctx context.Context) error {
		ok, err := c.redis.SetNX(ctx, key, token, ttl).Result()
		if err != nil {
			return err
		}

		acquired = ok

		return nil
	})

	if err != nil {
		// Redis is unreachable: fall back to working without coordination
		// rather than stopping the scraper altogether.
		c.report("cache: acquiring %s: %v", key, err)
		return func() {}, true
	}

	if !acquired {
		return func() {}, false
	}

	return func() {
		err := c.call(context.WithoutCancel(ctx), func(ctx context.Context) error {
			return releaseScript.Run(ctx, c.redis, []string{key}, token).Err()
		})

		if err != nil && !errors.Is(err, redis.ErrClosed) {
			c.report("cache: releasing %s: %v", key, err)
		}
	}, true
}
