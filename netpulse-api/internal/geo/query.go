package geo

import (
	"strconv"
	"strings"
)

const (
	MaxCells        = 250
	MaxIncidentRefs = 50
	MinSampleHealth = 3
)

const (
	LevelWorld   = "world"
	LevelCountry = "country"
	LevelRegion  = "region"
	LevelNetwork = "network"
	LevelService = "service"
)

const (
	LayerGlobal    = "global"
	LayerRegional  = "regional"
	LayerNetwork   = "network"
	LayerService   = "service"
	LayerIncidents = "incidents"
)

var DefaultLayers = []string{LayerGlobal, LayerRegional, LayerNetwork, LayerService, LayerIncidents}

type Query struct {
	Level   string
	Parent  string
	West    *float64
	South   *float64
	East    *float64
	North   *float64
	Layers  []string
	Search  string
	Service string
	Limit   int
}

func ParseQuery(values map[string]string) Query {
	limit, _ := strconv.Atoi(values["limit"])
	if limit < 1 || limit > MaxCells {
		limit = MaxCells
	}
	level := strings.ToLower(strings.TrimSpace(values["level"]))
	if !validLevel(level) {
		level = LevelWorld
	}
	layers := parseLayers(values["layers"])
	return Query{
		Level:   level,
		Parent:  strings.TrimSpace(values["parent"]),
		West:    parseBound(values["west"]),
		South:   parseBound(values["south"]),
		East:    parseBound(values["east"]),
		North:   parseBound(values["north"]),
		Layers:  layers,
		Search:  strings.TrimSpace(values["q"]),
		Service: strings.TrimSpace(values["service"]),
		Limit:   limit,
	}
}

func (q Query) HasViewport() bool {
	return q.West != nil && q.South != nil && q.East != nil && q.North != nil
}

func (q Query) HasLayer(layer string) bool {
	if len(q.Layers) == 0 {
		return true
	}
	for _, item := range q.Layers {
		if item == layer {
			return true
		}
	}
	return false
}

func validLevel(level string) bool {
	switch level {
	case LevelWorld, LevelCountry, LevelRegion, LevelNetwork, LevelService:
		return true
	default:
		return false
	}
}

func parseLayers(raw string) []string {
	if strings.TrimSpace(raw) == "" {
		return append([]string{}, DefaultLayers...)
	}
	seen := map[string]bool{}
	out := []string{}
	for _, part := range strings.Split(raw, ",") {
		layer := strings.ToLower(strings.TrimSpace(part))
		if !validLayer(layer) || seen[layer] {
			continue
		}
		seen[layer] = true
		out = append(out, layer)
	}
	if len(out) == 0 {
		return append([]string{}, DefaultLayers...)
	}
	return out
}

func validLayer(layer string) bool {
	switch layer {
	case LayerGlobal, LayerRegional, LayerNetwork, LayerService, LayerIncidents:
		return true
	default:
		return false
	}
}

func parseBound(raw string) *float64 {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return nil
	}
	value, err := strconv.ParseFloat(raw, 64)
	if err != nil {
		return nil
	}
	snapped := SnapDegree(value)
	return &snapped
}
