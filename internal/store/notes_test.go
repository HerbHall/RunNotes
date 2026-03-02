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
		Title:          "Config Notes",
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
	if note.Title != "Config Notes" {
		t.Errorf("Title = %q, want %q", note.Title, "Config Notes")
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

func TestCreate_MultiplePerContainer(t *testing.T) {
	s := newTestStore(t)
	ctx := context.Background()

	n1, err := s.Create(ctx, models.CreateNoteRequest{
		ContainerName: "my-container",
		Title:         "First Note",
		NoteContent:   "first",
	})
	if err != nil {
		t.Fatalf("first Create: %v", err)
	}

	n2, err := s.Create(ctx, models.CreateNoteRequest{
		ContainerName: "my-container",
		Title:         "Second Note",
		NoteContent:   "second",
	})
	if err != nil {
		t.Fatalf("second Create: %v", err)
	}

	if n1.ID == n2.ID {
		t.Errorf("expected different IDs, got %d for both", n1.ID)
	}

	notes, err := s.ListByContainer(ctx, "my-container")
	if err != nil {
		t.Fatalf("ListByContainer: %v", err)
	}
	if len(notes) != 2 {
		t.Errorf("len = %d, want 2", len(notes))
	}
}

func TestCreate_NilTags(t *testing.T) {
	s := newTestStore(t)
	ctx := context.Background()

	note, err := s.Create(ctx, models.CreateNoteRequest{
		ContainerName: "no-tags",
		Title:         "No Tags",
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

func TestGetByID(t *testing.T) {
	s := newTestStore(t)
	ctx := context.Background()

	created, err := s.Create(ctx, models.CreateNoteRequest{
		ContainerName: "my-redis",
		ContainerID:   "def456",
		Title:         "Redis Config",
		NoteContent:   "Cache server",
	})
	if err != nil {
		t.Fatalf("Create: %v", err)
	}

	note, err := s.GetByID(ctx, created.ID)
	if err != nil {
		t.Fatalf("GetByID: %v", err)
	}
	if note.ContainerName != "my-redis" {
		t.Errorf("ContainerName = %q, want %q", note.ContainerName, "my-redis")
	}
	if note.Title != "Redis Config" {
		t.Errorf("Title = %q, want %q", note.Title, "Redis Config")
	}
	if note.NoteContent != "Cache server" {
		t.Errorf("NoteContent = %q, want %q", note.NoteContent, "Cache server")
	}
}

func TestGetByID_NotFound(t *testing.T) {
	s := newTestStore(t)
	ctx := context.Background()

	_, err := s.GetByID(ctx, 999)
	if err != ErrNotFound {
		t.Errorf("err = %v, want ErrNotFound", err)
	}
}

func TestListByContainer(t *testing.T) {
	s := newTestStore(t)
	ctx := context.Background()

	_, _ = s.Create(ctx, models.CreateNoteRequest{ContainerName: "web", Title: "Note A", NoteContent: "first"})
	_, _ = s.Create(ctx, models.CreateNoteRequest{ContainerName: "web", Title: "Note B", NoteContent: "second"})
	_, _ = s.Create(ctx, models.CreateNoteRequest{ContainerName: "db", Title: "DB Note", NoteContent: "third"})

	notes, err := s.ListByContainer(ctx, "web")
	if err != nil {
		t.Fatalf("ListByContainer: %v", err)
	}
	if len(notes) != 2 {
		t.Errorf("len = %d, want 2", len(notes))
	}

	notes, err = s.ListByContainer(ctx, "db")
	if err != nil {
		t.Fatalf("ListByContainer: %v", err)
	}
	if len(notes) != 1 {
		t.Errorf("len = %d, want 1", len(notes))
	}
}

func TestListByContainer_Empty(t *testing.T) {
	s := newTestStore(t)
	ctx := context.Background()

	notes, err := s.ListByContainer(ctx, "nonexistent")
	if err != nil {
		t.Fatalf("ListByContainer: %v", err)
	}
	if notes == nil {
		t.Error("expected empty slice, not nil")
	}
	if len(notes) != 0 {
		t.Errorf("len = %d, want 0", len(notes))
	}
}

func TestList(t *testing.T) {
	s := newTestStore(t)
	ctx := context.Background()

	_, _ = s.Create(ctx, models.CreateNoteRequest{ContainerName: "a", Title: "Note A", NoteContent: "first"})
	_, _ = s.Create(ctx, models.CreateNoteRequest{ContainerName: "b", Title: "Note B", NoteContent: "second"})

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

	_, _ = s.Create(ctx, models.CreateNoteRequest{ContainerName: "unpinned", Title: "Unpinned", NoteContent: "test"})
	pinned, _ := s.Create(ctx, models.CreateNoteRequest{ContainerName: "pinned", Title: "Pinned", NoteContent: "test"})

	pinTrue := true
	_, _ = s.Update(ctx, pinned.ID, models.UpdateNoteRequest{Pinned: &pinTrue})

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

	_, _ = s.Create(ctx, models.CreateNoteRequest{ContainerName: "web-server", Title: "Web", NoteContent: "serves HTTP"})
	_, _ = s.Create(ctx, models.CreateNoteRequest{ContainerName: "db-server", Title: "DB", NoteContent: "stores data"})

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

func TestList_SearchByTitle(t *testing.T) {
	s := newTestStore(t)
	ctx := context.Background()

	_, _ = s.Create(ctx, models.CreateNoteRequest{ContainerName: "app", Title: "Deployment Guide", NoteContent: "steps"})
	_, _ = s.Create(ctx, models.CreateNoteRequest{ContainerName: "app", Title: "Config Notes", NoteContent: "settings"})

	notes, err := s.List(ctx, nil, "Deployment")
	if err != nil {
		t.Fatalf("List search title: %v", err)
	}
	if len(notes) != 1 {
		t.Fatalf("len = %d, want 1", len(notes))
	}
	if notes[0].Title != "Deployment Guide" {
		t.Errorf("Title = %q, want %q", notes[0].Title, "Deployment Guide")
	}
}

func TestList_SearchByTag(t *testing.T) {
	s := newTestStore(t)
	ctx := context.Background()

	_, _ = s.Create(ctx, models.CreateNoteRequest{
		ContainerName: "web-app",
		Title:         "Web App",
		NoteContent:   "serves pages",
		Tags:          []string{"web", "frontend"},
	})
	_, _ = s.Create(ctx, models.CreateNoteRequest{
		ContainerName: "pg-server",
		Title:         "PG Server",
		NoteContent:   "stores data",
		Tags:          []string{"database", "backend"},
	})

	notes, err := s.List(ctx, nil, "frontend")
	if err != nil {
		t.Fatalf("List search frontend: %v", err)
	}
	if len(notes) != 1 {
		t.Fatalf("len = %d, want 1", len(notes))
	}
	if notes[0].ContainerName != "web-app" {
		t.Errorf("ContainerName = %q, want %q", notes[0].ContainerName, "web-app")
	}

	notes, err = s.List(ctx, nil, "database")
	if err != nil {
		t.Fatalf("List search database: %v", err)
	}
	if len(notes) != 1 {
		t.Fatalf("len = %d, want 1", len(notes))
	}
	if notes[0].ContainerName != "pg-server" {
		t.Errorf("ContainerName = %q, want %q", notes[0].ContainerName, "pg-server")
	}
}

func TestUpdate_NoteContent(t *testing.T) {
	s := newTestStore(t)
	ctx := context.Background()

	created, _ := s.Create(ctx, models.CreateNoteRequest{ContainerName: "test", Title: "Test", NoteContent: "old"})

	updated := "new content"
	note, err := s.Update(ctx, created.ID, models.UpdateNoteRequest{NoteContent: &updated})
	if err != nil {
		t.Fatalf("Update: %v", err)
	}
	if note.NoteContent != "new content" {
		t.Errorf("NoteContent = %q, want %q", note.NoteContent, "new content")
	}
}

func TestUpdate_Title(t *testing.T) {
	s := newTestStore(t)
	ctx := context.Background()

	created, _ := s.Create(ctx, models.CreateNoteRequest{ContainerName: "test", Title: "Old Title", NoteContent: "data"})

	newTitle := "New Title"
	note, err := s.Update(ctx, created.ID, models.UpdateNoteRequest{Title: &newTitle})
	if err != nil {
		t.Fatalf("Update: %v", err)
	}
	if note.Title != "New Title" {
		t.Errorf("Title = %q, want %q", note.Title, "New Title")
	}
}

func TestUpdate_Pinned(t *testing.T) {
	s := newTestStore(t)
	ctx := context.Background()

	created, _ := s.Create(ctx, models.CreateNoteRequest{ContainerName: "test", Title: "Test", NoteContent: "data"})

	pinTrue := true
	note, err := s.Update(ctx, created.ID, models.UpdateNoteRequest{Pinned: &pinTrue})
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

	created, _ := s.Create(ctx, models.CreateNoteRequest{
		ContainerName: "test",
		Title:         "Test",
		NoteContent:   "data",
		Tags:          []string{"old"},
	})

	newTags := []string{"new", "tags"}
	note, err := s.Update(ctx, created.ID, models.UpdateNoteRequest{Tags: &newTags})
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
	_, err := s.Update(ctx, 999, models.UpdateNoteRequest{NoteContent: &content})
	if err != ErrNotFound {
		t.Errorf("err = %v, want ErrNotFound", err)
	}
}

func TestDelete(t *testing.T) {
	s := newTestStore(t)
	ctx := context.Background()

	created, _ := s.Create(ctx, models.CreateNoteRequest{ContainerName: "to-delete", Title: "Delete Me", NoteContent: "bye"})

	err := s.Delete(ctx, created.ID)
	if err != nil {
		t.Fatalf("Delete: %v", err)
	}

	_, err = s.GetByID(ctx, created.ID)
	if err != ErrNotFound {
		t.Errorf("after delete, GetByID err = %v, want ErrNotFound", err)
	}
}

func TestDelete_NotFound(t *testing.T) {
	s := newTestStore(t)
	ctx := context.Background()

	err := s.Delete(ctx, 999)
	if err != ErrNotFound {
		t.Errorf("err = %v, want ErrNotFound", err)
	}
}

func TestDeleteByContainer(t *testing.T) {
	s := newTestStore(t)
	ctx := context.Background()

	_, _ = s.Create(ctx, models.CreateNoteRequest{ContainerName: "web", Title: "Note 1", NoteContent: "a"})
	_, _ = s.Create(ctx, models.CreateNoteRequest{ContainerName: "web", Title: "Note 2", NoteContent: "b"})
	_, _ = s.Create(ctx, models.CreateNoteRequest{ContainerName: "db", Title: "DB Note", NoteContent: "c"})

	deleted, err := s.DeleteByContainer(ctx, "web")
	if err != nil {
		t.Fatalf("DeleteByContainer: %v", err)
	}
	if deleted != 2 {
		t.Errorf("deleted = %d, want 2", deleted)
	}

	notes, _ := s.ListByContainer(ctx, "web")
	if len(notes) != 0 {
		t.Errorf("web notes len = %d, want 0", len(notes))
	}

	notes, _ = s.ListByContainer(ctx, "db")
	if len(notes) != 1 {
		t.Errorf("db notes len = %d, want 1", len(notes))
	}
}

func TestDeleteByContainer_NoNotes(t *testing.T) {
	s := newTestStore(t)
	ctx := context.Background()

	deleted, err := s.DeleteByContainer(ctx, "nonexistent")
	if err != nil {
		t.Fatalf("DeleteByContainer: %v", err)
	}
	if deleted != 0 {
		t.Errorf("deleted = %d, want 0", deleted)
	}
}

func TestExportAll(t *testing.T) {
	s := newTestStore(t)
	ctx := context.Background()

	_, _ = s.Create(ctx, models.CreateNoteRequest{
		ContainerName: "web",
		Title:         "Web Notes",
		NoteContent:   "web server",
		Tags:          []string{"http"},
	})
	_, _ = s.Create(ctx, models.CreateNoteRequest{
		ContainerName: "db",
		Title:         "DB Notes",
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
		{ContainerName: "a", Title: "Note A", NoteContent: "note-a", Tags: []string{"tag1"}},
		{ContainerName: "b", Title: "Note B", NoteContent: "note-b"},
		{ContainerName: "c", Title: "Note C", NoteContent: "note-c"},
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

	created, err := s.Create(ctx, models.CreateNoteRequest{
		ContainerName: "existing",
		Title:         "My Note",
		NoteContent:   "old content",
		Tags:          []string{"old"},
	})
	if err != nil {
		t.Fatalf("Create: %v", err)
	}

	notes := []models.Note{
		{ContainerName: "existing", Title: "My Note", NoteContent: "new content", Pinned: true, Tags: []string{"new"}},
	}

	imported, err := s.ImportAll(ctx, notes)
	if err != nil {
		t.Fatalf("ImportAll: %v", err)
	}
	if imported != 1 {
		t.Errorf("imported = %d, want 1", imported)
	}

	got, err := s.GetByID(ctx, created.ID)
	if err != nil {
		t.Fatalf("GetByID: %v", err)
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

func TestImportAll_DefaultTitle(t *testing.T) {
	s := newTestStore(t)
	ctx := context.Background()

	notes := []models.Note{
		{ContainerName: "legacy", NoteContent: "no title provided"},
	}

	imported, err := s.ImportAll(ctx, notes)
	if err != nil {
		t.Fatalf("ImportAll: %v", err)
	}
	if imported != 1 {
		t.Errorf("imported = %d, want 1", imported)
	}

	all, _ := s.ListByContainer(ctx, "legacy")
	if len(all) != 1 {
		t.Fatalf("len = %d, want 1", len(all))
	}
	if all[0].Title != "Note" {
		t.Errorf("Title = %q, want %q", all[0].Title, "Note")
	}
}

func TestImportAll_SkipsEmptyName(t *testing.T) {
	s := newTestStore(t)
	ctx := context.Background()

	notes := []models.Note{
		{ContainerName: "", Title: "Skip Me", NoteContent: "should be skipped"},
		{ContainerName: "valid", Title: "Valid Note", NoteContent: "included"},
	}

	imported, err := s.ImportAll(ctx, notes)
	if err != nil {
		t.Fatalf("ImportAll: %v", err)
	}
	if imported != 1 {
		t.Errorf("imported = %d, want 1", imported)
	}
}
