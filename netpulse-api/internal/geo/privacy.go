package geo

import (
	"math"
	"strconv"
)

// SnapDegree rounds a coordinate to 1 degree (~100 km). Precise points are never returned.
func SnapDegree(value float64) float64 {
	return math.Round(value)
}

func CoarsePoint(lon, lat float64) (float64, float64, bool) {
	if lon < -180 || lon > 180 || lat < -90 || lat > 90 {
		return 0, 0, false
	}
	if math.IsNaN(lon) || math.IsNaN(lat) || math.IsInf(lon, 0) || math.IsInf(lat, 0) {
		return 0, 0, false
	}
	return SnapDegree(lon), SnapDegree(lat), true
}

func InViewport(lon, lat float64, q Query) bool {
	if !q.HasViewport() {
		return true
	}
	west, south, east, north := *q.West, *q.South, *q.East, *q.North
	if lat < south || lat > north {
		return false
	}
	if west <= east {
		return lon >= west && lon <= east
	}
	return lon >= west || lon <= east
}

func ClusterKey(lon, lat float64) string {
	return strconv.Itoa(int(SnapDegree(lon))) + ":" + strconv.Itoa(int(SnapDegree(lat)))
}
