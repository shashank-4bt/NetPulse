import type { MapCell, MapQuery } from "@/domain/map";
import { mapQueryString } from "@/features/map/query";

export function drillHref(cell: MapCell, query: MapQuery): string {
  if (cell.level === "service") {
    const slug = cell.id.startsWith("service:") ? cell.id.slice("service:".length) : cell.label;
    return `/service/${encodeURIComponent(slug)}`;
  }
  return `/map${mapQueryString(nextQuery(cell, query))}`;
}

export function incidentHref(id: string): string {
  return `/incident/${encodeURIComponent(id)}`;
}

export function nextQuery(cell: MapCell, query: MapQuery): MapQuery {
  const next: MapQuery = { ...query, select: cell.id };
  switch (cell.level) {
    case "world":
      return { ...next, level: "country", parent: "world" };
    case "country":
      return { ...next, level: "region", parent: cell.id, country: cell.label };
    case "region":
      return { ...next, level: "network", parent: cell.id, region: cell.label };
    case "network":
      return { ...next, level: "service", parent: cell.id, network: cell.label };
    default:
      return next;
  }
}

export function selectedCell(cells: MapCell[], select: string): MapCell | null {
  if (!select) {
    return null;
  }
  return cells.find((cell) => cell.id === select) ?? null;
}

export function drillLabel(cell: MapCell): string {
  switch (cell.level) {
    case "world":
      return "Open country summary";
    case "country":
      return "Open regional summary";
    case "region":
      return "Open network summary";
    case "network":
      return "Open service summary";
    case "service":
      return "Open service intelligence";
    default:
      return "Open details";
  }
}
