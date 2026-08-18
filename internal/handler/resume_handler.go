package handler

import (
	"net/http"

	"recruit/internal/model"
	"recruit/pkg/httpx"
)

func (s *Server) registerResumeRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/resumes", s.createResume)
	mux.HandleFunc("GET /api/resumes", s.listResumes)
	mux.HandleFunc("GET /api/resumes/{id}", s.getResume)
	mux.HandleFunc("PUT /api/resumes/{id}", s.updateResume)
	mux.HandleFunc("DELETE /api/resumes/{id}", s.deleteResume)
}

type createResumeRequest struct {
	CandidateID    string `json:"candidate_id"`
	Summary        string `json:"summary"`
	Education      string `json:"education"`
	WorkExperience string `json:"work_experience"`
}

func (s *Server) createResume(w http.ResponseWriter, r *http.Request) {
	var req createResumeRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	res, err := s.svc.CreateResume(model.Resume{
		CandidateID:    req.CandidateID,
		Summary:        req.Summary,
		Education:      req.Education,
		WorkExperience: req.WorkExperience,
	})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.Created(w, res)
}

func (s *Server) listResumes(w http.ResponseWriter, r *http.Request) {
	pp := httpx.ParsePagination(r, 20, s.maxPageSize())
	filter := model.ResumeFilter{
		CandidateID: r.URL.Query().Get("candidate_id"),
	}
	items, total, err := s.svc.ListResumes(filter, pp.Page, pp.Size)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, httpx.PageResult{
		Items:      items,
		Pagination: httpx.Pagination{Page: pp.Page, Size: pp.Size, Total: total},
	})
}

func (s *Server) getResume(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	res, err := s.svc.GetResume(id)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, res)
}

type updateResumeRequest struct {
	Summary        string `json:"summary"`
	Education      string `json:"education"`
	WorkExperience string `json:"work_experience"`
}

func (s *Server) updateResume(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	var req updateResumeRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	res, err := s.svc.UpdateResume(id, model.Resume{
		Summary:        req.Summary,
		Education:      req.Education,
		WorkExperience: req.WorkExperience,
	})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, res)
}

func (s *Server) deleteResume(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if err := s.svc.DeleteResume(id); err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.NoContent(w)
}
