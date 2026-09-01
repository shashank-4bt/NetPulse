import { describe, expect, it } from "vitest";

import { MAP_CELL_LIMIT, type MapCell } from "@/domain/map";
import { cellsToGeoJSON, plottableCells } from "@/features/map/geojson";
import { capCells, payloadHasRawMeasurements } from "@/features/map/performance";

function cell(id: string, lon: number | null = null, lat: number | null = null): MapCell {
  return {
    id,
    level: "region",
    label: id,
    parentId: null,
    lon,
    lat,
    sampleCount: 1,
    status: "insufficient_evidence",
    summary: "fixture",
    layer: "regional",
    childCount: 0,
  };
}

describe("map performance", () => {
  it("caps cells and never treats raw measurements as map data", () => {
    const cells = Array.from({ length: MAP_CELL_LIMIT + 10 }, (_, index) => cell(`c${index}`));
    const capped = capCells(cells);
    expect(capped.cells).toHaveLength(MAP_CELL_LIMIT);
    expect(capped.truncated).toBe(true);
    expect(payloadHasRawMeasurements({ cells: capped.cells })).toBe(false);
    expect(payloadHasRawMeasurements({ measurements: [{ id: "raw" }] })).toBe(true);
  });

  it("does not plot cells without coarse coordinates", () => {
    const geojson = cellsToGeoJSON([
      cell("label-only"),
      cell("precise", -0.12, 51.5),
      cell("coarse", 0, 52),
    ]);
    expect(plottableCells([cell("precise", -0.12, 51.5)])).toEqual([]);
    expect(geojson.features).toHaveLength(1);
    expect(geojson.features[0]?.properties.id).toBe("coarse");
  });
});
