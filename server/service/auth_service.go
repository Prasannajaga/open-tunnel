package service

import (
	"crypto/sha256"
	"database/sql"
	"encoding/hex"

	"opentunnel/server/config"
	"opentunnel/server/constants"
	"opentunnel/server/database"
	"opentunnel/server/utils"
)

type User struct {
	ID           int
	Username     string
	PasswordHash string
}

type AuthService struct {
	cfg *config.Config
	db  *sql.DB
}

func NewAuthService(cfg *config.Config) *AuthService {
	return &AuthService{
		cfg: cfg,
		db:  database.DB,
	}
}

func (s *AuthService) Authenticate(username, password string) (string, error) {
	var user User
	err := s.db.QueryRow(
		"SELECT id, username, password_hash FROM users WHERE username = $1",
		username,
	).Scan(&user.ID, &user.Username, &user.PasswordHash)

	if err == sql.ErrNoRows {
		return "", &AuthError{Message: constants.ErrInvalidCreds}
	}
	if err != nil {
		return "", err
	}

	if !verifyPassword(password, user.PasswordHash) {
		return "", &AuthError{Message: constants.ErrInvalidCreds}
	}

	token, err := utils.GenerateToken(username, s.cfg.JWTSecret, s.cfg.TokenExpiry)
	if err != nil {
		return "", err
	}

	return token, nil
}

func (s *AuthService) ValidateToken(tokenString string) (*utils.Claims, error) {
	claims, err := utils.ValidateToken(tokenString, s.cfg.JWTSecret)
	if err != nil {
		return nil, &AuthError{Message: constants.ErrTokenInvalid}
	}
	return claims, nil
}

func (s *AuthService) CreateUser(username, password string) error {
	hash := hashPassword(password)
	_, err := s.db.Exec(
		"INSERT INTO users (username, password_hash) VALUES ($1, $2)",
		username, hash,
	)
	return err
}

func hashPassword(password string) string {
	hash := sha256.Sum256([]byte(password))
	return hex.EncodeToString(hash[:])
}

func verifyPassword(password, hash string) bool {
	return hashPassword(password) == hash
}

type AuthError struct {
	Message string
}

func (e *AuthError) Error() string {
	return e.Message
}
