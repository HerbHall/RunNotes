package store

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/HerbHall/RunNotes/internal/models"
)

// ErrNotFound is returned when a note does not exist.
var ErrNotFound = errors.New("note not found")

// NoteStore provides CRUD operations for notes against a SQLite database.
type NoteStore struct {
	db *sql.DB
}

// NewNoteStore creates a NoteStore backed by the given database connection.
func NewNoteStore(db *sql.DB) *NoteStore {
	return &NoteStore{db: db}
}

// columnList is the standard SELECT column order for notes.
const columnList = "id, container_name, container_id, compose_project, compose_service, title, note_content, pinned, tags, created_at, updated_at"

// List returns all notes, optionally filtered by pinned status or search term.
// Results are ordered by pinned DESC, then updated_at DESC.
func (s *NoteStore) List(ctx context.Context, pinned *bool, search string) ([]models.Note, error) {
	query := "SELECT " + columnList + " FROM notes"
	var conditions []string
	var args []any

	if pinned != nil {
		conditions = append(conditions, "pinned = ?")
		if *pinned {
			args = append(args, 1)
		} else {
			args = append(args, 0)
		}
	}
	if search != "" {
		conditions = append(conditions, "(container_name LIKE ? OR title LIKE ? OR note_content LIKE ? OR tags LIKE ?)")
		like := "%" + search + "%"
		args = append(args, like, like, like, like)
	}

	if len(conditions) > 0 {
		query += " WHERE " + strings.Join(conditions, " AND ")
	}
	query += " ORDER BY pinned DESC, updated_at DESC"

	rows, err := s.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("query notes: %w", err)
	}
	defer func() { _ = rows.Close() }()

	var notes []models.Note
	for rows.Next() {
		n, err := scanNote(rows)
		if err != nil {
			return nil, err
		}
		notes = append(notes, n)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate notes: %w", err)
	}

	if notes == nil {
		notes = []models.Note{}
	}
	return notes, nil
}

// ListByContainer returns all notes for the given container name.
func (s *NoteStore) ListByContainer(ctx context.Context, name string) ([]models.Note, error) {
	rows, err := s.db.QueryContext(ctx,
		"SELECT "+columnList+" FROM notes WHERE container_name = ? ORDER BY pinned DESC, updated_at DESC",
		name)
	if err != nil {
		return nil, fmt.Errorf("query notes for container %q: %w", name, err)
	}
	defer func() { _ = rows.Close() }()

	var notes []models.Note
	for rows.Next() {
		n, err := scanNote(rows)
		if err != nil {
			return nil, err
		}
		notes = append(notes, n)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate notes for container %q: %w", name, err)
	}

	if notes == nil {
		notes = []models.Note{}
	}
	return notes, nil
}

// GetByID returns the note with the given ID, or ErrNotFound.
func (s *NoteStore) GetByID(ctx context.Context, id int64) (*models.Note, error) {
	row := s.db.QueryRowContext(ctx,
		"SELECT "+columnList+" FROM notes WHERE id = ?",
		id)

	n, err := scanNoteRow(row)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("get note %d: %w", id, err)
	}
	return &n, nil
}

