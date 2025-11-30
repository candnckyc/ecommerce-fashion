package main

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	
	"ecommerce-backend/internal/config"
	"ecommerce-backend/internal/database"
	"ecommerce-backend/internal/handlers"
	"ecommerce-backend/internal/middleware"
	"ecommerce-backend/internal/repository"
	"ecommerce-backend/internal/services"
)

func main() {
	log.Println("🚀 Starting E-Commerce API Server...")

	// Load configuration
	cfg := config.LoadConfig()

	// Connect to database
	db := database.Connect(cfg.DatabaseURL)
	defer db.Close()

	// Initialize router
	router := mux.NewRouter()

	// Apply CORS middleware
	router.Use(middleware.CORS(cfg.AllowedOrigins))

	// Setup routes
	setupRoutes(router, db, cfg)

	// Start server
	addr := ":" + cfg.Port
	log.Printf("✓ Server running on http://localhost%s", addr)
	log.Fatal(http.ListenAndServe(addr, router))
}

func setupRoutes(router *mux.Router, db *sql.DB, cfg *config.Config) {
	// Initialize repositories
	userRepo := repository.NewUserRepository(db)
	productRepo := repository.NewProductRepository(db)
	cartRepo := repository.NewCartRepository(db)
	orderRepo := repository.NewOrderRepository(db)

	// Initialize services
	authService := services.NewAuthService(userRepo, cfg.JWTSecret)
	productService := services.NewProductService(productRepo)
	cartService := services.NewCartService(cartRepo, productRepo)
	orderService := services.NewOrderService(orderRepo, cartRepo, productRepo)

	// Initialize handlers
	healthHandler := handlers.NewHealthHandler()
	authHandler := handlers.NewAuthHandler(authService)
	productHandler := handlers.NewProductHandler(productService)
	cartHandler := handlers.NewCartHandler(cartService)
	orderHandler := handlers.NewOrderHandler(orderService)
	adminHandler := handlers.NewAdminHandler(productService, orderService, userRepo)

	// Health check
	router.HandleFunc("/health", healthHandler.Check).Methods("GET", "OPTIONS")

	// API routes
	api := router.PathPrefix("/api").Subrouter()
	
	// Public routes (no auth required)
	api.HandleFunc("/health", healthHandler.Check).Methods("GET", "OPTIONS")
	api.HandleFunc("/auth/register", authHandler.Register).Methods("POST", "OPTIONS")
	api.HandleFunc("/auth/login", authHandler.Login).Methods("POST", "OPTIONS")
	
	// Product routes (public)
	api.HandleFunc("/products", productHandler.GetProducts).Methods("GET", "OPTIONS")
	api.HandleFunc("/products/{id}", productHandler.GetProduct).Methods("GET", "OPTIONS")
	api.HandleFunc("/brands", productHandler.GetBrands).Methods("GET", "OPTIONS")
	api.HandleFunc("/categories", productHandler.GetCategories).Methods("GET", "OPTIONS")

	// Protected routes (auth required)
	protected := api.PathPrefix("").Subrouter()
	protected.Use(middleware.AuthMiddleware(cfg.JWTSecret))
	
	// Auth protected routes
	protected.HandleFunc("/auth/me", authHandler.Me).Methods("GET", "OPTIONS")
	
	// Cart routes (protected)
	protected.HandleFunc("/cart", cartHandler.GetCart).Methods("GET", "OPTIONS")
	protected.HandleFunc("/cart", cartHandler.AddToCart).Methods("POST", "OPTIONS")
	protected.HandleFunc("/cart/{id}", cartHandler.UpdateCartItem).Methods("PUT", "OPTIONS")
	protected.HandleFunc("/cart/{id}", cartHandler.RemoveFromCart).Methods("DELETE", "OPTIONS")
	protected.HandleFunc("/cart/clear", cartHandler.ClearCart).Methods("DELETE", "OPTIONS")
	
	// Order routes (protected)
	protected.HandleFunc("/orders", orderHandler.CreateOrder).Methods("POST", "OPTIONS")
	protected.HandleFunc("/orders", orderHandler.GetOrders).Methods("GET", "OPTIONS")
	protected.HandleFunc("/orders/{id}", orderHandler.GetOrder).Methods("GET", "OPTIONS")
	
	// Address routes (protected)
	protected.HandleFunc("/addresses", orderHandler.GetAddresses).Methods("GET", "OPTIONS")
	protected.HandleFunc("/addresses", orderHandler.CreateAddress).Methods("POST", "OPTIONS")
	
	// Admin routes (protected + admin role check)
	admin := api.PathPrefix("/admin").Subrouter()
	admin.Use(middleware.AuthMiddleware(cfg.JWTSecret))
	admin.Use(middleware.AdminMiddleware())
	
	admin.HandleFunc("/stats", adminHandler.GetStats).Methods("GET", "OPTIONS")
	admin.HandleFunc("/orders", adminHandler.GetAllOrders).Methods("GET", "OPTIONS")
	admin.HandleFunc("/orders/{id}/status", adminHandler.UpdateOrderStatus).Methods("PUT", "OPTIONS")
	admin.HandleFunc("/products/{id}/toggle", adminHandler.ToggleProduct).Methods("PUT", "OPTIONS")
	admin.HandleFunc("/customers", adminHandler.GetAllCustomers).Methods("GET", "OPTIONS")
	
	log.Println("✓ Routes configured")
	log.Println("  POST /api/auth/register")
	log.Println("  POST /api/auth/login")
	log.Println("  GET  /api/auth/me (protected)")
	log.Println("  GET  /api/products")
	log.Println("  GET  /api/products/{id}")
	log.Println("  GET  /api/brands")
	log.Println("  GET  /api/categories")
	log.Println("  GET  /api/cart (protected)")
	log.Println("  POST /api/cart (protected)")
	log.Println("  PUT  /api/cart/{id} (protected)")
	log.Println("  DELETE /api/cart/{id} (protected)")
	log.Println("  GET  /api/admin/stats (admin)")
	log.Println("  GET  /api/admin/orders (admin)")
}