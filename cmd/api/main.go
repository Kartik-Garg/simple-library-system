package main

import (
	"log"
	"os"

	"github.com/Kartik-Garg/simple-library-go/internal/handler"
	"github.com/Kartik-Garg/simple-library-go/internal/repo"
	"github.com/Kartik-Garg/simple-library-go/internal/service"
	"github.com/joho/godotenv"
)

func main() {

	// Load env variables
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	// ✅ REAL FACTORY PATTERN
	// Factory decides which database implementation to use based on DB_TYPE env variable
	// User doesn't need to know about postgres, mongo, or inmemory implementations
	// Just change DB_TYPE in .env to switch databases - no code changes needed!
	bookRepo := repo.NewBookRepository()

	// Dependency Injection - inject repository into service
	bookService := service.NewBookService(bookRepo)

	// Dependency Injection - inject service into handler
	bookHandler := handler.NewBookHandler(bookService)

	// Setup router with handler
	router := handler.SetUpRouter(bookHandler)

	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server starting on port %s..", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatal("Failed to start server: ", err)
	}

}
