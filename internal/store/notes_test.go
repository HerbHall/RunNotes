package store

import (
	"context"
	"database/sql"
	"testing"

	"github.com/HerbHall/RunNotes/internal/database"
	"github.com/HerbHall/RunNotes/internal/models"
)

func newTestDB(t *testing.T) *sql.DB {
	t.Helper()
	db, err := database.Open(":memory:")
	if err != nil {
		t.Fatalf("open test db: %v", err)
	}
	t.Cleanup(func() { _ = db.Close() })
	return db
}

func newTestStore(t *testing.T) *NoteStore {
	t.Helper()
	return NewNoteStore(newTestDB(t))
}

func TestCreate(t *testing.T) {
	s := newTestStore(t)
	ctx := context.Background()

	note, err := s.Create(ctx, models.CreateNoteRequest{
		ContainerName:  "my-postgres",
		ContainerID:    "abc123",
		ComposeProject: "myapp",
		ComposeService: "db",
		NoteContent:    "Production database",
		Tags:           []string{"db", "prod"},
	})
	if err != nil {
		t.Fatalf("Create: %v", err)
	}

	if note.ID == 0 {
		t.Error("expected non-zero ID")
	}
	if note.ContainerName != "my-postgres" {
		t.Errorf("ContainerName = %q, want %q", note.ContainerName, "my-postgres")
	}
	if note.NoteContent != "Production database" {
		t.Errorf("NoteContent = %q, want %q", note.NoteContent, "Production database")
	}
	if note.Pinned {
		t.Error("expected Pinned = false")
	}
	if len(note.Tags) != 2 || note.Tags[0] != "db" || note.Tags[1] != "prod" {
		t.Errorf("Tags = %v, want [db prod]", note.Tags)
	}
	if note.CreatedAt.IsZero() {
		t.Error("expected non-zero CreatedAt")
	}
}

func TestCreate_DuplicateName(t *testing.T) {
	s := newTestStore(t)
	ctx := context.Background()

	_, err := s.Create(ctx, models.CreateNoteRequest{
		ContainerName: "my-container",
		NoteContent:   "first",
	})
	if err != nil {
		t.Fatalf("first Create: %v", err)
	}

	_, err = s.Create(ctx, models.CreateNoteRequest{
		ContainerName: "my-container",
		NoteContent:   "duplicate",
	})
	if err == nil {
		t.Fatal("expected error on duplicate container_name, got nil")
	}
}

func TestCreate_NilTags(t *testing.T) {
	s := newTestStore(t)
	ctx := context.Background()

	note, err := s.Create(ctx, models.CreateNoteRequest{
		ContainerName: "no-tags",
		NoteContent:   "test",
	})
	if err != nil {
		t.Fatalf("Create: %v", err)
	}
	if note.Tags == nil {
		t.Error("Tags should be empty slice, not nil")
	}
	if len(note.Tags) != 0 {
		t.Errorf("Tags length = %d, want 0", len(note.Tags))
	}
}

func TestGetByName(t *testing.T) {
	s := newTestStore(t)
	ctx := context.Background()

	_, err := s.Create(ctx, models.CreateNoteRequest{
		ContainerName: "my-redis",
		ContainerID:   "def456",
		NoteContent:   "Cache server",
	})
	if err != nil {
		t.Fatalf("Create: %v", err)
	}

	note, err := s.GetByName(ctx, "my-redis")
	if err != nil {
		t.Fatalf("GetByName: %v", err)
	}
	if note.ContainerName != "my-redis" {
		t.Errorf("ContainerName = %q, want %q", note.ContainerName, "my-redis")
	}
	if note.NoteContent != "Cache server" {
		t.Errorf("NoteContent = %q, want %q", note.NoteContent, "Cache server")
	}
}

func TestGetByName_NotFound(t *testing.T) {
	s := newTestStore(t)
	ctx := context.Background()

	_, err := s.GetByName(ctx, "nonexistent")
	if err != ErrNotFound {
		t.Errorf("err = %v, want ErrNotFound", err)
	}
}

