package handler

import (
	"net/http"

	"recruit/internal/model"
	"recruit/pkg/httpx"
)

func (s *Server) registerCandidateRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/candidates", s.createCandidate)
	mux.HandleFunc("GET /api/candidates", s.listCandidates)
	mux.HandleFunc("GET /api/candidates/{id}", s.getCandidate)
	mux.HandleFunc("PUT /api/candidates/{id}", s.updateCandidate)
	mux.HandleFunc("DELETE /api/candidates/{id}", s.deleteCandidate)
}

type createCandidateRequest struct {
	Name              string   `json:"name"`
	Email             string   `json:"email"`
	Phone             string   `json:"phone"`
	YearsOfExperience int      `json:"years_of_experience"`
	Skills            []string `json:"skills"`
}

func (s *Server) createCandidate(w http.ResponseWriter, r *http.Request) {
	var req createCandidateRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	c, err := s.svc.CreateCandidate(model.Candidate{
		Name:              req.Name,
		Email:             req.Email,
		Phone:             req.Phone,
		YearsOfExperience: req.YearsOfExperience,
		Skills:            req.Skills,
	})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.Created(w, c)
}

func (s *Server) listCandidates(w http.ResponseWriter, r *http.Request) {
	pp := httpx.ParsePagination(r, 20, s.maxPageSize())
	filter := model.CandidateFilter{
		Keyword: r.URL.Query().Get("q"),
		Status:  r.URL.Query().Get("status"),
	}
	items, total, err := s.svc.ListCandidates(filter, pp.Page, pp.Size)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, httpx.PageResult{
		Items:      items,
		Pagination: httpx.Pagination{Page: pp.Page, Size: pp.Size, Total: total},
	})
}

func (s *Server) getCandidate(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	c, err := s.svc.GetCandidate(id)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, c)
}

type updateCandidateRequest struct {
	Name              string   `json:"name"`
	Email             string   `json:"email"`
	Phone             string   `json:"phone"`
	YearsOfExperience int      `json:"years_of_experience"`
	Skills            []string `json:"skills"`
	Status            string   `json:"status"`
}

func (s *Server) updateCandidate(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	var req updateCandidateRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	c, err := s.svc.UpdateCandidate(id, model.Candidate{
		Name:              req.Name,
		Email:             req.Email,
		Phone:             req.Phone,
		YearsOfExperience: req.YearsOfExperience,
		Skills:            req.Skills,
		Status:            req.Status,
	})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, c)
}

func (s *Server) deleteCandidate(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if err := s.svc.DeleteCandidate(id); err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.NoContent(w)
}
