package repository

import (
	"time"

	"saasproject/internal/models"

	"gorm.io/gorm"
)

type RefreshTokenRepository struct {
	db *gorm.DB
}

func NewRefreshTokenRepository(db *gorm.DB) *RefreshTokenRepository {
	return &RefreshTokenRepository{db: db}
}

func (r *RefreshTokenRepository) Create(token *models.RefreshToken) error {
	return r.db.Create(token).Error
}

func (r *RefreshTokenRepository) Find(token string) (*models.RefreshToken, error) {
	var rt models.RefreshToken
	err := r.db.Where("token = ?", token).First(&rt).Error
	return &rt, err
}

func (r *RefreshTokenRepository) Delete(token string) error {
	return r.db.Where("token = ?", token).Delete(&models.RefreshToken{}).Error
}

func (r *RefreshTokenRepository) DeleteExpired() error {
	return r.db.Where("expires_at < ?", time.Now()).Delete(&models.RefreshToken{}).Error
}

func (r *RefreshTokenRepository) DeleteByUser(userID uint) error {
	return r.db.Where("user_id = ?", userID).
		Delete(&models.RefreshToken{}).Error
}
