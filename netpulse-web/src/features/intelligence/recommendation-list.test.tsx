import { render, screen } from "@testing-library/react";
import { describe, expect, it } from "vitest";

import { RecommendationList } from "@/features/intelligence/recommendation-list";

describe("RecommendationList", () => {
  it("shows structured fields and never offers an execute action", () => {
    render(
      <RecommendationList
        recommendations={[
          {
            id: "safe-rerun",
            action: "Re-run after workers exist",
            reason: "No facts yet",
            risk: "Acting now is misleading",
            expectedResult: "A measured report",
            verification: "Wait for measured facts",
            safetyClass: "dangerous",
            autoExecute: false,
            evidenceIds: [],
          },
        ]}
      />
    );

    expect(screen.getByText("Action")).toBeInTheDocument();
    expect(screen.getByText("Reason")).toBeInTheDocument();
    expect(screen.getByText("Risk")).toBeInTheDocument();
    expect(screen.getByText("Expected result")).toBeInTheDocument();
    expect(screen.getByText("Verification")).toBeInTheDocument();
    expect(screen.getByText("Dangerous — not executed")).toBeInTheDocument();
    expect(screen.getByText(/Possible security concern/i)).toBeInTheDocument();
    expect(screen.queryByRole("button", { name: /execute/i })).not.toBeInTheDocument();
  });
});
