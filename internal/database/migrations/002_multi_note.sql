-- +goose Up
-- Add title column to support multiple notes per container.
ALTER TABLE notes ADD COLUMN title TEXT NOT NULL DEFAULT '';

-- Populate titles from existing note content (first line, max 100 chars).
UPDATE notes SET title = CASE
    WHEN note_content = '' THEN 'Note'
    WHEN instr(note_content, char(10)) > 0
        THEN substr(note_content, 1, min(instr(note_content, char(10)) - 1, 100))
    ELSE substr(note_content, 1, 100)
END;

-- Drop the unique constraint to allow multiple notes per container.
DROP INDEX IF EXISTS idx_notes_container_name;

-- Create a non-unique index for efficient container-scoped queries.
CREATE INDEX IF NOT EXISTS idx_notes_container_name ON notes(container_name);

-- +goose Down
DROP INDEX IF EXISTS idx_notes_container_name;
CREATE UNIQUE INDEX IF NOT EXISTS idx_notes_container_name ON notes(container_name);
ALTER TABLE notes DROP COLUMN title;
