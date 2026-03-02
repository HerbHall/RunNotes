package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
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
		Title:         "App Notes",
		NoteContent:   "Test app",
		Tags:          []string{"test"},
	})

	if note.ContainerName != "my-app" {
		t.Errorf("ContainerName = %q, want %q", note.ContainerName, "my-app")
	}
	if note.Title != "App Notes" {
		t.Errorf("Title = %q, want %q", note.Title, "App Notes")
	}
	if note.NoteContent != "Test app" {
		t.Errorf("NoteContent = %q, want %q", note.NoteContent, "Test app")
	}
}

func TestHandleCreateNote_MultiplePerContainer(t *testing.T) {
	_, mux := newTestHandler(t)

	n1 := createNote(t, mux, models.CreateNoteRequest{
		ContainerName: "web",
		Title:         "First",
		NoteContent:   "one",
	})
	n2 := createNote(t, mux, models.CreateNoteRequest{
		ContainerName: "web",
		Title:         "Second",
		NoteContent:   "two",
	})

	if n1.ID == n2.ID {
		t.Errorf("expected different IDs, got %d for both", n1.ID)
	}
}

func TestHandleCreateNote_MissingName(t *testing.T) {
	_, mux := newTestHandler(t)

	body, _ := json.Marshal(models.CreateNoteRequest{Title: "No Container", NoteContent: "no name"})
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPost, "/notes", bytes.NewReader(body))
	mux.ServeHTTP(w, r)

	if w.Code != http.StatusBadRequest {
		t.Errorf("status = %d, want %d", w.Code, http.StatusBadRequest)
	}
}

func TestHandleCreateNote_MissingTitle(t *testing.T) {
	_, mux := newTestHandler(t)

	body, _ := json.Marshal(models.CreateNoteRequest{ContainerName: "web", NoteContent: "no title"})
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPost, "/notes", bytes.NewReader(body))
	mux.ServeHTTP(w, r)

	if w.Code != http.StatusBadRequest {
		t.Errorf("status = %d, want %d", w.Code, http.StatusBadRequest)
	}
}

func TestHandleListNotes(t *testing.T) {
	_, mux := newTestHandler(t)

	createNote(t, mux, models.CreateNoteRequest{ContainerName: "a", Title: "Note A", NoteContent: "first"})
	createNote(t, mux, models.CreateNoteRequest{ContainerName: "b", Title: "Note B", NoteContent: "second"})

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

func TestHandleListContainerNotes(t *testing.T) {
	_, mux := newTestHandler(t)

	createNote(t, mux, models.CreateNoteRequest{ContainerName: "web", Title: "Note 1", NoteContent: "a"})
	createNote(t, mux, models.CreateNoteRequest{ContainerName: "web", Title: "Note 2", NoteContent: "b"})
	createNote(t, mux, models.CreateNoteRequest{ContainerName: "db", Title: "DB Note", NoteContent: "c"})

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/notes/container/web", http.NoBody)
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

func TestHandleGetNote(t *testing.T) {
	_, mux := newTestHandler(t)

	created := createNote(t, mux, models.CreateNoteRequest{
		ContainerName: "web",
		Title:         "Web Notes",
		NoteContent:   "web server",
	})

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/notes/%d", created.ID), http.NoBody)
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
	if note.Title != "Web Notes" {
		t.Errorf("Title = %q, want %q", note.Title, "Web Notes")
	}
}

func TestHandleGetNote_NotFound(t *testing.T) {
	_, mux := newTestHandler(t)

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/notes/999", http.NoBody)
	mux.ServeHTTP(w, r)

	if w.Code != http.StatusNotFound {
		t.Errorf("status = %d, want %d", w.Code, http.StatusNotFound)
	}
}

func TestHandleUpdateNote(t *testing.T) {
	_, mux := newTestHandler(t)

	created := createNote(t, mux, models.CreateNoteRequest{
		ContainerName: "app",
		Title:         "App Notes",
		NoteContent:   "old",
	})

	updated := "new content"
	body, _ := json.Marshal(models.UpdateNoteRequest{NoteContent: &updated})
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPut, fmt.Sprintf("/notes/%d", created.ID), bytes.NewReader(body))
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
	r := httptest.NewRequest(http.MethodPut, "/notes/999", bytes.NewReader(body))
	mux.ServeHTTP(w, r)

	if w.Code != http.StatusNotFound {
		t.Errorf("status = %d, want %d", w.Code, http.StatusNotFound)
	}
}

func TestHandleDeleteNote(t *testing.T) {
	_, mux := newTestHandler(t)

	created := createNote(t, mux, models.CreateNoteRequest{
		ContainerName: "bye",
		Title:         "Delete Me",
		NoteContent:   "delete me",
	})

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/notes/%d", created.ID), http.NoBody)
	mux.ServeHTTP(w, r)

	if w.Code != http.StatusNoContent {
		t.Errorf("status = %d, want %d", w.Code, http.StatusNoContent)
	}

	// Confirm deleted.
	w2 := httptest.NewRecorder()
	r2 := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/notes/%d", created.ID), http.NoBody)
	mux.ServeHTTP(w2, r2)
	if w2.Code != http.StatusNotFound {
		t.Errorf("after delete, GET status = %d, want %d", w2.Code, http.StatusNotFound)
	}
}

func TestHandleDeleteNote_NotFound(t *testing.T) {
	_, mux := newTestHandler(t)

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodDelete, "/notes/999", http.NoBody)
	mux.ServeHTTP(w, r)

	if w.Code != http.StatusNotFound {
		t.Errorf("status = %d, want %d", w.Code, http.StatusNotFound)
	}
}

