package domain

import (
	"encoding/json"
	_ "github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

type RequestSecretForm struct {
	Email string `json:"email" validate:"required,email"`
}

type SignupForm struct {
	Username       string `json:"username" validate:"required,lowercase,alphanumunicode"`
	Email          string `json:"email" validate:"required,email"`
	FullName       string `json:"full_name"`
	Password       string `json:"password" validate:"required"`
	PasswordRepeat string `json:"password_repeat" validate:"required,eqfield=Password"`
	Captcha        string `json:"captcha" validate:"required"`
	EmailSecretCode string `json:"email_secret_code" validate:"required,len=6"`
}

type LoginForm struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
	Captcha  string `json:"captcha" validate:"required"`
}

type User struct {
	ID                 uuid.UUID `json:"id"`
	Username           string    `json:"username"`
	FullName           string    `json:"full_name"`
	AuthKey            string    `json:"auth_key"`
	PasswordHash       string    `json:"password_hash"`
	PasswordResetToken string    `json:"password_reset_token"`
	VerificationToken  string    `json:"verification_token"`
	Email              string    `json:"email"`
	Status             int       `json:"status"`
	CreatedAt          int       `json:"created_at"`
	UpdatedAt          int       `json:"updated_at"`
}

// FromJSON decode json to user struct
func (u *User) FromJSON(msg []byte) error {
	return json.Unmarshal(msg, u)
}

// ToJSON encode user struct to json
func (u *User) ToJSON() []byte {
	str, _ := json.Marshal(u)
	return str
}

// UserUseCase represent the user's use cases
type UserUseCase interface {
	RequestSecret(email string) (secret string, err error)
	Signup(s *SignupForm) (err error)
	Login(l *LoginForm) (u User, err error)
	Profile() error
}

// UserRepository represent the user's repository
type UserRepository interface {
	RequestSecret(email string) (secret string, err error)
	Signup(s *SignupForm) (err error)
	Login(l *LoginForm) (u User, err error)
}
