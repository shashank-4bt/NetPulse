"use client";

import { useEffect, useRef } from "react";
import {
  Map as MapLibreMap,
  NavigationControl,
  setWorkerUrl,
  type ErrorEvent,
  type GeoJSONSource,
  type MapLayerMouseEvent,
} from "maplibre-gl";
import "maplibre-gl/dist/maplibre-gl.css";

import type { MapCell, MapViewport } from "@/domain/map";
import { cellsToGeoJSON } from "@/features/map/geojson";
import { roundViewport } from "@/features/map/privacy";

const BASEMAP_STYLE = "https://tiles.openfreemap.org/styles/positron";
const SOURCE_ID = "netpulse-aggregates";

if (typeof window !== "undefined") {
  setWorkerUrl("/maplibre/maplibre-gl-worker.mjs");
}

type MapCanvasProps = {
  cells: MapCell[];
  selectedId: string;
  onSelect: (id: string) => void;
  onViewport: (viewport: MapViewport) => void;
  onLoadError: (message: string | null) => void;
  onReady: () => void;
};

export function MapCanvas({
  cells,
  selectedId,
  onSelect,
  onViewport,
  onLoadError,
  onReady,
}: MapCanvasProps) {
  const containerRef = useRef<HTMLDivElement>(null);
  const mapRef = useRef<MapLibreMap | null>(null);
  const onSelectRef = useRef(onSelect);
  const onViewportRef = useRef(onViewport);
  const onLoadErrorRef = useRef(onLoadError);
  const onReadyRef = useRef(onReady);

  useEffect(() => {
    onSelectRef.current = onSelect;
    onViewportRef.current = onViewport;
    onLoadErrorRef.current = onLoadError;
    onReadyRef.current = onReady;
  }, [onSelect, onViewport, onLoadError, onReady]);

  useEffect(() => {
    const node = containerRef.current;
    if (!node || mapRef.current) {
      return;
    }
    const map = new MapLibreMap({
      container: node,
      style: BASEMAP_STYLE,
      center: [0, 20],
      zoom: 1.15,
      cooperativeGestures: true,
    });
    mapRef.current = map;
    map.addControl(new NavigationControl({ showCompass: false }), "top-right");

    const emitViewport = () => {
      const bounds = map.getBounds();
      onViewportRef.current(
        roundViewport({
          west: bounds.getWest(),
          south: bounds.getSouth(),
          east: bounds.getEast(),
          north: bounds.getNorth(),
        })
      );
    };

    map.on("load", () => {
      map.addSource(SOURCE_ID, {
        type: "geojson",
        data: cellsToGeoJSON([]),
        cluster: true,
        clusterMaxZoom: 6,
        clusterRadius: 48,
      });
      map.addLayer({
        id: "np-clusters",
        type: "circle",
        source: SOURCE_ID,
        filter: ["has", "point_count"],
        paint: {
          "circle-color": "#64748b",
          "circle-radius": 16,
          "circle-stroke-width": 2,
          "circle-stroke-color": "#0f172a",
        },
      });
      map.addLayer({
        id: "np-cluster-count",
        type: "symbol",
        source: SOURCE_ID,
        filter: ["has", "point_count"],
        layout: {
          "text-field": ["concat", ["get", "point_count_abbreviated"], " locations"],
          "text-size": 11,
        },
        paint: {
          "text-color": "#0f172a",
        },
      });
      map.addLayer({
        id: "np-points",
        type: "circle",
        source: SOURCE_ID,
        filter: ["!", ["has", "point_count"]],
        paint: {
          "circle-color": [
            "match",
            ["get", "status"],
            "insufficient_evidence",
            "#ca8a04",
            "not_measured",
            "#94a3b8",
            "#64748b",
          ],
          "circle-radius": 8,
          "circle-stroke-width": 2,
          "circle-stroke-color": "#0f172a",
        },
      });
      map.addLayer({
        id: "np-point-labels",
        type: "symbol",
        source: SOURCE_ID,
        filter: ["!", ["has", "point_count"]],
        layout: {
          "text-field": [
            "concat",
            ["get", "label"],
            " · ",
            ["get", "status"],
            " · n=",
            ["to-string", ["get", "sampleCount"]],
          ],
          "text-size": 11,
          "text-offset": [0, 1.2],
          "text-anchor": "top",
        },
        paint: {
          "text-color": "#0f172a",
          "text-halo-color": "#f8fafc",
          "text-halo-width": 1,
        },
      });
      onReadyRef.current();
      emitViewport();
    });

    map.on("error", (event: ErrorEvent) => {
      const message = event.error?.message ?? "The basemap failed to load.";
      onLoadErrorRef.current(message);
    });

    map.on("moveend", emitViewport);

    map.on("click", "np-points", (event: MapLayerMouseEvent) => {
      const id = event.features?.[0]?.properties?.id;
      if (typeof id === "string") {
        onSelectRef.current(id);
      }
    });

    map.on("click", "np-clusters", (event: MapLayerMouseEvent) => {
      const feature = event.features?.[0];
      if (!feature || feature.geometry.type !== "Point") {
        return;
      }
      const source = map.getSource(SOURCE_ID) as GeoJSONSource | undefined;
      const clusterId = feature.properties?.cluster_id;
      if (!source || typeof clusterId !== "number") {
        return;
      }
      const coordinates = feature.geometry.coordinates as [number, number];
      void source.getClusterExpansionZoom(clusterId).then((zoom: number) => {
        map.easeTo({ center: coordinates, zoom });
      });
    });

    return () => {
      map.remove();
      mapRef.current = null;
    };
  }, []);

  useEffect(() => {
    const map = mapRef.current;
    const source = map?.getSource(SOURCE_ID) as GeoJSONSource | undefined;
    if (!source) {
      return;
    }
    source.setData(cellsToGeoJSON(cells));
  }, [cells]);

  useEffect(() => {
    const map = mapRef.current;
    if (!map || !selectedId) {
      return;
    }
    const plotted = cells.find((cell) => cell.id === selectedId);
    if (!plotted || plotted.lon == null || plotted.lat == null) {
      return;
    }
    map.easeTo({ center: [plotted.lon, plotted.lat], zoom: Math.max(map.getZoom(), 3) });
  }, [selectedId, cells]);

  return (
    <div
      ref={containerRef}
      className="np-map h-full min-h-[12rem] w-full"
      role="region"
      aria-label="Internet health map canvas. Use the data table for the same aggregates."
    />
  );
}
