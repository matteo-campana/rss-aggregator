package cache

import (
	"context"
	"log"
	"os"
	"strings"
	"sync"
	"testing"
	"time"
)

// A disabled client has to behave exactly like the application did before Redis
// existed: nothing cached, nothing limited, every lock granted.
func TestDisabledClientDegrades(t *testing.T) {
	client := Disabled()
	ctx := context.Background()

	if client.Enabled() {
		t.Error("Enabled() = true on a disabled client")
	}

	client.SetRaw(ctx, "k", []byte(`{"a":1}`))

	if _, ok := client.GetRaw(ctx, "k"); ok {
		t.Error("a disabled cache returned a hit")
	}

	if decision := client.Allow(ctx, "some-key"); !decision.Allowed {
		t.Error("a disabled limiter refused a request")
	}

	release, acquired := client.AcquireLock(ctx, "lock", time.Second)

	if !acquired {
		t.Error("a disabled locker refused the lock")
	}

	release()

	client.BumpVersion(ctx, "items")

	if got := client.Version(ctx, "items"); got != 0 {
		t.Errorf("version = %d, want 0", got)
	}
}

// A nil client turns up wherever an ApiConfig is built by hand, in tests among
// other places, and must not panic.
func TestNilClientDegrades(t *testing.T) {
	var client *Client

	ctx := context.Background()

	if client.Enabled() {
		t.Error("Enabled() = true on a nil client")
	}

	if decision := client.Allow(ctx, "key"); !decision.Allowed {
		t.Error("a nil limiter refused a request")
	}

	if _, acquired := client.AcquireLock(ctx, "lock", time.Second); !acquired {
		t.Error("a nil locker refused the lock")
	}
}

func TestConfigFromEnv(t *testing.T) {
	t.Setenv("REDIS_ENABLED", "true")
	t.Setenv("REDIS_ADDR", "redis:6379")
	t.Setenv("REDIS_DB", "3")
	t.Setenv("CACHE_TTL", "30s")
	t.Setenv("RATE_LIMIT_REQUESTS", "10")
	t.Setenv("RATE_LIMIT_WINDOW", "5s")

	config := ConfigFromEnv()

	if !config.Enabled || config.Addr != "redis:6379" || config.DB != 3 {
		t.Errorf("config = %+v, want the environment values", config)
	}

	if config.TTL != 30*time.Second || config.RateLimitRequests != 10 || config.RateLimitWindow != 5*time.Second {
		t.Errorf("config = %+v, want the environment values", config)
	}
}

func TestConfigFromEnvFallsBackOnInvalidValues(t *testing.T) {
	t.Setenv("REDIS_ENABLED", "yes")
	t.Setenv("REDIS_DB", "-1")
	t.Setenv("CACHE_TTL", "soon")
	t.Setenv("RATE_LIMIT_REQUESTS", "0")

	config := ConfigFromEnv()

	if config.Enabled {
		t.Error(`enabled = true, want false: only "true" enables Redis`)
	}

	if config.DB != 0 || config.TTL != defaultTTL || config.RateLimitRequests != defaultRateLimitRequests {
		t.Errorf("config = %+v, want the defaults", config)
	}
}

// The digest is what reaches Redis, so an API key is never stored in the clear.
func TestHashIdentity(t *testing.T) {
	hashed := hashIdentity("super-secret-api-key")

	if hashed == "super-secret-api-key" {
		t.Fatal("the identity was not hashed")
	}

	if len(hashed) != 64 {
		t.Errorf("digest length = %d, want 64", len(hashed))
	}

	if hashIdentity("super-secret-api-key") != hashed {
		t.Error("the digest is not stable")
	}

	if hashIdentity("another-key") == hashed {
		t.Error("two identities share a digest")
	}
}

// --- integration ------------------------------------------------------------

func integrationClient(t *testing.T) *Client {
	t.Helper()

	addr := os.Getenv("TEST_REDIS_ADDR")

	if addr == "" {
		t.Skip("TEST_REDIS_ADDR is not set")
	}

	client := New(Config{
		Enabled:           true,
		Addr:              addr,
		DB:                15,
		TTL:               time.Minute,
		RateLimitRequests: 3,
		RateLimitWindow:   2 * time.Second,
	})

	if !client.Enabled() {
		t.Fatal("the client is not enabled")
	}

	if err := client.redis.FlushDB(context.Background()).Err(); err != nil {
		t.Fatalf("flushing the test database: %v", err)
	}

	t.Cleanup(func() { client.Close() })

	return client
}

func TestCacheRoundTrip(t *testing.T) {
	client := integrationClient(t)
	ctx := context.Background()

	key := client.Key(ctx, "items", "items", "page=1")

	if _, ok := client.GetRaw(ctx, key); ok {
		t.Fatal("an empty cache returned a hit")
	}

	client.SetRaw(ctx, key, []byte(`{"total":3}`))

	payload, ok := client.GetRaw(ctx, key)

	if !ok {
		t.Fatal("the value just written was not found")
	}

	if string(payload) != `{"total":3}` {
		t.Errorf("payload = %s, want the stored JSON", payload)
	}
}

// Bumping the version has to retire the old keys without deleting anything.
func TestBumpVersionInvalidates(t *testing.T) {
	client := integrationClient(t)
	ctx := context.Background()

	key := client.Key(ctx, "items", "items", "page=1")
	client.SetRaw(ctx, key, []byte(`{"stale":true}`))

	client.BumpVersion(ctx, "items")

	if got := client.Version(ctx, "items"); got != 1 {
		t.Errorf("version = %d, want 1", got)
	}

	newKey := client.Key(ctx, "items", "items", "page=1")

	if newKey == key {
		t.Fatal("the key did not change after the bump")
	}

	if _, ok := client.GetRaw(ctx, newKey); ok {
		t.Error("the new key already has a value: the old entry was not retired")
	}
}

