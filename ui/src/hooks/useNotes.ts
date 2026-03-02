import { useState, useEffect, useCallback } from "react";
import type { Note, CreateNoteRequest, UpdateNoteRequest } from "../types";
import {
  listNotes,
  createNote as apiCreate,
  updateNote as apiUpdate,
  deleteNote as apiDelete,
  showToast,
} from "../api/client";

export function useNotes() {
  const [notes, setNotes] = useState<Note[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  const refresh = useCallback(async () => {
    try {
      setLoading(true);
      setError(null);
      const result = await listNotes();
      setNotes(result);
    } catch (err) {
      setError(err instanceof Error ? err.message : "Failed to load notes");
    } finally {
      setLoading(false);
    }
  }, []);

  useEffect(() => {
    refresh();
  }, [refresh]);

  const create = useCallback(
    async (req: CreateNoteRequest) => {
      const note = await apiCreate(req);
      showToast("success", "Note created");
      await refresh();
      return note;
    },
    [refresh],
  );

  const update = useCallback(
    async (name: string, req: UpdateNoteRequest) => {
      const note = await apiUpdate(name, req);
      showToast("success", "Note saved");
      await refresh();
      return note;
    },
    [refresh],
  );

  const remove = useCallback(
    async (name: string) => {
      await apiDelete(name);
      showToast("success", "Note deleted");
      await refresh();
    },
    [refresh],
  );

  const getNoteForContainer = useCallback(
    (containerName: string) =>
      notes.find((n) => n.container_name === containerName) ?? null,
    [notes],
  );

  return { notes, loading, error, refresh, create, update, remove, getNoteForContainer };
}
