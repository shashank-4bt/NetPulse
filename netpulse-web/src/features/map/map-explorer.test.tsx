import { render, screen } from "@testing-library/react";
import { describe, expect, it, vi } from "vitest";

import { DEFAULT_MAP_QUERY, emptyMapAggregates, type MapCell } from "@/domain/map";
import { MapExplorer } from "@/features/map/map-explorer";

vi.mock("next/dynamic", () => ({
  default: () => {
    function MapStub() {
      return <div>Map canvas stub</div>;
    }
    return MapStub;
  },
}));

vi.mock("next/navigation", () => ({
  useRouter: () => ({ push: vi.fn() }),
}));

const fixtureCell: MapCell = {
  id: "region:eu-west",
  level: "region",
  label: "eu-west",
  parentId: null,
  lon: null,
  lat: null,
  sampleCount: 3,
  status: "insufficient_evidence",
  summary: "Stored region label eu-west. Observed sample count 3.",
  layer: "regional",
  childCount: 1,
};

describe("MapExplorer", () => {
  it("loads an empty map without inventing health", () => {
    render(
      <MapExplorer
        query={DEFAULT_MAP_QUERY}
        initial={emptyMapAggregates("No coarse geographic aggregates are stored.")}
        state="empty"
        reason="No coarse geographic aggregates are stored."
        liveEnabled={false}
      />
    );
    expect(screen.getByText("Map canvas stub")).toBeInTheDocument();
    expect(screen.getByText("No aggregated cells match these filters.")).toBeInTheDocument();
    expect(screen.getByText("No map aggregates")).toBeInTheDocument();
    expect(screen.getByText(/table is the primary way/i)).toBeInTheDocument();
    expect(screen.queryByText("87%")).not.toBeInTheDocument();
    expect(screen.queryByText(/operational/i)).not.toBeInTheDocument();
  });

  it("shows a keyboard-accessible table alternative for stored cells", () => {
    render(
      <MapExplorer
        query={{ ...DEFAULT_MAP_QUERY, select: "region:eu-west" }}
        initial={{
          ...emptyMapAggregates("Aggregates are counted from stored incidents."),
          cells: [fixtureCell],
          totalSamples: 3,
          reason: "Aggregates are counted from stored incidents.",
        }}
        state="ready"
        reason="Aggregates are counted from stored incidents."
        liveEnabled={false}
      />
    );
    const links = screen.getAllByRole("link", { name: "eu-west" });
    expect(links[0]).toHaveAttribute("href", expect.stringContaining("/map"));
    expect(screen.getAllByText("Open network summary").length).toBeGreaterThan(0);
    expect(screen.getByText("Insufficient evidence")).toBeInTheDocument();
    expect(screen.getByText("Not plotted")).toBeInTheDocument();
  });

  it("renders unavailable and error states without filling the table", () => {
    const { rerender } = render(
      <MapExplorer
        query={DEFAULT_MAP_QUERY}
        initial={emptyMapAggregates("backend down")}
        state="unavailable"
        reason="backend down"
        liveEnabled={false}
      />
    );
    expect(screen.getByText("Map store unavailable")).toBeInTheDocument();
    rerender(
      <MapExplorer
        query={DEFAULT_MAP_QUERY}
        initial={emptyMapAggregates("store failed")}
        state="error"
        reason="store failed"
        liveEnabled={false}
      />
    );
    expect(screen.getByText("Map aggregates failed")).toBeInTheDocument();
  });
});
