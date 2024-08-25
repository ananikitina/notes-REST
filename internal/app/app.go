package app

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/ananikitina/notes-rest/internal/database"
	"github.com/ananikitina/notes-rest/internal/handlers"
)

func Start() {
	log.Println("Starting application...")

	// Connect to database
	_, err := database.ConnectDB()
	if err != nil {
		log.Fatalf("Could not connect to the database: %v", err)
	}
	defer database.DB.Close()

	// Set up router
	mux := http.NewServeMux()
	mux.HandleFunc("/notes", handlers.GetNotesHandler())
	mux.HandleFunc("/note", handlers.AddNoteHandler())

	// Create and configure the HTTP server
	server := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	// Start the HTTP server in a separate goroutine
	go func() {
		if err := server.ListenAndServe(); err != nil {
			log.Fatalf("Failed to start the server: %v", err)
		}
	}()

	log.Println("Server started on :8080")

	// Set up channel to listen for OS signals
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGTERM, syscall.SIGINT)
	<-signals

	log.Println("Shutdown signal received, shutting down server...")

	// Create a context with a timeout for graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exited gracefully.")
}
