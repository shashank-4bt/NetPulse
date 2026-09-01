package memory

import (
	"context"
	"sync"
	"time"

	"github.com/shashank-4bt/NetPulse/netpulse-api/internal/contract"
	"github.com/shashank-4bt/NetPulse/netpulse-api/internal/storage"
)

type Store struct {
	mu           sync.Mutex
	diagnoses    map[string]storage.DiagnosisRecord
	incidents    []contract.Incident
	queue        []storage.Job
	rate         map[string][]time.Time
	cache        map[string]cacheItem
	measurements map[string][]contract.Measurement
}

type cacheItem struct {
	value     string
	expiresAt time.Time
}

func New() *Store {
	return &Store{
		diagnoses:    map[string]storage.DiagnosisRecord{},
		incidents:    []contract.Incident{},
		queue:        []storage.Job{},
		rate:         map[string][]time.Time{},
		cache:        map[string]cacheItem{},
		measurements: map[string][]contract.Measurement{},
	}
}

func (s *Store) CreateDiagnosis(_ context.Context, rec storage.DiagnosisRecord) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.diagnoses[rec.Diagnosis.ID] = rec
	return nil
}

func (s *Store) GetDiagnosis(_ context.Context, id string) (*storage.DiagnosisRecord, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	rec, ok := s.diagnoses[id]
	if !ok {
		return nil, nil
	}
	copy := rec
	return &copy, nil
}

func (s *Store) UpdateDiagnosis(_ context.Context, rec storage.DiagnosisRecord) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.diagnoses[rec.Diagnosis.ID] = rec
	return nil
}

func (s *Store) ListIncidents(_ context.Context) ([]contract.Incident, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	out := make([]contract.Incident, len(s.incidents))
	copy(out, s.incidents)
	return out, nil
}

func (s *Store) GetIncident(_ context.Context, id string) (*contract.Incident, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, item := range s.incidents {
		if item.ID == id {
			copy := contract.NormalizeIncident(item)
			return &copy, nil
		}
	}
	return nil, nil
}

func (s *Store) ReplaceIncidents(items []contract.Incident) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.incidents = append([]contract.Incident{}, items...)
}

func (s *Store) ListServices(_ context.Context) ([]contract.Service, error) {
	return storage.ServiceCatalog(), nil
}

func (s *Store) GetService(_ context.Context, slug string) (*contract.Service, error) {
	for _, service := range storage.ServiceCatalog() {
		if service.Slug == slug {
			copy := service
			return &copy, nil
		}
	}
	return nil, nil
}

func (s *Store) Backend() string { return "memory" }

func (s *Store) Record(_ context.Context, diagnosisID string, measurements []contract.Measurement) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.measurements[diagnosisID] = measurements
	return nil
}

func (s *Store) Enqueue(_ context.Context, job storage.Job) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.queue = append(s.queue, job)
	return nil
}

func (s *Store) Dequeue(_ context.Context) (storage.Job, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if len(s.queue) == 0 {
		return storage.Job{}, false
	}
	job := s.queue[0]
	s.queue = s.queue[1:]
	return job, true
}

func (s *Store) Allow(key string, limitPerMin int) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	now := time.Now()
	window := now.Add(-time.Minute)
	events := s.rate[key]
	kept := events[:0]
	for _, event := range events {
		if event.After(window) {
			kept = append(kept, event)
		}
	}
	if len(kept) >= limitPerMin {
		s.rate[key] = kept
		return false
	}
	s.rate[key] = append(kept, now)
	return true
}

func (s *Store) Get(key string) (string, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	item, ok := s.cache[key]
	if !ok || time.Now().After(item.expiresAt) {
		delete(s.cache, key)
		return "", false
	}
	return item.value, true
}

func (s *Store) Set(key, value string, ttl time.Duration) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.cache[key] = cacheItem{value: value, expiresAt: time.Now().Add(ttl)}
}
