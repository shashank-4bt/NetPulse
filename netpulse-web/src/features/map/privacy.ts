import type { MapViewport } from "@/domain/map";

export function snapDegree(value: number): number {
  return Math.round(value);
}

export function isCoarseCoordinate(
  lon: number | null | undefined,
  lat: number | null | undefined
): boolean {
  if (lon == null || lat == null) {
    return false;
  }
  if (!Number.isFinite(lon) || !Number.isFinite(lat)) {
    return false;
  }
  if (lon < -180 || lon > 180 || lat < -90 || lat > 90) {
    return false;
  }
  return Number.isInteger(lon) && Number.isInteger(lat);
}

export function roundViewport(viewport: MapViewport): MapViewport {
  return {
    west: snapDegree(viewport.west),
    south: snapDegree(viewport.south),
    east: snapDegree(viewport.east),
    north: snapDegree(viewport.north),
  };
}

export function stripPreciseCoordinates<T extends { lon: number | null; lat: number | null }>(
  cell: T
): T {
  if (!isCoarseCoordinate(cell.lon, cell.lat)) {
    return { ...cell, lon: null, lat: null };
  }
  return cell;
}
