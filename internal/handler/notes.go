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
	mux.HandleFunc("GET /notes/export", h.HandleExportNotes)
	mux.HandleFunc("POST /notes/import", h.HandleImportNotes)
	mux.HandleFunc("GET /notes", h.HandleListNotes)
	mux.HandleFunc("POST /notes", h.HandleCreateNote)
	mux.HandleFunc("GET /notes/container/{name}", h.HandleListContainerNotes)
	mux.HandleFunc("DELETE /notes/container/{name}", h.HandleDeleteContainerNotes)
	mux.HandleFunc("GET /notes/{id}", h.HandleGetNote)
	mux.HandleFunc("PUT /notes/{id}", h.HandleUpdateNote)
	mux.HandleFunc("DELETE /notes/{id}", h.HandleDeleteNote)
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

// HandleListContainerNotes returns all notes for the given container name.
func (h *Handler) HandleListContainerNotes(w http.ResponseWriter, r *http.Request) {
	name := r.PathValue("name")

	notes, err := h.store.ListByContainer(r.Context(), name)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to list container notes")
		return
	}

	writeJSON(w, http.StatusOK, notes)
}

// HandleGetNote returns the note with the given ID.
func (h *Handler) HandleGetNote(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid note ID")
		return
	}

	note, err := h.store.GetByID(r.Context(), id)
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
		writeError(w, http.StatusInternalServerError, "failed to create note")
		return
	}

	writeJSON(w, http.StatusCreated, note)
}

// HandleUpdateNote applies partial updates to the note with the given ID.
func (h *Handler) HandleUpdateNote(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid note ID")
		return
	}

	var req models.UpdateNoteRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON body")
		return
	}

	if err := req.Validate(); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	note, err := h.store.Update(r.Context(), id, req)
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

// HandleExportNotes returns all notes as a JSON array.
func (h *Handler) HandleExportNotes(w http.ResponseWriter, r *http.Request) {
	notes, err := h.store.ExportAll(r.Context())
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to export notes")
		return
	}

	writeJSON(w, http.StatusOK, notes)
}

// HandleImportNotes imports notes from a JSON array in the request body.
func (h *Handler) HandleImportNotes(w http.ResponseWriter, r *http.Request) {
	var notes []models.Note
	if err := json.NewDecoder(r.Body).Decode(&notes); err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON body")
		return
	}

	imported, err := h.store.ImportAll(r.Context(), notes)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to import notes")
		return
	}

	writeJSON(w, http.StatusOK, map[string]int{"imported": imported})
}

// HandleDeleteNote deletes the note with the given ID.
func (h *Handler) HandleDeleteNote(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid note ID")
		return
	}

	err = h.store.Delete(r.Context(), id)
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

// HandleDeleteContainerNotes deletes all notes for the given container name.
func (h *Handler) HandleDeleteContainerNotes(w http.ResponseWriter, r *http.Request) {
	name := r.PathValue("name")

	deleted, err := h.store.DeleteByContainer(r.Context(), name)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to delete container notes")
		return
	}

	writeJSON(w, http.StatusOK, map[string]int64{"deleted": deleted})
}
