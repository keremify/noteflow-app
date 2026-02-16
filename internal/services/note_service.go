package services

import (
	"errors"
	"saasproject/internal/models"
	"saasproject/internal/repository"
)

type NoteService struct {
	Repo *repository.NoteRepository
}

func NewNoteService(r *repository.NoteRepository) *NoteService {
	return &NoteService{Repo: r}
}

func (s *NoteService) CreateNote(note *models.Note) error {
	return s.Repo.Create(note)
}

func (s *NoteService) ListNotes(userID uint, role string) ([]models.Note, error) {
	if role == "admin" {
		return s.Repo.GetAll()
	}
	return s.Repo.GetAllByUser(userID)
}

func (s *NoteService) UpdateNote(noteID, userID uint, role, title, content string) error {
	note, err := s.Repo.GetByID(noteID)
	if err != nil {
		return err
	}

	// owner control
	if note.UserID != userID && role != "admin" {
		return errors.New("Not eriÅŸim izni yok")
	}

	note.Title = title
	note.Content = content
	return s.Repo.Update(note)
}

func (s *NoteService) DeleteNote(noteID, userID uint, role string) error {
	note, err := s.Repo.GetByID(noteID)
	if err != nil {
		return err
	}

	if note.UserID != userID && role != "admin" { // for unauthorized access
		return errors.New("Not eriÅŸim izni yok")
	}

	return s.Repo.Delete(note)
}
