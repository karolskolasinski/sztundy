package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"time"

	"core/internal/httpapi"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

type config struct {
	databaseURL string
	httpAddr    string
}

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	cfg, err := loadConfig()
	if err != nil {
		return err
	}

	pool, err := runDatabaseConnection(cfg.databaseURL)
	if err != nil {
		return err
	}
	defer pool.Close()

	if err := runHTTPServer(cfg.httpAddr); err != nil {
		return err
	}

	return nil
}

func loadConfig() (config, error) {
	if err := godotenv.Load(); err != nil &&
		!errors.Is(err, os.ErrNotExist) {
		return config{}, fmt.Errorf("Load environment: %w", err)
	}

	cfg := config{
		databaseURL: os.Getenv("DATABASE_URL"),
		httpAddr:    os.Getenv("HTTP_ADDR"),
	}

	if cfg.databaseURL == "" {
		return config{}, errors.New("DATABASE_URL is not set")
	}

	if cfg.httpAddr == "" {
		return config{}, errors.New("HTTP_ADDR is not set")
	}

	return cfg, nil
}

func runDatabaseConnection(databaseURL string) (*pgxpool.Pool, error) {
	poolConfig, err := pgxpool.ParseConfig(databaseURL)
	if err != nil {
		return nil, fmt.Errorf("Parse PostgreSQL connection configuration: %w", err)
	}

	pool, err := pgxpool.NewWithConfig(context.Background(), poolConfig)
	if err != nil {
		return nil, fmt.Errorf("Create PostgreSQL connection pool: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("Connect to PostgreSQL: %w", err)
	}

	host := poolConfig.ConnConfig.Host
	port := poolConfig.ConnConfig.Port
	log.Printf("Connected to PostgreSQL at %s:%d", host, port)

	return pool, nil
}

func runHTTPServer(addr string) error {
	if _, _, err := net.SplitHostPort(addr); err != nil {
		return fmt.Errorf("Validate HTTP_ADDR %q: %w", addr, err)
	}

	server := &http.Server{
		Addr:              addr,
		Handler:           httpapi.New(),
		ReadHeaderTimeout: 5 * time.Second,
	}

	log.Printf("HTTP server listening on %s", server.Addr)

	if err := server.ListenAndServe(); err != nil &&
		!errors.Is(err, http.ErrServerClosed) {
		return fmt.Errorf("Start HTTP server: %w", err)
	}

	return nil
}
