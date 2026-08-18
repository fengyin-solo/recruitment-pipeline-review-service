package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"recruit/internal/config"
	"recruit/internal/model"
	"recruit/internal/service"
	"recruit/internal/store"
	"recruit/pkg/httpx"
	"recruit/pkg/logger"
)

func newTestServer() *Server {
	cfg := &config.Config{MaxPageSize: 100}
	log := logger.NewLevel(logger.LevelError)
	st := store.NewMemoryStore()
	svc := service.New(st, log, cfg)
	return NewServer(svc, log, cfg)
}

func parseResponse(t *testing.T, w *httptest.ResponseRecorder) httpx.Response {
	var resp httpx.Response
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("unmarshal response: %v", err)
	}
	return resp
}

func TestHandlerCreateJob(t *testing.T) {
	srv := newTestServer()
	body := `{"title":"Go Dev","department":"Tech","level":"P5","headcount":2,"salary_min":20000,"salary_max":40000}`
	req := httptest.NewRequest(http.MethodPost, "/api/jobs", strings.NewReader(body))
	w := httptest.NewRecorder()
	srv.Routes().ServeHTTP(w, req)
	if w.Code != http.StatusCreated {
		t.Errorf("expect 201, got %d: %s", w.Code, w.Body.String())
	}
	var resp struct {
		Code int       `json:"code"`
		Data model.Job `json:"data"`
	}
	json.Unmarshal(w.Body.Bytes(), &resp)
	if resp.Data.Title != "Go Dev" {
		t.Errorf("title mismatch: %s", resp.Data.Title)
	}
}

func TestHandlerCreateJobValidation(t *testing.T) {
	srv := newTestServer()
	body := `{"title":"","department":"Tech","level":"P5","headcount":2}`
	req := httptest.NewRequest(http.MethodPost, "/api/jobs", strings.NewReader(body))
	w := httptest.NewRecorder()
	srv.Routes().ServeHTTP(w, req)
	if w.Code != http.StatusBadRequest {
		t.Errorf("expect 400, got %d", w.Code)
	}
}

func TestHandlerListJobs(t *testing.T) {
	srv := newTestServer()
	srv.svc.CreateJob(model.Job{Title: "Go Dev", Department: "Tech", Level: "P5", Headcount: 1})
	srv.svc.CreateJob(model.Job{Title: "Java Dev", Department: "Tech", Level: "P5", Headcount: 1})
	req := httptest.NewRequest(http.MethodGet, "/api/jobs?page=1&size=10", nil)
	w := httptest.NewRecorder()
	srv.Routes().ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("expect 200, got %d", w.Code)
	}
	var resp struct {
		Code int `json:"code"`
		Data struct {
			Items      []model.Job      `json:"items"`
			Pagination httpx.Pagination `json:"pagination"`
		} `json:"data"`
	}
	json.Unmarshal(w.Body.Bytes(), &resp)
	if resp.Data.Pagination.Total != 2 {
		t.Errorf("expect total 2, got %d", resp.Data.Pagination.Total)
	}
	if len(resp.Data.Items) != 2 {
		t.Errorf("expect 2 items, got %d", len(resp.Data.Items))
	}
}

func TestHandlerListJobsFilter(t *testing.T) {
	srv := newTestServer()
	srv.svc.CreateJob(model.Job{Title: "Go Dev", Department: "Tech", Level: "P5", Headcount: 1, Status: model.JobStatusOpen})
	srv.svc.CreateJob(model.Job{Title: "Java Dev", Department: "Tech", Level: "P5", Headcount: 1, Status: model.JobStatusClosed})
	req := httptest.NewRequest(http.MethodGet, "/api/jobs?status=open", nil)
	w := httptest.NewRecorder()
	srv.Routes().ServeHTTP(w, req)
	var resp struct {
		Code int `json:"code"`
		Data struct {
			Items []model.Job `json:"items"`
			Total int         `json:"pagination"`
		} `json:"data"`
	}
	json.Unmarshal(w.Body.Bytes(), &resp)
	var pageResp struct {
		Code int `json:"code"`
		Data struct {
			Items      []model.Job      `json:"items"`
			Pagination httpx.Pagination `json:"pagination"`
		} `json:"data"`
	}
	json.Unmarshal(w.Body.Bytes(), &pageResp)
	if pageResp.Data.Pagination.Total != 1 {
		t.Errorf("expect total 1, got %d", pageResp.Data.Pagination.Total)
	}
}

