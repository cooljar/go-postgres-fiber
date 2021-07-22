package usecase

import (
	"github.com/cooljar/go-postgres-fiber/domain"
	"time"
)

type bookUsecase struct {
	bookRepo       domain.BookRepository
	contextTimeout time.Duration
}

// NewBookUsecase will create new an bookUsecase object representation of domain.BookUsecase interface
func NewBookUsecase(b domain.BookRepository, timeout time.Duration) domain.BookUsecase {
	return &bookUsecase{
		bookRepo:       b,
		contextTimeout: timeout,
	}
}

func (b *bookUsecase) Create(bd *domain.BookForm) (book domain.Book, err error) {
	/*ctx, cancel := context.WithTimeout(c, b.contextTimeout)
	defer cancel()*/

	//blokir jika judul buku sudah ada
	/*existedBook, _ := b.GetByTitle(ctx, bd.Title)
	if existedBook != (domain.Book{}) {
		return domain.ErrConflict
	}*/

	book, err = b.bookRepo.Create(bd)
	return
}

func (b *bookUsecase) Fetch(perPage, page int) (books []domain.Book, totalCount, pageCount, currentPage int, err error) {
	books, totalCount, pageCount, currentPage, err = b.bookRepo.Fetch(perPage, page)
	if err != nil {
		return
	}

	return
}

func (b *bookUsecase) GetByID(id int) (book domain.Book, err error) {
	book, err = b.bookRepo.GetByID(id)
	return
}

func (b *bookUsecase) Update(id int, bf *domain.BookForm) (book domain.Book, err error) {
	book, err = b.bookRepo.Update(id, bf)
	return
}

func (b *bookUsecase) Delete(id int) (rowsAffected int64, err error) {
	return b.bookRepo.Delete(id)
}
