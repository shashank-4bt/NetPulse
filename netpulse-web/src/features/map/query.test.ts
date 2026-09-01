import { describe, expect, it } from "vitest";

import { DEFAULT_MAP_QUERY } from "@/domain/map";
import {
  mapAggregatesQueryString,
  mapQueryString,
  parseMapQuery,
} from "@/features/map/query";

describe("map query", () => {
  it("defaults to world with every layer", () => {
    const query = parseMapQuery({});
    expect(query).toEqual(DEFAULT_MAP_QUERY);
  });

  it("parses repeated layer params and preserves drill crumbs", () => {
    const query = parseMapQuery({
      level: "network",
      parent: "region:eu-west",
      region: "eu-west",
      layers: ["regional", "network"],
      q: "AS64500",
    });
    expect(query.level).toBe("network");
    expect(query.layers).toEqual(["regional", "network"]);
    expect(mapQueryString(query)).toContain("parent=region%3Aeu-west");
  });

  it("includes viewport bounds on aggregate fetches", () => {
    const encoded = mapAggregatesQueryString(DEFAULT_MAP_QUERY, {
      west: -10.4,
      south: 40.4,
      east: 10.4,
      north: 60.4,
    });
    expect(encoded).toContain("west=-10.4");
    expect(encoded).toContain("level=world");
  });
});
