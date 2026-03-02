import { render, screen } from "@testing-library/react";
import userEvent from "@testing-library/user-event";
import { describe, it, expect, vi } from "vitest";
import { NoteEditor } from "./NoteEditor";
import type { Note, ContainerInfo } from "../types";

const mockContainer: ContainerInfo = {
  id: "abc123",
  name: "web-app",
  image: "nginx:latest",
  state: "running",
  status: "Up 2 hours",
};

const mockNote: Note = {
  id: 1,
  container_name: "web-app",
  container_id: "abc123",
  compose_project: "",
  compose_service: "",
  title: "Test Note",
  note_content: "This is a test note",
  pinned: false,
  tags: ["web", "frontend"],
  created_at: "2026-01-01T00:00:00Z",
  updated_at: "2026-01-01T00:00:00Z",
};

describe("NoteEditor", () => {
  it("shows note content when note exists", () => {
    render(
      <NoteEditor
        note={mockNote}
        container={mockContainer}
        onSave={vi.fn()}
        onDelete={vi.fn()}
        onBack={vi.fn()}
      />,
    );
    const textarea = screen.getByPlaceholderText("Write your note here... (supports Markdown)");
    expect(textarea).toHaveValue("This is a test note");
  });

  it("shows note title in editable field", () => {
    render(
      <NoteEditor
        note={mockNote}
        container={mockContainer}
        onSave={vi.fn()}
        onDelete={vi.fn()}
        onBack={vi.fn()}
      />,
    );
    const titleInput = screen.getByPlaceholderText("Note title");
    expect(titleInput).toHaveValue("Test Note");
  });

  it("shows tags as chips when note exists", () => {
    render(
      <NoteEditor
        note={mockNote}
        container={mockContainer}
        onSave={vi.fn()}
        onDelete={vi.fn()}
        onBack={vi.fn()}
      />,
    );
    expect(screen.getByText("web")).toBeInTheDocument();
    expect(screen.getByText("frontend")).toBeInTheDocument();
  });

  it("calls onSave with note ID and updated content", async () => {
    const user = userEvent.setup();
    const onSave = vi.fn().mockResolvedValue(undefined);
    render(
      <NoteEditor
        note={mockNote}
        container={mockContainer}
        onSave={onSave}
        onDelete={vi.fn()}
        onBack={vi.fn()}
      />,
    );

    const textarea = screen.getByPlaceholderText("Write your note here... (supports Markdown)");
    await user.clear(textarea);
    await user.type(textarea, "Updated content");

    const saveBtn = screen.getByRole("button", { name: "Save" });
    await user.click(saveBtn);

    expect(onSave).toHaveBeenCalledWith(1, {
      title: "Test Note",
      note_content: "Updated content",
      pinned: false,
      tags: ["web", "frontend"],
      container_id: "abc123",
    });
  });

  it("shows delete confirmation dialog", async () => {
    const user = userEvent.setup();
    render(
      <NoteEditor
        note={mockNote}
        container={mockContainer}
        onSave={vi.fn()}
        onDelete={vi.fn()}
        onBack={vi.fn()}
      />,
    );

    const deleteBtn = screen.getByLabelText("Delete note");
    await user.click(deleteBtn);

    expect(screen.getByText("Delete Note")).toBeInTheDocument();
    expect(
      screen.getByText(/Are you sure you want to delete/),
    ).toBeInTheDocument();
  });

  it("shows Edit and Preview toggle buttons", () => {
    render(
      <NoteEditor
        note={mockNote}
        container={mockContainer}
        onSave={vi.fn()}
        onDelete={vi.fn()}
        onBack={vi.fn()}
      />,
    );
    expect(screen.getByRole("button", { name: /Edit/i })).toBeInTheDocument();
    expect(screen.getByRole("button", { name: /Preview/i })).toBeInTheDocument();
  });

  it("switches to preview mode and hides the text field", async () => {
    const user = userEvent.setup();
    render(
      <NoteEditor
        note={{ ...mockNote, note_content: "# Hello Markdown" }}
        container={mockContainer}
        onSave={vi.fn()}
        onDelete={vi.fn()}
        onBack={vi.fn()}
      />,
    );

    const previewBtn = screen.getByRole("button", { name: /Preview/i });
    await user.click(previewBtn);

    expect(screen.queryByPlaceholderText("Write your note here... (supports Markdown)")).not.toBeInTheDocument();
    expect(screen.getByRole("heading", { level: 1 })).toHaveTextContent("Hello Markdown");
  });

  it("switches back to edit mode from preview", async () => {
    const user = userEvent.setup();
    render(
      <NoteEditor
        note={{ ...mockNote, note_content: "# Hello Markdown" }}
        container={mockContainer}
        onSave={vi.fn()}
        onDelete={vi.fn()}
        onBack={vi.fn()}
      />,
    );

    const previewBtn = screen.getByRole("button", { name: /Preview/i });
    await user.click(previewBtn);

    expect(screen.queryByPlaceholderText("Write your note here... (supports Markdown)")).not.toBeInTheDocument();

    const editBtn = screen.getByRole("button", { name: /Edit/i });
    await user.click(editBtn);

    expect(screen.getByPlaceholderText("Write your note here... (supports Markdown)")).toBeInTheDocument();
  });

  it("calls onBack when back button is clicked", async () => {
    const user = userEvent.setup();
    const onBack = vi.fn();
    render(
      <NoteEditor
        note={mockNote}
        container={mockContainer}
        onSave={vi.fn()}
        onDelete={vi.fn()}
        onBack={onBack}
      />,
    );

    await user.click(screen.getByLabelText("Back to note list"));
    expect(onBack).toHaveBeenCalledOnce();
  });
});
