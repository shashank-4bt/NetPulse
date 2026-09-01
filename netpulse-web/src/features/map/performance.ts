import { MAP_CELL_LIMIT, type MapAggregates, type MapCell } from "@/domain/map";

export function capCells(
  cells: MapCell[],
  limit = MAP_CELL_LIMIT
): { cells: MapCell[]; truncated: boolean } {
  if (cells.length <= limit) {
    return { cells, truncated: false };
  }
  return { cells: cells.slice(0, limit), truncated: true };
}

export function payloadHasRawMeasurements(value: unknown): boolean {
  if (!value || typeof value !== "object") {
    return false;
  }
  if ("measurements" in value) {
    return true;
  }
  const record = value as { cells?: unknown };
  if (!Array.isArray(record.cells)) {
    return false;
  }
  return record.cells.some(
    (cell) => cell && typeof cell === "object" && "measurements" in cell
  );
}

export function sanitizeAggregates(value: MapAggregates): MapAggregates {
  const capped = capCells(value.cells);
  return {
    ...value,
    cells: capped.cells,
    incidentRefs: value.incidentRefs.slice(0, 50),
    truncated: value.truncated || capped.truncated,
    limit: MAP_CELL_LIMIT,
  };
}
