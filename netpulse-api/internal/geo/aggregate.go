package geo

import (
	"fmt"
	"sort"
	"strings"

	"github.com/shashank-4bt/NetPulse/netpulse-api/internal/contract"
)

func Aggregate(incidents []contract.Incident, query Query) contract.MapAggregates {
	if query.Limit < 1 || query.Limit > MaxCells {
		query.Limit = MaxCells
	}
	matched := matchingIncidents(incidents, query)
	cells := buildCells(matched, query)
	cells = filterCells(cells, query)
	cells = Cluster(cells)
	truncated := false
	if len(cells) > query.Limit {
		cells = cells[:query.Limit]
		truncated = true
	}

	refs := []contract.MapIncidentRef{}
	if query.HasLayer(LayerIncidents) {
		refs = incidentRefs(matched)
		if len(refs) > MaxIncidentRefs {
			refs = refs[:MaxIncidentRefs]
			truncated = true
		}
	}

	total := 0
	hasCoords := false
	for _, item := range matched {
		total += item.SampleCount
	}
	for _, cell := range cells {
		if cell.Lon != nil && cell.Lat != nil {
			hasCoords = true
			break
		}
	}

	precision := "none"
	if hasCoords {
		precision = "degree"
	}

	reason := "No coarse geographic aggregates are stored. The map will not invent country health."
	if len(matched) > 0 {
		reason = "Aggregates are counted from stored incidents. Coordinates are omitted unless a coarse centroid is stored."
	}
	if truncated {
		reason += " Result limited to aggregated cells for this viewport and filter."
	}

	var parent *string
	if query.Parent != "" {
		parent = &query.Parent
	}

	return contract.NormalizeMap(contract.MapAggregates{
		Level:          query.Level,
		ParentID:       parent,
		Cells:          cells,
		IncidentRefs:   refs,
		TotalSamples:   total,
		Limit:          query.Limit,
		Truncated:      truncated,
		Precision:      precision,
		Reason:         reason,
		HasCoordinates: hasCoords,
	})
}

// Finalize applies privacy, viewport, clustering, and the cell cap to already-built cells.
// Tests use this with fixture centroids. Production aggregation does not invent coordinates.
func Finalize(cells []contract.MapCell, refs []contract.MapIncidentRef, query Query) contract.MapAggregates {
	if query.Limit < 1 || query.Limit > MaxCells {
		query.Limit = MaxCells
	}
	sanitized := make([]contract.MapCell, 0, len(cells))
	for _, cell := range cells {
		if snapped, ok := applyPrivacy(cell); ok {
			sanitized = append(sanitized, snapped)
		}
	}
	sanitized = filterCells(sanitized, query)
	sanitized = Cluster(sanitized)
	truncated := false
	if len(sanitized) > query.Limit {
		sanitized = sanitized[:query.Limit]
		truncated = true
	}
	if refs == nil {
		refs = []contract.MapIncidentRef{}
	}
	if len(refs) > MaxIncidentRefs {
		refs = refs[:MaxIncidentRefs]
		truncated = true
	}
	hasCoords := false
	total := 0
	for _, cell := range sanitized {
		total += cell.SampleCount
		if cell.Lon != nil && cell.Lat != nil {
			hasCoords = true
		}
	}
	precision := "none"
	if hasCoords {
		precision = "degree"
	}
	return contract.NormalizeMap(contract.MapAggregates{
		Level:          query.Level,
		Cells:          sanitized,
		IncidentRefs:   refs,
		TotalSamples:   total,
		Limit:          query.Limit,
		Truncated:      truncated,
		Precision:      precision,
		Reason:         "Fixture aggregates for tests. Production responses omit coordinates unless a coarse centroid is stored.",
		HasCoordinates: hasCoords,
	})
}

func Cluster(cells []contract.MapCell) []contract.MapCell {
	groups := map[string][]contract.MapCell{}
	unlocated := []contract.MapCell{}
	order := []string{}
	for _, cell := range cells {
		if cell.Lon == nil || cell.Lat == nil {
			unlocated = append(unlocated, cell)
			continue
		}
		key := ClusterKey(*cell.Lon, *cell.Lat)
		if _, ok := groups[key]; !ok {
			order = append(order, key)
		}
		groups[key] = append(groups[key], cell)
	}
	out := append([]contract.MapCell{}, unlocated...)
	for _, key := range order {
		group := groups[key]
		if len(group) == 1 {
			out = append(out, group[0])
			continue
		}
		out = append(out, mergeCluster(key, group))
	}
	sortCells(out)
	return out
}

func matchingIncidents(items []contract.Incident, query Query) []contract.Incident {
	out := []contract.Incident{}
	for _, item := range items {
		item = contract.NormalizeIncident(item)
		if query.Service != "" && !containsFold(item.AffectedServices, query.Service) && !strings.EqualFold(item.Scope, query.Service) {
			continue
		}
		if query.Search != "" && !incidentMatchesSearch(item, query.Search) {
			continue
		}
		if !matchesParent(item, query.Parent) {
			continue
		}
		out = append(out, item)
	}
	return out
}