func TestList(t *testing.T) {
	s := newTestStore(t)
	ctx := context.Background()

	_, _ = s.Create(ctx, models.CreateNoteRequest{ContainerName: "a", NoteContent: "first"})
	_, _ = s.Create(ctx, models.CreateNoteRequest{ContainerName: "b", NoteContent: "second"})

	notes, err := s.List(ctx, nil, "")
	if err != nil {
		t.Fatalf("List: %v", err)
	}
	if len(notes) != 2 {
		t.Errorf("len = %d, want 2", len(notes))
	}
}

func TestList_Empty(t *testing.T) {
	s := newTestStore(t)
	ctx := context.Background()

	notes, err := s.List(ctx, nil, "")
	if err != nil {
		t.Fatalf("List: %v", err)
	}
	if notes == nil {
		t.Error("expected empty slice, not nil")
	}
	if len(notes) != 0 {
		t.Errorf("len = %d, want 0", len(notes))
	}
}

func TestList_FilterPinned(t *testing.T) {
	s := newTestStore(t)
	ctx := context.Background()

	_, _ = s.Create(ctx, models.CreateNoteRequest{ContainerName: "unpinned", NoteContent: "test"})
	_, _ = s.Create(ctx, models.CreateNoteRequest{ContainerName: "pinned", NoteContent: "test"})

	pinTrue := true
	_, _ = s.Update(ctx, "pinned", models.UpdateNoteRequest{Pinned: &pinTrue})

	notes, err := s.List(ctx, &pinTrue, "")
	if err != nil {
		t.Fatalf("List pinned: %v", err)
	}
	if len(notes) != 1 {
		t.Fatalf("len = %d, want 1", len(notes))
	}
	if notes[0].ContainerName != "pinned" {
		t.Errorf("ContainerName = %q, want %q", notes[0].ContainerName, "pinned")
	}
}

func TestList_Search(t *testing.T) {
	s := newTestStore(t)
	ctx := context.Background()

	_, _ = s.Create(ctx, models.CreateNoteRequest{ContainerName: "web-server", NoteContent: "serves HTTP"})
	_, _ = s.Create(ctx, models.CreateNoteRequest{ContainerName: "db-server", NoteContent: "stores data"})

	notes, err := s.List(ctx, nil, "HTTP")
	if err != nil {
		t.Fatalf("List search: %v", err)
	}
	if len(notes) != 1 {
		t.Fatalf("len = %d, want 1", len(notes))
	}
	if notes[0].ContainerName != "web-server" {
		t.Errorf("ContainerName = %q, want %q", notes[0].ContainerName, "web-server")
	}
}

func TestUpdate_NoteContent(t *testing.T) {
	s := newTestStore(t)
	ctx := context.Background()

	_, _ = s.Create(ctx, models.CreateNoteRequest{ContainerName: "test", NoteContent: "old"})

	updated := "new content"
	note, err := s.Update(ctx, "test", models.UpdateNoteRequest{NoteContent: &updated})
	if err != nil {
		t.Fatalf("Update: %v", err)
	}
	if note.NoteContent != "new content" {
		t.Errorf("NoteContent = %q, want %q", note.NoteContent, "new content")
	}
}

func TestUpdate_Pinned(t *testing.T) {
	s := newTestStore(t)
	ctx := context.Background()

	_, _ = s.Create(ctx, models.CreateNoteRequest{ContainerName: "test", NoteContent: "data"})

	pinTrue := true
	note, err := s.Update(ctx, "test", models.UpdateNoteRequest{Pinned: &pinTrue})
	if err != nil {
		t.Fatalf("Update: %v", err)
	}
	if !note.Pinned {
		t.Error("expected Pinned = true")
	}
}

func TestUpdate_Tags(t *testing.T) {
	s := newTestStore(t)
	ctx := context.Background()

	_, _ = s.Create(ctx, models.CreateNoteRequest{
		ContainerName: "test",
		NoteContent:   "data",
		Tags:          []string{"old"},
	})

	newTags := []string{"new", "tags"}
	note, err := s.Update(ctx, "test", models.UpdateNoteRequest{Tags: &newTags})
	if err != nil {
		t.Fatalf("Update: %v", err)
	}
	if len(note.Tags) != 2 || note.Tags[0] != "new" || note.Tags[1] != "tags" {
		t.Errorf("Tags = %v, want [new tags]", note.Tags)
	}
}

