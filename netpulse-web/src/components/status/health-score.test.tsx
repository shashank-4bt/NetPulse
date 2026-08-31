import { render, screen } from "@testing-library/react";
import { describe, expect, it } from "vitest";

import { HealthScore } from "@/components/status/health-score";
import { INSUFFICIENT_EVIDENCE_LABEL } from "@/lib/design/taxonomy";

describe("HealthScore", () => {
  it("does not invent a score when value is null", () => {
    render(<HealthScore value={null} confidence={null} />);
    expect(screen.getByRole("status")).toHaveTextContent(
      INSUFFICIENT_EVIDENCE_LABEL
    );
    expect(screen.queryByText("/ 100")).not.toBeInTheDocument();
  });

  it("renders a provided score with confidence text", () => {
    render(<HealthScore value={72} confidence="medium" />);
    expect(screen.getByText("72")).toBeInTheDocument();
    expect(screen.getByText("Medium")).toBeInTheDocument();
  });
});
