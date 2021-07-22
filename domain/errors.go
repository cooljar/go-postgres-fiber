package domain

import "errors"

var (
	ErrSqlNoRowFoundToUpsert = errors.New("no row found to insert or update")
)

type DataValidationError struct {
	Field string
	Message string
}
func (m DataValidationError) Error() string {
	return m.Message
}

type NotFoundError struct {
	Code int
	Message string
}
func (nf NotFoundError) Error() string {
	return nf.Message
}