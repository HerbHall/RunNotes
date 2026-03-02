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
import ToggleButton from "@mui/material/ToggleButton";
import ToggleButtonGroup from "@mui/material/ToggleButtonGroup";
import Tooltip from "@mui/material/Tooltip";
import Typography from "@mui/material/Typography";
import ArrowBackIcon from "@mui/icons-material/ArrowBack";
import DeleteIcon from "@mui/icons-material/Delete";
import EditIcon from "@mui/icons-material/Edit";
import PushPinIcon from "@mui/icons-material/PushPin";
import PushPinOutlinedIcon from "@mui/icons-material/PushPinOutlined";
import VisibilityIcon from "@mui/icons-material/Visibility";
import type {
  Note,
  ContainerInfo,
  UpdateNoteRequest,
} from "../types";
import { MarkdownPreview } from "./MarkdownPreview";

interface NoteEditorProps {
  note: Note;
  container: ContainerInfo;
  onSave: (id: number, req: UpdateNoteRequest) => Promise<void>;
  onDelete: (id: number) => Promise<void>;
  onBack: () => void;
}

export function NoteEditor({
  note,
  container,
  onSave,
  onDelete,
  onBack,
}: NoteEditorProps) {
  const [title, setTitle] = useState(note.title);
  const [content, setContent] = useState(note.note_content);
  const [pinned, setPinned] = useState(note.pinned);
  const [tags, setTags] = useState<string[]>(note.tags);
  const [tagInput, setTagInput] = useState("");
  const [confirmOpen, setConfirmOpen] = useState(false);
  const [saving, setSaving] = useState(false);
  const [preview, setPreview] = useState(false);

  useEffect(() => {
    setTitle(note.title);
    setContent(note.note_content);
    setPinned(note.pinned);
    setTags(note.tags);
    setTagInput("");
    setPreview(false);
  }, [note]);

  const handleSave = useCallback(async () => {
    setSaving(true);
    try {
      const req: UpdateNoteRequest = {
        title,
        note_content: content,
        pinned,
        tags,
        container_id: container.id,
      };
      await onSave(note.id, req);
    } finally {
      setSaving(false);
    }
  }, [note.id, title, content, pinned, tags, container, onSave]);

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
    await onDelete(note.id);
    onBack();
  }, [onDelete, note.id, onBack]);

  const hasChanges =
    title !== note.title ||
    content !== note.note_content ||
    pinned !== note.pinned ||
    JSON.stringify(tags) !== JSON.stringify(note.tags);

  return (
    <Box sx={{ p: 2, height: "100%", display: "flex", flexDirection: "column" }} onKeyDown={handleKeyDown}>
      <Stack direction="row" alignItems="center" spacing={1} sx={{ mb: 2 }}>
        <Tooltip title="Back to note list">
          <IconButton onClick={onBack} size="small">
            <ArrowBackIcon />
          </IconButton>
        </Tooltip>
        <TextField
          variant="standard"
          value={title}
          onChange={(e) => setTitle(e.target.value)}
          placeholder="Note title"
          InputProps={{
            style: { fontSize: "1.25rem", fontWeight: 500 },
          }}
          sx={{ flex: 1 }}
        />
        <Typography variant="caption" color="text.secondary" sx={{ flexShrink: 0 }}>
          {note.container_name}
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

      <ToggleButtonGroup
        value={preview ? "preview" : "edit"}
        exclusive
        onChange={(_, val) => { if (val) setPreview(val === "preview"); }}
        size="small"
        sx={{ mb: 1 }}
      >
        <ToggleButton value="edit">
          <EditIcon sx={{ mr: 0.5 }} fontSize="small" />
          Edit
        </ToggleButton>
        <ToggleButton value="preview">
          <VisibilityIcon sx={{ mr: 0.5 }} fontSize="small" />
          Preview
        </ToggleButton>
      </ToggleButtonGroup>

      {preview ? (
        <MarkdownPreview content={content} />
      ) : (
        <TextField
          multiline
          minRows={4}
          maxRows={20}
          fullWidth
          placeholder="Write your note here... (supports Markdown)"
          value={content}
          onChange={(e) => setContent(e.target.value)}
          sx={{ mb: 2 }}
        />
      )}

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
            setTitle(note.title);
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
            Are you sure you want to delete &quot;{note.title}&quot;? This
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
