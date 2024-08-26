package usecases

import (
	"github.com/ananikitina/notes-rest/internal/domain"
	"github.com/ananikitina/notes-rest/internal/models"
)

// NoteUseCase represents the business logic for notes.
type NoteUseCase struct {
	noteRepo domain.NoteRepository
}

// NewNoteUseCase creates a new instance of NoteUseCase.
func NewNoteUseCase(noteRepo domain.NoteRepository) *NoteUseCase {
	return &NoteUseCase{noteRepo: noteRepo}
}

func (n *NoteUseCase) AddNote(note models.Note) error {
	return n.noteRepo.Add(note)
}

func (n *NoteUseCase) GetNotesByUserID(userID int) ([]models.Note, error) {
	return n.noteRepo.GetByUserID(userID)
}

func (n *NoteUseCase) GetAllNotes() ([]models.Note, error) {
	return n.noteRepo.GetAllNotes()
}