func buildCells(incidents []contract.Incident, query Query) []contract.MapCell {
	cells := []contract.MapCell{}
	grain := childGrain(query)
	if query.HasLayer(LayerGlobal) && showGlobal(query) {
		if cell, ok := worldCell(incidents); ok {
			cells = append(cells, cell)
		}
	}
	if query.HasLayer(LayerRegional) && showRegions(grain, query.Level) {
		cells = append(cells, uniqueSliceCells(incidents, LevelRegion, LayerRegional, regionValues, networkValues)...)
	}
	if query.HasLayer(LayerNetwork) && showNetworks(grain, query.Level) {
		cells = append(cells, uniqueSliceCells(incidents, LevelNetwork, LayerNetwork, networkValues, serviceValues)...)
	}
	if query.HasLayer(LayerService) && showServices(grain, query.Level) {
		cells = append(cells, uniqueSliceCells(incidents, LevelService, LayerService, serviceValues, func(contract.Incident) []string { return nil })...)
	}
	return cells
}

func childGrain(query Query) string {
	if query.Parent == "" || query.Parent == LevelWorld {
		return query.Level
	}
	kind, _ := splitID(query.Parent)
	switch kind {
	case LevelCountry:
		return LevelRegion
	case LevelRegion:
		return LevelNetwork
	case LevelNetwork:
		return LevelService
	case LevelService:
		return LevelService
	default:
		return query.Level
	}
}

func showGlobal(query Query) bool {
	return (query.Parent == "" || query.Parent == LevelWorld) && query.Level == LevelWorld
}

func showRegions(grain, level string) bool {
	return grain == LevelRegion || (grain == LevelWorld && level == LevelWorld)
}

func showNetworks(grain, level string) bool {
	return grain == LevelNetwork || (grain == LevelWorld && level == LevelWorld)
}

func showServices(grain, level string) bool {
	return grain == LevelService || (grain == LevelWorld && level == LevelWorld)
}

func worldCell(incidents []contract.Incident) (contract.MapCell, bool) {
	if len(incidents) == 0 {
		return contract.MapCell{}, false
	}
	samples := 0
	regions := map[string]struct{}{}
	for _, item := range incidents {
		samples += item.SampleCount
		for _, region := range item.Regions {
			if region != "" {
				regions[region] = struct{}{}
			}
		}
	}
	return contract.MapCell{
		ID:          "world",
		Level:       LevelWorld,
		Label:       "World",
		SampleCount: samples,
		Status:      cellStatus(samples),
		Summary:     "Observed sample count " + itoa(samples) + ". Country health is not assigned from incident labels.",
		Layer:       LayerGlobal,
		ChildCount:  len(regions),
	}, true
}

func uniqueSliceCells(
	incidents []contract.Incident,
	level string,
	layer string,
	values func(contract.Incident) []string,
	children func(contract.Incident) []string,
) []contract.MapCell {
	type bucket struct {
		samples  int
		children map[string]struct{}
	}
	buckets := map[string]*bucket{}
	order := []string{}
	for _, item := range incidents {
		for _, value := range values(item) {
			label := strings.TrimSpace(value)
			if label == "" {
				continue
			}
			if _, ok := buckets[label]; !ok {
				buckets[label] = &bucket{children: map[string]struct{}{}}
				order = append(order, label)
			}
			buckets[label].samples += item.SampleCount
			for _, child := range children(item) {
				child = strings.TrimSpace(child)
				if child != "" {
					buckets[label].children[child] = struct{}{}
				}
			}
		}
	}
	sort.Strings(order)
	out := make([]contract.MapCell, 0, len(order))
	for _, label := range order {
		bucket := buckets[label]
		out = append(out, contract.MapCell{
			ID:          level + ":" + label,
			Level:       level,
			Label:       label,
			SampleCount: bucket.samples,
			Status:      cellStatus(bucket.samples),
			Summary:     sliceSummary(level, label, bucket.samples),
			Layer:       layer,
			ChildCount:  len(bucket.children),
		})
	}
	return out
}

func sliceSummary(level, label string, samples int) string {
	switch level {
	case LevelRegion:
		return "Stored region label " + label + ". Observed sample count " + itoa(samples) + ". This is not a plotted location."
	case LevelNetwork:
		return "Stored network " + label + ". Observed sample count " + itoa(samples) + ". ASN health is not inferred beyond the sample."
	case LevelService:
		return "Stored service " + label + ". Observed sample count " + itoa(samples) + ". Open service intelligence for catalog context."
	default:
		return "Observed sample count " + itoa(samples) + "."
	}
}

func cellStatus(samples int) string {
	if samples <= 0 {
		return "not_measured"
	}
	return "insufficient_evidence"
}