func TestHandlerGetJob(t *testing.T) {
	srv := newTestServer()
	j, _ := srv.svc.CreateJob(model.Job{Title: "Go Dev", Department: "Tech", Level: "P5", Headcount: 1})
	req := httptest.NewRequest(http.MethodGet, "/api/jobs/"+j.ID, nil)
	w := httptest.NewRecorder()
	srv.Routes().ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("expect 200, got %d", w.Code)
	}
	var resp struct {
		Code int       `json:"code"`
		Data model.Job `json:"data"`
	}
	json.Unmarshal(w.Body.Bytes(), &resp)
	if resp.Data.ID != j.ID {
		t.Errorf("id mismatch")
	}
}

func TestHandlerGetJobNotFound(t *testing.T) {
	srv := newTestServer()
	req := httptest.NewRequest(http.MethodGet, "/api/jobs/notfound", nil)
	w := httptest.NewRecorder()
	srv.Routes().ServeHTTP(w, req)
	if w.Code != http.StatusNotFound {
		t.Errorf("expect 404, got %d", w.Code)
	}
}

func TestHandlerUpdateJob(t *testing.T) {
	srv := newTestServer()
	j, _ := srv.svc.CreateJob(model.Job{Title: "Go Dev", Department: "Tech", Level: "P5", Headcount: 1})
	body := `{"title":"Senior Go Dev"}`
	req := httptest.NewRequest(http.MethodPut, "/api/jobs/"+j.ID, strings.NewReader(body))
	w := httptest.NewRecorder()
	srv.Routes().ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("expect 200, got %d: %s", w.Code, w.Body.String())
	}
	var resp struct {
		Code int       `json:"code"`
		Data model.Job `json:"data"`
	}
	json.Unmarshal(w.Body.Bytes(), &resp)
	if resp.Data.Title != "Senior Go Dev" {
		t.Errorf("title mismatch: %s", resp.Data.Title)
	}
}

func TestHandlerDeleteJob(t *testing.T) {
	srv := newTestServer()
	j, _ := srv.svc.CreateJob(model.Job{Title: "Go Dev", Department: "Tech", Level: "P5", Headcount: 1})
	req := httptest.NewRequest(http.MethodDelete, "/api/jobs/"+j.ID, nil)
	w := httptest.NewRecorder()
	srv.Routes().ServeHTTP(w, req)
	if w.Code != http.StatusNoContent {
		t.Errorf("expect 204, got %d", w.Code)
	}
}

func TestHandlerBatchCloseJobs(t *testing.T) {
	srv := newTestServer()
	j1, _ := srv.svc.CreateJob(model.Job{Title: "Dev1", Department: "Tech", Level: "P5", Headcount: 1})
	j2, _ := srv.svc.CreateJob(model.Job{Title: "Dev2", Department: "Tech", Level: "P5", Headcount: 1})
	body := `{"ids":["` + j1.ID + `","` + j2.ID + `"]}`
	req := httptest.NewRequest(http.MethodPost, "/api/jobs/batch-close", strings.NewReader(body))
	w := httptest.NewRecorder()
	srv.Routes().ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("expect 200, got %d: %s", w.Code, w.Body.String())
	}
}

func TestHandlerCreateCandidate(t *testing.T) {
	srv := newTestServer()
	body := `{"name":"Alice","email":"alice@example.com","phone":"123","years_of_experience":3}`
	req := httptest.NewRequest(http.MethodPost, "/api/candidates", strings.NewReader(body))
	w := httptest.NewRecorder()
	srv.Routes().ServeHTTP(w, req)
	if w.Code != http.StatusCreated {
		t.Errorf("expect 201, got %d: %s", w.Code, w.Body.String())
	}
}

