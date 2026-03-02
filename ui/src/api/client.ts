import { createDockerDesktopClient } from "@docker/extension-api-client";
import type {
  Note,
  CreateNoteRequest,
  UpdateNoteRequest,
  ImportResult,
  ContainerInfo,
} from "../types";

const ddClient = createDockerDesktopClient();

export async function listNotes(params?: {
  pinned?: boolean;
  search?: string;
}): Promise<Note[]> {
  const query = new URLSearchParams();
  if (params?.pinned !== undefined) query.set("pinned", String(params.pinned));
  if (params?.search) query.set("search", params.search);
  const qs = query.toString();
  const url = qs ? `/notes?${qs}` : "/notes";
  const result = await ddClient.extension.vm?.service?.get(url);
  return result as Note[];
}

export async function getNote(name: string): Promise<Note> {
  const result = await ddClient.extension.vm?.service?.get(
    `/notes/${encodeURIComponent(name)}`,
  );
  return result as Note;
}

export async function createNote(req: CreateNoteRequest): Promise<Note> {
  const result = await ddClient.extension.vm?.service?.post("/notes", req);
  return result as Note;
}

export async function updateNote(
  name: string,
  req: UpdateNoteRequest,
): Promise<Note> {
  const result = await ddClient.extension.vm?.service?.put(
    `/notes/${encodeURIComponent(name)}`,
    req,
  );
  return result as Note;
}

export async function deleteNote(name: string): Promise<void> {
  await ddClient.extension.vm?.service?.delete(
    `/notes/${encodeURIComponent(name)}`,
  );
}

export async function exportNotes(): Promise<Note[]> {
  const result = await ddClient.extension.vm?.service?.get("/notes/export");
  return result as Note[];
}

export async function importNotes(notes: Note[]): Promise<ImportResult> {
  const result = await ddClient.extension.vm?.service?.post(
    "/notes/import",
    notes,
  );
  return result as ImportResult;
}

export async function listContainers(): Promise<ContainerInfo[]> {
  const result = await ddClient.docker.listContainers({ all: true });
  return (result as Array<Record<string, unknown>>).map((c) => ({
    id: c.Id as string,
    name: ((c.Names as string[])?.[0] ?? "").replace(/^\//, ""),
    image: c.Image as string,
    state: c.State as string,
    status: c.Status as string,
    composeProject: (c.Labels as Record<string, string>)?.[
      "com.docker.compose.project"
    ],
    composeService: (c.Labels as Record<string, string>)?.[
      "com.docker.compose.service"
    ],
  }));
}

export function showToast(
  type: "success" | "warning" | "error",
  msg: string,
): void {
  ddClient.desktopUI.toast[type](msg);
}
