package data

import (
	"database/sql"
	"errors"
	"log/slog"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID 		 	int64  		`json:"id"`
	CreatedAt 	string 		`json:"created_at"`
	Email 	 	string 		`json:"email"`
	Password 	password 	`json:"-"`
	Activated 	bool   		`json:"activated"`
	Version 	int    		`json:"-"`
}

type password struct {
	Plaintext *string
	Hash []byte
}

type UserModel struct {
	DB *sql.DB
	logger *slog.Logger
}

func (m *UserModel) Insert(user *User) error {

	query := `
	INSERT INTO users (email, password_hash, activated, version) 
	VALUES (?, ?, ?, ?) 
	RETURNING id, created_at, version;`

	err := m.DB.QueryRow(query, user.Email, user.Password.Hash, user.Activated, user.Version).Scan(&user.ID, &user.CreatedAt, &user.Version)
	if err != nil {
		// Check for duplicate email
		return err
	}

	return nil
}

func (m *UserModel) GetByEmail(email string) (*User, error) {

	user := User{}

	query := `
	SELECT id, created_at, email, password_hash, activated, version
	FROM users
	WHERE email = $1;`

	err := m.DB.QueryRow(query, email).Scan(&user.ID, &user.CreatedAt, &user.Email, &user.Password.Hash, &user.Activated, &user.Version)
	if err != nil {
		switch {
			case errors.Is(err, sql.ErrNoRows):
				return nil, errors.New("no record found")
			default:
				return nil, err
		}
	}

	return &user, nil
}

func (m *UserModel) Update(user User) error {

	query := `
	UPDATE users
	SET email = ?, password_hash = ?, activated = ?, version = version + 1
	WHERE id = ? AND version = ?
	RETURNING version;`

	args := []any{
		user.Email,
		user.Password.Hash,
		user.Activated,
		user.ID,
		user.Version,
	}

	err := m.DB.QueryRow(query, args...).Scan(&user.Version)
	if err != nil {
		switch {
			// Check for duplicate email
			case errors.Is(err, sql.ErrNoRows):
				return errors.New("no record found")
			default:
				return err
		}
	}

	return nil
}

func (p* password) Set(plaintextPassword string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(plaintextPassword), 12)
	if err != nil {
		return err
	}

	p.Plaintext = &plaintextPassword
	p.Hash = hash

	return nil
}

func (p *password) Matches(plaintextPassword string) (bool, error){
	err := bcrypt.CompareHashAndPassword(p.Hash, []byte(plaintextPassword))
	if err != nil {
		switch {
			case errors.Is(err, bcrypt.ErrMismatchedHashAndPassword):
				return false, nil
			default:
				return false, err
		}
	}
	return true, nil
}