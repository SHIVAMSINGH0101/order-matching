package main

import (
	"log"
	"net/http"

	"github.com/SHIVAMSINGH0101/go-demo/internal/config"
	"github.com/SHIVAMSINGH0101/go-demo/internal/database"
	"github.com/SHIVAMSINGH0101/go-demo/internal/handlers"
	"github.com/SHIVAMSINGH0101/go-demo/internal/repository"
	"github.com/SHIVAMSINGH0101/go-demo/internal/services"
	"github.com/gorilla/mux"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatal("Failed to load configuration:", err)
	}

	// Initialize database connection
	db, err := database.NewConnection(cfg.Database)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer db.Close()

	router := mux.NewRouter()

	api := router.PathPrefix("/api/v1").Subrouter()

	// Initialize repository layer
	routeRepo := repository.NewOrderRepository(db)

	// Initialize service layer
	routeService := services.NewOrderService(routeRepo)

	// Initialize handlers
	orderHandler := handlers.NewOrderHandler(routeService)
	orderHandler.RegisterOrderHandlers(api)

	// Health check
	router.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	}).Methods("GET")

	log.Printf("Server starting on port %s", cfg.Server.Port)
	log.Fatal(http.ListenAndServe(":"+cfg.Server.Port, router))
}
