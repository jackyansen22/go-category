package main

import (
	"log"
	"net/http"
	"os"

	"github.com/jackyansen22/crud-category/internal/config"
	"github.com/jackyansen22/crud-category/internal/database"
	"github.com/jackyansen22/crud-category/internal/handler"
	"github.com/jackyansen22/crud-category/internal/repository"
	"github.com/jackyansen22/crud-category/internal/service"
)

func main() {
	log.Println("🚀 CATEGORY API STARTED (RAILWAY)")

	// ===== PORT (Railway first, fallback local) =====
	port := os.Getenv("PORT")
	if port == "" {
		port = os.Getenv("APP_PORT")
	}
	if port == "" {
		port = "8080"
	}

	// ===== CONFIG =====
	cfg := config.Load()

	// ===== ROUTE ROOT =====
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{
			"service":"Category API",
			"status":"running",
			"health":"/health",
			"categories":"/categories"
		}`))
	})

	// ===== HEALTH =====
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"OK","message":"API Running"}`))
	})

	// ===== DATABASE =====
	db, err := database.Connect(cfg.DBUrl)
	if err != nil {
		log.Println("❌ DB connection failed:", err)
		log.Println("⚠️ Category routes DISABLED")
	} else {
		log.Println("✅ Connected to database")

		repo := repository.NewCategoryRepository(db)
		svc := service.NewCategoryService(repo)
		h := handler.NewCategoryHandler(svc)

		http.HandleFunc("/categories", h.Categories)
		http.HandleFunc("/categories/", h.CategoryByID)

		productRepo := repository.NewProductRepository(db)
		productSvc := service.NewProductService(productRepo)
		productHandler := handler.NewProductHandler(productSvc)

		//http.HandleFunc("/api/produk", productHandler.Products)
		//http.HandleFunc("/api/produk/", productHandler.ProductByID)
		http.HandleFunc("/product", productHandler.Products)
		http.HandleFunc("/product/", productHandler.ProductByID)

		// =====================
		// Transaction (Checkout)
		// =====================
		transactionRepo := repository.NewTransactionRepository(db)
		transactionService := service.NewTransactionService(transactionRepo)
		transactionHandler := handler.NewTransactionHandler(transactionService)

		// POST /api/checkout
		http.HandleFunc("/checkout", transactionHandler.Checkout)

		// Transactions (read only)
		http.HandleFunc("/transactions", transactionHandler.GetAll)
		http.HandleFunc("/transactions/", transactionHandler.GetByID)

	}

	log.Println("🌐 Listening on :" + port)

	log.Fatal(
		http.ListenAndServe(
			":"+port,
			handler.RecoverMiddleware(http.DefaultServeMux),
		),
	)
}
