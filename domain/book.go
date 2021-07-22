package domain

import (
	"encoding/json"
)

// BookForm form for create book
type BookForm struct {
	Title   string  `json:"title" validate:"required"`
	Author  string  `json:"author" validate:"required"`
	Content string  `json:"content" validate:"required"`
	Price   float64 `json:"price" validate:"gte=1"`
	Rating  int     `json:"rating" validate:"lte=5"`
}

// Book the book model
type Book struct {
	ID        int     `json:"id"`
	Title     string  `json:"title"`
	Author    string  `json:"author"`
	Content   string  `json:"content"`
	Price     float64 `json:"price"`
	CreatedAt int     `json:"created_at"`
	UpdatedAt int     `json:"updated_at"`
	Rating    int     `json:"rating"`
}

// FromJSON decode json to book struct
func (b *Book) FromJSON(msg []byte) error {
	return json.Unmarshal(msg, b)
}

// ToJSON encode book struct to json
func (b *Book) ToJSON() []byte {
	str, _ := json.Marshal(b)
	return str
}

// BookUsecase represent the book's use cases
type BookUsecase interface {
	Create(b *BookForm) (book Book, err error)
	Fetch(perPage, page int) (books []Book, totalCount, pageCount, currentPage int, err error)
	GetByID(id int) (Book, error)
	Update(id int, b *BookForm) (book Book, err error)
	Delete(id int) (rowsAffected int64, err error)
}

// BookRepository represent the book's repository
type BookRepository interface {
	Create(b *BookForm) (book Book, err error)
	Fetch(perPage, page int) (books []Book, totalCount, pageCount, currentPage int, err error)
	GetByID(id int) (Book, error)
	Update(id int, b *BookForm) (book Book, err error)
	Delete(id int) (rowsAffected int64, err error)
}
