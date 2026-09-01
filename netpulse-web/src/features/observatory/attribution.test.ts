import { describe, expect, it } from "vitest";

import {
  incidentTitle,
  isCausalOverclaim,
} from "@/features/observatory/attribution";

describe("incident attribution", () => {
  it("prefers observed-failure language over ISP blame", () => {
    expect(
      incidentTitle({
        evidenceClass: "inferred_hypothesis",
        isolatedLayer: "isp",
        namedNetwork: "AS123",
      })
    ).toBe("Elevated connectivity failures observed");
    expect(isCausalOverclaim("ISP X caused the outage")).toBe(true);
  });

  it("may name a network only after measured isolation", () => {
    const title = incidentTitle({
      evidenceClass: "measured_fact",
      isolatedLayer: "isp",
      namedNetwork: "AS64500",
    });
    expect(title).toBe("Measured failures isolated to network AS64500");
    expect(isCausalOverclaim(title)).toBe(false);
  });
});
