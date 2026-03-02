import { render, screen } from "@testing-library/react";
import userEvent from "@testing-library/user-event";
import { describe, it, expect, vi } from "vitest";
import { SearchBar } from "./SearchBar";

describe("SearchBar", () => {
  it("renders with the given value", () => {
    render(<SearchBar value="hello" onChange={vi.fn()} />);
    const input = screen.getByPlaceholderText("Search containers...");
    expect(input).toHaveValue("hello");
  });

  it("calls onChange when user types", async () => {
    const user = userEvent.setup();
    const onChange = vi.fn();
    render(<SearchBar value="" onChange={onChange} />);

    const input = screen.getByPlaceholderText("Search containers...");
    await user.type(input, "test");

    expect(onChange).toHaveBeenCalled();
    expect(onChange).toHaveBeenCalledWith("t");
  });
});
