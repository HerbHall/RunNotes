import { render, screen } from "@testing-library/react";
import userEvent from "@testing-library/user-event";
import { describe, it, expect, vi } from "vitest";
import { NoteListView } from "./NoteListView";
import type { Note, ContainerInfo } from "../types";

const mockContainer: ContainerInfo = {
  id: "abc123",
  name: "web-app",
  image: "nginx:latest",
  state: "running",
  status: "Up 2 hours",
};

const mockNotes: Note[] = [
  {
    id: 1,
    container_name: "web-app",
    container_id: "abc123",
    compose_project: "",
    compose_service: "",
    title: "Config Notes",
    note_content: "Nginx config details",
    pinned: true,
    tags: ["config"],
    created_at: "2026-01-01T00:00:00Z",
    updated_at: "2026-01-02T00:00:00Z",
  },
  {
    id: 2,
    container_name: "web-app",
    container_id: "abc123",
    compose_project: "",
    compose_service: "",
    title: "Deployment Notes",
    note_content: "Steps for deploying",
    pinned: false,
    tags: [],
    created_at: "2026-01-01T00:00:00Z",
    updated_at: "2026-01-03T00:00:00Z",
  },
];

describe("NoteListView", () => {
  it("shows container name as heading", () => {
    render(
      <NoteListView
        containerName="web-app"
        container={mockContainer}
        notes={mockNotes}
        onSelectNote={vi.fn()}
        onCreateNote={vi.fn()}
      />,
    );
    expect(screen.getByText("web-app")).toBeInTheDocument();
  });

  it("shows 'Add Note' button", () => {
    render(
      <NoteListView
        containerName="web-app"
        container={mockContainer}
        notes={mockNotes}
        onSelectNote={vi.fn()}
        onCreateNote={vi.fn()}
      />,
    );
    expect(screen.getByRole("button", { name: /Add Note/i })).toBeInTheDocument();
  });

  it("lists note titles", () => {
    render(
      <NoteListView
        containerName="web-app"
        container={mockContainer}
        notes={mockNotes}
        onSelectNote={vi.fn()}
        onCreateNote={vi.fn()}
      />,
    );
    expect(screen.getByText("Config Notes")).toBeInTheDocument();
    expect(screen.getByText("Deployment Notes")).toBeInTheDocument();
  });

  it("shows empty state when no notes", () => {
    render(
      <NoteListView
        containerName="web-app"
        container={mockContainer}
        notes={[]}
        onSelectNote={vi.fn()}
        onCreateNote={vi.fn()}
      />,
    );
    expect(screen.getByText("No notes yet for this container")).toBeInTheDocument();
  });

  it("calls onSelectNote when a note is clicked", async () => {
    const user = userEvent.setup();
    const onSelectNote = vi.fn();
    render(
      <NoteListView
        containerName="web-app"
        container={mockContainer}
        notes={mockNotes}
        onSelectNote={onSelectNote}
        onCreateNote={vi.fn()}
      />,
    );
    await user.click(screen.getByText("Config Notes"));
    expect(onSelectNote).toHaveBeenCalledWith(1);
  });

  it("shows title input when Add Note is clicked", async () => {
    const user = userEvent.setup();
    render(
      <NoteListView
        containerName="web-app"
        container={mockContainer}
        notes={mockNotes}
        onSelectNote={vi.fn()}
        onCreateNote={vi.fn()}
      />,
    );
    await user.click(screen.getByRole("button", { name: /Add Note/i }));
    expect(screen.getByPlaceholderText("Note title...")).toBeInTheDocument();
  });

  it("shows pin icon for pinned notes", () => {
    const { container } = render(
      <NoteListView
        containerName="web-app"
        container={mockContainer}
        notes={mockNotes}
        onSelectNote={vi.fn()}
        onCreateNote={vi.fn()}
      />,
    );
    const pinIcons = container.querySelectorAll("[data-testid='PushPinIcon']");
    expect(pinIcons).toHaveLength(1);
  });
});
