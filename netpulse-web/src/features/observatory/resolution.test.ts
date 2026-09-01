import { describe, expect, it } from "vitest";

import { canMarkResolved } from "@/features/observatory/resolution";

describe("incident resolution policy", () => {
  it("never marks resolved from one recovered measurement", () => {
    const decision = canMarkResolved({
      recoverySampleCount: 1,
      identifiedCause: true,
    });
    expect(decision.ok).toBe(false);
  });

  it("requires an identified cause in addition to recoveries", () => {
    expect(
      canMarkResolved({ recoverySampleCount: 4, identifiedCause: false }).ok
    ).toBe(false);
  });
});
