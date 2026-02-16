package repository

import (
	"gorm.io/gorm"
	"saasproject/internal/models"
)

type NoteRepository struct {
	DB *gorm.DB
}

func NewNoteRepository(db *gorm.DB) *NoteRepository {
	return &NoteRepository{DB: db}
}

// create note
func (r *NoteRepository) Create(note *models.Note) error {
	return r.DB.Create(note).Error
}

// list (just specified user's notes)
func (r *NoteRepository) GetAllByUser(userID uint) ([]models.Note, error) {
	var notes []models.Note
	err := r.DB.Where("user_id = ?", userID).Find(&notes).Error
	return notes, err
}

// list all notes (admin)
func (r *NoteRepository) GetAll() ([]models.Note, error) {
	var notes []models.Note
	err := r.DB.Find(&notes).Error
	return notes, err
}

// get by id
func (r *NoteRepository) GetByID(id uint) (*models.Note, error) {
	var note models.Note
	err := r.DB.First(&note, id).Error
	return &note, err
}

// update notes
func (r *NoteRepository) Update(note *models.Note) error {
	return r.DB.Save(note).Error
}

// delete notes
func (r *NoteRepository) Delete(note *models.Note) error {
	return r.DB.Delete(note).Error
}