func TestHandlerCreateCandidateDuplicate(t *testing.T) {
	srv := newTestServer()
	body := `{"name":"Alice","email":"alice@example.com","phone":"123"}`
	req := httptest.NewRequest(http.MethodPost, "/api/candidates", strings.NewReader(body))
	w := httptest.NewRecorder()
	srv.Routes().ServeHTTP(w, req)
	req2 := httptest.NewRequest(http.MethodPost, "/api/candidates", strings.NewReader(body))
	w2 := httptest.NewRecorder()
	srv.Routes().ServeHTTP(w2, req2)
	if w2.Code != http.StatusBadRequest {
		t.Errorf("expect 400, got %d", w2.Code)
	}
}

func TestHandlerListCandidates(t *testing.T) {
	srv := newTestServer()
	srv.svc.CreateCandidate(model.Candidate{Name: "Alice", Email: "alice@example.com", Phone: "123"})
	srv.svc.CreateCandidate(model.Candidate{Name: "Bob", Email: "bob@example.com", Phone: "456"})
	req := httptest.NewRequest(http.MethodGet, "/api/candidates?keyword=ali", nil)
	w := httptest.NewRecorder()
	srv.Routes().ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("expect 200, got %d", w.Code)
	}
	var resp struct {
		Code int `json:"code"`
		Data struct {
			Items      []model.Candidate `json:"items"`
			Pagination httpx.Pagination  `json:"pagination"`
		} `json:"data"`
	}
	json.Unmarshal(w.Body.Bytes(), &resp)
	if resp.Data.Pagination.Total != 1 {
		t.Errorf("expect total 1, got %d", resp.Data.Pagination.Total)
	}
}

func TestHandlerGetCandidate(t *testing.T) {
	srv := newTestServer()
	c, _ := srv.svc.CreateCandidate(model.Candidate{Name: "Alice", Email: "alice@example.com", Phone: "123"})
	req := httptest.NewRequest(http.MethodGet, "/api/candidates/"+c.ID, nil)
	w := httptest.NewRecorder()
	srv.Routes().ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("expect 200, got %d", w.Code)
	}
}

func TestHandlerUpdateCandidate(t *testing.T) {
	srv := newTestServer()
	c, _ := srv.svc.CreateCandidate(model.Candidate{Name: "Alice", Email: "alice@example.com", Phone: "123"})
	body := `{"name":"Alice Updated"}`
	req := httptest.NewRequest(http.MethodPut, "/api/candidates/"+c.ID, strings.NewReader(body))
	w := httptest.NewRecorder()
	srv.Routes().ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("expect 200, got %d: %s", w.Code, w.Body.String())
	}
}

func TestHandlerDeleteCandidate(t *testing.T) {
	srv := newTestServer()
	c, _ := srv.svc.CreateCandidate(model.Candidate{Name: "Alice", Email: "alice@example.com", Phone: "123"})
	req := httptest.NewRequest(http.MethodDelete, "/api/candidates/"+c.ID, nil)
	w := httptest.NewRecorder()
	srv.Routes().ServeHTTP(w, req)
	if w.Code != http.StatusNoContent {
		t.Errorf("expect 204, got %d", w.Code)
	}
}

func TestHandlerCreateResume(t *testing.T) {
	srv := newTestServer()
	c, _ := srv.svc.CreateCandidate(model.Candidate{Name: "Bob", Email: "bob@example.com", Phone: "789"})
	body := `{"candidate_id":"` + c.ID + `","summary":"good","education":"BS"}`
	req := httptest.NewRequest(http.MethodPost, "/api/resumes", strings.NewReader(body))
	w := httptest.NewRecorder()
	srv.Routes().ServeHTTP(w, req)
	if w.Code != http.StatusCreated {
		t.Errorf("expect 201, got %d: %s", w.Code, w.Body.String())
	}
}

