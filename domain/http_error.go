package domain

import (
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

// HTTPError example
type HTTPError struct {
	Field    string    `json:"field"`
	Message string `json:"message"`
}

// NewHttpError custom error response
func NewHttpError(ctx *fiber.Ctx, err error) error {
	status := fiber.StatusInternalServerError

	if _, ok1 := err.(DataValidationError); ok1 {
		status = fiber.StatusUnprocessableEntity
		return ctx.Status(status).JSON(HTTPError{Field: err.(DataValidationError).Field, Message: err.Error()})
	}else if _, ok2 := err.(validator.ValidationErrors); ok2 {
		var fields []HTTPError
		for _, err := range err.(validator.ValidationErrors) {
			/*fmt.Println(err.Namespace())
			fmt.Println(err.Field())
			fmt.Println(err.StructNamespace())
			fmt.Println(err.StructField())
			fmt.Println(err.Tag())
			fmt.Println(err.ActualTag())
			fmt.Println(err.Kind())
			fmt.Println(err.Type())
			fmt.Println(err.Value())
			fmt.Println(err.Param())*/

			fields = append(fields, HTTPError{Field: err.Field(), Message: err.Tag()})
		}
		return ctx.Status(fiber.StatusBadRequest).JSON(fields)
	}

	return ctx.Status(status).JSON(fiber.Map{"code":"", "message":err.Error()})
}
