import { render, screen } from "@testing-library/react";
import { describe, it, expect } from "vitest";
import App from "./App";

describe("App", () => {
  it("renders the heading", () => {
    render(<App />);
    expect(screen.getByText("RunNotes")).toBeInTheDocument();
  });

  it("shows loading state initially", () => {
    render(<App />);
    expect(screen.getByRole("progressbar")).toBeInTheDocument();
  });

  it("shows placeholder when no container is selected", () => {
    render(<App />);
    expect(
      screen.getByText("Select a container to view or add notes"),
    ).toBeInTheDocument();
  });
});
