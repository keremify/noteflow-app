package models

import "time"

type Note struct {
	ID           uint   `gorm:"primaryKey"`
	UserID       uint   `gorm:"not null;index"`
	Title        string `gorm:"not null"`
	Content      string `gorm:"type:text"`
	ReminderAt   *time.Time
	ReminderSent bool `gorm:"default:false"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
}
