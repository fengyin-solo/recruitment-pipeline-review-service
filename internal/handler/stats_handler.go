package handler

import (
	"net/http"

	"recruit/pkg/httpx"
)

func (s *Server) registerStatsRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /api/stats/departments", s.departmentStats)
	mux.HandleFunc("GET /api/stats/funnel", s.jobFunnelStats)
}

func (s *Server) departmentStats(w http.ResponseWriter, r *http.Request) {
	stats, err := s.svc.DepartmentStats()
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, stats)
}

func (s *Server) jobFunnelStats(w http.ResponseWriter, r *http.Request) {
	stats, err := s.svc.JobFunnelStats()
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, stats)
}
