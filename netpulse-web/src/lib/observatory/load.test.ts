import { describe, expect, it } from "vitest";

import { emptyServiceIntelligence } from "@/features/observatory/empty-intelligence";
import { parseOutageQuery } from "@/features/observatory/filter-incidents";

describe("observatory load helpers", () => {
  it("keeps unmeasured intelligence off live scores", () => {
    const intel = emptyServiceIntelligence();
    expect(intel.currentState).toBe("not_measured");
    expect(intel.health).toBeNull();
    expect(intel.availability.sampleCount).toBe(0);
    expect(intel.availability.measured).toBe(false);
    expect(intel.regionalHealth).toEqual([]);
  });

  it("parses outage filters from search params", () => {
    const query = parseOutageQuery({
      service: "youtube",
      severity: "high",
      page: "2",
      q: "elevated",
    });
    expect(query.service).toBe("youtube");
    expect(query.page).toBe(2);
    expect(query.q).toBe("elevated");
  });
});
