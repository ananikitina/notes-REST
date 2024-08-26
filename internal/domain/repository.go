package domain

import "github.com/ananikitina/notes-rest/internal/models"

// NoteRepository defines the contract for note-related database operations.
type NoteRepository interface {
	Add(note models.Note) error
	GetByUserID(userID int) ([]models.Note, error)
}
