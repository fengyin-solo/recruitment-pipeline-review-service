package handler

import (
	"net/http"
	"time"

	"recruit/internal/model"
	"recruit/pkg/httpx"
)

func (s *Server) registerInterviewRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/interviews", s.createInterview)
	mux.HandleFunc("GET /api/interviews", s.listInterviews)
	mux.HandleFunc("GET /api/interviews/{id}", s.getInterview)
	mux.HandleFunc("PUT /api/interviews/{id}", s.updateInterview)
	mux.HandleFunc("DELETE /api/interviews/{id}", s.deleteInterview)
	mux.HandleFunc("POST /api/interviews/batch-update-status", s.batchUpdateInterviewStatus)
}

type createInterviewRequest struct {
	JobID       string    `json:"job_id"`
	CandidateID string    `json:"candidate_id"`
	Round       int       `json:"round"`
	Interviewer string    `json:"interviewer"`
	ScheduledAt time.Time `json:"scheduled_at"`
}

func (s *Server) createInterview(w http.ResponseWriter, r *http.Request) {
	var req createInterviewRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	iv, err := s.svc.CreateInterview(model.Interview{
		JobID:       req.JobID,
		CandidateID: req.CandidateID,
		Round:       req.Round,
		Interviewer: req.Interviewer,
		ScheduledAt: req.ScheduledAt,
	})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.Created(w, iv)
}

func (s *Server) listInterviews(w http.ResponseWriter, r *http.Request) {
	pp := httpx.ParsePagination(r, 20, s.maxPageSize())
	filter := model.InterviewFilter{
		JobID:       r.URL.Query().Get("job_id"),
		CandidateID: r.URL.Query().Get("candidate_id"),
		Status:      r.URL.Query().Get("status"),
	}
	items, total, err := s.svc.ListInterviews(filter, pp.Page, pp.Size)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, httpx.PageResult{
		Items:      items,
		Pagination: httpx.Pagination{Page: pp.Page, Size: pp.Size, Total: total},
	})
}

func (s *Server) getInterview(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	iv, err := s.svc.GetInterview(id)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, iv)
}

type updateInterviewRequest struct {
	Round       int       `json:"round"`
	Interviewer string    `json:"interviewer"`
	ScheduledAt time.Time `json:"scheduled_at"`
	Status      string    `json:"status"`
	Feedback    string    `json:"feedback"`
	Score       int       `json:"score"`
}

func (s *Server) updateInterview(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	var req updateInterviewRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	iv, err := s.svc.UpdateInterview(id, model.Interview{
		Round:       req.Round,
		Interviewer: req.Interviewer,
		ScheduledAt: req.ScheduledAt,
		Status:      req.Status,
		Feedback:    req.Feedback,
		Score:       req.Score,
	})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, iv)
}

func (s *Server) deleteInterview(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if err := s.svc.DeleteInterview(id); err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.NoContent(w)
}

type batchUpdateInterviewStatusRequest struct {
	IDs    []string `json:"ids"`
	Status string   `json:"status"`
}

func (s *Server) batchUpdateInterviewStatus(w http.ResponseWriter, r *http.Request) {
	var req batchUpdateInterviewStatusRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	count, err := s.svc.BatchUpdateInterviewStatus(req.IDs, req.Status)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, map[string]interface{}{"updated_count": count})
}