func TestDifferentQueriesUseDifferentKeys(t *testing.T) {
	client := integrationClient(t)
	ctx := context.Background()

	first := client.Key(ctx, "items", "items", "page=1")
	second := client.Key(ctx, "items", "items", "page=2")

	if first == second {
		t.Error("two different query strings share a cache key")
	}
}

func TestRateLimitAllowsThenRefuses(t *testing.T) {
	client := integrationClient(t)
	ctx := context.Background()

	for attempt := 1; attempt <= 3; attempt++ {
		decision := client.Allow(ctx, "a-key")

		if !decision.Allowed {
			t.Fatalf("request %d refused, want the first three allowed", attempt)
		}

		if decision.Remaining != 3-attempt {
			t.Errorf("request %d: remaining = %d, want %d", attempt, decision.Remaining, 3-attempt)
		}
	}

	decision := client.Allow(ctx, "a-key")

	if decision.Allowed {
		t.Error("the fourth request was allowed, want it refused")
	}

	if decision.RetryAfter <= 0 {
		t.Errorf("retry after = %s, want a positive delay", decision.RetryAfter)
	}

	// The allowance is per identity.
	if other := client.Allow(ctx, "another-key"); !other.Allowed {
		t.Error("a different identity was refused: the counters are shared")
	}
}

func TestRateLimitWindowExpires(t *testing.T) {
	client := integrationClient(t)
	ctx := context.Background()

	for i := 0; i < 4; i++ {
		client.Allow(ctx, "expiring")
	}

	if client.Allow(ctx, "expiring").Allowed {
		t.Fatal("still allowed after exceeding the limit")
	}

	time.Sleep(2100 * time.Millisecond)

	if !client.Allow(ctx, "expiring").Allowed {
		t.Error("still refused after the window elapsed")
	}
}

func TestLockIsExclusiveAndReleasable(t *testing.T) {
	client := integrationClient(t)
	ctx := context.Background()

	release, acquired := client.AcquireLock(ctx, "lock:feed:1", time.Minute)

	if !acquired {
		t.Fatal("the first lock was refused")
	}

	if _, second := client.AcquireLock(ctx, "lock:feed:1", time.Minute); second {
		t.Error("the lock was granted twice")
	}

	release()

	releaseAgain, third := client.AcquireLock(ctx, "lock:feed:1", time.Minute)

	if !third {
		t.Error("the lock was not granted after being released")
	}

	releaseAgain()
}

// Releasing must only remove our own token: a lock that expired and was taken
// by somebody else has to survive our release.
func TestReleaseDoesNotFreeSomebodyElsesLock(t *testing.T) {
	client := integrationClient(t)
	ctx := context.Background()

	release, acquired := client.AcquireLock(ctx, "lock:feed:2", 500*time.Millisecond)

	if !acquired {
		t.Fatal("the first lock was refused")
	}

	time.Sleep(700 * time.Millisecond)

	// The lock expired and another worker took it.
	otherRelease, taken := client.AcquireLock(ctx, "lock:feed:2", time.Minute)

	if !taken {
		t.Fatal("the expired lock was not available again")
	}

	// The first worker finishes and releases: it must not free the new holder.
	release()

	if _, stolen := client.AcquireLock(ctx, "lock:feed:2", time.Minute); stolen {
		t.Error("the stale release freed a lock held by somebody else")
	}

	otherRelease()
}

func TestOnlyOneWorkerGetsTheLock(t *testing.T) {
	client := integrationClient(t)
	ctx := context.Background()

	var (
		waitGroup sync.WaitGroup
		mutex     sync.Mutex
		granted   int
		releases  []func()
	)

	for i := 0; i < 8; i++ {
		waitGroup.Add(1)

		go func() {
			defer waitGroup.Done()

			release, acquired := client.AcquireLock(ctx, "lock:contended", time.Minute)

			if !acquired {
				return
			}

			mutex.Lock()
			granted++
			releases = append(releases, release)
			mutex.Unlock()
		}()
	}

	waitGroup.Wait()

	if granted != 1 {
		t.Errorf("%d workers got the lock, want exactly 1", granted)
	}

	for _, release := range releases {
		release()
	}
}

// A Redis that is down fails on every operation; the log must not drown in it.
func TestErrorReportingIsThrottled(t *testing.T) {
	var output strings.Builder

	log.SetOutput(&output)
	t.Cleanup(func() { log.SetOutput(os.Stderr) })

	client := &Client{}

	for i := 0; i < 50; i++ {
		client.report("cache: something failed %d", i)
	}

	lines := strings.Count(output.String(), "\n")

	if lines != 1 {
		t.Errorf("logged %d lines for 50 failures, want 1 within the cooldown", lines)
	}

	if client.suppressedErrors != 49 {
		t.Errorf("suppressed = %d, want 49", client.suppressedErrors)
	}

	// Once the cooldown has passed, the next line reports the backlog.
	client.lastErrorLoggedAt = time.Now().Add(-2 * errorCooldown)
	client.report("cache: something failed again")

	if !strings.Contains(output.String(), "49 further errors suppressed") {
		t.Errorf("the backlog was never reported:\n%s", output.String())
	}
}
