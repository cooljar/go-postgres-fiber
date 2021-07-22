package http

import (
	"github.com/cooljar/go-postgres-fiber/domain"
	"github.com/cooljar/go-postgres-fiber/utils"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt"
)

type UserHandler struct {
	UserUsecase domain.UserUseCase
	Validate *validator.Validate
}

// NewUserHandler will initialize the /user resources endpoint
func NewUserHandler(app *fiber.App, validator *validator.Validate, userUseCase domain.UserUseCase, rPublic, rPrivate fiber.Router) {
	handler := &UserHandler{
		UserUsecase: userUseCase,
		Validate: validator,
	}

	rUser := rPublic.Group("/user")
	rUser.Post("/email-validation-secret", handler.RequestSecret)
	rUser.Post("/signup", handler.Signup)
	rUser.Post("/login", handler.Login)

	rUserPrivate := rPrivate.Group("/user")
	rUserPrivate.Get("/", handler.Profile)
}

// RequestSecret func for send secret code to specified email address.
// @Summary request email validation secret code
// @Description Sending secret code to specified email address.
// @Tags User
// @Accept mpfd
// @Param email formData string true "Destination Email address"
// @Produce json
// @Success 200 {object} domain.JSONResult{data=domain.EmailValidation} "Description"
// @Failure 422 {object} []domain.HTTPError
// @Failure 400 {object} domain.HTTPError
// @Failure 404 {object} domain.HTTPError
// @Failure 500 {object} domain.HTTPError
// @Router /v1/user/email-validation-secret [post]
func (uh *UserHandler) RequestSecret(c *fiber.Ctx) error {
	model := new(domain.RequestSecretForm)
	model.Email = c.FormValue("email")

	// Validate form input
	err := uh.Validate.Struct(model)
	if err != nil {
		return domain.NewHttpError(c, err)
	}

	secretCode, err := uh.UserUsecase.RequestSecret(model.Email)
	if err != nil {
		return domain.NewHttpError(c, err)
	}

	modelResp := domain.EmailValidation{Email: model.Email, SecretCode: secretCode}

	return c.JSON(domain.JSONResult{Data: modelResp, Message: "Success"})
}

// Signup func for signup.
// @Summary user signup
// @Description Signup to get JWT token.
// @Tags User
// @Accept json
// @Produce json
// @Param login form body domain.SignupForm true "Fill form"
// @success 200 {object} domain.JSONResult{data=string} "Registration Success"
// @Failure 400 {object} domain.HTTPError
// @Failure 404 {object} domain.HTTPError
// @Failure 422 {array} domain.HTTPError
// @Failure 500 {object} domain.HTTPError
// @Router /v1/user/signup [post]
func (uh *UserHandler) Signup(c *fiber.Ctx) error {
	var signupResponse domain.JSONResult

	signupForm := new(domain.SignupForm)

	//  Parse body into application struct
	if err := c.BodyParser(signupForm); err != nil {
		return domain.NewHttpError(c, err)
	}

	// Validate form input
	err := uh.Validate.Struct(signupForm)
	if err != nil {
		return domain.NewHttpError(c, err)
	}

	err = uh.UserUsecase.Signup(signupForm)
	if err != nil {
		return domain.NewHttpError(c, err)
	}

	signupResponse.Message = "Registration Success"
	signupResponse.Data = "Registration Success, continue to login"
	return c.JSON(signupResponse)
}

// Login func for login.
// @Summary user login
// @Description Login to get JWT token.
// @Tags User
// @Accept mpfd
// @Param email formData string true "Email address"
// @Param password formData string true "Password"
// @Param captcha formData string true "Captcha"
// @Produce json
// @Success 200 {object} domain.JSONResult{data=string} "Login Success, JWT Token provided"
// @Failure 400 {object} domain.HTTPError
// @Failure 422 {object} []domain.HTTPError
// @Failure 404 {object} domain.HTTPError
// @Failure 500 {object} domain.HTTPError
// @Router /v1/user/login [post]
func (uh *UserHandler) Login(c *fiber.Ctx) error {
	loginForm := new(domain.LoginForm)
	loginForm.Email = c.FormValue("email")
	loginForm.Password = c.FormValue("password")
	loginForm.Captcha = c.FormValue("captcha")

	// Validate form input
	err := uh.Validate.Struct(loginForm)
	if err != nil {
		return domain.NewHttpError(c, err)
	}

	user, err := uh.UserUsecase.Login(loginForm)
	if err != nil {
		return domain.NewHttpError(c, err)
	}

	var loginResponse domain.JSONResult
	loginResponse.Message = "Login Success, JWT Token provided"

	// Generate a new Access token.
	loginResponse.Data, err = utils.GenerateNewAccessToken(&user)
	if err != nil {
		return domain.NewHttpError(c, err)
	}

	return c.JSON(loginResponse)
}

// Profile godoc
// @Summary Show an profile
// @Description Show an user profile
// @Tags User
// @Produce  json
// @Success 200 {object} interface{}
// @Failure 404 {object} domain.HTTPError
// @Failure 500 {object} domain.HTTPError
// @Security ApiKeyAuth
// @Router /v1/auth/user [get]
func (uh *UserHandler) Profile(c *fiber.Ctx) error {
	user := c.Locals("jwt").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	//email := claims["email"].(string)

	return c.JSON(claims)
}
