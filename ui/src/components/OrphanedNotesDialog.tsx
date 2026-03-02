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
  onDeleteContainer: (name: string) => Promise<void>;
  onDeleteAll: () => Promise<void>;
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
  onDeleteContainer,
  onDeleteAll,
}: OrphanedNotesDialogProps) {
  const [confirmOpen, setConfirmOpen] = useState(false);

  // Group orphaned notes by container name
  const grouped = new Map<string, Note[]>();
  for (const note of orphanedNotes) {
    const existing = grouped.get(note.container_name) ?? [];
    existing.push(note);
    grouped.set(note.container_name, existing);
  }

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
              {Array.from(grouped.entries()).map(([containerName, notes]) => (
                <ListItem
                  key={containerName}
                  secondaryAction={
                    <Tooltip title={`Delete all notes for ${containerName}`}>
                      <IconButton
                        edge="end"
                        aria-label={`Delete notes for ${containerName}`}
                        color="error"
                        onClick={() => onDeleteContainer(containerName)}
                      >
                        <DeleteIcon />
                      </IconButton>
                    </Tooltip>
                  }
                  sx={{ borderBottom: 1, borderColor: "divider" }}
                >
                  <ListItemText
                    primary={
                      notes.length > 1
                        ? `${containerName} (${notes.length} notes)`
                        : containerName
                    }
                    secondary={
                      <>
                        {notes.map((n) => n.title).join(", ")}
                        {" — "}
                        {formatDate(notes[0].updated_at)}
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
