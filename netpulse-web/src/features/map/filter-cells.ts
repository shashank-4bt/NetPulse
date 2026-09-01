import type { MapCell, MapLayer, MapQuery } from "@/domain/map";

export function filterCells(
  cells: MapCell[],
  query: Pick<MapQuery, "layers" | "q">
): MapCell[] {
  const layers = new Set(query.layers);
  const needle = query.q.trim().toLowerCase();
  return cells.filter((cell) => {
    if (layers.size > 0 && !layers.has(cell.layer as MapLayer)) {
      return false;
    }
    if (!needle) {
      return true;
    }
    return `${cell.label} ${cell.summary} ${cell.id}`.toLowerCase().includes(needle);
  });
}

export function filterIncidentRefs(
  refs: { id: string; title: string; coarseRegion: string }[],
  query: Pick<MapQuery, "layers" | "q">
): typeof refs {
  if (!query.layers.includes("incidents")) {
    return [];
  }
  const needle = query.q.trim().toLowerCase();
  if (!needle) {
    return refs;
  }
  return refs.filter((item) =>
    `${item.title} ${item.coarseRegion} ${item.id}`.toLowerCase().includes(needle)
  );
}