// Create inserts a new note and returns it with generated fields populated.
func (s *NoteStore) Create(ctx context.Context, req models.CreateNoteRequest) (*models.Note, error) {
	tags := req.Tags
	if tags == nil {
		tags = []string{}
	}
	tagsJSON, err := json.Marshal(tags)
	if err != nil {
		return nil, fmt.Errorf("marshal tags: %w", err)
	}

	now := time.Now().UTC().Format(time.RFC3339)
	result, err := s.db.ExecContext(ctx,
		"INSERT INTO notes (container_name, container_id, compose_project, compose_service, title, note_content, tags, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)",
		req.ContainerName, req.ContainerID, req.ComposeProject, req.ComposeService,
		req.Title, req.NoteContent, string(tagsJSON), now, now)
	if err != nil {
		return nil, fmt.Errorf("insert note: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("last insert id: %w", err)
	}

	return &models.Note{
		ID:             id,
		ContainerName:  req.ContainerName,
		ContainerID:    req.ContainerID,
		ComposeProject: req.ComposeProject,
		ComposeService: req.ComposeService,
		Title:          req.Title,
		NoteContent:    req.NoteContent,
		Pinned:         false,
		Tags:           tags,
		CreatedAt:      parseTime(now),
		UpdatedAt:      parseTime(now),
	}, nil
}

// Update applies partial updates to the note identified by ID.
// Only non-nil fields in the request are changed.
func (s *NoteStore) Update(ctx context.Context, id int64, req models.UpdateNoteRequest) (*models.Note, error) {
	var setClauses []string
	var args []any

	if req.Title != nil {
		setClauses = append(setClauses, "title = ?")
		args = append(args, *req.Title)
	}
	if req.NoteContent != nil {
		setClauses = append(setClauses, "note_content = ?")
		args = append(args, *req.NoteContent)
	}
	if req.Pinned != nil {
		setClauses = append(setClauses, "pinned = ?")
		if *req.Pinned {
			args = append(args, 1)
		} else {
			args = append(args, 0)
		}
	}
	if req.Tags != nil {
		tagsJSON, err := json.Marshal(*req.Tags)
		if err != nil {
			return nil, fmt.Errorf("marshal tags: %w", err)
		}
		setClauses = append(setClauses, "tags = ?")
		args = append(args, string(tagsJSON))
	}
	if req.ContainerID != "" {
		setClauses = append(setClauses, "container_id = ?")
		args = append(args, req.ContainerID)
	}

	now := time.Now().UTC().Format(time.RFC3339)
	setClauses = append(setClauses, "updated_at = ?")
	args = append(args, now)

	args = append(args, id)
	query := "UPDATE notes SET " + strings.Join(setClauses, ", ") + " WHERE id = ?"

	result, err := s.db.ExecContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("update note %d: %w", id, err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return nil, fmt.Errorf("rows affected: %w", err)
	}
	if rowsAffected == 0 {
		return nil, ErrNotFound
	}

	return s.GetByID(ctx, id)
}

// Delete removes the note identified by ID.
func (s *NoteStore) Delete(ctx context.Context, id int64) error {
	result, err := s.db.ExecContext(ctx,
		"DELETE FROM notes WHERE id = ?", id)
	if err != nil {
		return fmt.Errorf("delete note %d: %w", id, err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("rows affected: %w", err)
	}
	if rowsAffected == 0 {
		return ErrNotFound
	}
	return nil
}

// DeleteByContainer removes all notes for the given container name.
// Returns the number of notes deleted.
func (s *NoteStore) DeleteByContainer(ctx context.Context, name string) (int64, error) {
	result, err := s.db.ExecContext(ctx,
		"DELETE FROM notes WHERE container_name = ?", name)
	if err != nil {
		return 0, fmt.Errorf("delete notes for container %q: %w", name, err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("rows affected: %w", err)
	}
	return rowsAffected, nil
}

// scanNote scans a note from a sql.Rows iterator.
func scanNote(rows *sql.Rows) (models.Note, error) {
	var n models.Note
	var pinnedInt int
	var tagsJSON, createdStr, updatedStr string

	err := rows.Scan(&n.ID, &n.ContainerName, &n.ContainerID,
		&n.ComposeProject, &n.ComposeService, &n.Title, &n.NoteContent,
		&pinnedInt, &tagsJSON, &createdStr, &updatedStr)
	if err != nil {
		return n, fmt.Errorf("scan note: %w", err)
	}

	n.Pinned = pinnedInt == 1
	n.CreatedAt = parseTime(createdStr)
	n.UpdatedAt = parseTime(updatedStr)

	if err := json.Unmarshal([]byte(tagsJSON), &n.Tags); err != nil {
		n.Tags = []string{}
	}
	if n.Tags == nil {
		n.Tags = []string{}
	}

	return n, nil
}

// scanNoteRow scans a note from a sql.Row.
func scanNoteRow(row *sql.Row) (models.Note, error) {
	var n models.Note
	var pinnedInt int
	var tagsJSON, createdStr, updatedStr string

	err := row.Scan(&n.ID, &n.ContainerName, &n.ContainerID,
		&n.ComposeProject, &n.ComposeService, &n.Title, &n.NoteContent,
		&pinnedInt, &tagsJSON, &createdStr, &updatedStr)
	if err != nil {
		return n, err
	}

	n.Pinned = pinnedInt == 1
	n.CreatedAt = parseTime(createdStr)
	n.UpdatedAt = parseTime(updatedStr)

	if err := json.Unmarshal([]byte(tagsJSON), &n.Tags); err != nil {
		n.Tags = []string{}
	}
	if n.Tags == nil {
		n.Tags = []string{}
	}

	return n, nil
}

// ExportAll returns all notes for export.
func (s *NoteStore) ExportAll(ctx context.Context) ([]models.Note, error) {
	return s.List(ctx, nil, "")
}

// ImportAll upserts notes from an import payload inside a transaction.
// Notes with an empty ContainerName are skipped. Returns the count of
// successfully imported (created or updated) notes.
func (s *NoteStore) ImportAll(ctx context.Context, notes []models.Note) (imported int, err error) {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return 0, fmt.Errorf("begin tx: %w", err)
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback()
		}
	}()

	for i := range notes {
		n := &notes[i]
		if n.ContainerName == "" {
			continue
		}

		// Default title for imports missing one (backward compat).
		if n.Title == "" {
			n.Title = "Note"
		}

		tags := n.Tags
		if tags == nil {
			tags = []string{}
		}
		tagsJSON, marshalErr := json.Marshal(tags)
		if marshalErr != nil {
			err = fmt.Errorf("marshal tags for %q: %w", n.ContainerName, marshalErr)
			return 0, err
		}

		now := time.Now().UTC().Format(time.RFC3339)
		pinnedVal := 0
		if n.Pinned {
			pinnedVal = 1
		}

		// Check if a note with the same container_name and title already exists.
		var existingID int64
		row := tx.QueryRowContext(ctx,
			"SELECT id FROM notes WHERE container_name = ? AND title = ?", n.ContainerName, n.Title)
		scanErr := row.Scan(&existingID)

		if scanErr == nil {
			// Exists: update.
			_, execErr := tx.ExecContext(ctx,
				"UPDATE notes SET container_id = ?, compose_project = ?, compose_service = ?, note_content = ?, pinned = ?, tags = ?, updated_at = ? WHERE id = ?",
				n.ContainerID, n.ComposeProject, n.ComposeService,
				n.NoteContent, pinnedVal, string(tagsJSON), now, existingID)
			if execErr != nil {
				err = fmt.Errorf("update note %q/%q: %w", n.ContainerName, n.Title, execErr)
				return 0, err
			}
		} else if errors.Is(scanErr, sql.ErrNoRows) {
			// Does not exist: insert.
			_, execErr := tx.ExecContext(ctx,
				"INSERT INTO notes (container_name, container_id, compose_project, compose_service, title, note_content, pinned, tags, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)",
				n.ContainerName, n.ContainerID, n.ComposeProject, n.ComposeService,
				n.Title, n.NoteContent, pinnedVal, string(tagsJSON), now, now)
			if execErr != nil {
				err = fmt.Errorf("insert note %q/%q: %w", n.ContainerName, n.Title, execErr)
				return 0, err
			}
		} else {
			err = fmt.Errorf("check existing note %q/%q: %w", n.ContainerName, n.Title, scanErr)
			return 0, err
		}
		imported++
	}

	if commitErr := tx.Commit(); commitErr != nil {
		err = fmt.Errorf("commit tx: %w", commitErr)
		return 0, err
	}
	return imported, nil
}

func parseTime(s string) time.Time {
	t, _ := time.Parse(time.RFC3339, s)
	return t
}
