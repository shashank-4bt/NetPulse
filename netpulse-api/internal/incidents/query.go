package incidents

import (
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/shashank-4bt/NetPulse/netpulse-api/internal/contract"
)

type Query struct {
	Service  string
	Region   string
	Network  string
	Severity string
	Status   string
	Search   string
	Sort     string
	Time     string
	Page     int
	PageSize int
	Now      time.Time
}

type Result struct {
	Items []contract.Incident
	Page  contract.Page
}

func ParseQuery(values map[string]string, now time.Time) Query {
	page, _ := strconv.Atoi(values["page"])
	size, _ := strconv.Atoi(values["pageSize"])
	if page < 1 {
		page = 1
	}
	if size < 1 || size > 100 {
		size = 20
	}
	return Query{
		Service:  strings.TrimSpace(values["service"]),
		Region:   strings.TrimSpace(values["region"]),
		Network:  strings.TrimSpace(values["network"]),
		Severity: strings.TrimSpace(values["severity"]),
		Status:   strings.TrimSpace(values["status"]),
		Search:   strings.TrimSpace(values["q"]),
		Sort:     strings.TrimSpace(values["sort"]),
		Time:     strings.TrimSpace(values["time"]),
		Page:     page,
		PageSize: size,
		Now:      now,
	}
}

func Filter(items []contract.Incident, query Query) Result {
	now := query.Now
	if now.IsZero() {
		now = time.Now().UTC()
	}
	filtered := make([]contract.Incident, 0, len(items))
	for _, item := range items {
		item = contract.NormalizeIncident(item)
		if query.Service != "" && !containsFold(item.AffectedServices, query.Service) && !strings.EqualFold(item.Scope, query.Service) {
			continue
		}
		if query.Region != "" && !containsFold(item.Regions, query.Region) {
			continue
		}
		if query.Network != "" && !containsFold(item.Networks, query.Network) {
			continue
		}
		if query.Severity != "" && !strings.EqualFold(item.Severity, query.Severity) {
			continue
		}
		if query.Status != "" && !strings.EqualFold(item.Status, query.Status) {
			continue
		}
		if query.Search != "" && !matchesSearch(item, query.Search) {
			continue
		}
		if !inTimeWindow(item.StartedAt, query.Time, now) {
			continue
		}
		filtered = append(filtered, item)
	}

	sortIncidents(filtered, query.Sort)

	total := len(filtered)
	start := (query.Page - 1) * query.PageSize
	if start > total {
		start = total
	}
	end := start + query.PageSize
	if end > total {
		end = total
	}
	return Result{
		Items: filtered[start:end],
		Page: contract.Page{
			Number: query.Page,
			Size:   query.PageSize,
			Total:  total,
		},
	}
}

func matchesSearch(item contract.Incident, raw string) bool {
	q := strings.ToLower(raw)
	hay := strings.ToLower(strings.Join([]string{
		item.Title,
		item.Scope,
		item.Status,
		item.Severity,
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

func inTimeWindow(startedAt, window string, now time.Time) bool {
	if window == "" || window == "all" {
		return true
	}
	started, err := time.Parse(time.RFC3339, startedAt)
	if err != nil {
		return false
	}
	switch window {
	case "24h":
		return !started.Before(now.Add(-24 * time.Hour))
	case "7d":
		return !started.Before(now.Add(-7 * 24 * time.Hour))
	case "30d":
		return !started.Before(now.Add(-30 * 24 * time.Hour))
	default:
		return true
	}
}

func sortIncidents(items []contract.Incident, key string) {
	sort.SliceStable(items, func(i, j int) bool {
		a, b := items[i], items[j]
		switch key {
		case "started_asc":
			return a.StartedAt < b.StartedAt
		case "severity":
			return severityRank(a.Severity) > severityRank(b.Severity)
		case "status":
			return a.Status < b.Status
		case "updated_desc":
			return a.LastUpdatedAt > b.LastUpdatedAt
		default:
			return a.StartedAt > b.StartedAt
		}
	})
}

func severityRank(value string) int {
	switch strings.ToLower(value) {
	case "critical":
		return 4
	case "high":
		return 3
	case "moderate":
		return 2
	case "informational":
		return 1
	default:
		return 0
	}
}
