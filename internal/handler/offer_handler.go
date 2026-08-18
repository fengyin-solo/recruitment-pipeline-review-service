package handler

import (
	"net/http"
	"time"

	"recruit/internal/model"
	"recruit/pkg/httpx"
)

func (s *Server) registerOfferRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/offers", s.createOffer)
	mux.HandleFunc("GET /api/offers", s.listOffers)
	mux.HandleFunc("GET /api/offers/{id}", s.getOffer)
	mux.HandleFunc("PUT /api/offers/{id}", s.updateOffer)
	mux.HandleFunc("DELETE /api/offers/{id}", s.deleteOffer)
}

type createOfferRequest struct {
	JobID       string    `json:"job_id"`
	CandidateID string    `json:"candidate_id"`
	Salary      int       `json:"salary"`
	ExpiresAt   time.Time `json:"expires_at"`
}

func (s *Server) createOffer(w http.ResponseWriter, r *http.Request) {
	var req createOfferRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	o, err := s.svc.CreateOffer(model.Offer{
		JobID:       req.JobID,
		CandidateID: req.CandidateID,
		Salary:      req.Salary,
		ExpiresAt:   req.ExpiresAt,
	})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.Created(w, o)
}

func (s *Server) listOffers(w http.ResponseWriter, r *http.Request) {
	pp := httpx.ParsePagination(r, 20, s.maxPageSize())
	filter := model.OfferFilter{
		JobID:       r.URL.Query().Get("job_id"),
		CandidateID: r.URL.Query().Get("candidate_id"),
		Status:      r.URL.Query().Get("status"),
	}
	items, total, err := s.svc.ListOffers(filter, pp.Page, pp.Size)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, httpx.PageResult{
		Items:      items,
		Pagination: httpx.Pagination{Page: pp.Page, Size: pp.Size, Total: total},
	})
}

func (s *Server) getOffer(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	o, err := s.svc.GetOffer(id)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, o)
}

type updateOfferRequest struct {
	Salary    int       `json:"salary"`
	Status    string    `json:"status"`
	ExpiresAt time.Time `json:"expires_at"`
}

func (s *Server) updateOffer(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	var req updateOfferRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	o, err := s.svc.UpdateOffer(id, model.Offer{
		Salary:    req.Salary,
		Status:    req.Status,
		ExpiresAt: req.ExpiresAt,
	})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, o)
}

func (s *Server) deleteOffer(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if err := s.svc.DeleteOffer(id); err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.NoContent(w)
}
