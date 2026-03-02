import { useRef, useState } from "react";
import Button from "@mui/material/Button";
import Dialog from "@mui/material/Dialog";
import DialogActions from "@mui/material/DialogActions";
import DialogContent from "@mui/material/DialogContent";
import DialogContentText from "@mui/material/DialogContentText";
import DialogTitle from "@mui/material/DialogTitle";
import DownloadIcon from "@mui/icons-material/Download";
import UploadIcon from "@mui/icons-material/Upload";
import { exportNotes, importNotes, showToast } from "../api/client";
import type { Note } from "../types";

interface ExportImportButtonsProps {
  onRefresh: () => Promise<void>;
}

export function ExportImportButtons({ onRefresh }: ExportImportButtonsProps) {
  const fileInputRef = useRef<HTMLInputElement>(null);
  const [confirmOpen, setConfirmOpen] = useState(false);
  const [pendingNotes, setPendingNotes] = useState<Note[]>([]);

  const handleExport = async () => {
    try {
      const notes = await exportNotes();
      const blob = new Blob([JSON.stringify(notes, null, 2)], {
        type: "application/json",
      });
      const url = URL.createObjectURL(blob);
      const a = document.createElement("a");
      a.href = url;
      a.download = "runnotes-export.json";
      document.body.appendChild(a);
      a.click();
      document.body.removeChild(a);
      URL.revokeObjectURL(url);
      showToast("success", `Exported ${notes.length} notes`);
    } catch {
      showToast("error", "Failed to export notes");
    }
  };

  const handleFileSelect = (event: React.ChangeEvent<HTMLInputElement>) => {
    const file = event.target.files?.[0];
    if (!file) return;

    const reader = new FileReader();
    reader.onload = (e) => {
      try {
        const parsed = JSON.parse(e.target?.result as string) as Note[];
        if (!Array.isArray(parsed)) {
          showToast("error", "Invalid file: expected a JSON array of notes");
          return;
        }
        setPendingNotes(parsed);
        setConfirmOpen(true);
      } catch {
        showToast("error", "Invalid file: could not parse JSON");
      }
    };
    reader.readAsText(file);

    // Reset the input so the same file can be selected again.
    event.target.value = "";
  };

  const handleConfirmImport = async () => {
    setConfirmOpen(false);
    try {
      const result = await importNotes(pendingNotes);
      showToast("success", `Imported ${result.imported} notes`);
      await onRefresh();
    } catch {
      showToast("error", "Failed to import notes");
    } finally {
      setPendingNotes([]);
    }
  };

  const handleCancelImport = () => {
    setConfirmOpen(false);
    setPendingNotes([]);
  };

  return (
    <>
      <Button
        size="small"
        startIcon={<DownloadIcon />}
        onClick={handleExport}
      >
        Export
      </Button>
      <Button
        size="small"
        startIcon={<UploadIcon />}
        onClick={() => fileInputRef.current?.click()}
      >
        Import
      </Button>
      <input
        ref={fileInputRef}
        type="file"
        accept=".json"
        style={{ display: "none" }}
        onChange={handleFileSelect}
      />
      <Dialog open={confirmOpen} onClose={handleCancelImport}>
        <DialogTitle>Import Notes</DialogTitle>
        <DialogContent>
          <DialogContentText>
            Import {pendingNotes.length} notes? Existing notes with matching
            container names and titles will be updated.
          </DialogContentText>
        </DialogContent>
        <DialogActions>
          <Button onClick={handleCancelImport}>Cancel</Button>
          <Button onClick={handleConfirmImport} variant="contained">
            Import
          </Button>
        </DialogActions>
      </Dialog>
    </>
  );
}
