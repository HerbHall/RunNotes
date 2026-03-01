package database

import (
	"context"
	"testing"
)

func TestOpen_InMemory(t *testing.T) {
	db, err := Open(":memory:")
	if err != nil {
		t.Fatalf("Open(:memory:): %v", err)
	}
	defer func() { _ = db.Close() }()

	// Verify WAL mode is active.
	var journalMode string
	err = db.QueryRowContext(context.Background(), "PRAGMA journal_mode").Scan(&journalMode)
	if err != nil {
		t.Fatalf("query journal_mode: %v", err)
	}
	// In-memory databases use "memory" journal mode, not "wal".
	if journalMode != "memory" && journalMode != "wal" {
		t.Errorf("journal_mode = %q, want memory or wal", journalMode)
	}

	// Verify foreign keys are enabled.
	var fk int
	err = db.QueryRowContext(context.Background(), "PRAGMA foreign_keys").Scan(&fk)
	if err != nil {
		t.Fatalf("query foreign_keys: %v", err)
	}
	if fk != 1 {
		t.Errorf("foreign_keys = %d, want 1", fk)
	}
}

func TestOpen_CreatesNotesTable(t *testing.T) {
	db, err := Open(":memory:")
	if err != nil {
		t.Fatalf("Open(:memory:): %v", err)
	}
	defer func() { _ = db.Close() }()

	// Verify notes table exists with expected columns.
	rows, err := db.QueryContext(context.Background(), "PRAGMA table_info(notes)")
	if err != nil {
		t.Fatalf("table_info: %v", err)
	}
	defer func() { _ = rows.Close() }()

	columns := make(map[string]bool)
	for rows.Next() {
		var cid int
		var name, ctype string
		var notnull, pk int
		var dfltValue *string
		if err := rows.Scan(&cid, &name, &ctype, &notnull, &dfltValue, &pk); err != nil {
			t.Fatalf("scan column: %v", err)
		}
		columns[name] = true
	}
	if err := rows.Err(); err != nil {
		t.Fatalf("rows iteration: %v", err)
	}

	expected := []string{
		"id", "container_name", "container_id",
		"compose_project", "compose_service",
		"note_content", "pinned", "tags",
		"created_at", "updated_at",
	}
	for _, col := range expected {
		if !columns[col] {
			t.Errorf("missing column %q in notes table", col)
		}
	}
}

func TestOpen_Idempotent(t *testing.T) {
	db, err := Open(":memory:")
	if err != nil {
		t.Fatalf("first Open: %v", err)
	}
	_ = db.Close()

	// Opening again should not fail (IF NOT EXISTS).
	db2, err := Open(":memory:")
	if err != nil {
		t.Fatalf("second Open: %v", err)
	}
	_ = db2.Close()
}

func TestOpen_UniqueIndex(t *testing.T) {
	db, err := Open(":memory:")
	if err != nil {
		t.Fatalf("Open(:memory:): %v", err)
	}
	defer func() { _ = db.Close() }()

	ctx := context.Background()
	_, err = db.ExecContext(ctx,
		"INSERT INTO notes (container_name, note_content) VALUES (?, ?)",
		"test-container", "first note")
	if err != nil {
		t.Fatalf("first insert: %v", err)
	}

	_, err = db.ExecContext(ctx,
		"INSERT INTO notes (container_name, note_content) VALUES (?, ?)",
		"test-container", "duplicate note")
	if err == nil {
		t.Fatal("expected unique constraint error on duplicate container_name, got nil")
	}
}
