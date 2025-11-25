package inmemory

import (
	"context"
	"errors"
	"sync"

	"github.com/Kartik-Garg/simple-library-go/internal/domain"
)

// bookRepository is an in-memory implementation of domain.BookRepository
// Useful for testing or development without a database
type bookRepository struct {
	mu      sync.RWMutex          // Protects concurrent access
	books   map[uint]*domain.Book // In-memory storage
	nextID  uint                  // Auto-increment ID
}

// NewBookRepository creates a new in-memory repository
func NewBookRepository() domain.BookRepository {
	return &bookRepository{
		books:  make(map[uint]*domain.Book),
		nextID: 1,
	}
}

// Create adds a new book to the in-memory store
func (r *bookRepository) Create(ctx context.Context, book *domain.Book) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	// Assign auto-increment ID
	book.ID = r.nextID
	r.nextID++

	// Store a copy to avoid external modifications
	bookCopy := *book
	r.books[book.ID] = &bookCopy

	return nil
}

// FindByID retrieves a book by ID from the in-memory store
func (r *bookRepository) FindByID(ctx context.Context, id uint) (*domain.Book, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	book, exists := r.books[id]
	if !exists {
		return nil, errors.New("book not found")
	}

	// Return a copy to avoid external modifications
	bookCopy := *book
	return &bookCopy, nil
}

// FindAll retrieves all books from the in-memory store
func (r *bookRepository) FindAll(ctx context.Context) ([]*domain.Book, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	books := make([]*domain.Book, 0, len(r.books))
	for _, book := range r.books {
		// Return copies to avoid external modifications
		bookCopy := *book
		books = append(books, &bookCopy)
	}

	return books, nil
}

// Update modifies an existing book in the in-memory store
func (r *bookRepository) Update(ctx context.Context, book *domain.Book) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.books[book.ID]; !exists {
		return errors.New("book not found")
	}

	// Store a copy to avoid external modifications
	bookCopy := *book
	r.books[book.ID] = &bookCopy

	return nil
}

// Delete removes a book from the in-memory store
func (r *bookRepository) Delete(ctx context.Context, id uint) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.books[id]; !exists {
		return errors.New("book not found")
	}

	delete(r.books, id)
	return nil
}

