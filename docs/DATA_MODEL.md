# Data Model

Container note schema for RunNotes. See [ARCHITECTURE.md](ARCHITECTURE.md) for the full system design.

## Design Principles

- **Name-based persistence** — notes are keyed by container name, not container ID, so they survive `docker-compose down && up` cycles
- **Dual-key lookup** — container name is the stable lookup key; container ID is refreshed metadata
- **MVP simplicity** — single table, JSON tags, no premature normalization
- **No NULLs** — default empty strings avoid `sql.NullString` complexity in Go

## SQLite Schema

```sql
CREATE TABLE IF NOT EXISTS notes (
    id              INTEGER PRIMARY KEY AUTOINCREMENT,
    container_name  TEXT    NOT NULL,
    container_id    TEXT    NOT NULL DEFAULT '',
    compose_project TEXT    NOT NULL DEFAULT '',
    compose_service TEXT    NOT NULL DEFAULT '',
    note_content    TEXT    NOT NULL DEFAULT '',
    pinned          INTEGER NOT NULL DEFAULT 0,
    tags            TEXT    NOT NULL DEFAULT '[]',
    created_at      TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%SZ', 'now')),
    updated_at      TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%SZ', 'now'))
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_notes_container_name ON notes(container_name);
```

## Field Reference

| Field | Type | Constraints | Purpose |
|-------|------|-------------|---------|
| `id` | INTEGER | PK, autoincrement | Stable internal identifier |
| `container_name` | TEXT | NOT NULL, unique index | Primary lookup key (e.g., `my-postgres`) |
| `container_id` | TEXT | NOT NULL, default `''` | Last known Docker container ID |
| `compose_project` | TEXT | NOT NULL, default `''` | Compose project name if applicable |
| `compose_service` | TEXT | NOT NULL, default `''` | Compose service name if applicable |
| `note_content` | TEXT | NOT NULL, default `''` | User-authored note text |
| `pinned` | INTEGER | NOT NULL, default `0` | 0 = normal, 1 = pinned to top |
| `tags` | TEXT | NOT NULL, default `'[]'` | JSON array of tag strings |
| `created_at` | TEXT | NOT NULL, default now | RFC3339 timestamp |
| `updated_at` | TEXT | NOT NULL, default now | RFC3339 timestamp |

### Why integer PK instead of container_name as PK

An auto-increment integer PK provides:

- Stable references if a container is renamed (the note ID doesn't change)
- Simpler foreign keys for future tables (history, attachments)
- SQLite rowid optimization (integer PK aliases the rowid, zero overhead)

The unique index on `container_name` enforces one note per container name.

### Why tags as JSON array instead of a separate table

A normalized `note_tags` many-to-many table would be standard relational design, but it's over-engineered for this MVP:

- Most notes will have 0-3 tags
- SQLite's `json_each()` function supports querying tags when needed
- Migration to a separate table is straightforward later
- One table means simpler CRUD, simpler backup, simpler export

## Dual-Key Lookup

The dual-key system handles the container recreation scenario:

1. User creates note for container `my-postgres` (ID: `abc123`)
2. User runs `docker-compose down && docker-compose up`
3. Container `my-postgres` is recreated with new ID `def456`
4. Frontend fetches container list from Docker API — sees `my-postgres` with ID `def456`
5. Frontend requests note for `my-postgres` from backend
6. Backend finds the note by `container_name` (the stable key)
7. Backend updates `container_id` to `def456` (refresh metadata)
8. Note content is preserved — the user sees their note unchanged

The `container_id` field serves two purposes:

- **Disambiguation** — if two containers share a name across different Compose projects (future enhancement)
- **Audit trail** — tracking which container instance the note was last associated with

## Go Backend Types

These structs live in `internal/models/note.go`.

```go
package models

import "time"

// Note represents a user-authored annotation attached to a Docker container.
type Note struct {
    ID             int64     `json:"id"`
    ContainerName  string    `json:"container_name"`
    ContainerID    string    `json:"container_id"`
    ComposeProject string    `json:"compose_project,omitempty"`
    ComposeService string    `json:"compose_service,omitempty"`
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
    NoteContent    string   `json:"note_content"`
    Tags           []string `json:"tags,omitempty"`
}

// UpdateNoteRequest is the payload for updating an existing note.
// Pointer fields distinguish "not provided" from "set to empty".
type UpdateNoteRequest struct {
    NoteContent *string   `json:"note_content,omitempty"`
    Pinned      *bool     `json:"pinned,omitempty"`
    Tags        *[]string `json:"tags,omitempty"`
    ContainerID string    `json:"container_id,omitempty"`
}
```

### Notes on Go types

- `Tags` is `[]string` in Go, serialized to/from JSON TEXT in SQLite via `json.Marshal`/`json.Unmarshal`
- `UpdateNoteRequest` uses pointers for optional fields — PATCH semantics where `nil` means "don't change" and zero value means "set to empty"
- `Pinned` is `bool` in Go, mapped to INTEGER 0/1 in SQLite
- `omitempty` on `ComposeProject`/`ComposeService` keeps JSON clean for standalone containers

## TypeScript Frontend Types

These interfaces will live in `ui/src/types.ts` when the frontend is scaffolded (issue #2).

```typescript
/** A user-authored annotation attached to a Docker container. */
export interface Note {
  id: number;
  container_name: string;
  container_id: string;
  compose_project: string;
  compose_service: string;
  note_content: string;
  pinned: boolean;
  tags: string[];
  created_at: string;
  updated_at: string;
}

/** Payload for creating a new note. */
export interface CreateNoteRequest {
  container_name: string;
  container_id: string;
  compose_project?: string;
  compose_service?: string;
  note_content: string;
  tags?: string[];
}

/** Payload for updating an existing note. */
export interface UpdateNoteRequest {
  note_content?: string;
  pinned?: boolean;
  tags?: string[];
  container_id?: string;
}

/** Container info from the Docker API, used for the container list view. */
export interface ContainerInfo {
  id: string;
  name: string;
  image: string;
  state: string;
  status: string;
  compose_project?: string;
  compose_service?: string;
}
```

### Notes on TypeScript types

- Timestamps are `string` (ISO 8601), not `Date` — JSON serialization sends strings; parse to `Date` in the UI layer when needed for display
- `ContainerInfo` represents Docker API data, not stored data — used for the container list view and correlating containers with notes
- Optional fields use TypeScript's `?` syntax, matching Go's `omitempty` behavior

## Future Considerations

These are explicitly deferred beyond MVP. Documented here so the schema can be evaluated against future needs:

- **Note versioning/history** — a `note_history` table with FK to `notes.id`, storing previous content and timestamps
- **Attachments** — an `attachments` table for files linked to notes (screenshots, configs)
- **Shared/team notes** — multi-user support would require user identity and access control
- **Compose-project-level notes** — notes attached to a Compose project rather than individual containers
- **Note templates** — predefined note structures for common container types (databases, web servers)
