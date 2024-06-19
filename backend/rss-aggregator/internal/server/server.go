package server

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"rss-aggregator/internal/database"

	_ "github.com/jackc/pgx/v5/stdlib"
	_ "github.com/joho/godotenv/autoload"
)

type ApiConfig struct {
	queries *database.Queries
	conn    *sql.DB
	port    int
}

func NewServer() *http.Server {

	port, _ := strconv.Atoi(os.Getenv("PORT"))

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

	queries := database.New(conn)

	apiCfg := &ApiConfig{
		queries: queries,
		conn:    conn,
		port:    port,
	}

	// Declare Server config
	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", apiCfg.port),
		Handler:      apiCfg.RegisterRoutes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	return server
}
