package models

type Tag struct {
	ID     uint   `gorm:"primaryKey"`
	UserID uint   `gorm:"not null;index"`
	Name   string `gorm:"not null;uniqueIndex"`
}

type NoteTag struct {
	ID     uint `gorm:"primaryKey"`
	NoteID uint
	TagID  uint
}
