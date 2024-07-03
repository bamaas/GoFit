package data

import (
	"crypto/sha256"
	"database/sql"
	"errors"
	"log/slog"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID 		 	int64  		`json:"id"`
	CreatedAt 	string 		`json:"created_at"`
	Email 	 	string 		`json:"email"`
	Password 	password 	`json:"-"`
	Activated 	bool   		`json:"activated"`
	Version 	int    		`json:"-"`
	Goal 	 	string 		`json:"goal"`
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
		// TODO: Check for duplicate email
		return err
	}

	return nil
}

func (m *UserModel) Update(user *User) error {
	
	q := `
	UPDATE users
	SET email = ?, password_hash = ?, activated = ?, version = version + 1, goal = ?
	WHERE id = ?
	RETURNING version;
	`

	args := []any{
		user.Email,
		user.Password.Hash,
		user.Activated,
		user.Goal,
	}

	return m.DB.QueryRow(q, args...).Scan(&user.Version)
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
				return nil, ErrRecordNotFound
			default:
				return nil, err
		}
	}

	return &user, nil
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

func (m *UserModel) GetForToken(tokenScope, tokenPlaintext string) (*User, error) {

	tokenHash := sha256.Sum256([]byte(tokenPlaintext))

	query := `
	SELECT users.id, users.created_at, users.email, users.password_hash, users.activated, users.version 
	FROM users
	INNER JOIN tokens
	ON users.id = tokens.user_id
	WHERE tokens.hash = $1 
	AND tokens.scope = $2 
	AND tokens.expiry > $3
	`
	args := []any{tokenHash[:], tokenScope, time.Now()}

	var user User
	err := m.DB.QueryRow(query, args...).Scan(
		&user.ID, 
		&user.CreatedAt, 
		&user.Email, 
		&user.Password.Hash, 
		&user.Activated, 
		&user.Version,
	)
	if err != nil {
		switch {
			case errors.Is(err, sql.ErrNoRows):
				return nil, ErrRecordNotFound
			default:
				return nil, err
		}
	}

	return &user, nil
}