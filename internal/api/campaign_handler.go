package api

import (
	"encoding/json"
	"net/http"
)

func (s *Server) health(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]any{"status": "ok", "phase": s.plant.Campaign.Snapshot().Phase, "started": s.plant.Status()["started"]})
}

func (s *Server) operations(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, s.plant.Status())
}

func (s *Server) equipment(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, s.plant.Snapshot().Column)
}

func writeJSON(w http.ResponseWriter, status int, value any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(value)
}
