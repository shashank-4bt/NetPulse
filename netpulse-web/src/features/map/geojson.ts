import type { MapCell } from "@/domain/map";
import { isCoarseCoordinate } from "@/features/map/privacy";

export type MapFeatureCollection = {
  type: "FeatureCollection";
  features: Array<{
    type: "Feature";
    id: string;
    geometry: { type: "Point"; coordinates: [number, number] };
    properties: {
      id: string;
      label: string;
      status: string;
      sampleCount: number;
      level: string;
    };
  }>;
};

export function plottableCells(cells: MapCell[]): MapCell[] {
  return cells.filter((cell) => isCoarseCoordinate(cell.lon, cell.lat));
}

export function cellsToGeoJSON(cells: MapCell[]): MapFeatureCollection {
  return {
    type: "FeatureCollection",
    features: plottableCells(cells).map((cell) => ({
      type: "Feature" as const,
      id: cell.id,
      geometry: {
        type: "Point" as const,
        coordinates: [cell.lon as number, cell.lat as number],
      },
      properties: {
        id: cell.id,
        label: cell.label,
        status: cell.status,
        sampleCount: cell.sampleCount,
        level: String(cell.level),
      },
    })),
  };
}
