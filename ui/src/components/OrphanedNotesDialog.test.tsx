import { render, screen, within } from "@testing-library/react";
import userEvent from "@testing-library/user-event";
import { describe, it, expect, vi } from "vitest";
import { OrphanedNotesDialog } from "./OrphanedNotesDialog";
import type { Note } from "../types";

const mockOrphanedNotes: Note[] = [
  {
    id: 1,
    container_name: "old-web-server",
    container_id: "abc123",
    compose_project: "",
    compose_service: "",
    title: "Server Config",
    note_content: "This container was running nginx for the old site configuration",
    pinned: false,
    tags: [],
    created_at: "2026-01-15T10:00:00Z",
    updated_at: "2026-02-20T14:30:00Z",
  },
  {
    id: 2,
    container_name: "temp-database",
    container_id: "def456",
    compose_project: "",
    compose_service: "",
    title: "Migration Notes",
    note_content: "Temporary postgres for migration testing",
    pinned: false,
    tags: [],
    created_at: "2026-01-10T08:00:00Z",
    updated_at: "2026-02-18T09:15:00Z",
  },
];

describe("OrphanedNotesDialog", () => {
  it("shows empty state when no orphaned notes", () => {
    render(
      <OrphanedNotesDialog
        open={true}
        onClose={vi.fn()}
        orphanedNotes={[]}
        onDeleteContainer={vi.fn()}
        onDeleteAll={vi.fn()}
      />,
    );
    expect(screen.getByText("No orphaned notes found")).toBeInTheDocument();
    expect(screen.queryByText("Delete All")).not.toBeInTheDocument();
  });

  it("shows orphaned notes grouped by container name", () => {
    render(
      <OrphanedNotesDialog
        open={true}
        onClose={vi.fn()}
        orphanedNotes={mockOrphanedNotes}
        onDeleteContainer={vi.fn()}
        onDeleteAll={vi.fn()}
      />,
    );
    expect(screen.getByText("old-web-server")).toBeInTheDocument();
    expect(screen.getByText("temp-database")).toBeInTheDocument();
  });

  it("shows note titles in secondary text", () => {
    render(
      <OrphanedNotesDialog
        open={true}
        onClose={vi.fn()}
        orphanedNotes={mockOrphanedNotes}
        onDeleteContainer={vi.fn()}
        onDeleteAll={vi.fn()}
      />,
    );
    expect(screen.getByText(/Server Config/)).toBeInTheDocument();
    expect(screen.getByText(/Migration Notes/)).toBeInTheDocument();
  });

  it("calls onDeleteContainer when individual delete button is clicked", async () => {
    const user = userEvent.setup();
    const onDeleteContainer = vi.fn().mockResolvedValue(undefined);
    render(
      <OrphanedNotesDialog
        open={true}
        onClose={vi.fn()}
        orphanedNotes={mockOrphanedNotes}
        onDeleteContainer={onDeleteContainer}
        onDeleteAll={vi.fn()}
      />,
    );
    const deleteButton = screen.getByLabelText(
      "Delete notes for old-web-server",
    );
    await user.click(deleteButton);
    expect(onDeleteContainer).toHaveBeenCalledWith("old-web-server");
  });

  it("shows confirmation dialog when Delete All is clicked", async () => {
    const user = userEvent.setup();
    render(
      <OrphanedNotesDialog
        open={true}
        onClose={vi.fn()}
        orphanedNotes={mockOrphanedNotes}
        onDeleteContainer={vi.fn()}
        onDeleteAll={vi.fn()}
      />,
    );
    await user.click(screen.getByText("Delete All"));
    expect(
      screen.getByText("Delete 2 orphaned notes? This cannot be undone."),
    ).toBeInTheDocument();
  });

  it("calls onDeleteAll when confirmation is accepted", async () => {
    const user = userEvent.setup();
    const onDeleteAll = vi.fn().mockResolvedValue(undefined);
    render(
      <OrphanedNotesDialog
        open={true}
        onClose={vi.fn()}
        orphanedNotes={mockOrphanedNotes}
        onDeleteContainer={vi.fn()}
        onDeleteAll={onDeleteAll}
      />,
    );
    await user.click(screen.getByText("Delete All"));
    const confirmDialog = screen.getByText("Confirm Deletion").closest<HTMLElement>('[role="dialog"]')!;
    await user.click(within(confirmDialog).getByText("Delete"));
    expect(onDeleteAll).toHaveBeenCalledOnce();
  });

  it("does not call onDeleteAll when confirmation is cancelled", async () => {
    const user = userEvent.setup();
    const onDeleteAll = vi.fn().mockResolvedValue(undefined);
    render(
      <OrphanedNotesDialog
        open={true}
        onClose={vi.fn()}
        orphanedNotes={mockOrphanedNotes}
        onDeleteContainer={vi.fn()}
        onDeleteAll={onDeleteAll}
      />,
    );
    await user.click(screen.getByText("Delete All"));
    const confirmDialog = screen.getByText("Confirm Deletion").closest<HTMLElement>('[role="dialog"]')!;
    await user.click(within(confirmDialog).getByText("Cancel"));
    expect(onDeleteAll).not.toHaveBeenCalled();
  });

  it("calls onClose when Close button is clicked", async () => {
    const user = userEvent.setup();
    const onClose = vi.fn();
    render(
      <OrphanedNotesDialog
        open={true}
        onClose={onClose}
        orphanedNotes={mockOrphanedNotes}
        onDeleteContainer={vi.fn()}
        onDeleteAll={vi.fn()}
      />,
    );
    await user.click(screen.getByText("Close"));
    expect(onClose).toHaveBeenCalledOnce();
  });

  it("does not render dialog content when closed", () => {
    render(
      <OrphanedNotesDialog
        open={false}
        onClose={vi.fn()}
        orphanedNotes={mockOrphanedNotes}
        onDeleteContainer={vi.fn()}
        onDeleteAll={vi.fn()}
      />,
    );
    expect(screen.queryByText("Orphaned Notes")).not.toBeInTheDocument();
  });
});
