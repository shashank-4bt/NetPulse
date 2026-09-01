import { render, screen } from "@testing-library/react";
import { describe, expect, it } from "vitest";

import { CONFIDENCE_CAVEAT } from "@/domain/diagnostic";
import { emptyConfidence } from "@/features/intelligence/confidence";
import { ConfidencePanel } from "@/features/intelligence/confidence-panel";

describe("ConfidencePanel", () => {
  it("never implies certainty when no value is supplied", () => {
    render(
      <ConfidencePanel
        confidence={emptyConfidence()}
        evidence={[]}
        alternatives={[]}
      />
    );

    expect(screen.getByText(CONFIDENCE_CAVEAT)).toBeInTheDocument();
    expect(screen.getByText("No numeric value")).toBeInTheDocument();
    expect(screen.queryByText("certain")).not.toBeInTheDocument();
    expect(screen.queryByText("87%")).not.toBeInTheDocument();
  });
});
