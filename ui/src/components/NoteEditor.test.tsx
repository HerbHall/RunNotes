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
  note_content: "This is a test note",
  pinned: false,
  tags: ["web", "frontend"],
  created_at: "2026-01-01T00:00:00Z",
  updated_at: "2026-01-01T00:00:00Z",
};

describe("NoteEditor", () => {
  it("shows 'Add Note' button when note is null", () => {
    render(
      <NoteEditor
        containerName="web-app"
        container={mockContainer}
        note={null}
        onSave={vi.fn()}
        onDelete={vi.fn()}
      />,
    );
    expect(screen.getByText("No note yet for web-app")).toBeInTheDocument();
    expect(screen.getByRole("button", { name: "Add Note" })).toBeInTheDocument();
  });

  it("shows note content when note exists", () => {
    render(
      <NoteEditor
        containerName="web-app"
        container={mockContainer}
        note={mockNote}
        onSave={vi.fn()}
        onDelete={vi.fn()}
      />,
    );
    const textarea = screen.getByPlaceholderText("Write your note here... (supports Markdown)");
    expect(textarea).toHaveValue("This is a test note");
  });

  it("shows tags as chips when note exists", () => {
    render(
      <NoteEditor
        containerName="web-app"
        container={mockContainer}
        note={mockNote}
        onSave={vi.fn()}
        onDelete={vi.fn()}
      />,
    );
    expect(screen.getByText("web")).toBeInTheDocument();
    expect(screen.getByText("frontend")).toBeInTheDocument();
  });

  it("calls onSave with updated content", async () => {
    const user = userEvent.setup();
    const onSave = vi.fn().mockResolvedValue(undefined);
    render(
      <NoteEditor
        containerName="web-app"
        container={mockContainer}
        note={mockNote}
        onSave={onSave}
        onDelete={vi.fn()}
      />,
    );

    const textarea = screen.getByPlaceholderText("Write your note here... (supports Markdown)");
    await user.clear(textarea);
    await user.type(textarea, "Updated content");

    const saveBtn = screen.getByRole("button", { name: "Save" });
    await user.click(saveBtn);

    expect(onSave).toHaveBeenCalledWith("web-app", {
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
        containerName="web-app"
        container={mockContainer}
        note={mockNote}
        onSave={vi.fn()}
        onDelete={vi.fn()}
      />,
    );

    const deleteBtn = screen.getByLabelText("Delete note");
    await user.click(deleteBtn);

    expect(screen.getByText("Delete Note")).toBeInTheDocument();
    expect(
      screen.getByText(/Are you sure you want to delete/),
    ).toBeInTheDocument();
  });

  it("shows Edit and Preview toggle buttons when note exists", () => {
    render(
      <NoteEditor
        containerName="web-app"
        container={mockContainer}
        note={mockNote}
        onSave={vi.fn()}
        onDelete={vi.fn()}
      />,
    );
    expect(screen.getByRole("button", { name: /Edit/i })).toBeInTheDocument();
    expect(screen.getByRole("button", { name: /Preview/i })).toBeInTheDocument();
  });

  it("switches to preview mode and hides the text field", async () => {
    const user = userEvent.setup();
    render(
      <NoteEditor
        containerName="web-app"
        container={mockContainer}
        note={{ ...mockNote, note_content: "# Hello Markdown" }}
        onSave={vi.fn()}
        onDelete={vi.fn()}
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
        containerName="web-app"
        container={mockContainer}
        note={{ ...mockNote, note_content: "# Hello Markdown" }}
        onSave={vi.fn()}
        onDelete={vi.fn()}
      />,
    );

    const previewBtn = screen.getByRole("button", { name: /Preview/i });
    await user.click(previewBtn);

    expect(screen.queryByPlaceholderText("Write your note here... (supports Markdown)")).not.toBeInTheDocument();

    const editBtn = screen.getByRole("button", { name: /Edit/i });
    await user.click(editBtn);

    expect(screen.getByPlaceholderText("Write your note here... (supports Markdown)")).toBeInTheDocument();
  });
});
