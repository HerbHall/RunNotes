import { render, screen } from "@testing-library/react";
import userEvent from "@testing-library/user-event";
import { describe, it, expect, vi } from "vitest";
import { ContainerList } from "./ContainerList";
import type { ContainerInfo } from "../types";

const mockContainers: ContainerInfo[] = [
  {
    id: "abc123",
    name: "web-app",
    image: "nginx:latest",
    state: "running",
    status: "Up 2 hours",
  },
  {
    id: "def456",
    name: "database",
    image: "postgres:16",
    state: "exited",
    status: "Exited (0) 5 minutes ago",
  },
];

describe("ContainerList", () => {
  it("renders a list of containers", () => {
    render(
      <ContainerList
        containers={mockContainers}
        selectedName={null}
        onSelect={vi.fn()}
        noteCounts={new Map()}
        loading={false}
      />,
    );
    expect(screen.getByText("web-app")).toBeInTheDocument();
    expect(screen.getByText("database")).toBeInTheDocument();
  });

  it("shows empty state message when no containers", () => {
    render(
      <ContainerList
        containers={[]}
        selectedName={null}
        onSelect={vi.fn()}
        noteCounts={new Map()}
        loading={false}
      />,
    );
    expect(screen.getByText("No containers found")).toBeInTheDocument();
  });

  it("shows loading spinner", () => {
    render(
      <ContainerList
        containers={[]}
        selectedName={null}
        onSelect={vi.fn()}
        noteCounts={new Map()}
        loading={true}
      />,
    );
    expect(screen.getByRole("progressbar")).toBeInTheDocument();
  });

  it("calls onSelect when a container is clicked", async () => {
    const user = userEvent.setup();
    const onSelect = vi.fn();
    render(
      <ContainerList
        containers={mockContainers}
        selectedName={null}
        onSelect={onSelect}
        noteCounts={new Map()}
        loading={false}
      />,
    );
    await user.click(screen.getByText("web-app"));
    expect(onSelect).toHaveBeenCalledWith("web-app");
  });

  it("shows note indicator for containers with notes", () => {
    const { container } = render(
      <ContainerList
        containers={mockContainers}
        selectedName={null}
        onSelect={vi.fn()}
        noteCounts={new Map([["web-app", 1]])}
        loading={false}
      />,
    );
    const noteIcons = container.querySelectorAll("[data-testid='DescriptionIcon']");
    expect(noteIcons).toHaveLength(1);
  });

  it("shows badge with count for containers with multiple notes", () => {
    render(
      <ContainerList
        containers={mockContainers}
        selectedName={null}
        onSelect={vi.fn()}
        noteCounts={new Map([["web-app", 3]])}
        loading={false}
      />,
    );
    expect(screen.getByText("3")).toBeInTheDocument();
  });
});
