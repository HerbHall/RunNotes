package main

import (
	"context"
	"flag"
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
		log.Fatalf("open database: %v", err)
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
		ln, err = net.Listen("tcp", addr)
		if err != nil {
			log.Fatalf("listen tcp %s: %v", addr, err)
		}
		log.Printf("DEV MODE: listening on http://localhost%s", addr)
	} else {
		_ = os.RemoveAll(*socketPath)
		ln, err = net.Listen("unix", *socketPath)
		if err != nil {
			log.Fatalf("listen unix %s: %v", *socketPath, err)
		}
		log.Printf("listening on unix socket %s", *socketPath)
	}

	// Graceful shutdown on SIGTERM/SIGINT.
	done := make(chan os.Signal, 1)
	signal.Notify(done, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		if err := srv.Serve(ln); err != nil && err != http.ErrServerClosed {
			log.Fatalf("serve: %v", err)
		}
	}()

	<-done
	log.Println("shutting down...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("shutdown: %v", err)
	}
	log.Println("stopped")
}
