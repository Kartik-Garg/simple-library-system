package repo

import (
	"fmt"
	"log"
	"os"

	"github.com/Kartik-Garg/simple-library-go/internal/domain"
	"github.com/Kartik-Garg/simple-library-go/internal/infra/db"
	"github.com/Kartik-Garg/simple-library-go/internal/repo/inmemory"
	"github.com/Kartik-Garg/simple-library-go/internal/repo/postgres"
)

// DatabaseType represents the type of database to use
type DatabaseType string

const (
	Postgres DatabaseType = "postgres"
	InMemory DatabaseType = "inmemory"
	// MongoDB  DatabaseType = "mongo"  // Just showing extensibility here, as we are factory pattern here
	// in future if we decide to use some other type of DB, it will be easier to plug it in.
)

// NewBookRepository uses Factory Pattern
// It decides which implementation to create based on configuration
// The caller doesn't need to know about specific implementations
func NewBookRepository() domain.BookRepository {
	// Read database type from environment variable
	dbType := DatabaseType(os.Getenv("DB_TYPE"))

	// Default to postgres if not specified
	if dbType == "" {
		dbType = Postgres
		log.Println("DB_TYPE not set, defaulting to postgres")
	}

	log.Printf("Creating repository for database type: %s", dbType)

	// Factory decides which implementation to create
	switch dbType {
	case Postgres:
		return createPostgresRepository()

	case InMemory:
		return createInMemoryRepository()

	// case MongoDB:
	//     return createMongoRepository()  // Future extensibility

	default:
		panic(fmt.Sprintf("unsupported database type: %s. Supported types: postgres, inmemory", dbType))
	}
}

// createPostgresRepository creates a PostgreSQL repository
// lets keep this private so only the factory can call it
func createPostgresRepository() domain.BookRepository {
	// Create database connection
	database := db.NewPostgresDB()

	// Auto-migrate schema
	if err := database.AutoMigrate(&domain.Book{}); err != nil {
		log.Fatal("Failed to migrate database:", err)
	}

	// Create and return repository
	return postgres.NewBookRepository(database)
}

// createInMemoryRepository creates an in-memory repository
// Useful for testing or development without a database
func createInMemoryRepository() domain.BookRepository {
	log.Println("Using in-memory repository (no database required)")
	return inmemory.NewBookRepository()
}

// Future: createMongoRepository for MongoDB implementation
// func createMongoRepository() domain.BookRepository {
//     database := db.NewMongoDB()
//     return mongo.NewBookRepository(database)
// }
