package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/HerbHall/RunNotes/internal/models"
	"github.com/HerbHall/RunNotes/internal/store"
)

// Handler provides HTTP handlers for note CRUD operations.
type Handler struct {
	store *store.NoteStore
}

// NewHandler creates a Handler backed by the given NoteStore.
func NewHandler(s *store.NoteStore) *Handler {
	return &Handler{store: s}
}

// RegisterRoutes registers all note routes on the given ServeMux.
func (h *Handler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /notes", h.HandleListNotes)
	mux.HandleFunc("GET /notes/{name}", h.HandleGetNote)
	mux.HandleFunc("POST /notes", h.HandleCreateNote)
	mux.HandleFunc("PUT /notes/{name}", h.HandleUpdateNote)
	mux.HandleFunc("DELETE /notes/{name}", h.HandleDeleteNote)
}

// HandleListNotes returns all notes, with optional ?pinned= and ?search= filters.
func (h *Handler) HandleListNotes(w http.ResponseWriter, r *http.Request) {
	var pinned *bool
	if p := r.URL.Query().Get("pinned"); p != "" {
		b, err := strconv.ParseBool(p)
		if err != nil {
			writeError(w, http.StatusBadRequest, "invalid pinned parameter")
			return
		}
		pinned = &b
	}

	search := r.URL.Query().Get("search")

	notes, err := h.store.List(r.Context(), pinned, search)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to list notes")
		return
	}

	writeJSON(w, http.StatusOK, notes)
}

// HandleGetNote returns the note for the given container name.
func (h *Handler) HandleGetNote(w http.ResponseWriter, r *http.Request) {
	name := r.PathValue("name")

	note, err := h.store.GetByName(r.Context(), name)
	if errors.Is(err, store.ErrNotFound) {
		writeError(w, http.StatusNotFound, "note not found")
		return
	}
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to get note")
		return
	}

	writeJSON(w, http.StatusOK, note)
}

// HandleCreateNote creates a new note from the JSON request body.
func (h *Handler) HandleCreateNote(w http.ResponseWriter, r *http.Request) {
	var req models.CreateNoteRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON body")
		return
	}

	if err := req.Validate(); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	note, err := h.store.Create(r.Context(), req)
	if err != nil {
		// Check for unique constraint violation.
		if isUniqueViolation(err) {
			writeError(w, http.StatusConflict, "note already exists for this container")
			return
		}
		writeError(w, http.StatusInternalServerError, "failed to create note")
		return
	}

	writeJSON(w, http.StatusCreated, note)
}

// HandleUpdateNote applies partial updates to the note for the given container name.
func (h *Handler) HandleUpdateNote(w http.ResponseWriter, r *http.Request) {
	name := r.PathValue("name")

	var req models.UpdateNoteRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON body")
		return
	}

	if err := req.Validate(); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	note, err := h.store.Update(r.Context(), name, req)
	if errors.Is(err, store.ErrNotFound) {
		writeError(w, http.StatusNotFound, "note not found")
		return
	}
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to update note")
		return
	}

	writeJSON(w, http.StatusOK, note)
}

// HandleDeleteNote deletes the note for the given container name.
func (h *Handler) HandleDeleteNote(w http.ResponseWriter, r *http.Request) {
	name := r.PathValue("name")

	err := h.store.Delete(r.Context(), name)
	if errors.Is(err, store.ErrNotFound) {
		writeError(w, http.StatusNotFound, "note not found")
		return
	}
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to delete note")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// isUniqueViolation checks if an error is a SQLite unique constraint violation.
func isUniqueViolation(err error) bool {
	if err == nil {
		return false
	}
	// modernc.org/sqlite returns errors containing "UNIQUE constraint failed".
	return contains(err.Error(), "UNIQUE constraint failed")
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && searchString(s, substr)
}

func searchString(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