func TestHandlerCreateResumeForeignKey(t *testing.T) {
	srv := newTestServer()
	body := `{"candidate_id":"notfound","summary":"good","education":"BS"}`
	req := httptest.NewRequest(http.MethodPost, "/api/resumes", strings.NewReader(body))
	w := httptest.NewRecorder()
	srv.Routes().ServeHTTP(w, req)
	if w.Code != http.StatusBadRequest {
		t.Errorf("expect 400, got %d", w.Code)
	}
}

func TestHandlerListResumes(t *testing.T) {
	srv := newTestServer()
	c, _ := srv.svc.CreateCandidate(model.Candidate{Name: "Bob", Email: "bob@example.com", Phone: "789"})
	srv.svc.CreateResume(model.Resume{CandidateID: c.ID, Summary: "s", Education: "BS"})
	req := httptest.NewRequest(http.MethodGet, "/api/resumes?candidate_id="+c.ID, nil)
	w := httptest.NewRecorder()
	srv.Routes().ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("expect 200, got %d", w.Code)
	}
}

func TestHandlerGetResume(t *testing.T) {
	srv := newTestServer()
	c, _ := srv.svc.CreateCandidate(model.Candidate{Name: "Bob", Email: "bob@example.com", Phone: "789"})
	r, _ := srv.svc.CreateResume(model.Resume{CandidateID: c.ID, Summary: "s", Education: "BS"})
	req := httptest.NewRequest(http.MethodGet, "/api/resumes/"+r.ID, nil)
	w := httptest.NewRecorder()
	srv.Routes().ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("expect 200, got %d", w.Code)
	}
}

func TestHandlerUpdateResume(t *testing.T) {
	srv := newTestServer()
	c, _ := srv.svc.CreateCandidate(model.Candidate{Name: "Bob", Email: "bob@example.com", Phone: "789"})
	r, _ := srv.svc.CreateResume(model.Resume{CandidateID: c.ID, Summary: "s", Education: "BS"})
	body := `{"summary":"updated"}`
	req := httptest.NewRequest(http.MethodPut, "/api/resumes/"+r.ID, strings.NewReader(body))
	w := httptest.NewRecorder()
	srv.Routes().ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("expect 200, got %d: %s", w.Code, w.Body.String())
	}
}

func TestHandlerDeleteResume(t *testing.T) {
	srv := newTestServer()
	c, _ := srv.svc.CreateCandidate(model.Candidate{Name: "Bob", Email: "bob@example.com", Phone: "789"})
	r, _ := srv.svc.CreateResume(model.Resume{CandidateID: c.ID, Summary: "s", Education: "BS"})
	req := httptest.NewRequest(http.MethodDelete, "/api/resumes/"+r.ID, nil)
	w := httptest.NewRecorder()
	srv.Routes().ServeHTTP(w, req)
	if w.Code != http.StatusNoContent {
		t.Errorf("expect 204, got %d", w.Code)
	}
}

func TestHandlerCreateInterview(t *testing.T) {
	srv := newTestServer()
	j, _ := srv.svc.CreateJob(model.Job{Title: "Dev", Department: "Tech", Level: "P5", Headcount: 1})
	c, _ := srv.svc.CreateCandidate(model.Candidate{Name: "Tom", Email: "tom@example.com", Phone: "000"})
	body := `{"job_id":"` + j.ID + `","candidate_id":"` + c.ID + `","round":1,"interviewer":"M1","scheduled_at":"` + time.Now().Add(time.Hour).Format(time.RFC3339) + `"}`
	req := httptest.NewRequest(http.MethodPost, "/api/interviews", strings.NewReader(body))
	w := httptest.NewRecorder()
	srv.Routes().ServeHTTP(w, req)
	if w.Code != http.StatusCreated {
		t.Errorf("expect 201, got %d: %s", w.Code, w.Body.String())
	}
}