func TestHandleDeleteContainerNotes(t *testing.T) {
	_, mux := newTestHandler(t)

	createNote(t, mux, models.CreateNoteRequest{ContainerName: "web", Title: "Note 1", NoteContent: "a"})
	createNote(t, mux, models.CreateNoteRequest{ContainerName: "web", Title: "Note 2", NoteContent: "b"})

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodDelete, "/notes/container/web", http.NoBody)
	mux.ServeHTTP(w, r)

	if w.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d; body: %s", w.Code, http.StatusOK, w.Body.String())
	}

	var result map[string]int64
	if err := json.NewDecoder(w.Body).Decode(&result); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if result["deleted"] != 2 {
		t.Errorf("deleted = %d, want 2", result["deleted"])
	}
}

func TestHandleExportNotes(t *testing.T) {
	_, mux := newTestHandler(t)

	createNote(t, mux, models.CreateNoteRequest{ContainerName: "x", Title: "X Note", NoteContent: "one"})
	createNote(t, mux, models.CreateNoteRequest{ContainerName: "y", Title: "Y Note", NoteContent: "two"})

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/notes/export", http.NoBody)
	mux.ServeHTTP(w, r)

	if w.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d; body: %s", w.Code, http.StatusOK, w.Body.String())
	}

	var notes []models.Note
	if err := json.NewDecoder(w.Body).Decode(&notes); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if len(notes) != 2 {
		t.Errorf("len = %d, want 2", len(notes))
	}
}

func TestHandleExportNotes_Empty(t *testing.T) {
	_, mux := newTestHandler(t)

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/notes/export", http.NoBody)
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

func TestHandleImportNotes(t *testing.T) {
	_, mux := newTestHandler(t)

	payload := []models.Note{
		{ContainerName: "svc-a", Title: "SVC A", NoteContent: "alpha", Tags: []string{"a"}},
		{ContainerName: "svc-b", Title: "SVC B", NoteContent: "beta"},
	}
	body, _ := json.Marshal(payload)

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPost, "/notes/import", bytes.NewReader(body))
	mux.ServeHTTP(w, r)

	if w.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d; body: %s", w.Code, http.StatusOK, w.Body.String())
	}

	var result map[string]int
	if err := json.NewDecoder(w.Body).Decode(&result); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if result["imported"] != 2 {
		t.Errorf("imported = %d, want 2", result["imported"])
	}
}

func TestHandleImportNotes_InvalidJSON(t *testing.T) {
	_, mux := newTestHandler(t)

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPost, "/notes/import", bytes.NewReader([]byte("not-json")))
	mux.ServeHTTP(w, r)

	if w.Code != http.StatusBadRequest {
		t.Errorf("status = %d, want %d", w.Code, http.StatusBadRequest)
	}
}
