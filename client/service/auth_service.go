package service

import (
	"opentunnel/client/utils"
)

type AuthService struct{}

func NewAuthService() *AuthService {
	return &AuthService{}
}

func (s *AuthService) SaveToken(token string) error {
	return utils.SaveToken(token)
}

func (s *AuthService) LoadToken() (string, error) {
	return utils.LoadToken()
}

func (s *AuthService) IsAuthenticated() bool {
	return utils.TokenExists()
}
