package geo

import (
	"math"
	"testing"

	"github.com/shashank-4bt/NetPulse/netpulse-api/internal/contract"
)

func TestSnapDegreeDropsPreciseLocations(t *testing.T) {
	lon, lat, ok := CoarsePoint(-0.1276, 51.5074)
	if !ok {
		t.Fatal("expected valid coarse point")
	}
	if lon != 0 || lat != 52 {
		t.Fatalf("got %v,%v want 0,52", lon, lat)
	}
	if _, _, ok := CoarsePoint(200, 10); ok {
		t.Fatal("out of range longitude must be rejected")
	}
	if _, _, ok := CoarsePoint(0, math.NaN()); ok {
		t.Fatal("NaN must be rejected")
	}
}

func TestViewportFiltersOnlyCoarsePoints(t *testing.T) {
	west, south, east, north := -10.0, 40.0, 10.0, 60.0
	q := Query{West: &west, South: &south, East: &east, North: &north}
	if !InViewport(0, 52, q) {
		t.Fatal("point inside viewport")
	}
	if InViewport(20, 52, q) {
		t.Fatal("point outside viewport")
	}
}

func TestFinalizeNeverKeepsSubDegreeCoordinates(t *testing.T) {
	lon, lat := -0.1276, 51.5074
	agg := Finalize([]contract.MapCell{{
		ID:          "fixture:london",
		Level:       LevelRegion,
		Label:       "fixture-london",
		Lon:         &lon,
		Lat:         &lat,
		SampleCount: 4,
		Status:      "insufficient_evidence",
		Layer:       LayerRegional,
	}}, nil, Query{Level: LevelRegion, Limit: MaxCells, Layers: DefaultLayers})
	if len(agg.Cells) != 1 {
		t.Fatalf("cells %d", len(agg.Cells))
	}
	if agg.Cells[0].Lon == nil || agg.Cells[0].Lat == nil {
		t.Fatal("coarse centroid should remain")
	}
	if *agg.Cells[0].Lon != 0 || *agg.Cells[0].Lat != 52 {
		t.Fatalf("expected snapped degree, got %v,%v", *agg.Cells[0].Lon, *agg.Cells[0].Lat)
	}
	if agg.Precision != "degree" {
		t.Fatalf("precision %s", agg.Precision)
	}
}