func TestHandlerCreateInterviewForeignKey(t *testing.T) {
	srv := newTestServer()
	body := `{"job_id":"notfound","candidate_id":"notfound","round":1,"interviewer":"M1","scheduled_at":"` + time.Now().Add(time.Hour).Format(time.RFC3339) + `"}`
	req := httptest.NewRequest(http.MethodPost, "/api/interviews", strings.NewReader(body))
	w := httptest.NewRecorder()
	srv.Routes().ServeHTTP(w, req)
	if w.Code != http.StatusBadRequest {
		t.Errorf("expect 400, got %d", w.Code)
	}
}

func TestHandlerListInterviews(t *testing.T) {
	srv := newTestServer()
	j, _ := srv.svc.CreateJob(model.Job{Title: "Dev", Department: "Tech", Level: "P5", Headcount: 1})
	c, _ := srv.svc.CreateCandidate(model.Candidate{Name: "Tom", Email: "tom@example.com", Phone: "000"})
	srv.svc.CreateInterview(model.Interview{JobID: j.ID, CandidateID: c.ID, Round: 1, Interviewer: "M1", ScheduledAt: time.Now().Add(time.Hour)})
	req := httptest.NewRequest(http.MethodGet, "/api/interviews?job_id="+j.ID, nil)
	w := httptest.NewRecorder()
	srv.Routes().ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("expect 200, got %d", w.Code)
	}
}

func TestHandlerGetInterview(t *testing.T) {
	srv := newTestServer()
	j, _ := srv.svc.CreateJob(model.Job{Title: "Dev", Department: "Tech", Level: "P5", Headcount: 1})
	c, _ := srv.svc.CreateCandidate(model.Candidate{Name: "Tom", Email: "tom@example.com", Phone: "000"})
	iv, _ := srv.svc.CreateInterview(model.Interview{JobID: j.ID, CandidateID: c.ID, Round: 1, Interviewer: "M1", ScheduledAt: time.Now().Add(time.Hour)})
	req := httptest.NewRequest(http.MethodGet, "/api/interviews/"+iv.ID, nil)
	w := httptest.NewRecorder()
	srv.Routes().ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("expect 200, got %d", w.Code)
	}
}

func TestHandlerUpdateInterviewStateMachine(t *testing.T) {
	srv := newTestServer()
	j, _ := srv.svc.CreateJob(model.Job{Title: "Dev", Department: "Tech", Level: "P5", Headcount: 1})
	c, _ := srv.svc.CreateCandidate(model.Candidate{Name: "Tom", Email: "tom@example.com", Phone: "000"})
	iv, _ := srv.svc.CreateInterview(model.Interview{JobID: j.ID, CandidateID: c.ID, Round: 1, Interviewer: "M1", ScheduledAt: time.Now().Add(time.Hour)})
	body := `{"status":"completed","score":85,"feedback":"good"}`
	req := httptest.NewRequest(http.MethodPut, "/api/interviews/"+iv.ID, strings.NewReader(body))
	w := httptest.NewRecorder()
	srv.Routes().ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("expect 200, got %d: %s", w.Code, w.Body.String())
	}
	var resp struct {
		Code int             `json:"code"`
		Data model.Interview `json:"data"`
	}
	json.Unmarshal(w.Body.Bytes(), &resp)
	if resp.Data.Status != model.InterviewStatusCompleted {
		t.Errorf("expect completed, got %s", resp.Data.Status)
	}
}

func TestHandlerUpdateInterviewInvalidTransition(t *testing.T) {
	srv := newTestServer()
	j, _ := srv.svc.CreateJob(model.Job{Title: "Dev", Department: "Tech", Level: "P5", Headcount: 1})
	c, _ := srv.svc.CreateCandidate(model.Candidate{Name: "Tom", Email: "tom@example.com", Phone: "000"})
	iv, _ := srv.svc.CreateInterview(model.Interview{JobID: j.ID, CandidateID: c.ID, Round: 1, Interviewer: "M1", ScheduledAt: time.Now().Add(time.Hour)})
	body1 := `{"status":"cancelled"}`
	req1 := httptest.NewRequest(http.MethodPut, "/api/interviews/"+iv.ID, strings.NewReader(body1))
	w1 := httptest.NewRecorder()
	srv.Routes().ServeHTTP(w1, req1)
	if w1.Code != http.StatusOK {
		t.Fatalf("expect 200 for cancel, got %d", w1.Code)
	}
	body2 := `{"status":"completed"}`
	req2 := httptest.NewRequest(http.MethodPut, "/api/interviews/"+iv.ID, strings.NewReader(body2))
	w2 := httptest.NewRecorder()
	srv.Routes().ServeHTTP(w2, req2)
	if w2.Code != http.StatusBadRequest {
		t.Errorf("expect 400 for invalid transition, got %d", w2.Code)
	}
}

