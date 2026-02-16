package models

import "time"

type RefreshToken struct {
	ID        uint   `gorm:"primaryKey"`
	UserID    uint   `gorm:"index"`
	Token     string `gorm:"uniqueIndex"`
	UserAgent string
	IP        string
	ExpiresAt time.Time
	CreatedAt time.Time
}
