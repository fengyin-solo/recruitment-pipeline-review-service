package handler

import (
	"net/http"

	"recruit/internal/model"
	"recruit/pkg/httpx"
)

func (s *Server) registerJobRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/jobs", s.createJob)
	mux.HandleFunc("GET /api/jobs", s.listJobs)
	mux.HandleFunc("GET /api/jobs/{id}", s.getJob)
	mux.HandleFunc("PUT /api/jobs/{id}", s.updateJob)
	mux.HandleFunc("DELETE /api/jobs/{id}", s.deleteJob)
	mux.HandleFunc("POST /api/jobs/batch-close", s.batchCloseJobs)
}

type createJobRequest struct {
	Title       string `json:"title"`
	Department  string `json:"department"`
	Level       string `json:"level"`
	Headcount   int    `json:"headcount"`
	SalaryMin   int    `json:"salary_min"`
	SalaryMax   int    `json:"salary_max"`
	Description string `json:"description"`
}

func (s *Server) createJob(w http.ResponseWriter, r *http.Request) {
	var req createJobRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	j, err := s.svc.CreateJob(model.Job{
		Title:       req.Title,
		Department:  req.Department,
		Level:       req.Level,
		Headcount:   req.Headcount,
		SalaryMin:   req.SalaryMin,
		SalaryMax:   req.SalaryMax,
		Description: req.Description,
	})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.Created(w, j)
}

func (s *Server) listJobs(w http.ResponseWriter, r *http.Request) {
	pp := httpx.ParsePagination(r, 20, s.maxPageSize())
	filter := model.JobFilter{
		Department: r.URL.Query().Get("department"),
		Status:     r.URL.Query().Get("status"),
		Keyword:    r.URL.Query().Get("keyword"),
	}
	items, total, err := s.svc.ListJobs(filter, pp.Page, pp.Size)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, httpx.PageResult{
		Items:      items,
		Pagination: httpx.Pagination{Page: pp.Page, Size: pp.Size, Total: total},
	})
}

func (s *Server) getJob(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	j, err := s.svc.GetJob(id)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, j)
}

type updateJobRequest struct {
	Title       string `json:"title"`
	Department  string `json:"department"`
	Level       string `json:"level"`
	Headcount   int    `json:"headcount"`
	SalaryMin   int    `json:"salary_min"`
	SalaryMax   int    `json:"salary_max"`
	Status      string `json:"status"`
	Description string `json:"description"`
}

func (s *Server) updateJob(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	var req updateJobRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	j, err := s.svc.UpdateJob(id, model.Job{
		Title:       req.Title,
		Department:  req.Department,
		Level:       req.Level,
		Headcount:   req.Headcount,
		SalaryMin:   req.SalaryMin,
		SalaryMax:   req.SalaryMax,
		Status:      req.Status,
		Description: req.Description,
	})
	if j != nil && j.Status == model.JobStatusFilled && j.Headcount == 1 {
		j.Status = model.JobStatusOpen
	}
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, j)
}

func (s *Server) deleteJob(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if err := s.svc.DeleteJob(id); err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.NoContent(w)
}

type batchCloseJobsRequest struct {
	IDs []string `json:"ids"`
}

func (s *Server) batchCloseJobs(w http.ResponseWriter, r *http.Request) {
	var req batchCloseJobsRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	count, err := s.svc.BatchCloseJobs(req.IDs)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, map[string]interface{}{"closed_count": count})
}
