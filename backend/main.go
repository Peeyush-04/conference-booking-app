package main

import (
	"backend/handler"
	"backend/middleware"
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	// Connect to PostgreSQL
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		log.Fatal("DATABASE_URL not set in .env or environment")
	}

	dbpool, err := pgxpool.New(context.Background(), dbURL)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v", err)
	}
	defer dbpool.Close()

	// Initialize Router
	r := chi.NewRouter()

	// Global Middleware
	r.Use(middleware.RecoverMiddleware)
	r.Use(middleware.LoggingMiddleware)
	r.Use(middleware.RequestIDMiddleware)
	r.Use(middleware.RateLimitMiddleware)
	r.Use(middleware.CORSMiddleware)

	// Register Handlers
	handler.NewUserHandler(dbpool).RegisterRoutes(r)
	handler.NewAuthHandler(dbpool).RegisterRoutes(r)
	handler.NewConferenceHandler(dbpool).RegisterRoutes(r)
	handler.NewBookingHandler(dbpool).RegisterRoutes(r)
	handler.NewTicketHandler(dbpool).RegisterRoutes(r)

	// Run Server with Graceful Shutdown
	srv := &http.Server{
		Addr:         ":8080",
		Handler:      r,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Graceful shutdown
	go func() {
		log.Println("Server started on http://localhost:8080")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server failed: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit

	log.Println("Shutting down server...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exited properly")
}
