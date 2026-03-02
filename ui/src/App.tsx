import { useState, useMemo, useEffect, useRef, useCallback } from "react";
import Box from "@mui/material/Box";
import Divider from "@mui/material/Divider";
import IconButton from "@mui/material/IconButton";
import Tooltip from "@mui/material/Tooltip";
import Typography from "@mui/material/Typography";
import PushPinIcon from "@mui/icons-material/PushPin";
import PushPinOutlinedIcon from "@mui/icons-material/PushPinOutlined";
import RefreshIcon from "@mui/icons-material/Refresh";
import { ContainerList } from "./components/ContainerList";
import { ExportImportButtons } from "./components/ExportImportButtons";
import { NoteEditor } from "./components/NoteEditor";
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
    refresh: refreshNotes,
    getNoteForContainer,
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
          return update(note.container_name, { container_id: container.id });
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
  const [search, setSearch] = useState("");
  const [pinFilter, setPinFilter] = useState(false);

  const noteNames = useMemo(
    () => new Set(notes.map((n) => n.container_name)),
    [notes],
  );

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
  const selectedNote = selectedName ? getNoteForContainer(selectedName) : null;

  const handleSave = async (
    name: string,
    req: CreateNoteRequest | UpdateNoteRequest,
  ) => {
    if ("container_name" in req) {
      await create(req);
    } else {
      await update(name, req);
    }
  };

  const handleDelete = async (name: string) => {
    await remove(name);
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
            onSelect={setSelectedName}
            noteNames={noteNames}
            loading={containersLoading || notesLoading}
          />
        </Box>

        <Divider orientation="vertical" flexItem />

        <Box sx={{ flex: 1, overflowY: "auto" }}>
          {selectedContainer != null ? (
            <NoteEditor
              containerName={selectedName!}
              container={selectedContainer}
              note={selectedNote}
              onSave={handleSave}
              onDelete={handleDelete}
            />
          ) : (
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
          )}
        </Box>
      </Box>
    </Box>
  );
}
