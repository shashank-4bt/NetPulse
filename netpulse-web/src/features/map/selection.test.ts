import { describe, expect, it } from "vitest";

import { DEFAULT_MAP_QUERY, type MapCell } from "@/domain/map";
import { drillHref, drillLabel, nextQuery } from "@/features/map/selection";

const region: MapCell = {
  id: "region:eu-west",
  level: "region",
  label: "eu-west",
  parentId: null,
  lon: null,
  lat: null,
  sampleCount: 3,
  status: "insufficient_evidence",
  summary: "Stored region label",
  layer: "regional",
  childCount: 1,
};

describe("map selection", () => {
  it("drills region to a network summary", () => {
    const next = nextQuery(region, DEFAULT_MAP_QUERY);
    expect(next.level).toBe("network");
    expect(next.parent).toBe("region:eu-west");
    expect(drillLabel(region)).toBe("Open network summary");
    expect(drillHref(region, DEFAULT_MAP_QUERY)).toContain("parent=region%3Aeu-west");
  });

  it("opens service intelligence for a service cell", () => {
    const href = drillHref(
      { ...region, id: "service:youtube", level: "service", label: "youtube" },
      DEFAULT_MAP_QUERY
    );
    expect(href).toBe("/service/youtube");
    expect(drillLabel({ ...region, level: "service" })).toBe("Open service intelligence");
  });

  it("opens a country summary from the world cell", () => {
    const next = nextQuery({ ...region, id: "world", level: "world", label: "World" }, DEFAULT_MAP_QUERY);
    expect(next.level).toBe("country");
    expect(next.parent).toBe("world");
  });
});
