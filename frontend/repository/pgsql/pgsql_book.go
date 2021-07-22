package pgsql

import (
	"context"
	"github.com/cooljar/go-postgres-fiber/domain"
	"github.com/jackc/pgx/v4/pgxpool"
	"time"
)

type pgsqlBookRepository struct {
	Conn *pgxpool.Pool
}

// NewPgsqlBookRepository will create an object that represent the application.Repository interface
func NewPgsqlBookRepository(conn *pgxpool.Pool) domain.BookRepository {
	return &pgsqlBookRepository{Conn: conn}
}

func (m *pgsqlBookRepository) Create(b *domain.BookForm) (domain.Book, error) {
	var book domain.Book
	var id int

	ts := time.Now().Unix()
	qStr := `insert into books (title, content, author, price, rating, created_at, updated_at) VALUES ($1,$2,$3,$4,$5,$6,$7) returning id`
	err := m.Conn.QueryRow(context.Background(), qStr, b.Title, b.Content, b.Author, b.Price, b.Rating, ts, ts).Scan(&id)
	if err != nil {
		return book, err
	}

	book.ID = id
	book.Title = b.Title
	book.Author = b.Author
	book.Content = b.Content
	book.Price = b.Price
	book.CreatedAt = int(ts)
	book.UpdatedAt = int(ts)
	book.Rating = b.Rating

	return book, nil
}

func (m *pgsqlBookRepository) Fetch(perPage, page int) (books []domain.Book, totalCount, pageCount, currentPage int, err error) {
	err = m.Conn.QueryRow(context.Background(),"SELECT COUNT(*) FROM books").Scan(&totalCount)
	if err != nil {
		return
	}

	pageCount = totalCount / perPage
	if page > pageCount {
		page = pageCount
	}

	offset := perPage * (page-1)
	qStr := `SELECT id,title,author,content,price,created_at,updated_at,rating FROM books ORDER BY $1 DESC LIMIT $2 OFFSET $3`
	rows, err := m.Conn.Query(context.Background(), qStr, "created_at", perPage, offset)
	if err != nil {
		return
	}

	for rows.Next() {
		var b domain.Book
		var title, author, content string
		var id, rating, createdAt, updatedAt int
		var price float64
		err = rows.Scan(&id, &title, &author, &content, &price, &createdAt, &updatedAt, &rating)
		if err != nil {
			return
		}

		b.ID = id
		b.Title = title
		b.Author = author
		b.Content = content
		b.Price = price
		b.CreatedAt = createdAt
		b.UpdatedAt = updatedAt
		b.Rating = rating
		books = append(books, b)
	}

	return books, totalCount, pageCount, page, nil
}

func (m *pgsqlBookRepository) GetByID(bookId int) (domain.Book, error) {
	b := domain.Book{}
	var title, author, content string
	var id, rating, createdAt, updatedAt int
	var price float64
	err := m.Conn.QueryRow(context.Background(), "SELECT id,title,author,content,price,created_at,updated_at,rating FROM books WHERE id=$1 LIMIT 1", bookId).Scan(&id, &title, &author, &content, &price, &createdAt, &updatedAt, &rating)
	if err != nil {
		return b, err
	}

	b.ID = id
	b.Title = title
	b.Author = author
	b.Content = content
	b.Price = price
	b.CreatedAt = createdAt
	b.UpdatedAt = updatedAt
	b.Rating = rating

	return b, nil
}

func (m *pgsqlBookRepository) Update(bookId int, b *domain.BookForm) (book domain.Book, err error) {
	var title, author, content string
	var id, rating, createdAt, updatedAt int
	var price float64

	err = m.Conn.QueryRow(context.Background(), "SELECT id,title,author,content,price,created_at,updated_at,rating FROM books WHERE id=$1 LIMIT 1", bookId).Scan(&id, &title, &author, &content, &price, &createdAt, &updatedAt, &rating)
	if err != nil {
		return
	}

	book = domain.Book{
		id,
		title,
		author,
		content,
		price,
		createdAt,
		updatedAt,
		rating,
	}

	book.Title = b.Title
	book.Author = b.Author
	book.Content = b.Content
	book.Price = b.Price
	book.Rating = b.Rating

	qCmd := `UPDATE books SET title=$1, author=$2, content=$3, price=$4, rating=$5 WHERE id=$6`
	_, err = m.Conn.Exec(context.Background(), qCmd, book.Title, book.Author, book.Content, book.Price, book.Rating, book.ID)
	if err != nil {
		return
	}

	return
}

func (m *pgsqlBookRepository) Delete(id int) (rowsAffected int64, err error) {
	qCmd := `DELETE FROM books WHERE id=$1`
	res, err := m.Conn.Exec(context.Background(), qCmd, id)
	if err != nil {
		return 0, err
	}

	return res.RowsAffected(), nil
}
