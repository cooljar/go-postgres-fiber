package http

import (
	"github.com/cooljar/go-postgres-fiber/domain"
	"github.com/cooljar/go-postgres-fiber/utils"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"strconv"
)

// BookHandler  represent the httphandler for book
type BookHandler struct {
	BookUsecase domain.BookUsecase
	Validate *validator.Validate
}

func NewBookHandler(app *fiber.App, bookUseCase domain.BookUsecase, rPublic, rPrivate fiber.Router) {
	handler := &BookHandler{
		BookUsecase: bookUseCase,
		Validate: utils.NewValidator(),
	}

	rBook := rPublic.Group("/book")
	rBook.Get("/", handler.FetchBooks)
	rBook.Get("/:id", handler.GetByID)

	rAuthBook := rPrivate.Group("/book")
	rAuthBook.Post("/", handler.Create)
	rAuthBook.Put("/:id", handler.Update)
	rAuthBook.Delete("/:id", handler.Delete)
}

// Create func for creates a new book.
// @Summary create a new book
// @Description Create a new book.
// @Tags Book
// @Accept json
// @Produce json
// @Param book body domain.BookForm true "Add Book"
// @Success 200 {object} domain.JSONResult{data=[]domain.Book,message=string} "Description"
// @Failure 422 {object} []domain.HTTPError
// @Failure 400 {object} domain.HTTPError
// @Failure 404 {object} domain.HTTPError
// @Failure 500 {object} domain.HTTPError
// @Security ApiKeyAuth
// @Router /v1/auth/book [post]
func (b *BookHandler) Create(c *fiber.Ctx) error {
	// Instantiate new Book struct
	book := new(domain.BookForm)

	//  Parse body into application struct
	if err := c.BodyParser(book); err != nil {
		return domain.NewHttpError(c, err)
	}

	// Validate form input
	err := b.Validate.Struct(book)
	if err != nil {
		return domain.NewHttpError(c, err)
	}

	newBook, err := b.BookUsecase.Create(book)
	if err != nil {
		return domain.NewHttpError(c, err)
	}

	return c.JSON(domain.JSONResult{
		Data: newBook,
		Message: "Success",
	})
}

// FetchBooks func gets all exists books.
// @Description Get all exists books.
// @Summary get all exists books
// @Tags Book
// @Produce json
// @Param title query string false "search by title"
// @Param page query string false "page to display, default to 1"
// @Param perPage query string false "num of records per page, default to 20"
// @Success 200 {object} domain.JSONResult{data=[]domain.Book,meta=domain.JSONResultMeta,message=string} "Description"
// @Router /v1/book [get]
func (b *BookHandler) FetchBooks(c *fiber.Ctx) error {
	perPage, err := strconv.Atoi(c.Query("perPage"))
	if err != nil {
		perPage = 20
	}

	page, err := strconv.Atoi(c.Query("page"))
	if err != nil {
		page = 1
	}

	books, totalCount, pageCount, currentPage, err := b.BookUsecase.Fetch(perPage, page)

	if err != nil {
		return domain.NewHttpError(c, err)
	}

	return c.JSON(domain.JSONResult{Data: books, Message: "Success", Meta: domain.JSONResultMeta{TotalCount: totalCount, PageCount: pageCount, CurrentPage: currentPage, PerPage: perPage}})
}

// GetByID godoc
// @Summary Show an book
// @Description get book by ID
// @Tags Book
// @Produce  json
// @Param id path int true "Book ID"
// @Success 200 {object} domain.JSONResult{data=domain.Book,message=string}
// @Failure 400 {object} domain.HTTPError
// @Failure 404 {object} domain.HTTPError
// @Failure 500 {object} domain.HTTPError
// @Router /v1/book/{id} [get]
func (b *BookHandler) GetByID(c *fiber.Ctx) error {
	idBook, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return domain.NewHttpError(c, err)
	}

	book, err := b.BookUsecase.GetByID(idBook)
	if err != nil {
		return domain.NewHttpError(c, err)
	}

	return c.JSON(domain.JSONResult{Data: book, Message: "Success"})
}

// Update func for update a selected book.
// @Summary update a selected book
// @Description Update a selected book.
// @Tags Book
// @Accept json
// @Produce json
// @Param id path int true "Book ID"
// @Param book body domain.BookForm true "Update Book"
// @Success 200 {object} domain.JSONResult{data=[]domain.Book,message=string} "Description"
// @Failure 422 {object} []domain.HTTPError
// @Failure 400 {object} domain.HTTPError
// @Failure 404 {object} domain.HTTPError
// @Failure 500 {object} domain.HTTPError
// @Security ApiKeyAuth
// @Router /v1/auth/book/{id} [put]
func (b *BookHandler) Update(c *fiber.Ctx) error {
	idBook, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return domain.NewHttpError(c, err)
	}

	// Instantiate new Book struct
	book := new(domain.BookForm)

	//  Parse body into application struct
	if err := c.BodyParser(book); err != nil {
		return domain.NewHttpError(c, err)
	}

	// Validate form input
	err = b.Validate.Struct(book)
	if err != nil {
		return domain.NewHttpError(c, err)
	}

	updatedBook, err := b.BookUsecase.Update(idBook, book)
	if err != nil {
		return domain.NewHttpError(c, err)
	}

	return c.JSON(domain.JSONResult{Data: updatedBook, Message: "Success"})
}

// Delete godoc
// @Summary Delete an book
// @Description Delete by book ID
// @Tags Book
// @Accept  json
// @Produce  json
// @Param  id path int true "Book ID" Format(int64)
// @Success 200 {object} domain.JSONResult{data=string,message=string} "Description"
// @Security ApiKeyAuth
// @Failure 400 {object} domain.HTTPError
// @Failure 500 {object} domain.HTTPError
// @Router /v1/auth/book/{id} [delete]
func (b *BookHandler) Delete(c *fiber.Ctx) error {
	idBook, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return domain.NewHttpError(c, err)
	}

	rowsAffected, err := b.BookUsecase.Delete(idBook)
	if err != nil {
		return domain.NewHttpError(c, err)
	}

	if rowsAffected < 1 {
		return domain.NewHttpError(c, err)
	}

	return c.JSON(domain.JSONResult{Data: "deleted", Message: "Success"})
}
