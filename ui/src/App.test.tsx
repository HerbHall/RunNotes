import { render, screen } from "@testing-library/react";
import { describe, it, expect } from "vitest";
import App from "./App";

describe("App", () => {
  it("renders the heading", () => {
    render(<App />);
    expect(screen.getByText("RunNotes")).toBeInTheDocument();
  });

  it("renders the description", () => {
    render(<App />);
    expect(
      screen.getByText(
        "Attach notes and annotations to your Docker containers.",
      ),
    ).toBeInTheDocument();
  });
});