func TestHandlerDeleteInterview(t *testing.T) {
	srv := newTestServer()
	j, _ := srv.svc.CreateJob(model.Job{Title: "Dev", Department: "Tech", Level: "P5", Headcount: 1})
	c, _ := srv.svc.CreateCandidate(model.Candidate{Name: "Tom", Email: "tom@example.com", Phone: "000"})
	iv, _ := srv.svc.CreateInterview(model.Interview{JobID: j.ID, CandidateID: c.ID, Round: 1, Interviewer: "M1", ScheduledAt: time.Now().Add(time.Hour)})
	req := httptest.NewRequest(http.MethodDelete, "/api/interviews/"+iv.ID, nil)
	w := httptest.NewRecorder()
	srv.Routes().ServeHTTP(w, req)
	if w.Code != http.StatusNoContent {
		t.Errorf("expect 204, got %d", w.Code)
	}
}

func TestHandlerBatchUpdateInterviewStatus(t *testing.T) {
	srv := newTestServer()
	j, _ := srv.svc.CreateJob(model.Job{Title: "Dev", Department: "Tech", Level: "P5", Headcount: 1})
	c, _ := srv.svc.CreateCandidate(model.Candidate{Name: "Tom", Email: "tom@example.com", Phone: "000"})
	iv1, _ := srv.svc.CreateInterview(model.Interview{JobID: j.ID, CandidateID: c.ID, Round: 1, Interviewer: "M1", ScheduledAt: time.Now().Add(time.Hour)})
	iv2, _ := srv.svc.CreateInterview(model.Interview{JobID: j.ID, CandidateID: c.ID, Round: 2, Interviewer: "M2", ScheduledAt: time.Now().Add(2 * time.Hour)})
	body := `{"ids":["` + iv1.ID + `","` + iv2.ID + `"],"status":"completed"}`
	req := httptest.NewRequest(http.MethodPost, "/api/interviews/batch-update-status", strings.NewReader(body))
	w := httptest.NewRecorder()
	srv.Routes().ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("expect 200, got %d: %s", w.Code, w.Body.String())
	}
}

func TestHandlerCreateOffer(t *testing.T) {
	srv := newTestServer()
	j, _ := srv.svc.CreateJob(model.Job{Title: "Dev", Department: "Tech", Level: "P5", Headcount: 1})
	c, _ := srv.svc.CreateCandidate(model.Candidate{Name: "Tom", Email: "tom@example.com", Phone: "000"})
	body := `{"job_id":"` + j.ID + `","candidate_id":"` + c.ID + `","salary":30000,"expires_at":"` + time.Now().Add(7*24*time.Hour).Format(time.RFC3339) + `"}`
	req := httptest.NewRequest(http.MethodPost, "/api/offers", strings.NewReader(body))
	w := httptest.NewRecorder()
	srv.Routes().ServeHTTP(w, req)
	if w.Code != http.StatusCreated {
		t.Errorf("expect 201, got %d: %s", w.Code, w.Body.String())
	}
}

func TestHandlerCreateOfferForeignKey(t *testing.T) {
	srv := newTestServer()
	body := `{"job_id":"notfound","candidate_id":"notfound","salary":30000,"expires_at":"` + time.Now().Add(7*24*time.Hour).Format(time.RFC3339) + `"}`
	req := httptest.NewRequest(http.MethodPost, "/api/offers", strings.NewReader(body))
	w := httptest.NewRecorder()
	srv.Routes().ServeHTTP(w, req)
	if w.Code != http.StatusBadRequest {
		t.Errorf("expect 400, got %d", w.Code)
	}
}

