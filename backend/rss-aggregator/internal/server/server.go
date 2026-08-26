package server

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"rss-aggregator/internal/cache"
	"rss-aggregator/internal/database"

	_ "github.com/jackc/pgx/v5/stdlib"
	_ "github.com/joho/godotenv/autoload"
)

const defaultPort = 8080

type ApiConfig struct {
	queries *database.Queries
	conn    *sql.DB
	cache   *cache.Client
	port    int
}

// NewApiConfig opens the database connection and builds the API configuration.
func NewApiConfig() *ApiConfig {

	port, err := strconv.Atoi(os.Getenv("PORT"))
	if err != nil || port <= 0 {
		log.Printf("PORT is not a valid port number, falling back to %d", defaultPort)
		port = defaultPort
	}

	var (
		db_database = os.Getenv("DB_DATABASE")
		db_password = os.Getenv("DB_PASSWORD")
		db_username = os.Getenv("DB_USERNAME")
		db_port     = os.Getenv("DB_PORT")
		db_host     = os.Getenv("DB_HOST")
	)

	connStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		db_username, db_password, db_host, db_port, db_database)

	conn, err := sql.Open("pgx", connStr)
	if err != nil {
		log.Fatal("Error connecting to DB: ", err)
	}

	return &ApiConfig{
		queries: database.New(conn),
		conn:    conn,
		cache:   cache.New(cache.ConfigFromEnv()),
		port:    port,
	}
}

// Queries exposes the generated queries to the other components, the
// background scraper in particular.
func (apiCfg *ApiConfig) Queries() *database.Queries {
	return apiCfg.queries
}

// Cache exposes the Redis client to the other components, the scraper in
// particular. It is never nil: a disabled client is still usable.
func (apiCfg *ApiConfig) Cache() *cache.Client {
	if apiCfg.cache == nil {
		apiCfg.cache = cache.Disabled()
	}

	return apiCfg.cache
}

// Close releases the database and Redis connections.
func (apiCfg *ApiConfig) Close() error {
	if apiCfg.cache != nil {
		apiCfg.cache.Close()
	}

	if apiCfg.conn == nil {
		return nil
	}

	return apiCfg.conn.Close()
}

// Server builds the HTTP server serving this configuration.
func (apiCfg *ApiConfig) Server() *http.Server {
	return &http.Server{
		Addr:         fmt.Sprintf(":%d", apiCfg.port),
		Handler:      apiCfg.RegisterRoutes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}
}

// NewServer builds a ready-to-serve HTTP server with its own configuration.
func NewServer() *http.Server {
	return NewApiConfig().Server()
}
