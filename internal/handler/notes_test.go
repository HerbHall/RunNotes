package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/HerbHall/RunNotes/internal/database"
	"github.com/HerbHall/RunNotes/internal/models"
	"github.com/HerbHall/RunNotes/internal/store"
)

func newTestHandler(t *testing.T) (*Handler, *http.ServeMux) {
	t.Helper()
	db, err := database.Open(":memory:")
	if err != nil {
		t.Fatalf("open test db: %v", err)
	}
	t.Cleanup(func() { _ = db.Close() })

	s := store.NewNoteStore(db)
	h := NewHandler(s)
	mux := http.NewServeMux()
	h.RegisterRoutes(mux)
	return h, mux
}

func createNote(t *testing.T, mux *http.ServeMux, req models.CreateNoteRequest) *models.Note {
	t.Helper()
	body, _ := json.Marshal(req)
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPost, "/notes", bytes.NewReader(body))
	mux.ServeHTTP(w, r)
	if w.Code != http.StatusCreated {
		t.Fatalf("POST /notes status = %d, want %d; body: %s", w.Code, http.StatusCreated, w.Body.String())
	}
	var note models.Note
	if err := json.NewDecoder(w.Body).Decode(&note); err != nil {
		t.Fatalf("decode create response: %v", err)
	}
	return &note
}

func TestHandleCreateNote(t *testing.T) {
	_, mux := newTestHandler(t)

	note := createNote(t, mux, models.CreateNoteRequest{
		ContainerName: "my-app",
		ContainerID:   "abc123",
		NoteContent:   "Test app",
		Tags:          []string{"test"},
	})

	if note.ContainerName != "my-app" {
		t.Errorf("ContainerName = %q, want %q", note.ContainerName, "my-app")
	}
	if note.NoteContent != "Test app" {
		t.Errorf("NoteContent = %q, want %q", note.NoteContent, "Test app")
	}
}

func TestHandleCreateNote_MissingName(t *testing.T) {
	_, mux := newTestHandler(t)

	body, _ := json.Marshal(models.CreateNoteRequest{NoteContent: "no name"})
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPost, "/notes", bytes.NewReader(body))
	mux.ServeHTTP(w, r)

	if w.Code != http.StatusBadRequest {
		t.Errorf("status = %d, want %d", w.Code, http.StatusBadRequest)
	}
}

func TestHandleCreateNote_Duplicate(t *testing.T) {
	_, mux := newTestHandler(t)

	createNote(t, mux, models.CreateNoteRequest{ContainerName: "dup", NoteContent: "first"})

	body, _ := json.Marshal(models.CreateNoteRequest{ContainerName: "dup", NoteContent: "second"})
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPost, "/notes", bytes.NewReader(body))
	mux.ServeHTTP(w, r)

	if w.Code != http.StatusConflict {
		t.Errorf("status = %d, want %d", w.Code, http.StatusConflict)
	}
}

func TestHandleListNotes(t *testing.T) {
	_, mux := newTestHandler(t)

	createNote(t, mux, models.CreateNoteRequest{ContainerName: "a", NoteContent: "first"})
	createNote(t, mux, models.CreateNoteRequest{ContainerName: "b", NoteContent: "second"})

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/notes", http.NoBody)
	mux.ServeHTTP(w, r)

	if w.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", w.Code, http.StatusOK)
	}

	var notes []models.Note
	if err := json.NewDecoder(w.Body).Decode(&notes); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if len(notes) != 2 {
		t.Errorf("len = %d, want 2", len(notes))
	}
}

func TestHandleListNotes_Empty(t *testing.T) {
	_, mux := newTestHandler(t)

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/notes", http.NoBody)
	mux.ServeHTTP(w, r)

	if w.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", w.Code, http.StatusOK)
	}

	var notes []models.Note
	if err := json.NewDecoder(w.Body).Decode(&notes); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if len(notes) != 0 {
		t.Errorf("len = %d, want 0", len(notes))
	}
}

func TestHandleGetNote(t *testing.T) {
	_, mux := newTestHandler(t)

	createNote(t, mux, models.CreateNoteRequest{ContainerName: "web", NoteContent: "web server"})

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/notes/web", http.NoBody)
	mux.ServeHTTP(w, r)

	if w.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", w.Code, http.StatusOK)
	}

	var note models.Note
	if err := json.NewDecoder(w.Body).Decode(&note); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if note.ContainerName != "web" {
		t.Errorf("ContainerName = %q, want %q", note.ContainerName, "web")
	}
}

func TestHandleGetNote_NotFound(t *testing.T) {
	_, mux := newTestHandler(t)

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/notes/nonexistent", http.NoBody)
	mux.ServeHTTP(w, r)

	if w.Code != http.StatusNotFound {
		t.Errorf("status = %d, want %d", w.Code, http.StatusNotFound)
	}
}

func TestHandleUpdateNote(t *testing.T) {
	_, mux := newTestHandler(t)

	createNote(t, mux, models.CreateNoteRequest{ContainerName: "app", NoteContent: "old"})

	updated := "new content"
	body, _ := json.Marshal(models.UpdateNoteRequest{NoteContent: &updated})
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPut, "/notes/app", bytes.NewReader(body))
	mux.ServeHTTP(w, r)

	if w.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d; body: %s", w.Code, http.StatusOK, w.Body.String())
	}

	var note models.Note
	if err := json.NewDecoder(w.Body).Decode(&note); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if note.NoteContent != "new content" {
		t.Errorf("NoteContent = %q, want %q", note.NoteContent, "new content")
	}
}

func TestHandleUpdateNote_NotFound(t *testing.T) {
	_, mux := newTestHandler(t)

	content := "data"
	body, _ := json.Marshal(models.UpdateNoteRequest{NoteContent: &content})
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPut, "/notes/missing", bytes.NewReader(body))
	mux.ServeHTTP(w, r)

	if w.Code != http.StatusNotFound {
		t.Errorf("status = %d, want %d", w.Code, http.StatusNotFound)
	}
}

func TestHandleDeleteNote(t *testing.T) {
	_, mux := newTestHandler(t)

	createNote(t, mux, models.CreateNoteRequest{ContainerName: "bye", NoteContent: "delete me"})

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodDelete, "/notes/bye", http.NoBody)
	mux.ServeHTTP(w, r)

	if w.Code != http.StatusNoContent {
		t.Errorf("status = %d, want %d", w.Code, http.StatusNoContent)
	}

	// Confirm deleted.
	w2 := httptest.NewRecorder()
	r2 := httptest.NewRequest(http.MethodGet, "/notes/bye", http.NoBody)
	mux.ServeHTTP(w2, r2)
	if w2.Code != http.StatusNotFound {
		t.Errorf("after delete, GET status = %d, want %d", w2.Code, http.StatusNotFound)
	}
}

func TestHandleDeleteNote_NotFound(t *testing.T) {
	_, mux := newTestHandler(t)

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodDelete, "/notes/nope", http.NoBody)
	mux.ServeHTTP(w, r)

	if w.Code != http.StatusNotFound {
		t.Errorf("status = %d, want %d", w.Code, http.StatusNotFound)
	}
}
