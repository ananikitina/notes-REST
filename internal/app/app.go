package app

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/ananikitina/notes-rest/internal/config"
	"github.com/ananikitina/notes-rest/internal/database"
	"github.com/ananikitina/notes-rest/internal/handlers"
	"github.com/ananikitina/notes-rest/internal/middleware"
	"github.com/ananikitina/notes-rest/internal/repository"
	"github.com/ananikitina/notes-rest/internal/services"
	"github.com/ananikitina/notes-rest/internal/usecases"
	"github.com/go-chi/chi"
)

func Start() {
	log.Println("Starting application...")

	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load configuration:%v", err)
	}

	// Connect to database
	_, err = database.ConnectDB()
	if err != nil {
		log.Fatalf("Could not connect to the database: %v", err)
	}
	defer database.DB.Close()

	//Initialize the JWT service
	jwtService := services.NewJWTService(cfg)

	//Initialize the User repository, use case and handler
	userRepo := repository.NewUserRepository(database.DB)
	userUseCase := usecases.NewUserUseCase(userRepo)
	userHandler := handlers.NewUserHandler(userUseCase, jwtService)

	//Initialize the Note repository, use case and handler
	noteRepo := repository.NewNoteRepository(database.DB)
	noteUseCase := usecases.NewNoteUseCase(noteRepo)
	noteHandler := handlers.NewNoteHandler(noteUseCase)

	// Set up the router
	r := chi.NewRouter()

	//Public routes
	r.Post("/register", userHandler.RegisterHandler())
	r.Post("/login", userHandler.LoginHandler())

	// Protected routes
	r.Group(func(r chi.Router) {
		r.Use(middleware.AuthMiddleware(jwtService))
		r.Get("/notes", noteHandler.GetNotesHandler())
		r.Post("/note", noteHandler.AddNoteHandler())

		r.Group(func(r chi.Router) {
			r.Use(middleware.AdminOnlyMiddleware)
			r.Get("/allnotes", noteHandler.GetAllNotesHandler())
		})
	})

	// Create and configure the HTTP server
	server := &http.Server{
		Addr:    ":8080",
		Handler: r,
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
