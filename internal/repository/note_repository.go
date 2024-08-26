package repository

import (
	"database/sql"

	"github.com/ananikitina/notes-rest/internal/domain"
	"github.com/ananikitina/notes-rest/internal/models"
)

// noteRepository is an implementation of the NoteRepository interface.
type noteRepository struct {
	DB *sql.DB
}

// NewNoteRepository creates a new note repository with the given database connection.
func NewNoteRepository(DB *sql.DB) domain.NoteRepository {
	return &noteRepository{DB: DB}
}

// AddNote inserts a new note into the database.
func (n *noteRepository) Add(note models.Note) error {
	_, err := n.DB.Exec(`
	INSERT INTO notes (content, user_id)  
	VALUES ($1,$2)
	RETURNING id;
`, note.Content, note.UserID)
	return err
}

// GetNotesByUserID retrieves notes from specified user.
func (n *noteRepository) GetByUserID(userID int) ([]models.Note, error) {
	rows, err := n.DB.Query(
		"SELECT id, content, user_id, created_at FROM notes WHERE user_id = $1;",
		userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var notes []models.Note
	for rows.Next() {
		var note models.Note
		if err := rows.Scan(&note.ID, &note.Content, &note.UserID, &note.CreatedAt); err != nil {
			return nil, err
		}
		notes = append(notes, note)
	}
	return notes, rows.Err()
}
