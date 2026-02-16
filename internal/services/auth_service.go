package services

import (
	"errors"
	"saasproject/internal/models"
	"saasproject/internal/repository"
	"saasproject/internal/utils"
	"time"
)

type AuthService struct {
	userRepo         *repository.UserRepository
	refreshTokenRepo *repository.RefreshTokenRepository
}

func NewAuthService(userRepo *repository.UserRepository, refreshTokenRepo *repository.RefreshTokenRepository) *AuthService {
	return &AuthService{userRepo: userRepo, refreshTokenRepo: refreshTokenRepo}
}

// Register
func (s *AuthService) Register(name, email, password string) error {
	_, err := s.userRepo.FindByEmail(email)
	if err == nil {
		return errors.New("user already exists")
	}

	hashedPassword, err := utils.HashPassword(password)
	if err != nil {
		return err
	}

	user := &models.User{
		Name:     name,
		Email:    email,
		Password: hashedPassword,
	}

	return s.userRepo.Create(user)
}

// Login
func (s *AuthService) Login(email, password, userAgent, ip string) (string, string, error) {
	user, err := s.userRepo.FindByEmail(email)
	if err != nil {
		return "", "", errors.New("Geçersiz giriş bilgileri")
	}

	if !utils.CheckPasswordHash(password, user.Password) {
		return "", "", errors.New("Geçersiz giriş bilgileri")
	}

	accessToken, err := utils.GenerateToken(user.ID, user.Role)
	if err != nil {
		return "", "", err
	}

	refreshToken, err := utils.GenerateRefreshToken()
	if err != nil {
		return "", "", err
	}

	rt := &models.RefreshToken{
		UserID:    user.ID,
		Token:     refreshToken,
		UserAgent: userAgent,
		IP:        ip,
		ExpiresAt: time.Now().Add(7 * 24 * time.Hour),
	}

	if err := s.refreshTokenRepo.Create(rt); err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}

func (s *AuthService) Refresh(refreshToken string) (string, string, error) {
	rt, err := s.refreshTokenRepo.Find(refreshToken)
	if err != nil || rt.ExpiresAt.Before(time.Now()) {
		return "", "", errors.New("Refresh token geçersiz veya süresi dolmuş")
	}

	user, err := s.userRepo.FindByID(rt.UserID)
	if err != nil {
		return "", "", errors.New("Kullanıcı bulunamadı")
	}

	if err := s.refreshTokenRepo.Delete(refreshToken); err != nil {
		return "", "", errors.New("failed to revoke refresh token")
	}

	accessToken, err := utils.GenerateToken(user.ID, user.Role)
	if err != nil {
		return "", "", err
	}

	newRefreshToken, err := utils.GenerateRefreshToken()
	if err != nil {
		return "", "", err
	}

	err = s.refreshTokenRepo.Create(&models.RefreshToken{
		UserID:    user.ID,
		Token:     newRefreshToken,
		ExpiresAt: time.Now().Add(7 * 24 * time.Hour),
	})

	if err != nil {
		return "", "", err
	}

	return accessToken, newRefreshToken, nil
}

func (s *AuthService) Logout(refreshToken string) error {
	return s.refreshTokenRepo.Delete(refreshToken)
}

func (s *AuthService) LogoutAll(userID uint) error {
	return s.refreshTokenRepo.DeleteByUser(userID)
}
