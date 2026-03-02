package models

import (
	"errors"
	"time"
)

// Note represents a user-authored annotation attached to a Docker container.
type Note struct {
	ID             int64     `json:"id"`
	ContainerName  string    `json:"container_name"`
	ContainerID    string    `json:"container_id"`
	ComposeProject string    `json:"compose_project,omitempty"`
	ComposeService string    `json:"compose_service,omitempty"`
	Title          string    `json:"title"`
	NoteContent    string    `json:"note_content"`
	Pinned         bool      `json:"pinned"`
	Tags           []string  `json:"tags"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

// CreateNoteRequest is the payload for creating a new note.
type CreateNoteRequest struct {
	ContainerName  string   `json:"container_name"`
	ContainerID    string   `json:"container_id"`
	ComposeProject string   `json:"compose_project,omitempty"`
	ComposeService string   `json:"compose_service,omitempty"`
	Title          string   `json:"title"`
	NoteContent    string   `json:"note_content"`
	Tags           []string `json:"tags,omitempty"`
}

// UpdateNoteRequest is the payload for updating an existing note.
// Pointer fields distinguish "not provided" from "set to empty".
type UpdateNoteRequest struct {
	Title       *string   `json:"title,omitempty"`
	NoteContent *string   `json:"note_content,omitempty"`
	Pinned      *bool     `json:"pinned,omitempty"`
	Tags        *[]string `json:"tags,omitempty"`
	ContainerID string    `json:"container_id,omitempty"`
}

// Validate checks that required fields are present on CreateNoteRequest.
func (r *CreateNoteRequest) Validate() error {
	if r.ContainerName == "" {
		return errors.New("container_name is required")
	}
	if r.Title == "" {
		return errors.New("title is required")
	}
	return nil
}

// Validate checks that at least one field is set on UpdateNoteRequest.
func (r *UpdateNoteRequest) Validate() error {
	if r.Title == nil && r.NoteContent == nil && r.Pinned == nil && r.Tags == nil && r.ContainerID == "" {
		return errors.New("at least one field must be provided")
	}
	return nil
}
