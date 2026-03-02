import { cleanup, render, screen, waitFor } from "@testing-library/react";
import userEvent from "@testing-library/user-event";
import { describe, it, expect, vi, beforeEach, afterEach } from "vitest";
import { ExportImportButtons } from "./ExportImportButtons";

vi.mock("../api/client", () => ({
  exportNotes: vi.fn(),
  importNotes: vi.fn(),
  showToast: vi.fn(),
}));

import { exportNotes, importNotes, showToast } from "../api/client";

const mockExportNotes = vi.mocked(exportNotes);
const mockImportNotes = vi.mocked(importNotes);
const mockShowToast = vi.mocked(showToast);

describe("ExportImportButtons", () => {
  const onRefresh = vi.fn().mockResolvedValue(undefined);

  beforeEach(() => {
    vi.clearAllMocks();
  });

  afterEach(() => {
    cleanup();
  });

  it("renders Export and Import buttons", () => {
    render(<ExportImportButtons onRefresh={onRefresh} />);
    expect(
      screen.getByRole("button", { name: /export/i }),
    ).toBeInTheDocument();
    expect(
      screen.getByRole("button", { name: /import/i }),
    ).toBeInTheDocument();
  });

  it("calls exportNotes on export click", async () => {
    const user = userEvent.setup();
    mockExportNotes.mockResolvedValueOnce([]);

    // Mock only the anchor element creation, not all createElement calls
    const mockAnchor = document.createElement("a");
    mockAnchor.click = vi.fn();
    const origCreateElement = document.createElement.bind(document);
    vi.spyOn(document, "createElement").mockImplementation((tag: string) => {
      if (tag === "a") return mockAnchor;
      return origCreateElement(tag);
    });
    globalThis.URL.createObjectURL = vi.fn().mockReturnValue("blob:fake");
    globalThis.URL.revokeObjectURL = vi.fn();

    render(<ExportImportButtons onRefresh={onRefresh} />);
    await user.click(screen.getByRole("button", { name: /export/i }));

    await waitFor(() => {
      expect(mockExportNotes).toHaveBeenCalled();
    });
    expect(mockShowToast).toHaveBeenCalledWith("success", "Exported 0 notes");

    vi.restoreAllMocks();
  });

  it("shows confirmation dialog after file selection", async () => {
    const user = userEvent.setup();
    const notesJson = JSON.stringify([
      { container_name: "a", note_content: "note-a", tags: [] },
      { container_name: "b", note_content: "note-b", tags: [] },
    ]);

    render(<ExportImportButtons onRefresh={onRefresh} />);

    const file = new File([notesJson], "notes.json", {
      type: "application/json",
    });
    const input = document.querySelector(
      'input[type="file"]',
    ) as HTMLInputElement;
    await user.upload(input, file);

    expect(await screen.findByText("Import Notes")).toBeInTheDocument();
    expect(screen.getByText(/Import 2 notes/)).toBeInTheDocument();
  });

  it("imports notes on confirm and refreshes", async () => {
    const user = userEvent.setup();
    mockImportNotes.mockResolvedValueOnce({ imported: 2 });

    const notesJson = JSON.stringify([
      { container_name: "a", note_content: "note-a", tags: [] },
      { container_name: "b", note_content: "note-b", tags: [] },
    ]);

    render(<ExportImportButtons onRefresh={onRefresh} />);

    const file = new File([notesJson], "notes.json", {
      type: "application/json",
    });
    const input = document.querySelector(
      'input[type="file"]',
    ) as HTMLInputElement;
    await user.upload(input, file);

    // Wait for dialog to appear
    await screen.findByText("Import Notes");

    // Find the dialog's Import button (the contained variant in DialogActions)
    const dialog = screen.getByRole("dialog");
    const dialogImportBtn = dialog.querySelector(
      "button.MuiButton-contained",
    ) as HTMLButtonElement;
    await user.click(dialogImportBtn);

    await waitFor(() => {
      expect(mockImportNotes).toHaveBeenCalled();
    });
    await waitFor(() => {
      expect(onRefresh).toHaveBeenCalled();
    });
    expect(mockShowToast).toHaveBeenCalledWith("success", "Imported 2 notes");
  });

  it("closes dialog on cancel", async () => {
    const user = userEvent.setup();
    const notesJson = JSON.stringify([
      { container_name: "a", note_content: "note-a", tags: [] },
    ]);

    render(<ExportImportButtons onRefresh={onRefresh} />);

    const file = new File([notesJson], "notes.json", {
      type: "application/json",
    });
    const input = document.querySelector(
      'input[type="file"]',
    ) as HTMLInputElement;
    await user.upload(input, file);

    await screen.findByText("Import Notes");
    await user.click(screen.getByRole("button", { name: "Cancel" }));

    await waitFor(() => {
      expect(screen.queryByText("Import Notes")).not.toBeInTheDocument();
    });
  });
});
