// cmd/server/main.go
package main

import (
	"context"
	"log"
	"time"
	"os"
	"hesh-core/internal/db"
	"hesh-core/internal/repository"
)

func main() {
	// Create root context for application
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	dbURL := os.Getenv("DATABASE_URL")

	if dbURL == "" {
		// Fallback
		dbURL = "postgres://admin:password@localhost:5433/hesh_core_db?sslmode=disable"
		log.Println("DATABASE_URL not set, falling back to local host connection")
	}
	

	log.Println("Waking up hesh-core-service... Allocating connection pool...")

	
	dbPool, err := db.Connect(ctx, dbURL)
	if err != nil {
		log.Fatalf("Fatal: Could not connect to database container: %v", err)
	}
	defer dbPool.Close() // Safely close all connection pipelines when the service stops

	log.Println("Database connection pool allocated successfully!")

	// 4. HAND OUT THE SHORTCUTS (Pointers)
	// We pass the exact same 'dbPool' memory reference into every constructor
	userRepo := repository.NewPostgresUserRepository(dbPool)
	merchantRepo := repository.NewPostgresMerchantRepository(dbPool)
	programRepo := repository.NewPostgresProgramRepository(dbPool)
	transactionRepo := repository.NewPostgresTransactionRepository(dbPool)

	log.Printf("Repositories wired up successfully. UserRepo memory link: %p", userRepo)

	// pass repo to route handlers
	_ = merchantRepo
	_ = programRepo
	_ = transactionRepo
	_ = userRepo
}