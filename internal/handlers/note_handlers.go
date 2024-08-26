package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/ananikitina/notes-rest/internal/middleware"
	"github.com/ananikitina/notes-rest/internal/models"
	"github.com/ananikitina/notes-rest/internal/services"
	"github.com/ananikitina/notes-rest/internal/usecases"
	_ "github.com/lib/pq"
)

type NoteHandler struct {
	noteUseCase *usecases.NoteUseCase
}

func NewNoteHandler(noteUseCase *usecases.NoteUseCase) *NoteHandler {
	return &NoteHandler{noteUseCase: noteUseCase}
}

func (n *NoteHandler) AddNoteHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		userID := r.Context().Value(middleware.UserIDKey).(int)

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

		note.UserID = userID
		note.CreatedAt = time.Now()

		if err := n.noteUseCase.AddNote(note); err != nil {
			fmt.Println(err)
			http.Error(w, "Failed to add note", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(note)
	}
}

func (n *NoteHandler) GetNotesHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		userID := r.Context().Value(middleware.UserIDKey).(int)
		notes, err := n.noteUseCase.GetNotesByUserID(userID)
		if err != nil {
			fmt.Println(err)
			http.Error(w, "Failed to fetch notes", http.StatusInternalServerError)
			return
		}
		// if err := rows.Err(); err != nil {
		// 	http.Error(w, "Error during rows iteration", http.StatusInternalServerError)
		// 	return
		// }

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(notes)
	}
}
