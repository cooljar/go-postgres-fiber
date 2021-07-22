package pgsql

import (
	"context"
	"github.com/cooljar/go-postgres-fiber/domain"
	"github.com/cooljar/go-postgres-fiber/utils"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v4/pgxpool"
	"golang.org/x/crypto/bcrypt"
	"math/rand"
	"strconv"
	"time"
)

type pgsqlUserRepository struct {
	Conn *pgxpool.Pool
}

// NewPgsqlUserRepository will create an object that represent the user Repository interface
func NewPgsqlUserRepository(conn *pgxpool.Pool) domain.UserRepository {
	return &pgsqlUserRepository{Conn: conn}
}

func (ur *pgsqlUserRepository) RequestSecret(email string) (secret string, err error) {
	var emailExist bool
	emailExist, err = isExistByEmail(ur.Conn, email)

	if emailExist {
		return secret, domain.DataValidationError{Field: "email", Message: "invalid email address, already registered"}
	}

	min := 100000
	max := 999999
	randNum := rand.Intn(rand.Intn(max - min) + min)
	secret = strconv.Itoa(randNum)

	qStr := `INSERT INTO email_validation (email,secret_code) VALUES($1,$2) ON CONFLICT ON CONSTRAINT email_validation_pkey DO UPDATE SET secret_code = EXCLUDED.secret_code`
	commandTag, err := ur.Conn.Exec(context.Background(), qStr, email, secret)
	if err != nil {
		return secret,err
	}
	if commandTag.RowsAffected() != 1 {
		return secret, domain.ErrSqlNoRowFoundToUpsert
	}

	return
}

func (ur *pgsqlUserRepository) Signup(sf *domain.SignupForm) (err error) {
	var secretCodeExist, usernameExists, emailExists bool
	var verificationToken, passwordResetToken, authKey, passwordHash string

	// periksa apakah secret code email validation valid
	qSecret := `SELECT EXISTS(SELECT 1 FROM email_validation WHERE email = $1 AND secret_code = $2)`
	err = ur.Conn.QueryRow(context.Background(), qSecret, sf.Email, sf.EmailSecretCode).Scan(&secretCodeExist)
	if err != nil {
		return
	}
	if !secretCodeExist {
		err = domain.DataValidationError{Field: "email_secret_code", Message: "Invalid secret code"}
		return
	}

	usernameExists, err = isExistByUsername(ur.Conn, sf.Username)
	if err != nil {
		return
	}
	if usernameExists {
		err = domain.DataValidationError{Field: "username", Message: "username already taken"}
		return
	}

	emailExists, err = isExistByEmail(ur.Conn, sf.Email)
	if err != nil {
		return
	}
	if emailExists {
		err = domain.DataValidationError{Field: "email", Message: "email address already taken"}
		return
	}

	ts := time.Now().Unix()
	now := int(ts)

	passwordHash, err = generatePasswordHash(sf.Password)
	if err != nil {
		return
	}

	verificationToken, err = utils.GenerateRandomString(16)
	if err != nil {
		return
	}

	passwordResetToken, err = utils.GenerateRandomString(16)
	if err != nil {
		return
	}

	authKey, err = utils.GenerateRandomString(16)
	if err != nil {
		return
	}

	verificationToken = verificationToken + "_" + strconv.Itoa(now)
	passwordResetToken = passwordResetToken + "_" + strconv.Itoa(now)

	user := domain.User{
		uuid.New(),
		sf.Username,
		sf.FullName,
		authKey,
		passwordHash,
		passwordResetToken,
		verificationToken,
		sf.Email,
		domain.ConstUserStatusActive,
		now,
		now,
	}

	qStr := `insert into "user" (id, username, full_name, auth_key, password_hash, password_reset_token, verification_token, email, status, created_at, updated_at) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11) returning id`
	_, err = ur.Conn.Exec(context.Background(), qStr, user.ID, user.Username, user.FullName, user.AuthKey, user.PasswordHash, user.PasswordResetToken, user.VerificationToken, user.Email, user.Status, user.CreatedAt, user.UpdatedAt)
	if err != nil {
		return
	}

	return
}

func (ur *pgsqlUserRepository) Login(lf *domain.LoginForm) (u domain.User, err error) {
	var id uuid.UUID
	var username, email, fullName, passwordHash string

	qStr := `SELECT id, username, email, full_name, password_hash FROM "user" WHERE email = $1 AND status = $2`
	err = ur.Conn.QueryRow(context.Background(), qStr, lf.Email, domain.ConstUserStatusActive).Scan(&id, &username, &email, &fullName, &passwordHash)
	if err != nil {
		return
	}

	u.ID = id
	u.Username = username
	u.Email = email
	u.FullName = fullName
	u.PasswordHash = passwordHash

	return u, nil
}

func isExistByUsername(conn *pgxpool.Pool, username string) (exist bool, err error) {
	qStr := `SELECT EXISTS(SELECT 1 FROM "user" WHERE username = $1)`
	err = conn.QueryRow(context.Background(), qStr, username).Scan(&exist)
	if err != nil {
		return false, err
	}

	return
}

func isExistByEmail(conn *pgxpool.Pool, email string) (exist bool, err error) {
	qStr := `SELECT EXISTS(SELECT 1 FROM "user" WHERE email = $1)`
	err = conn.QueryRow(context.Background(), qStr, email).Scan(&exist)
	if err != nil {
		return false, err
	}

	return
}

func generatePasswordHash(password string) (string, error) {
	bytesPwd, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytesPwd), err
}
