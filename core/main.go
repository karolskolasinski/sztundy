//package main
//
//import (
//	"context"
//	"errors"
//	"fmt"
//	"log"
//	"net/http"
//	"os"
//	"os/signal"
//	"syscall"
//	"time"
//
//	"github.com/jackc/pgx/v5/pgxpool"
//	"github.com/joho/godotenv"
//
//	"core/internal/httpapi"
//)
//
//func main() {
//	if err := run(); err != nil {
//		log.Fatal(err)
//	}
//}
//
//func run() error {
//	if err := godotenv.Load(); err != nil && !errors.Is(err, os.ErrNotExist) {
//		return fmt.Errorf("load .env: %w", err)
//	}
//
//	databaseURL := os.Getenv("DATABASE_URL")
//	if databaseURL == "" {
//		return errors.New("DATABASE_URL is not set")
//	}
//
//	addr := os.Getenv("HTTP_ADDR")
//	if addr == "" {
//		addr = "127.0.0.1:8080"
//	}
//
//	ctx, stop := signal.NotifyContext(
//		context.Background(),
//		os.Interrupt,
//		syscall.SIGTERM,
//	)
//	defer stop()
//
//	pool, err := pgxpool.New(ctx, databaseURL)
//	if err != nil {
//		return fmt.Errorf("create PostgreSQL connection pool: %w", err)
//	}
//	defer pool.Close()
//
//	pingCtx, cancelPing := context.WithTimeout(ctx, 5*time.Second)
//	err = pool.Ping(pingCtx)
//	cancelPing()
//
//	if err != nil {
//		return fmt.Errorf("connect to PostgreSQL: %w", err)
//	}
//
//	log.Println("PostgreSQL connection established")
//
//	server := &http.Server{
//		Addr:              addr,
//		Handler:           httpapi.New(),
//		ReadHeaderTimeout: 5 * time.Second,
//		ReadTimeout:       10 * time.Second,
//		WriteTimeout:      15 * time.Second,
//		IdleTimeout:       60 * time.Second,
//	}
//
//	serverErrors := make(chan error, 1)
//
//	go func() {
//		log.Printf("HTTP server listening on %s", server.Addr)
//		serverErrors <- server.ListenAndServe()
//	}()
//
//	select {
//	case err := <-serverErrors:
//		if !errors.Is(err, http.ErrServerClosed) {
//			return fmt.Errorf("HTTP server: %w", err)
//		}
//		return nil
//
//	case <-ctx.Done():
//		shutdownCtx, cancelShutdown := context.WithTimeout(
//			context.Background(),
//			10*time.Second,
//		)
//		defer cancelShutdown()
//
//		if err := server.Shutdown(shutdownCtx); err != nil {
//			_ = server.Close()
//			return fmt.Errorf("shutdown HTTP server: %w", err)
//		}
//
//		err := <-serverErrors
//		if !errors.Is(err, http.ErrServerClosed) {
//			return fmt.Errorf("HTTP server after shutdown: %w", err)
//		}
//
//		return nil
//	}
//}

package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	"core/internal/httpapi"
)

func main() {
	if err := runHTTPServer(); err != nil {
		log.Fatal(err)
	}
}

func runHTTPServer() error {
	server := &http.Server{
		Addr:              ":3000",
		Handler:           httpapi.New(),
		ReadHeaderTimeout: 5 * time.Second,
	}

	log.Printf("HTTP server listening on %s", server.Addr)

	if err := server.ListenAndServe(); err != nil &&
		!errors.Is(err, http.ErrServerClosed) {
		return fmt.Errorf("HTTP server: %w", err)
	}

	return nil
}
