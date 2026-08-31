import { describe, expect, it } from "vitest";

import {
  CONFIDENCE_LEVELS,
  OPERATIONAL_STATUSES,
  SEVERITIES,
} from "@/domain/display";
import {
  CONFIDENCE_LABEL,
  SEVERITY_LABEL,
  STATUS_LABEL,
} from "@/lib/design/taxonomy";

describe("display taxonomy", () => {
  it("maps every status, severity, and confidence to a visible label", () => {
    expect(OPERATIONAL_STATUSES.map((status) => STATUS_LABEL[status])).toEqual([
      "Operational",
      "Degraded",
      "Investigating",
      "Major Incident",
      "Unknown",
    ]);
    expect(SEVERITIES.map((severity) => SEVERITY_LABEL[severity])).toEqual([
      "Informational",
      "Moderate",
      "High",
      "Critical",
    ]);
    expect(
      CONFIDENCE_LEVELS.map((confidence) => CONFIDENCE_LABEL[confidence])
    ).toEqual(["Low", "Medium", "High", "Very High"]);
  });
});
