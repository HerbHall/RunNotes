import { useState, useEffect, useCallback } from "react";
import type { Note, CreateNoteRequest, UpdateNoteRequest } from "../types";
import {
  listNotes,
  createNote as apiCreate,
  updateNote as apiUpdate,
  deleteNote as apiDelete,
  deleteContainerNotes as apiDeleteContainer,
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
    async (id: number, req: UpdateNoteRequest) => {
      const note = await apiUpdate(id, req);
      showToast("success", "Note saved");
      await refresh();
      return note;
    },
    [refresh],
  );

  const remove = useCallback(
    async (id: number) => {
      await apiDelete(id);
      showToast("success", "Note deleted");
      await refresh();
    },
    [refresh],
  );

  const removeAllForContainer = useCallback(
    async (name: string) => {
      await apiDeleteContainer(name);
      showToast("success", "Notes deleted");
      await refresh();
    },
    [refresh],
  );

  const getNotesForContainer = useCallback(
    (containerName: string) =>
      notes.filter((n) => n.container_name === containerName),
    [notes],
  );

  return {
    notes,
    loading,
    error,
    refresh,
    create,
    update,
    remove,
    removeAllForContainer,
    getNotesForContainer,
  };
}
