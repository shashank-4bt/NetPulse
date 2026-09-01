import { describe, expect, it } from "vitest";

import { filterCells, filterIncidentRefs } from "@/features/map/filter-cells";
import type { MapCell } from "@/domain/map";

const cells: MapCell[] = [
  {
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
  },
  {
    id: "network:AS64500",
    level: "network",
    label: "AS64500",
    parentId: null,
    lon: null,
    lat: null,
    sampleCount: 3,
    status: "insufficient_evidence",
    summary: "Stored network",
    layer: "network",
    childCount: 1,
  },
];

describe("filterCells", () => {
  it("hides layers that are turned off", () => {
    const result = filterCells(cells, { layers: ["regional"], q: "" });
    expect(result.map((cell) => cell.id)).toEqual(["region:eu-west"]);
  });

  it("searches labels without inventing rows", () => {
    expect(filterCells(cells, { layers: ["network"], q: "github" })).toEqual([]);
    expect(filterCells([], { layers: ["regional"], q: "" })).toEqual([]);
  });

  it("hides incident refs when the incidents layer is off", () => {
    const refs = [{ id: "1", title: "Elevated connectivity failures observed", coarseRegion: "eu-west" }];
    expect(filterIncidentRefs(refs, { layers: ["regional"], q: "" })).toEqual([]);
    expect(filterIncidentRefs(refs, { layers: ["incidents"], q: "elevated" })).toHaveLength(1);
  });
});
