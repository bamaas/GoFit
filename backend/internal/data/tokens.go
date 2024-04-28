package data

import (
	"database/sql"
	"log/slog"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base32"
	"time"
)

const (
	ScopeActivation = "activation"
	ScopeAuthentication = "authentication"
)

type Token struct {
	Plaintext 	string 		`json:"token"`
	Hash 		[]byte 		`json:"-"`
	UserID 		int64 		`json:"-"`
	Expiry 		time.Time 	`json:"expiry"`
	Scope 		string 		`json:"-"`
}

type TokenModel struct {
	DB *sql.DB
	logger *slog.Logger
}

func generateToken(userID int64, ttl time.Duration, scope string) (*Token, error) {
	token := &Token{
		UserID: userID,
		Expiry: time.Now().Add(ttl),
		Scope: scope,
	}

	randomBytes := make([]byte, 16)
	_, err := rand.Read(randomBytes)
	if err != nil {
		return nil, err
	}

	// Generate plaintext token
	token.Plaintext = base32.StdEncoding.WithPadding(base32.NoPadding).EncodeToString(randomBytes)

	// Generate hash of the token
	hash := sha256.Sum256([]byte(token.Plaintext))
	token.Hash = hash[:]

	return token, nil
}

func (m *TokenModel) Insert(token *Token) error {

	query := `
	INSERT INTO tokens (hash, user_id, expiry, scope)
	VALUES (?, ?, ?, ?);`

	_, err := m.DB.Exec(query, token.Hash, token.UserID, token.Expiry, token.Scope)
	if err != nil {
		return err
	}

	return nil
}

func (m *TokenModel) New(UserID int64, ttl time.Duration, scope string) (*Token, error) {

	token, err := generateToken(UserID, ttl, scope)
	if err != nil {
		return nil, err
	}

	err = m.Insert(token)		// TODO: Verify this impelmentation
	if err != nil {
		return nil, err
	}

	return token, nil
}

func (m *TokenModel) DeleteAllForUser(scope string, userID int64) error {

	query := `
	DELETE FROM tokens
	WHERE user_id = ? AND scope = ?;`

	_, err := m.DB.Exec(query, userID, scope)
	if err != nil {
		return err
	}

	return nil
}