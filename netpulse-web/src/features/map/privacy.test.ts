import { describe, expect, it } from "vitest";

import { isCoarseCoordinate, stripPreciseCoordinates } from "@/features/map/privacy";

describe("map privacy", () => {
  it("rejects sub-degree coordinates", () => {
    expect(isCoarseCoordinate(-0.1276, 51.5074)).toBe(false);
    expect(isCoarseCoordinate(0, 52)).toBe(true);
  });

  it("strips precise points instead of plotting them", () => {
    const stripped = stripPreciseCoordinates({
      lon: -0.1276,
      lat: 51.5074,
    });
    expect(stripped.lon).toBeNull();
    expect(stripped.lat).toBeNull();
  });
});
