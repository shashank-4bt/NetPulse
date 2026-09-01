import { describe, expect, it } from "vitest";

import { getPublicHealthSnapshot } from "@/lib/api/public-health";

describe("getPublicHealthSnapshot", () => {
  it("never invents live health or incidents", async () => {
    const snapshot = await getPublicHealthSnapshot();
    expect(snapshot.state).toBe("unavailable");
    expect(snapshot.incidents).toEqual([]);
    expect(snapshot.regions).toEqual([]);
    expect(snapshot.measuredAt).toBeNull();
  });
});
