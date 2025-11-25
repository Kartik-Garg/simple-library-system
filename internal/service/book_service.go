package service

import (
	"context"
	"errors"

	"github.com/Kartik-Garg/simple-library-go/internal/domain"
)

// Q) We are taking interface here, which allows us to do DI?
// What exactly is happening here? is it DI or something else?
type BookService struct {
	repo domain.BookRepository
}

// constructor of the service. We are also doing DI here
func NewBookService(repo domain.BookRepository) *BookService {
	return &BookService{
		repo: repo,
	}
}

// Services usually handles business logic and also interacts with repo layers
// So the functionality of services is declared here. That is, what exactly can a service do (use cases)

func (s *BookService) CreateBook(ctx context.Context, book *domain.Book) error {
	// first we do validation
	if err := book.Validate(); err != nil {
		return err
	}

	// set availablity as well
	book.Available = book.Quantity

	return s.repo.Create(ctx, book)
}

func (s *BookService) GetBook(ctx context.Context, id uint) (*domain.Book, error) {
	book, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, errors.New("book not found")
	}
	return book, nil
}

func (s *BookService) GetAllBooks(ctx context.Context) ([]*domain.Book, error) {
	return s.repo.FindAll(ctx)
}

func (s *BookService) UpdateBook(ctx context.Context, id uint, updatedBook *domain.Book) error {
	// Check if book exists
	existing, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return errors.New("book not found")
	}

	// Update fields
	// Can we update only few fields and not all
	existing.Title = updatedBook.Title
	existing.Author = updatedBook.Author
	existing.Quantity = updatedBook.Quantity
	existing.Available = updatedBook.Available

	// Validate
	// After making the local changes, before making transaction with db, we should validate
	// if the given values are actually valid/legal or not
	if err := existing.Validate(); err != nil {
		return err
	}

	return s.repo.Update(ctx, existing)
}

func (s *BookService) DeleteBook(ctx context.Context, id uint) error {
	// Check if book is even present
	if _, err := s.repo.FindByID(ctx, id); err != nil {
		return errors.New("book not present or already deleted")
	}

	return s.repo.Delete(ctx, id)
}