func TestHandlerListOffers(t *testing.T) {
	srv := newTestServer()
	j, _ := srv.svc.CreateJob(model.Job{Title: "Dev", Department: "Tech", Level: "P5", Headcount: 1})
	c, _ := srv.svc.CreateCandidate(model.Candidate{Name: "Tom", Email: "tom@example.com", Phone: "000"})
	srv.svc.CreateOffer(model.Offer{JobID: j.ID, CandidateID: c.ID, Salary: 30000, ExpiresAt: time.Now().Add(7 * 24 * time.Hour)})
	req := httptest.NewRequest(http.MethodGet, "/api/offers?job_id="+j.ID, nil)
	w := httptest.NewRecorder()
	srv.Routes().ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("expect 200, got %d", w.Code)
	}
}

func TestHandlerGetOffer(t *testing.T) {
	srv := newTestServer()
	j, _ := srv.svc.CreateJob(model.Job{Title: "Dev", Department: "Tech", Level: "P5", Headcount: 1})
	c, _ := srv.svc.CreateCandidate(model.Candidate{Name: "Tom", Email: "tom@example.com", Phone: "000"})
	o, _ := srv.svc.CreateOffer(model.Offer{JobID: j.ID, CandidateID: c.ID, Salary: 30000, ExpiresAt: time.Now().Add(7 * 24 * time.Hour)})
	req := httptest.NewRequest(http.MethodGet, "/api/offers/"+o.ID, nil)
	w := httptest.NewRecorder()
	srv.Routes().ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("expect 200, got %d", w.Code)
	}
}

func TestHandlerUpdateOfferStateMachine(t *testing.T) {
	srv := newTestServer()
	j, _ := srv.svc.CreateJob(model.Job{Title: "Dev", Department: "Tech", Level: "P5", Headcount: 1})
	c, _ := srv.svc.CreateCandidate(model.Candidate{Name: "Tom", Email: "tom@example.com", Phone: "000"})
	o, _ := srv.svc.CreateOffer(model.Offer{JobID: j.ID, CandidateID: c.ID, Salary: 30000, ExpiresAt: time.Now().Add(7 * 24 * time.Hour)})
	body := `{"status":"accepted"}`
	req := httptest.NewRequest(http.MethodPut, "/api/offers/"+o.ID, strings.NewReader(body))
	w := httptest.NewRecorder()
	srv.Routes().ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("expect 200, got %d: %s", w.Code, w.Body.String())
	}
	var resp struct {
		Code int         `json:"code"`
		Data model.Offer `json:"data"`
	}
	json.Unmarshal(w.Body.Bytes(), &resp)
	if resp.Data.Status != model.OfferStatusAccepted {
		t.Errorf("expect accepted, got %s", resp.Data.Status)
	}
	job, _ := srv.svc.GetJob(j.ID)
	if job.Status != model.JobStatusFilled {
		t.Errorf("expect job filled, got %s", job.Status)
	}
}

func TestHandlerUpdateOfferInvalidTransition(t *testing.T) {
	srv := newTestServer()
	j, _ := srv.svc.CreateJob(model.Job{Title: "Dev", Department: "Tech", Level: "P5", Headcount: 1})
	c, _ := srv.svc.CreateCandidate(model.Candidate{Name: "Tom", Email: "tom@example.com", Phone: "000"})
	o, _ := srv.svc.CreateOffer(model.Offer{JobID: j.ID, CandidateID: c.ID, Salary: 30000, ExpiresAt: time.Now().Add(7 * 24 * time.Hour)})
	body1 := `{"status":"declined"}`
	req1 := httptest.NewRequest(http.MethodPut, "/api/offers/"+o.ID, strings.NewReader(body1))
	w1 := httptest.NewRecorder()
	srv.Routes().ServeHTTP(w1, req1)
	if w1.Code != http.StatusOK {
		t.Fatalf("expect 200 for decline, got %d", w1.Code)
	}
	body2 := `{"status":"accepted"}`
	req2 := httptest.NewRequest(http.MethodPut, "/api/offers/"+o.ID, strings.NewReader(body2))
	w2 := httptest.NewRecorder()
	srv.Routes().ServeHTTP(w2, req2)
	if w2.Code != http.StatusBadRequest {
		t.Errorf("expect 400 for invalid transition, got %d", w2.Code)
	}
}

