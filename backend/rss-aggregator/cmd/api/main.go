package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"rss-aggregator/internal/scraper"
	"rss-aggregator/internal/server"
)

const shutdownTimeout = 10 * time.Second

func main() {

	apiCfg := server.NewApiConfig()
	defer apiCfg.Close()

	// The context is cancelled on SIGINT/SIGTERM: it stops the scraper and
	// triggers the graceful shutdown of the HTTP server.
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	if config := scraper.ConfigFromEnv(); config.Enabled {
		go scraper.New(apiCfg.Queries(), config).Run(ctx)
	} else {
		log.Print("scraper: disabled, set SCRAPER_ENABLED=true to enable it")
	}

	httpServer := apiCfg.Server()

	go func() {
		log.Printf("listening on %s", httpServer.Addr)

		if err := httpServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("cannot start server: %s", err)
		}
	}()

	<-ctx.Done()

	log.Print("shutting down")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	if err := httpServer.Shutdown(shutdownCtx); err != nil {
		log.Printf("shutdown: %v", err)
	}
}
