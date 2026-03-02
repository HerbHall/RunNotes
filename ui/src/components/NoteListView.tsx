import { useState, useCallback } from "react";
import Box from "@mui/material/Box";
import Button from "@mui/material/Button";
import IconButton from "@mui/material/IconButton";
import List from "@mui/material/List";
import ListItemButton from "@mui/material/ListItemButton";
import ListItemText from "@mui/material/ListItemText";
import Stack from "@mui/material/Stack";
import TextField from "@mui/material/TextField";
import Tooltip from "@mui/material/Tooltip";
import Typography from "@mui/material/Typography";
import AddIcon from "@mui/icons-material/Add";
import PushPinIcon from "@mui/icons-material/PushPin";
import type { Note, ContainerInfo, CreateNoteRequest } from "../types";

function truncate(text: string, maxLength: number): string {
  if (text.length <= maxLength) return text;
  return text.slice(0, maxLength) + "...";
}

function formatDate(dateStr: string): string {
  try {
    return new Date(dateStr).toLocaleDateString();
  } catch {
    return dateStr;
  }
}

interface NoteListViewProps {
  containerName: string;
  container: ContainerInfo;
  notes: Note[];
  onSelectNote: (id: number) => void;
  onCreateNote: (req: CreateNoteRequest) => Promise<void>;
}

export function NoteListView({
  containerName,
  container,
  notes,
  onSelectNote,
  onCreateNote,
}: NoteListViewProps) {
  const [creating, setCreating] = useState(false);
  const [newTitle, setNewTitle] = useState("");
  const [saving, setSaving] = useState(false);

  const handleCreate = useCallback(async () => {
    const title = newTitle.trim();
    if (!title) return;
    setSaving(true);
    try {
      await onCreateNote({
        container_name: containerName,
        container_id: container.id,
        compose_project: container.composeProject,
        compose_service: container.composeService,
        title,
        note_content: "",
      });
      setNewTitle("");
      setCreating(false);
    } finally {
      setSaving(false);
    }
  }, [newTitle, containerName, container, onCreateNote]);

  return (
    <Box sx={{ p: 2 }}>
      <Stack direction="row" alignItems="center" spacing={1} sx={{ mb: 2 }}>
        <Typography variant="h6" sx={{ flex: 1 }}>
          {containerName}
        </Typography>
        <Button
          variant="contained"
          size="small"
          startIcon={<AddIcon />}
          onClick={() => setCreating(true)}
          disabled={creating}
        >
          Add Note
        </Button>
      </Stack>

      {creating && (
        <Stack direction="row" spacing={1} sx={{ mb: 2 }}>
          <TextField
            size="small"
            placeholder="Note title..."
            value={newTitle}
            onChange={(e) => setNewTitle(e.target.value)}
            onKeyDown={(e) => {
              if (e.key === "Enter") {
                e.preventDefault();
                handleCreate();
              }
              if (e.key === "Escape") {
                setCreating(false);
                setNewTitle("");
              }
            }}
            autoFocus
            sx={{ flex: 1 }}
          />
          <Button
            variant="contained"
            size="small"
            onClick={handleCreate}
            disabled={saving || !newTitle.trim()}
          >
            {saving ? "Creating..." : "Create"}
          </Button>
          <Button
            size="small"
            onClick={() => {
              setCreating(false);
              setNewTitle("");
            }}
          >
            Cancel
          </Button>
        </Stack>
      )}

      {notes.length === 0 ? (
        <Box sx={{ py: 4, textAlign: "center" }}>
          <Typography variant="body1" color="text.secondary" gutterBottom>
            No notes yet for this container
          </Typography>
          {!creating && (
            <Button
              variant="outlined"
              startIcon={<AddIcon />}
              onClick={() => setCreating(true)}
            >
              Add Note
            </Button>
          )}
        </Box>
      ) : (
        <List disablePadding>
          {notes.map((note) => (
            <ListItemButton
              key={note.id}
              onClick={() => onSelectNote(note.id)}
              sx={{ borderBottom: 1, borderColor: "divider" }}
            >
              <ListItemText
                primary={note.title}
                secondary={
                  <>
                    {note.note_content
                      ? truncate(note.note_content, 80)
                      : "Empty note"}
                    {" — "}
                    {formatDate(note.updated_at)}
                  </>
                }
                primaryTypographyProps={{ variant: "subtitle2" }}
                secondaryTypographyProps={{
                  variant: "body2",
                  color: "text.secondary",
                }}
              />
              {note.pinned && (
                <Tooltip title="Pinned">
                  <IconButton size="small" sx={{ ml: 1 }} disabled>
                    <PushPinIcon fontSize="small" color="primary" />
                  </IconButton>
                </Tooltip>
              )}
            </ListItemButton>
          ))}
        </List>
      )}
    </Box>
  );
}