func TestHandlerDeleteOffer(t *testing.T) {
	srv := newTestServer()
	j, _ := srv.svc.CreateJob(model.Job{Title: "Dev", Department: "Tech", Level: "P5", Headcount: 1})
	c, _ := srv.svc.CreateCandidate(model.Candidate{Name: "Tom", Email: "tom@example.com", Phone: "000"})
	o, _ := srv.svc.CreateOffer(model.Offer{JobID: j.ID, CandidateID: c.ID, Salary: 30000, ExpiresAt: time.Now().Add(7 * 24 * time.Hour)})
	req := httptest.NewRequest(http.MethodDelete, "/api/offers/"+o.ID, nil)
	w := httptest.NewRecorder()
	srv.Routes().ServeHTTP(w, req)
	if w.Code != http.StatusNoContent {
		t.Errorf("expect 204, got %d", w.Code)
	}
}

func TestHandlerStatsDepartments(t *testing.T) {
	srv := newTestServer()
	srv.svc.CreateJob(model.Job{Title: "Dev", Department: "Tech", Level: "P5", Headcount: 1})
	srv.svc.CreateJob(model.Job{Title: "PM", Department: "Product", Level: "P5", Headcount: 1})
	req := httptest.NewRequest(http.MethodGet, "/api/stats/departments", nil)
	w := httptest.NewRecorder()
	srv.Routes().ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("expect 200, got %d", w.Code)
	}
}

func TestHandlerStatsFunnel(t *testing.T) {
	srv := newTestServer()
	srv.svc.CreateJob(model.Job{Title: "Dev", Department: "Tech", Level: "P5", Headcount: 1})
	req := httptest.NewRequest(http.MethodGet, "/api/stats/funnel", nil)
	w := httptest.NewRecorder()
	srv.Routes().ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("expect 200, got %d", w.Code)
	}
}

func TestHandlerBadRequestBody(t *testing.T) {
	srv := newTestServer()
	req := httptest.NewRequest(http.MethodPost, "/api/jobs", bytes.NewReader([]byte("not json")))
	w := httptest.NewRecorder()
	srv.Routes().ServeHTTP(w, req)
	if w.Code != http.StatusBadRequest {
		t.Errorf("expect 400, got %d", w.Code)
	}
}

func TestHandlerMethodNotFound(t *testing.T) {
	srv := newTestServer()
	req := httptest.NewRequest(http.MethodGet, "/api/unknown", nil)
	w := httptest.NewRecorder()
	srv.Routes().ServeHTTP(w, req)
	if w.Code != http.StatusNotFound {
		t.Errorf("expect 404, got %d", w.Code)
	}
}

func TestHandlerRecoveryMiddleware(t *testing.T) {
	srv := newTestServer()
	mux := http.NewServeMux()
	mux.HandleFunc("GET /panic", func(w http.ResponseWriter, r *http.Request) {
		panic("intentional panic")
	})
	handler := srv.recoveryMiddleware(mux)
	req := httptest.NewRequest(http.MethodGet, "/panic", nil)
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)
	if w.Code != http.StatusInternalServerError {
		t.Errorf("expect 500, got %d", w.Code)
	}
}

func TestHandlerLoggingMiddleware(t *testing.T) {
	srv := newTestServer()
	mux := http.NewServeMux()
	mux.HandleFunc("GET /ok", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	handler := srv.loggingMiddleware(mux)
	req := httptest.NewRequest(http.MethodGet, "/ok", nil)
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Errorf("expect 200, got %d", w.Code)
	}
}
