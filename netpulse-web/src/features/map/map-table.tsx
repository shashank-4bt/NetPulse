import Link from "next/link";

import {
  Table,
  TableBody,
  TableCaption,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table";
import { LayerStatusBadge } from "@/components/status/layer-status-badge";
import { Badge } from "@/components/ui/badge";
import type { LayerStatus } from "@/domain/display";
import { LAYER_STATUSES } from "@/domain/display";
import {
  MAP_LAYER_LABELS,
  MAP_LEVEL_LABELS,
  type MapAggregates,
  type MapLayer,
  type MapQuery,
} from "@/domain/map";
import { drillHref, drillLabel, incidentHref } from "@/features/map/selection";

type MapTableProps = {
  aggregates: MapAggregates;
  query: MapQuery;
};

export function MapTable({ aggregates, query }: MapTableProps) {
  return (
    <Table>
      <TableCaption className="mt-0 caption-top text-left">
        Equivalent important information for the map. Sample counts are observed, not
        population impact. Coordinates are omitted unless a coarse centroid is stored.
      </TableCaption>
      <TableHeader>
        <TableRow>
          <TableHead>Label</TableHead>
          <TableHead>Level</TableHead>
          <TableHead>Layer</TableHead>
          <TableHead>Status</TableHead>
          <TableHead>Samples</TableHead>
          <TableHead>Location</TableHead>
          <TableHead>Open</TableHead>
        </TableRow>
      </TableHeader>
      <TableBody>
        {aggregates.cells.length === 0 ? (
          <TableRow>
            <TableCell colSpan={7} className="text-muted-foreground">
              No aggregated cells match these filters.
            </TableCell>
          </TableRow>
        ) : (
          aggregates.cells.map((cell) => (
            <TableRow
              key={cell.id}
              data-state={query.select === cell.id ? "selected" : undefined}
            >
              <TableCell>
                <Link
                  href={drillHref(cell, query)}
                  className="font-medium hover:underline focus-visible:outline-none focus-visible:ring-3 focus-visible:ring-ring/50"
                >
                  {cell.label}
                </Link>
              </TableCell>
              <TableCell>{levelLabel(cell.level)}</TableCell>
              <TableCell>{layerLabel(cell.layer)}</TableCell>
              <TableCell>
                {isLayerStatus(cell.status) ? (
                  <LayerStatusBadge status={cell.status} />
                ) : (
                  <Badge variant="outline">{cell.status || "Unclassified"}</Badge>
                )}
              </TableCell>
              <TableCell className="font-mono text-xs">{cell.sampleCount}</TableCell>
              <TableCell className="font-mono text-xs text-muted-foreground">
                {cell.lon != null && cell.lat != null
                  ? `${cell.lon}, ${cell.lat} (1°)`
                  : "Not plotted"}
              </TableCell>
              <TableCell>
                <Link
                  href={drillHref(cell, query)}
                  className="text-sm hover:underline"
                >
                  {drillLabel(cell)}
                </Link>
              </TableCell>
            </TableRow>
          ))
        )}
        {aggregates.incidentRefs.map((item) => (
          <TableRow key={`incident:${item.id}`}>
            <TableCell>
              <Link href={incidentHref(item.id)} className="hover:underline">
                {item.title || item.id}
              </Link>
            </TableCell>
            <TableCell>Incident</TableCell>
            <TableCell>{MAP_LAYER_LABELS.incidents}</TableCell>
            <TableCell>
              <Badge variant="outline">Stored record</Badge>
            </TableCell>
            <TableCell className="font-mono text-xs">{item.sampleCount}</TableCell>
            <TableCell className="text-muted-foreground">
              {item.coarseRegion || "No coarse region"}
            </TableCell>
            <TableCell>
              <Link href={incidentHref(item.id)} className="text-sm hover:underline">
                Open incident
              </Link>
            </TableCell>
          </TableRow>
        ))}
      </TableBody>
    </Table>
  );
}

function isLayerStatus(value: string): value is LayerStatus {
  return (LAYER_STATUSES as readonly string[]).includes(value);
}

function levelLabel(level: string): string {
  return (MAP_LEVEL_LABELS as Record<string, string>)[level] ?? level;
}

function layerLabel(layer: string): string {
  return (MAP_LAYER_LABELS as Record<MapLayer, string>)[layer as MapLayer] ?? layer;
}
