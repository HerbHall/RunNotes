import { useState, useEffect, useCallback } from "react";
import Box from "@mui/material/Box";
import Button from "@mui/material/Button";
import Chip from "@mui/material/Chip";
import Dialog from "@mui/material/Dialog";
import DialogActions from "@mui/material/DialogActions";
import DialogContent from "@mui/material/DialogContent";
import DialogContentText from "@mui/material/DialogContentText";
import DialogTitle from "@mui/material/DialogTitle";
import IconButton from "@mui/material/IconButton";
import Stack from "@mui/material/Stack";
import TextField from "@mui/material/TextField";
import Tooltip from "@mui/material/Tooltip";
import Typography from "@mui/material/Typography";
import PushPinIcon from "@mui/icons-material/PushPin";
import PushPinOutlinedIcon from "@mui/icons-material/PushPinOutlined";
import DeleteIcon from "@mui/icons-material/Delete";
import type {
  Note,
  ContainerInfo,
  CreateNoteRequest,
  UpdateNoteRequest,
} from "../types";

interface NoteEditorProps {
  containerName: string;
  container: ContainerInfo;
  note: Note | null;
  onSave: (
    name: string,
    req: CreateNoteRequest | UpdateNoteRequest,
  ) => Promise<void>;
  onDelete: (name: string) => Promise<void>;
}

export function NoteEditor({
  containerName,
  container,
  note,
  onSave,
  onDelete,
}: NoteEditorProps) {
  const [content, setContent] = useState(note?.note_content ?? "");
  const [pinned, setPinned] = useState(note?.pinned ?? false);
  const [tags, setTags] = useState<string[]>(note?.tags ?? []);
  const [tagInput, setTagInput] = useState("");
  const [confirmOpen, setConfirmOpen] = useState(false);
  const [saving, setSaving] = useState(false);

  useEffect(() => {
    setContent(note?.note_content ?? "");
    setPinned(note?.pinned ?? false);
    setTags(note?.tags ?? []);
    setTagInput("");
  }, [note, containerName]);

  const handleSave = useCallback(async () => {
    setSaving(true);
    try {
      if (note) {
        const req: UpdateNoteRequest = {
          note_content: content,
          pinned,
          tags,
          container_id: container.id,
        };
        await onSave(containerName, req);
      } else {
        const req: CreateNoteRequest = {
          container_name: containerName,
          container_id: container.id,
          compose_project: container.composeProject,
          compose_service: container.composeService,
          note_content: content,
          tags,
        };
        await onSave(containerName, req);
      }
    } finally {
      setSaving(false);
    }
  }, [note, content, pinned, tags, containerName, container, onSave]);

  const handleKeyDown = useCallback(
    (e: React.KeyboardEvent) => {
      if ((e.ctrlKey || e.metaKey) && e.key === "s") {
        e.preventDefault();
        handleSave();
      }
    },
    [handleSave],
  );

  const handleAddTag = useCallback(() => {
    const trimmed = tagInput.trim();
    if (trimmed && !tags.includes(trimmed)) {
      setTags((prev) => [...prev, trimmed]);
    }
    setTagInput("");
  }, [tagInput, tags]);

  const handleRemoveTag = useCallback((tag: string) => {
    setTags((prev) => prev.filter((t) => t !== tag));
  }, []);

  const handleDelete = useCallback(async () => {
    setConfirmOpen(false);
    await onDelete(containerName);
  }, [onDelete, containerName]);

  const hasChanges =
    note != null &&
    (content !== note.note_content ||
      pinned !== note.pinned ||
      JSON.stringify(tags) !== JSON.stringify(note.tags));

  if (note == null) {
    return (
      <Box sx={{ p: 3, textAlign: "center" }}>
        <Typography variant="body1" color="text.secondary" gutterBottom>
          No note yet for {containerName}
        </Typography>
        <Button variant="contained" onClick={handleSave} disabled={saving}>
          Add Note
        </Button>
      </Box>
    );
  }

  return (
    <Box sx={{ p: 2, height: "100%", display: "flex", flexDirection: "column" }} onKeyDown={handleKeyDown}>
      <Stack direction="row" alignItems="center" spacing={1} sx={{ mb: 2 }}>
        <Typography variant="h6" sx={{ flex: 1 }}>
          {containerName}
        </Typography>
        <Tooltip title={pinned ? "Unpin note" : "Pin note"}>
          <IconButton onClick={() => setPinned(!pinned)} size="small">
            {pinned ? (
              <PushPinIcon color="primary" />
            ) : (
              <PushPinOutlinedIcon />
            )}
          </IconButton>
        </Tooltip>
        <Tooltip title="Delete note">
          <IconButton
            onClick={() => setConfirmOpen(true)}
            size="small"
            color="error"
          >
            <DeleteIcon />
          </IconButton>
        </Tooltip>
      </Stack>

      <TextField
        multiline
        minRows={4}
        maxRows={20}
        fullWidth
        placeholder="Write your note here..."
        value={content}
        onChange={(e) => setContent(e.target.value)}
        sx={{ mb: 2 }}
      />

      <Stack direction="row" spacing={1} flexWrap="wrap" sx={{ mb: 1 }}>
        {tags.map((tag) => (
          <Chip
            key={tag}
            label={tag}
            size="small"
            onDelete={() => handleRemoveTag(tag)}
            sx={{ mb: 0.5 }}
          />
        ))}
      </Stack>

      <Stack direction="row" spacing={1} sx={{ mb: 2 }}>
        <TextField
          size="small"
          placeholder="Add tag..."
          value={tagInput}
          onChange={(e) => setTagInput(e.target.value)}
          onKeyDown={(e) => {
            if (e.key === "Enter") {
              e.preventDefault();
              handleAddTag();
            }
          }}
          sx={{ width: 160 }}
        />
        <Button size="small" variant="outlined" onClick={handleAddTag}>
          Add
        </Button>
      </Stack>

      <Stack direction="row" spacing={1}>
        <Button
          variant="contained"
          onClick={handleSave}
          disabled={saving || !hasChanges}
        >
          {saving ? "Saving..." : "Save"}
        </Button>
        <Button
          variant="outlined"
          disabled={!hasChanges}
          onClick={() => {
            setContent(note.note_content);
            setPinned(note.pinned);
            setTags(note.tags);
          }}
        >
          Cancel
        </Button>
      </Stack>

      <Dialog open={confirmOpen} onClose={() => setConfirmOpen(false)}>
        <DialogTitle>Delete Note</DialogTitle>
        <DialogContent>
          <DialogContentText>
            Are you sure you want to delete the note for {containerName}? This
            action cannot be undone.
          </DialogContentText>
        </DialogContent>
        <DialogActions>
          <Button onClick={() => setConfirmOpen(false)}>Cancel</Button>
          <Button onClick={handleDelete} color="error" variant="contained">
            Delete
          </Button>
        </DialogActions>
      </Dialog>
    </Box>
  );
}
