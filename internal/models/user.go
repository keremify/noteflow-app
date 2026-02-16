package models

import "time"

type User struct {
	ID        uint   `gorm:"primaryKey"`
	Name      string `gorm:"not null"`
	Email     string `gorm:"uniqueIndex;not null"`
	Password  string `gorm:"not null"`
	Premium   bool   `gorm:"default:false"`
	Role      string `gorm:"default:'user'"`
	CreatedAt time.Time
}