func TestUpdate_NotFound(t *testing.T) {
	s := newTestStore(t)
	ctx := context.Background()

	content := "data"
	_, err := s.Update(ctx, "nonexistent", models.UpdateNoteRequest{NoteContent: &content})
	if err != ErrNotFound {
		t.Errorf("err = %v, want ErrNotFound", err)
	}
}

func TestDelete(t *testing.T) {
	s := newTestStore(t)
	ctx := context.Background()

	_, _ = s.Create(ctx, models.CreateNoteRequest{ContainerName: "to-delete", NoteContent: "bye"})

	err := s.Delete(ctx, "to-delete")
	if err != nil {
		t.Fatalf("Delete: %v", err)
	}

	_, err = s.GetByName(ctx, "to-delete")
	if err != ErrNotFound {
		t.Errorf("after delete, GetByName err = %v, want ErrNotFound", err)
	}
}

func TestDelete_NotFound(t *testing.T) {
	s := newTestStore(t)
	ctx := context.Background()

	err := s.Delete(ctx, "nonexistent")
	if err != ErrNotFound {
		t.Errorf("err = %v, want ErrNotFound", err)
	}
}

func TestExportAll(t *testing.T) {
	s := newTestStore(t)
	ctx := context.Background()

	_, _ = s.Create(ctx, models.CreateNoteRequest{
		ContainerName: "web",
		NoteContent:   "web server",
		Tags:          []string{"http"},
	})
	_, _ = s.Create(ctx, models.CreateNoteRequest{
		ContainerName: "db",
		NoteContent:   "database",
		Tags:          []string{"sql"},
	})

	notes, err := s.ExportAll(ctx)
	if err != nil {
		t.Fatalf("ExportAll: %v", err)
	}
	if len(notes) != 2 {
		t.Errorf("len = %d, want 2", len(notes))
	}
}

func TestImportAll(t *testing.T) {
	s := newTestStore(t)
	ctx := context.Background()

	notes := []models.Note{
		{ContainerName: "a", NoteContent: "note-a", Tags: []string{"tag1"}},
		{ContainerName: "b", NoteContent: "note-b"},
		{ContainerName: "c", NoteContent: "note-c"},
	}

	imported, err := s.ImportAll(ctx, notes)
	if err != nil {
		t.Fatalf("ImportAll: %v", err)
	}
	if imported != 3 {
		t.Errorf("imported = %d, want 3", imported)
	}

	all, err := s.List(ctx, nil, "")
	if err != nil {
		t.Fatalf("List: %v", err)
	}
	if len(all) != 3 {
		t.Errorf("len = %d, want 3", len(all))
	}
}

func TestImportAll_Upsert(t *testing.T) {
	s := newTestStore(t)
	ctx := context.Background()

	_, err := s.Create(ctx, models.CreateNoteRequest{
		ContainerName: "existing",
		NoteContent:   "old content",
		Tags:          []string{"old"},
	})
	if err != nil {
		t.Fatalf("Create: %v", err)
	}

	notes := []models.Note{
		{ContainerName: "existing", NoteContent: "new content", Pinned: true, Tags: []string{"new"}},
	}

	imported, err := s.ImportAll(ctx, notes)
	if err != nil {
		t.Fatalf("ImportAll: %v", err)
	}
	if imported != 1 {
		t.Errorf("imported = %d, want 1", imported)
	}

	got, err := s.GetByName(ctx, "existing")
	if err != nil {
		t.Fatalf("GetByName: %v", err)
	}
	if got.NoteContent != "new content" {
		t.Errorf("NoteContent = %q, want %q", got.NoteContent, "new content")
	}
	if !got.Pinned {
		t.Error("expected Pinned = true after upsert")
	}
	if len(got.Tags) != 1 || got.Tags[0] != "new" {
		t.Errorf("Tags = %v, want [new]", got.Tags)
	}
}

func TestImportAll_SkipsEmptyName(t *testing.T) {
	s := newTestStore(t)
	ctx := context.Background()

	notes := []models.Note{
		{ContainerName: "", NoteContent: "should be skipped"},
		{ContainerName: "valid", NoteContent: "included"},
	}

	imported, err := s.ImportAll(ctx, notes)
	if err != nil {
		t.Fatalf("ImportAll: %v", err)
	}
	if imported != 1 {
		t.Errorf("imported = %d, want 1", imported)
	}
}
