package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	customhttp "student-cert-service/internal/delivery/http"
	"student-cert-service/internal/infrastructure"
	"student-cert-service/internal/repository"
	"student-cert-service/internal/usecase"
)

func main() {
	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		log.Fatal("DB_URL environment variable is required")
	}
	rabbitURL := os.Getenv("RABBITMQ_URL")
	if rabbitURL == "" {
		log.Fatal("RABBITMQ_URL environment variable is required")
	}
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		log.Fatal("JWT_SECRET environment variable is required")
	}

	// Wait for DB and RabbitMQ to be ready (Docker compose healthcheck should handle this mostly, but good practice)
	time.Sleep(2 * time.Second)

	repo, err := repository.NewPostgresRepository(dbURL)
	if err != nil {
		log.Fatalf("failed to init repository: %v", err)
	}

	rabbitClient, err := infrastructure.NewRabbitMQClient(rabbitURL)
	if err != nil {
		log.Fatalf("failed to init rabbitmq: %v", err)
	}
	defer rabbitClient.Close()

	authUC := usecase.NewAuthUseCase(repo, jwtSecret)
	requestsUC := usecase.NewRequestsUseCase(repo, rabbitClient)
	handlers := customhttp.NewHandlers(authUC, requestsUC)

	mux := http.NewServeMux()

	// Metrics
	mux.Handle("GET /metrics", promhttp.Handler())

	// Public routes
	mux.HandleFunc("POST /register", handlers.Register)
	mux.HandleFunc("POST /login", handlers.Login)

	// Protected routes
	authMiddleware := customhttp.AuthMiddleware(jwtSecret)

	// In Go 1.22, we can wrap specific routes, but since we are using a standard mux, 
	// we will wrap the handler functions for protected routes.
	mux.Handle("POST /request", authMiddleware(http.HandlerFunc(handlers.CreateRequest)))
	mux.Handle("GET /requests", authMiddleware(http.HandlerFunc(handlers.GetRequests)))

	// Apply CORS middleware to the whole mux
	handlerWithCORS := customhttp.CORSMiddleware(mux)

	server := &http.Server{
		Addr:    ":8080",
		Handler: handlerWithCORS,
	}

	go func() {
		log.Println("Starting server on :8080")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("failed to start server: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}
	log.Println("Server exiting")
}