func incidentRefs(items []contract.Incident) []contract.MapIncidentRef {
	out := make([]contract.MapIncidentRef, 0, len(items))
	for _, item := range items {
		region := ""
		if len(item.Regions) > 0 {
			region = item.Regions[0]
		}
		out = append(out, contract.MapIncidentRef{
			ID:           item.ID,
			Title:        item.Title,
			CoarseRegion: region,
			SampleCount:  item.SampleCount,
		})
	}
	return out
}

func filterCells(cells []contract.MapCell, query Query) []contract.MapCell {
	out := []contract.MapCell{}
	for _, cell := range cells {
		if query.Search != "" && !strings.Contains(strings.ToLower(cell.Label+" "+cell.Summary), strings.ToLower(query.Search)) {
			continue
		}
		if cell.Lon != nil && cell.Lat != nil && !InViewport(*cell.Lon, *cell.Lat, query) {
			continue
		}
		out = append(out, cell)
	}
	return out
}

func applyPrivacy(cell contract.MapCell) (contract.MapCell, bool) {
	if cell.Lon == nil || cell.Lat == nil {
		return cell, true
	}
	lon, lat, ok := CoarsePoint(*cell.Lon, *cell.Lat)
	if !ok {
		cell.Lon = nil
		cell.Lat = nil
		return cell, true
	}
	cell.Lon = &lon
	cell.Lat = &lat
	return cell, true
}

func mergeCluster(key string, group []contract.MapCell) contract.MapCell {
	samples := 0
	children := 0
	labels := []string{}
	var lon, lat float64
	layer := group[0].Layer
	level := group[0].Level
	for i, cell := range group {
		samples += cell.SampleCount
		children += cell.ChildCount
		labels = append(labels, cell.Label)
		if i == 0 && cell.Lon != nil && cell.Lat != nil {
			lon = *cell.Lon
			lat = *cell.Lat
		}
		if cell.Layer != layer {
			layer = LayerRegional
		}
		if cell.Level != level {
			level = LevelRegion
		}
	}
	return contract.MapCell{
		ID:          "cluster:" + key,
		Level:       level,
		Label:       itoa(len(group)) + " coarse locations",
		Lon:         &lon,
		Lat:         &lat,
		SampleCount: samples,
		Status:      cellStatus(samples),
		Summary:     "Clustered to 1-degree precision. Labels: " + strings.Join(labels, ", ") + ". Precise individual locations are not exposed.",
		Layer:       layer,
		ChildCount:  children,
	}
}

func matchesParent(item contract.Incident, parent string) bool {
	if parent == "" || parent == LevelWorld {
		return true
	}
	kind, value := splitID(parent)
	switch kind {
	case LevelCountry:
		return false
	case LevelRegion:
		return containsFold(item.Regions, value)
	case LevelNetwork:
		return containsFold(item.Networks, value)
	case LevelService:
		return containsFold(item.AffectedServices, value) || strings.EqualFold(item.Scope, value)
	default:
		return containsFold(item.Regions, parent) || containsFold(item.Networks, parent) ||
			containsFold(item.AffectedServices, parent) || strings.EqualFold(item.Scope, parent)
	}
}

func splitID(raw string) (string, string) {
	kind, value, ok := strings.Cut(raw, ":")
	if !ok {
		return "", raw
	}
	return strings.ToLower(kind), value
}

func regionValues(item contract.Incident) []string  { return item.Regions }
func networkValues(item contract.Incident) []string { return item.Networks }
func serviceValues(item contract.Incident) []string {
	if len(item.AffectedServices) > 0 {
		return item.AffectedServices
	}
	if item.Scope != "" {
		return []string{item.Scope}
	}
	return nil
}

func incidentMatchesSearch(item contract.Incident, raw string) bool {
	q := strings.ToLower(raw)
	hay := strings.ToLower(strings.Join([]string{
		item.Title,
		item.Scope,
		strings.Join(item.AffectedServices, " "),
		strings.Join(item.Regions, " "),
		strings.Join(item.Networks, " "),
	}, " "))
	return strings.Contains(hay, q)
}

func containsFold(values []string, want string) bool {
	for _, value := range values {
		if strings.EqualFold(value, want) {
			return true
		}
	}
	return false
}

func sortCells(cells []contract.MapCell) {
	sort.SliceStable(cells, func(i, j int) bool {
		if cells[i].Level != cells[j].Level {
			return levelRank(cells[i].Level) < levelRank(cells[j].Level)
		}
		if cells[i].SampleCount != cells[j].SampleCount {
			return cells[i].SampleCount > cells[j].SampleCount
		}
		return cells[i].Label < cells[j].Label
	})
}

func levelRank(level string) int {
	switch level {
	case LevelWorld:
		return 0
	case LevelCountry:
		return 1
	case LevelRegion:
		return 2
	case LevelNetwork:
		return 3
	case LevelService:
		return 4
	default:
		return 5
	}
}

func itoa(value int) string {
	return fmt.Sprintf("%d", value)
}
