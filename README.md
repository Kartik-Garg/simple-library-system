# Library Management System

A simple REST API for managing a library's book inventory, built with Go and PostgreSQL.

## What it does

This is a basic CRUD application where you can:
- Add new books to the library
- View all books or get details of a specific book
- Update book information
- Delete books from the system

Each book has a title, author, total quantity, and available count.

## Tech Stack

- **Go 1.25** - Backend language
- **Gin** - HTTP web framework
- **PostgreSQL** - Database
- **GORM** - ORM for database operations
- **Docker** - Containerization

## Project Structure

```
.
├── cmd/api/              # Application entry point
│   └── main.go
├── internal/
│   ├── domain/           # Core business logic and entities
│   ├── handler/          # HTTP handlers (controllers)
│   ├── service/          # Business logic layer
│   ├── repo/             # Data access layer
│   │   ├── postgres/     # PostgreSQL implementation
│   │   └── inmemory/     # In-memory implementation (for testing)
│   └── infra/db/         # Database connection setup
├── Dockerfile
├── docker-compose.yml
└── .env
```

I followed a clean architecture approach with clear separation between layers. The repository pattern makes it easy to swap database implementations, by just changing an environment variable.

## Getting Started

### Prerequisites

- Docker and Docker Compose installed

### Running the application

1. Clone the repo
```bash
git clone <your-repo-url>
cd simple-library-system
```

2. Start everything with Docker Compose
```bash
docker-compose up --build
```

This will:
- Build the Go application
- Start PostgreSQL database
- Run the API server on port 8080

3. Check if it's running
```bash
curl http://localhost:8080/health
```

You should see: `{"service":"library-api","status":"healthy"}`

### Without Docker (local development)

If you want to run it locally without Docker:

1. Make sure PostgreSQL is running and update `.env` with your database credentials

2. Install dependencies
```bash
go mod download
```

3. Run the application
```bash
go run cmd/api/main.go
```

## API Endpoints

Base URL: `http://localhost:8080/api/v1`

### Create a book
```bash
POST /books
Content-Type: application/json

{
  "title": "The Great Gatsby",
  "author": "F. Scott Fitzgerald",
  "quantity": 5
}
```

### Get all books
```bash
GET /books
```

### Get a specific book
```bash
GET /books/:id
```

### Update a book
```bash
PUT /books/:id
Content-Type: application/json

{
  "title": "The Great Gatsby",
  "author": "F. Scott Fitzgerald",
  "quantity": 10,
  "available": 8
}
```

### Delete a book
```bash
DELETE /books/:id
```

## Configuration

The app uses environment variables for configuration. Check `.env` file:

- `DB_TYPE` - Set to `postgres` or `inmemory`
- `DB_HOST`, `DB_PORT`, `DB_USER`, `DB_PASSWORD`, `DB_NAME` - PostgreSQL connection details
- `APP_PORT` - Port where the API runs (default: 8080)
- `GIN_MODE` - Set to `debug` or `release`

## Architecture & Design Patterns

### Domain-Driven Design (DDD)

I structured this project using DDD principles to keep the code organized and maintainable:

**Domain Layer** (`internal/domain/`)
- Contains the core business entities (`Book`) and business rules
- Defines the `BookRepository` interface - this is key to DDD's dependency inversion
- Has validation logic (`Validate()`) and domain methods (`IsAvailable()`)
- No dependencies on external frameworks or databases

**Application Layer** (`internal/service/`)
- Orchestrates business operations
- Uses the repository interface (doesn't care about implementation details)
- Handles business logic like setting `available = quantity` when creating books

**Infrastructure Layer** (`internal/repo/`, `internal/infra/`)
- Implements the repository interface with actual database code
- PostgreSQL implementation using GORM
- In-memory implementation for testing
- Database connection setup

**Interface Layer** (`internal/handler/`)
- HTTP handlers that expose the API
- Converts HTTP requests to domain operations
- Returns appropriate HTTP responses

### Design Patterns Used

**1. Repository Pattern** (`internal/domain/repository.go`, `internal/repo/`)

The repository abstracts data access. The interface lives in the domain layer, implementations live in infrastructure.

```go
// Domain defines what it needs
type BookRepository interface {
    Create(ctx context.Context, book *Book) error
    FindByID(ctx context.Context, id uint) (*Book, error)
    // ...
}

// Infrastructure provides it
type bookRepository struct {
    db *gorm.DB
}
```

Why? The domain doesn't depend on GORM or PostgreSQL. I can swap databases without touching business logic.

**2. Factory Pattern** (`internal/repo/factory.go`)

The factory creates the right repository implementation based on configuration:

```go
func NewBookRepository() domain.BookRepository {
    dbType := os.Getenv("DB_TYPE")
    switch dbType {
    case "postgres":
        return createPostgresRepository()
    case "inmemory":
        return createInMemoryRepository()
    }
}
```

Why? Change `DB_TYPE` in `.env` and you get a different implementation. No code changes needed.

**3. Dependency Injection** (`cmd/api/main.go`)

Dependencies are injected through constructors:

```go
bookRepo := repo.NewBookRepository()
bookService := service.NewBookService(bookRepo)
bookHandler := handler.NewBookHandler(bookService)
```

Why? Makes testing easier (can inject mocks) and keeps components loosely coupled.

**4. Service Layer Pattern** (`internal/service/`)

Business logic sits between handlers and repositories:

```go
func (s *BookService) CreateBook(ctx context.Context, book *domain.Book) error {
    if err := book.Validate(); err != nil {
        return err
    }
    book.Available = book.Quantity  // Business rule
    return s.repo.Create(ctx, book)
}
```

Why? Keeps business logic out of HTTP handlers and makes it reusable.

### Why This Architecture?

**Separation of Concerns**
Each layer has a single responsibility. Domain layer doesn't know about HTTP or databases. Handlers don't know about database queries.

**Testability**
Can test business logic without a database. Can test handlers without real services. Just inject mocks.

**Flexibility**
Want to add a MongoDB implementation? Just create a new repository. Want to add gRPC alongside REST? Add new handlers. Core logic stays untouched.

**Maintainability**
When requirements change, you know exactly where to look. Database query issue? Check repository. Business rule change? Check service layer.

## Stopping the application

```bash
# Stop containers
docker-compose down

# Stop and remove all data
docker-compose down -v
```

## Notes

- The in-memory implementation is useful for quick testing without spinning up a database
- Health check endpoint (`/health`) is used by Docker to monitor the application
- Database migrations happen automatically on startup using GORM's AutoMigrate
- The API uses proper HTTP status codes (200, 201, 400, 404, 500)

## Future improvements

Some things I'd add if I had more time:
- Unit and integration tests
- Pagination for the GET all books endpoint
- More comprehensive error handling

