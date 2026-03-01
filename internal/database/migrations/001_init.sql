-- +goose Up
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

-- +goose Down
DROP TABLE IF EXISTS notes;
