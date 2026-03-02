import { render, screen } from "@testing-library/react";
import { describe, it, expect } from "vitest";
import { MarkdownPreview } from "./MarkdownPreview";

describe("MarkdownPreview", () => {
  it("renders a markdown heading as an h1 element", () => {
    render(<MarkdownPreview content="# Hello World" />);
    const heading = screen.getByRole("heading", { level: 1 });
    expect(heading).toHaveTextContent("Hello World");
  });

  it("renders a markdown h2 heading", () => {
    render(<MarkdownPreview content="## Section Title" />);
    const heading = screen.getByRole("heading", { level: 2 });
    expect(heading).toHaveTextContent("Section Title");
  });

  it("renders bold text", () => {
    render(<MarkdownPreview content="This is **bold** text" />);
    const bold = screen.getByText("bold");
    expect(bold.tagName).toBe("STRONG");
  });

  it("renders a code block", () => {
    render(<MarkdownPreview content={"```\nconsole.log('hi')\n```"} />);
    expect(screen.getByText("console.log('hi')")).toBeInTheDocument();
  });

  it("renders plain text without errors", () => {
    render(<MarkdownPreview content="Just plain text here" />);
    expect(screen.getByText("Just plain text here")).toBeInTheDocument();
  });

  it("renders an unordered list", () => {
    render(<MarkdownPreview content={"- item one\n- item two"} />);
    expect(screen.getByText("item one")).toBeInTheDocument();
    expect(screen.getByText("item two")).toBeInTheDocument();
  });

  it("renders empty content without errors", () => {
    const { container } = render(<MarkdownPreview content="" />);
    expect(container).toBeInTheDocument();
  });
});
