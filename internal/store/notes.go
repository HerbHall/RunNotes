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

// List returns all notes, optionally filtered by pinned status or search term.
// Results are ordered by pinned DESC, then updated_at DESC.
func (s *NoteStore) List(ctx context.Context, pinned *bool, search string) ([]models.Note, error) {
	query := "SELECT id, container_name, container_id, compose_project, compose_service, note_content, pinned, tags, created_at, updated_at FROM notes"
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
		conditions = append(conditions, "(container_name LIKE ? OR note_content LIKE ?)")
		like := "%" + search + "%"
		args = append(args, like, like)
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

// GetByName returns the note for the given container name, or ErrNotFound.
func (s *NoteStore) GetByName(ctx context.Context, name string) (*models.Note, error) {
	row := s.db.QueryRowContext(ctx,
		"SELECT id, container_name, container_id, compose_project, compose_service, note_content, pinned, tags, created_at, updated_at FROM notes WHERE container_name = ?",
		name)

	n, err := scanNoteRow(row)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("get note %q: %w", name, err)
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
		"INSERT INTO notes (container_name, container_id, compose_project, compose_service, note_content, tags, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?)",
		req.ContainerName, req.ContainerID, req.ComposeProject, req.ComposeService,
		req.NoteContent, string(tagsJSON), now, now)
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
		NoteContent:    req.NoteContent,
		Pinned:         false,
		Tags:           tags,
		CreatedAt:      parseTime(now),
		UpdatedAt:      parseTime(now),
	}, nil
}

// Update applies partial updates to the note identified by container name.
// Only non-nil fields in the request are changed.
func (s *NoteStore) Update(ctx context.Context, name string, req models.UpdateNoteRequest) (*models.Note, error) {
	var setClauses []string
	var args []any

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

	args = append(args, name)
	query := "UPDATE notes SET " + strings.Join(setClauses, ", ") + " WHERE container_name = ?"

	result, err := s.db.ExecContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("update note %q: %w", name, err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return nil, fmt.Errorf("rows affected: %w", err)
	}
	if rowsAffected == 0 {
		return nil, ErrNotFound
	}

	return s.GetByName(ctx, name)
}

// Delete removes the note identified by container name.
func (s *NoteStore) Delete(ctx context.Context, name string) error {
	result, err := s.db.ExecContext(ctx,
		"DELETE FROM notes WHERE container_name = ?", name)
	if err != nil {
		return fmt.Errorf("delete note %q: %w", name, err)
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

// scanNote scans a note from a sql.Rows iterator.
func scanNote(rows *sql.Rows) (models.Note, error) {
	var n models.Note
	var pinnedInt int
	var tagsJSON, createdStr, updatedStr string

	err := rows.Scan(&n.ID, &n.ContainerName, &n.ContainerID,
		&n.ComposeProject, &n.ComposeService, &n.NoteContent,
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
		&n.ComposeProject, &n.ComposeService, &n.NoteContent,
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

func parseTime(s string) time.Time {
	t, _ := time.Parse(time.RFC3339, s)
	return t
}
