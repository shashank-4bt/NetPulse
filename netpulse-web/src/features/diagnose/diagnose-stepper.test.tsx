import { render, screen } from "@testing-library/react";
import { describe, expect, it } from "vitest";

import { DiagnoseStepper } from "@/features/diagnose/diagnose-stepper";

describe("DiagnoseStepper", () => {
  it("shows step counts and omits duration and percentages when they are not supplied", () => {
    render(
      <DiagnoseStepper
        steps={[
          {
            id: "initializing",
            label: "Initializing",
            state: "complete",
            durationMs: null,
            note: null,
          },
          {
            id: "dns",
            label: "DNS",
            state: "unavailable",
            durationMs: null,
            note: "Not run",
          },
          {
            id: "complete",
            label: "Complete",
            state: "complete",
            durationMs: 12,
            note: null,
          },
        ]}
      />
    );

    expect(screen.getByText("2 complete · 1 unavailable · 3 total")).toBeInTheDocument();
    expect(screen.getByText("12 ms")).toBeInTheDocument();
    expect(screen.queryByText("%")).not.toBeInTheDocument();
    expect(screen.queryByText("Failed")).not.toBeInTheDocument();
  });
});
