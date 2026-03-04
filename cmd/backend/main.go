package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/HerbHall/RunNotes/internal/database"
	"github.com/HerbHall/RunNotes/internal/handler"
	"github.com/HerbHall/RunNotes/internal/store"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	socketPath := flag.String("socket", "/run/guest-services/backend.sock", "Unix socket path")
	dbPath := flag.String("db", "/data/runnotes.db", "SQLite database path")
	flag.Parse()

	// Dev mode: use TCP instead of Unix socket, local DB path.
	devMode := os.Getenv("ENV_MODE") == "dev"
	if devMode {
		if *dbPath == "/data/runnotes.db" {
			*dbPath = "./runnotes.db"
		}
	}

	log.Printf("opening database at %s", *dbPath)
	db, err := database.Open(*dbPath)
	if err != nil {
		return fmt.Errorf("open database: %w", err)
	}
	defer func() { _ = db.Close() }()

	noteStore := store.NewNoteStore(db)
	h := handler.NewHandler(noteStore)

	mux := http.NewServeMux()
	h.RegisterRoutes(mux)

	srv := &http.Server{
		Handler:      mux,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	// Start listener.
	var ln net.Listener
	if devMode {
		addr := ":3001"
		ln, err = net.Listen("tcp", addr) //nolint:gosec // G102: dev-mode only, not exposed in production
		if err != nil {
			return fmt.Errorf("listen tcp %s: %w", addr, err)
		}
		log.Printf("DEV MODE: listening on http://localhost%s", addr)
	} else {
		_ = os.RemoveAll(*socketPath)
		ln, err = net.Listen("unix", *socketPath)
		if err != nil {
			return fmt.Errorf("listen unix %s: %w", *socketPath, err)
		}
		log.Printf("listening on unix socket %s", *socketPath)
	}

	// Graceful shutdown on SIGTERM/SIGINT.
	done := make(chan os.Signal, 1)
	signal.Notify(done, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		if serveErr := srv.Serve(ln); serveErr != nil && serveErr != http.ErrServerClosed {
			log.Printf("serve: %v", serveErr)
		}
	}()

	<-done
	log.Println("shutting down...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		return fmt.Errorf("shutdown: %w", err)
	}
	log.Println("stopped")
	return nil
}
