package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/ananikitina/notes-rest/internal/database"
	"github.com/ananikitina/notes-rest/internal/models"
	"github.com/ananikitina/notes-rest/internal/services"
	_ "github.com/lib/pq"
)

func AddNoteHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var note models.Note
		if err := json.NewDecoder(r.Body).Decode(&note); err != nil {
			http.Error(w, "Invalid input", http.StatusBadRequest)
			return
		}

		spellChecker := services.NewYandexSpellChecker()

		// Spell check
		spellErrors, err := spellChecker.Check(note.Content)
		if err != nil {
			fmt.Println(err)
			http.Error(w, "Failed to check spelling", http.StatusInternalServerError)
			return
		}

		// Return errors if found
		if len(spellErrors) > 0 {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]interface{}{
				"message": "Spelling errors found",
				"errors":  spellErrors,
			})
			return
		}

		note.CreatedAt = time.Now()

		_, err = database.DB.Exec(`
		INSERT INTO notes (content, user_id)  
		VALUES ($1,$2)
		RETURNING id;
	`, note.Content, note.UserID)

		if err != nil {
			fmt.Println(err)
			http.Error(w, "Failed to add note", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(note)
	}
}

func GetNotesHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		rows, err := database.DB.Query("SELECT id, content, user_id, created_at FROM notes;")
		if err != nil {
			fmt.Println(err)
			http.Error(w, "Failed to fetch notes", http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		var notes []models.Note
		for rows.Next() {
			var note models.Note
			if err := rows.Scan(&note.ID, &note.Content, &note.UserID, &note.CreatedAt); err != nil {
				http.Error(w, "Failed to scan notes", http.StatusInternalServerError)
				return
			}
			notes = append(notes, note)
		}

		if err := rows.Err(); err != nil {
			http.Error(w, "Error during rows iteration", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(notes)
	}
}

// // NoteHandler handles HTTP requests related to notes.
// type NoteHandler struct {
// 	noteUsecase *usecase.NoteUsecase
// }

// // NewNoteHandler creates a new NoteHandler.
// func NewNoteHandler(noteUsecase *usecase.NoteUsecase) http.Handler {
// 	return &NoteHandler{
// 		noteUsecase: noteUsecase,
// 	}
// }

// // ServeHTTP handles HTTP requests for notes.
// func (h *NoteHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
// 	switch r.Method {
// 	case http.MethodPost:
// 		h.createNote(w, r)
// 	case http.MethodGet:
// 		h.getAllNotes(w, r)
// 	default:
// 		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
// 	}
// }
