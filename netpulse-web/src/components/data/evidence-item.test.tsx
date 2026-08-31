import { render, screen } from "@testing-library/react";
import { describe, expect, it } from "vitest";

import { EvidenceItem } from "@/components/data/evidence-item";

describe("EvidenceItem", () => {
  it("labels measured facts distinctly from inferences", () => {
    const { rerender } = render(
      <EvidenceItem
        evidenceClass="measured_fact"
        title="TLS handshake completed"
        body="Fixture body"
      />
    );
    expect(screen.getByText("Measured fact")).toBeInTheDocument();

    rerender(
      <EvidenceItem
        evidenceClass="inferred_hypothesis"
        title="Possible congestion"
        body="Fixture body"
      />
    );
    expect(screen.getByText("Inferred hypothesis")).toBeInTheDocument();
    expect(screen.queryByText("Measured fact")).not.toBeInTheDocument();
  });
});
