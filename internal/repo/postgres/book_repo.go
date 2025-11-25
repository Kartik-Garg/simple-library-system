package postgres

import (
	"context"

	"github.com/Kartik-Garg/simple-library-go/internal/domain"
	"gorm.io/gorm"
)

// keeping it private, so constructor needs to be called
type bookRepository struct {
	db *gorm.DB
}

// actual implementation of interfaces
func NewBookRepository(db *gorm.DB) domain.BookRepository {
	return &bookRepository{
		db: db,
	}
}

func (r *bookRepository) Create(ctx context.Context, book *domain.Book) error {
	return r.db.WithContext(ctx).Create(book).Error
}

func (r *bookRepository) FindByID(ctx context.Context, id uint) (*domain.Book, error) {
	var book domain.Book
	err := r.db.WithContext(ctx).First(&book, id).Error
	if err != nil {
		return nil, err
	}
	return &book, nil
}

func (r *bookRepository) FindAll(ctx context.Context) ([]*domain.Book, error) {
	var books []*domain.Book
	err := r.db.WithContext(ctx).Find(&books).Error
	return books, err
}

func (r *bookRepository) Update(ctx context.Context, book *domain.Book) error {
	return r.db.WithContext(ctx).Save(book).Error
}

func (r *bookRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&domain.Book{}, id).Error
}
