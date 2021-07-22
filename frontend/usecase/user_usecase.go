package usecase

import (
	"fmt"
	"github.com/cooljar/go-postgres-fiber/domain"
	"github.com/cooljar/go-postgres-fiber/utils"
	"golang.org/x/crypto/bcrypt"
	"path"
	"strings"
	"time"
)

type userUsecase struct {
	userRepo       domain.UserRepository
	contextTimeout time.Duration
}

// NewUserUsecase will create new an userUsecase object representation of domain.UserUsecase interface
func NewUserUsecase(u domain.UserRepository, timeout time.Duration) domain.UserUseCase {
	return &userUsecase{
		userRepo:       u,
		contextTimeout: timeout,
	}
}

func (uu *userUsecase) RequestSecret(email string) (secret string, err error) {
	secret, err = uu.userRepo.RequestSecret(email)
	if err != nil {
		return
	}

	var viewPath = path.Join("assets", "html", "email-validation.html")
	var dataEmail = map[string]interface{}{
		"email":  email,
		"secret": secret,
	}

	mailer := utils.NewMailer(
		"Email Validation Secret Code At - Cooljar Apps",
		email,
		viewPath,
		dataEmail,
	)

	go func(mailer *utils.Mailer) {
		err = mailer.SendMail()
		if err != nil {
			fmt.Println(err)
		}
	}(mailer)

	return
}

func (uu *userUsecase) Signup(sf *domain.SignupForm) (err error) {
	sf.Email = strings.ToLower(sf.Email)
	sf.Username = strings.ToLower(sf.Username)

	err = uu.userRepo.Signup(sf)

	return
}

func (uu *userUsecase) Login(lf *domain.LoginForm) (u domain.User, err error) {
	u, err = uu.userRepo.Login(lf)
	if err != nil {
		return
	}

	if validatePassword(lf.Password, u.PasswordHash) == false {
		err = domain.DataValidationError{Field: "password", Message: "invalid password"}
		return
	}

	return
}

func (uu *userUsecase) Profile() error {
	return nil
}

func validatePassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
