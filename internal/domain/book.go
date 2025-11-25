package domain

import (
	"errors"
	"time"
)

// Book entity
type Book struct {
	// json for API response and gorm is for db columns
	ID        uint      `json:"id" gorm:"primaryKey"`
	Title     string    `json:"title" gorm:"not null"`
	Author    string    `json:"author" gorm:"not null"`
	Available int       `json:"available" gorm:"default:1"`
	Quantity  int       `json:"quantity" gorm:"default:1"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Entity basic methods/validations
func (b *Book) Validate() error {
	if b.Title == "" {
		return errors.New("title is mandatory")
	}
	if b.Author == "" {
		return errors.New("author is mandatory")
	}
	if b.Available < 0 {
		return errors.New("availability can not be negative")
	}
	if b.Available > b.Quantity {
		return errors.New("quantity is always more or equal to availability")
	}
	return nil
}

// what basic functions can this entity/model have? depending on its structure?

// we can check if book is available
func (b *Book) IsAvailable() bool {
	return b.Available > 0
}

// Another thing which DDD mentions is that, we define what we need and not how it is actually implemented
// So book domain needs a way to do CRUD operations. So, we need to define an interface here, for books in domain
// We can have repo interface defined in domain - This will also allow us to use repository design pattern
