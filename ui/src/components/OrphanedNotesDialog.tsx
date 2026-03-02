import { useState } from "react";
import Box from "@mui/material/Box";
import Button from "@mui/material/Button";
import Dialog from "@mui/material/Dialog";
import DialogActions from "@mui/material/DialogActions";
import DialogContent from "@mui/material/DialogContent";
import DialogContentText from "@mui/material/DialogContentText";
import DialogTitle from "@mui/material/DialogTitle";
import IconButton from "@mui/material/IconButton";
import List from "@mui/material/List";
import ListItem from "@mui/material/ListItem";
import ListItemText from "@mui/material/ListItemText";
import Tooltip from "@mui/material/Tooltip";
import Typography from "@mui/material/Typography";
import DeleteIcon from "@mui/icons-material/Delete";
import type { Note } from "../types";

interface OrphanedNotesDialogProps {
  open: boolean;
  onClose: () => void;
  orphanedNotes: Note[];
  onDelete: (name: string) => Promise<void>;
  onDeleteAll: () => Promise<void>;
}

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

export function OrphanedNotesDialog({
  open,
  onClose,
  orphanedNotes,
  onDelete,
  onDeleteAll,
}: OrphanedNotesDialogProps) {
  const [confirmOpen, setConfirmOpen] = useState(false);

  const handleDeleteAll = async () => {
    setConfirmOpen(false);
    await onDeleteAll();
  };

  return (
    <>
      <Dialog open={open} onClose={onClose} maxWidth="sm" fullWidth>
        <DialogTitle>Orphaned Notes</DialogTitle>
        <DialogContent>
          <Typography variant="body2" color="text.secondary" sx={{ mb: 2 }}>
            These notes belong to containers that no longer exist.
          </Typography>
          {orphanedNotes.length === 0 ? (
            <Box sx={{ py: 2, textAlign: "center" }}>
              <Typography variant="body2" color="text.secondary">
                No orphaned notes found
              </Typography>
            </Box>
          ) : (
            <List disablePadding>
              {orphanedNotes.map((note) => (
                <ListItem
                  key={note.container_name}
                  secondaryAction={
                    <Tooltip title="Delete note">
                      <IconButton
                        edge="end"
                        aria-label={`Delete note for ${note.container_name}`}
                        color="error"
                        onClick={() => onDelete(note.container_name)}
                      >
                        <DeleteIcon />
                      </IconButton>
                    </Tooltip>
                  }
                  sx={{ borderBottom: 1, borderColor: "divider" }}
                >
                  <ListItemText
                    primary={note.container_name}
                    secondary={
                      <>
                        {note.note_content
                          ? truncate(note.note_content, 100)
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
                </ListItem>
              ))}
            </List>
          )}
        </DialogContent>
        <DialogActions>
          <Button onClick={onClose}>Close</Button>
          {orphanedNotes.length > 0 && (
            <Button
              color="error"
              variant="contained"
              onClick={() => setConfirmOpen(true)}
            >
              Delete All
            </Button>
          )}
        </DialogActions>
      </Dialog>

      <Dialog open={confirmOpen} onClose={() => setConfirmOpen(false)}>
        <DialogTitle>Confirm Deletion</DialogTitle>
        <DialogContent>
          <DialogContentText>
            Delete {orphanedNotes.length} orphaned note
            {orphanedNotes.length !== 1 ? "s" : ""}? This cannot be undone.
          </DialogContentText>
        </DialogContent>
        <DialogActions>
          <Button onClick={() => setConfirmOpen(false)}>Cancel</Button>
          <Button color="error" variant="contained" onClick={handleDeleteAll}>
            Delete
          </Button>
        </DialogActions>
      </Dialog>
    </>
  );
}
