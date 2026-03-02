import { useState, useMemo, useEffect, useRef, useCallback } from "react";
import Badge from "@mui/material/Badge";
import Box from "@mui/material/Box";
import Divider from "@mui/material/Divider";
import IconButton from "@mui/material/IconButton";
import Tooltip from "@mui/material/Tooltip";
import Typography from "@mui/material/Typography";
import PushPinIcon from "@mui/icons-material/PushPin";
import PushPinOutlinedIcon from "@mui/icons-material/PushPinOutlined";
import RefreshIcon from "@mui/icons-material/Refresh";
import WarningAmberIcon from "@mui/icons-material/WarningAmber";
import { ContainerList } from "./components/ContainerList";
import { ExportImportButtons } from "./components/ExportImportButtons";
import { NoteEditor } from "./components/NoteEditor";
import { NoteListView } from "./components/NoteListView";
import { OrphanedNotesDialog } from "./components/OrphanedNotesDialog";
import { SearchBar } from "./components/SearchBar";
import { useContainers } from "./hooks/useContainers";
import { useNotes } from "./hooks/useNotes";
import type { CreateNoteRequest, UpdateNoteRequest } from "./types";

export default function App() {
  const {
    containers,
    loading: containersLoading,
    refresh: refreshContainers,
  } = useContainers();
  const {
    notes,
    loading: notesLoading,
    create,
    update,
    remove,
    removeAllForContainer,
    refresh: refreshNotes,
    getNotesForContainer,
  } = useNotes();

  const reconciledRef = useRef(false);

  // Auto-reconcile stale container IDs after data loads
  useEffect(() => {
    if (containersLoading || notesLoading || reconciledRef.current) return;
    const stale = notes.filter((note) => {
      const container = containers.find(
        (c) => c.name === note.container_name,
      );
      return container && container.id !== note.container_id;
    });
    if (stale.length > 0) {
      reconciledRef.current = true;
      Promise.all(
        stale.map((note) => {
          const container = containers.find(
            (c) => c.name === note.container_name,
          )!;
          return update(note.id, { container_id: container.id });
        }),
      );
    }
  }, [containers, notes, containersLoading, notesLoading, update]);

  const handleRefresh = useCallback(async () => {
    reconciledRef.current = false;
    await refreshContainers();
    await refreshNotes();
  }, [refreshContainers, refreshNotes]);

  const [selectedName, setSelectedName] = useState<string | null>(null);
  const [selectedNoteId, setSelectedNoteId] = useState<number | null>(null);
  const [search, setSearch] = useState("");
  const [pinFilter, setPinFilter] = useState(false);
  const [orphanDialogOpen, setOrphanDialogOpen] = useState(false);

  const handleSelectContainer = useCallback((name: string) => {
    setSelectedName(name);
    setSelectedNoteId(null);
  }, []);

  const orphanedNotes = useMemo(
    () =>
      notes.filter(
        (n) => !containers.some((c) => c.name === n.container_name),
      ),
    [notes, containers],
  );

  const noteCounts = useMemo(() => {
    const counts = new Map<string, number>();
    for (const n of notes) {
      counts.set(n.container_name, (counts.get(n.container_name) ?? 0) + 1);
    }
    return counts;
  }, [notes]);

  const pinnedNames = useMemo(
    () => new Set(notes.filter((n) => n.pinned).map((n) => n.container_name)),
    [notes],
  );

  const filteredContainers = useMemo(() => {
    let list = containers;
    if (search) {
      const term = search.toLowerCase();
      list = list.filter((c) => c.name.toLowerCase().includes(term));
    }
    if (pinFilter) {
      list = list.filter((c) => pinnedNames.has(c.name));
    }
    return list;
  }, [containers, search, pinFilter, pinnedNames]);

  const selectedContainer = containers.find((c) => c.name === selectedName) ?? null;
  const containerNotes = selectedName ? getNotesForContainer(selectedName) : [];
  const selectedNote = selectedNoteId != null
    ? containerNotes.find((n) => n.id === selectedNoteId) ?? null
    : null;

  const handleCreateNote = async (req: CreateNoteRequest) => {
    const note = await create(req);
    setSelectedNoteId(note.id);
  };

  const handleUpdateNote = async (id: number, req: UpdateNoteRequest) => {
    await update(id, req);
  };

  const handleDeleteNote = async (id: number) => {
    await remove(id);
    setSelectedNoteId(null);
  };

  return (
    <Box sx={{ display: "flex", flexDirection: "column", height: "100vh" }}>
      <Box
        sx={{
          display: "flex",
          alignItems: "center",
          px: 2,
          py: 1,
          gap: 2,
          borderBottom: 1,
          borderColor: "divider",
        }}
      >
        <Typography variant="h6" sx={{ flexShrink: 0 }}>
          RunNotes
        </Typography>
        <ExportImportButtons onRefresh={handleRefresh} />
        <Box sx={{ flex: 1 }} />
        <SearchBar value={search} onChange={setSearch} />
        <Tooltip title={pinFilter ? "Show all containers" : "Show pinned only"}>
          <IconButton
            onClick={() => setPinFilter(!pinFilter)}
            size="small"
            color={pinFilter ? "primary" : "default"}
          >
            {pinFilter ? <PushPinIcon /> : <PushPinOutlinedIcon />}
          </IconButton>
        </Tooltip>
        {orphanedNotes.length > 0 && (
          <Tooltip title={`${orphanedNotes.length} orphaned notes`}>
            <IconButton
              onClick={() => setOrphanDialogOpen(true)}
              size="small"
              color="warning"
            >
              <Badge badgeContent={orphanedNotes.length} color="warning">
                <WarningAmberIcon />
              </Badge>
            </IconButton>
          </Tooltip>
        )}
        <Tooltip title="Refresh container list">
          <span>
            <IconButton
              onClick={handleRefresh}
              size="small"
              disabled={containersLoading || notesLoading}
            >
              <RefreshIcon />
            </IconButton>
          </span>
        </Tooltip>
      </Box>

      <Box sx={{ display: "flex", flex: 1, overflow: "hidden" }}>
        <Box
          sx={{
            width: 320,
            flexShrink: 0,
            borderRight: 1,
            borderColor: "divider",
            overflowY: "auto",
          }}
        >
          <ContainerList
            containers={filteredContainers}
            selectedName={selectedName}
            onSelect={handleSelectContainer}
            noteCounts={noteCounts}
            loading={containersLoading || notesLoading}
          />
        </Box>

        <Divider orientation="vertical" flexItem />

        <Box sx={{ flex: 1, overflowY: "auto" }}>
          {selectedContainer == null ? (
            <Box
              sx={{
                display: "flex",
                alignItems: "center",
                justifyContent: "center",
                height: "100%",
              }}
            >
              <Typography variant="body1" color="text.secondary">
                Select a container to view or add notes
              </Typography>
            </Box>
          ) : selectedNote != null ? (
            <NoteEditor
              note={selectedNote}
              container={selectedContainer}
              onSave={handleUpdateNote}
              onDelete={handleDeleteNote}
              onBack={() => setSelectedNoteId(null)}
            />
          ) : (
            <NoteListView
              containerName={selectedName!}
              container={selectedContainer}
              notes={containerNotes}
              onSelectNote={setSelectedNoteId}
              onCreateNote={handleCreateNote}
            />
          )}
        </Box>
      </Box>

      <OrphanedNotesDialog
        open={orphanDialogOpen}
        onClose={() => setOrphanDialogOpen(false)}
        orphanedNotes={orphanedNotes}
        onDeleteContainer={async (name) => {
          await removeAllForContainer(name);
        }}
        onDeleteAll={async () => {
          const containerNames = new Set(orphanedNotes.map((n) => n.container_name));
          for (const name of containerNames) {
            await removeAllForContainer(name);
          }
          setOrphanDialogOpen(false);
        }}
      />
    </Box>
  );
}
