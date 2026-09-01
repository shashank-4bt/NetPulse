package api

import (
	"encoding/json"
	"log/slog"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/shashank-4bt/NetPulse/netpulse-api/internal/config"
	"github.com/shashank-4bt/NetPulse/netpulse-api/internal/contract"
	"github.com/shashank-4bt/NetPulse/netpulse-api/internal/diagnostics"
	"github.com/shashank-4bt/NetPulse/netpulse-api/internal/storage"
)

type Server struct {
	Cfg         config.Config
	Log         *slog.Logger
	Diagnostics *diagnostics.Service
	Limiter     storage.RateLimiter
	StorageInfo map[string]string
}

func (s *Server) Handler() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /v1/health", s.health)
	mux.HandleFunc("POST /v1/diagnoses", s.createDiagnosis)
	mux.HandleFunc("GET /v1/diagnoses/{id}", s.getDiagnosis)
	mux.HandleFunc("GET /v1/services", s.listServices)
	mux.HandleFunc("GET /v1/services/{slug}", s.getService)
	mux.HandleFunc("GET /v1/incidents", s.listIncidents)
	return s.middleware(mux)
}

func (s *Server) middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if origin := s.Cfg.CORSOrigin; origin != "" {
			w.Header().Set("Access-Control-Allow-Origin", origin)
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
			w.Header().Set("Access-Control-Allow-Methods", "GET,POST,OPTIONS")
		}
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}

func (s *Server) health(w http.ResponseWriter, _ *http.Request) {
	write(w, http.StatusOK, contract.Envelope{
		OK: true,
		Health: &contract.Health{
			Status:  "ok",
			Version: s.Cfg.EngineVersion,
			Storage: s.StorageInfo,
		},
	})
}

func (s *Server) createDiagnosis(w http.ResponseWriter, r *http.Request) {
	ip, _, _ := net.SplitHostPort(r.RemoteAddr)
	if ip == "" {
		ip = r.RemoteAddr
	}
	if s.Limiter != nil && !s.Limiter.Allow("diag:"+ip, s.Cfg.RateLimitPerMin) {
		write(w, http.StatusTooManyRequests, contract.Envelope{Error: &contract.APIError{Code: "rate_limited", Message: "Too many diagnose requests"}})
		return
	}

	var body struct {
		Target string `json:"target"`
	}
	if err := json.NewDecoder(http.MaxBytesReader(w, r.Body, 8<<10)).Decode(&body); err != nil {
		write(w, http.StatusBadRequest, contract.Envelope{Error: &contract.APIError{Code: "validation_error", Message: "Expected JSON with a target"}})
		return
	}
	diag, apiErr, status := s.Diagnostics.Create(r.Context(), body.Target)
	if apiErr != nil {
		write(w, status, contract.Envelope{Error: apiErr})
		return
	}
	write(w, status, contract.Envelope{OK: true, Diagnosis: diag})
}

func (s *Server) getDiagnosis(w http.ResponseWriter, r *http.Request) {
	diag, apiErr, status := s.Diagnostics.Get(r.Context(), r.PathValue("id"))
	if apiErr != nil {
		write(w, status, contract.Envelope{Error: apiErr})
		return
	}
	write(w, status, contract.Envelope{OK: true, Diagnosis: diag})
}

func (s *Server) listServices(w http.ResponseWriter, r *http.Request) {
	items, apiErr, status := s.Diagnostics.ListServices(r.Context())
	if apiErr != nil {
		write(w, status, contract.Envelope{Error: apiErr})
		return
	}
	write(w, status, contract.Envelope{OK: true, Services: items})
}

func (s *Server) getService(w http.ResponseWriter, r *http.Request) {
	item, apiErr, status := s.Diagnostics.GetService(r.Context(), r.PathValue("slug"))
	if apiErr != nil {
		write(w, status, contract.Envelope{Error: apiErr})
		return
	}
	write(w, status, contract.Envelope{OK: true, Service: item})
}

func (s *Server) listIncidents(w http.ResponseWriter, r *http.Request) {
	items, apiErr, status := s.Diagnostics.ListIncidents(r.Context())
	if apiErr != nil {
		write(w, status, contract.Envelope{Error: apiErr})
		return
	}
	write(w, status, contract.Envelope{OK: true, Incidents: items})
}

func write(w http.ResponseWriter, status int, body contract.Envelope) {
	if body.Incidents == nil {
		body.Incidents = []contract.Incident{}
	}
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(body)
}

func ClientIP(r *http.Request) string {
	if forwarded := r.Header.Get("X-Forwarded-For"); forwarded != "" {
		return strings.TrimSpace(strings.Split(forwarded, ",")[0])
	}
	return r.RemoteAddr
}

func NewHTTPServer(addr string, handler http.Handler) *http.Server {
	return &http.Server{
		Addr:              addr,
		Handler:           handler,
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      15 * time.Second,
		IdleTimeout:       30 * time.Second,
	}
}
