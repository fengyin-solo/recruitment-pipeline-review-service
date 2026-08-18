package service

import (
	"testing"
	"time"

	"recruit/internal/config"
	"recruit/internal/model"
	"recruit/internal/store"
	"recruit/pkg/logger"
)

func newTestService() *Service {
	cfg := &config.Config{MaxPageSize: 100}
	log := logger.NewLevel(logger.LevelError)
	return New(store.NewMemoryStore(), log, cfg)
}

func TestServiceJobCRUD(t *testing.T) {
	svc := newTestService()
	j, err := svc.CreateJob(model.Job{Title: "Go Engineer", Department: "Tech", Level: "P5", Headcount: 2, SalaryMin: 20000, SalaryMax: 40000})
	if err != nil {
		t.Fatalf("create job: %v", err)
	}
	if j.ID == "" {
		t.Error("job id empty")
	}
	if j.Status != model.JobStatusOpen {
		t.Errorf("expect default open, got %s", j.Status)
	}
	got, err := svc.GetJob(j.ID)
	if err != nil {
		t.Fatalf("get job: %v", err)
	}
	if got.Title != "Go Engineer" {
		t.Error("title mismatch")
	}
	items, total, err := svc.ListJobs(model.JobFilter{}, 1, 10)
	if err != nil {
		t.Fatalf("list jobs: %v", err)
	}
	if total != 1 || len(items) != 1 {
		t.Errorf("list jobs count mismatch: total=%d len=%d", total, len(items))
	}
	updated, err := svc.UpdateJob(j.ID, model.Job{Title: "Senior Go Engineer"})
	if err != nil {
		t.Fatalf("update job: %v", err)
	}
	if updated.Title != "Senior Go Engineer" {
		t.Error("update title mismatch")
	}
	if err := svc.DeleteJob(j.ID); err != nil {
		t.Fatalf("delete job: %v", err)
	}
	if _, err := svc.GetJob(j.ID); err == nil {
		t.Error("expect error after delete")
	}
	if _, err := svc.GetJob(""); err == nil {
		t.Error("expect error for empty id")
	}
	if _, err := svc.UpdateJob("", model.Job{}); err == nil {
		t.Error("expect error for empty id update")
	}
	if err := svc.DeleteJob(""); err == nil {
		t.Error("expect error for empty id delete")
	}
}

func TestServiceJobValidation(t *testing.T) {
	svc := newTestService()
	tests := []struct {
		name string
		job  model.Job
	}{
		{"empty title", model.Job{Title: "", Department: "Tech", Level: "P5", Headcount: 1}},
		{"empty department", model.Job{Title: "Dev", Department: "", Level: "P5", Headcount: 1}},
		{"empty level", model.Job{Title: "Dev", Department: "Tech", Level: "", Headcount: 1}},
		{"zero headcount", model.Job{Title: "Dev", Department: "Tech", Level: "P5", Headcount: 0}},
		{"negative salary", model.Job{Title: "Dev", Department: "Tech", Level: "P5", Headcount: 1, SalaryMin: -1}},
		{"invalid salary range", model.Job{Title: "Dev", Department: "Tech", Level: "P5", Headcount: 1, SalaryMin: 50000, SalaryMax: 30000}},
		{"invalid status", model.Job{Title: "Dev", Department: "Tech", Level: "P5", Headcount: 1, Status: "invalid"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if _, err := svc.CreateJob(tt.job); err == nil {
				t.Error("expect validation error")
			}
		})
	}
}

func TestServiceJobFilter(t *testing.T) {
	svc := newTestService()
	svc.CreateJob(model.Job{Title: "Go Engineer", Department: "Tech", Level: "P4", Headcount: 1, Status: model.JobStatusOpen, Description: "backend"})
	svc.CreateJob(model.Job{Title: "Java Engineer", Department: "Tech", Level: "P5", Headcount: 1, Status: model.JobStatusClosed, Description: "backend"})
	svc.CreateJob(model.Job{Title: "PM", Department: "Product", Level: "P5", Headcount: 1, Status: model.JobStatusOpen, Description: "product"})
	tests := []struct {
		name   string
		filter model.JobFilter
		want   int
	}{
		{"department tech", model.JobFilter{Department: "Tech"}, 2},
		{"status open", model.JobFilter{Status: model.JobStatusOpen}, 2},
		{"department+status", model.JobFilter{Department: "Tech", Status: model.JobStatusOpen}, 1},
		{"keyword engineer", model.JobFilter{Keyword: "engineer"}, 2},
		{"keyword backend", model.JobFilter{Keyword: "backend"}, 2},
		{"keyword pm", model.JobFilter{Keyword: "pm"}, 1},
		{"no match", model.JobFilter{Keyword: "flutter"}, 0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, total, err := svc.ListJobs(tt.filter, 1, 10)
			if err != nil {
				t.Fatalf("list jobs: %v", err)
			}
			if total != tt.want {
				t.Errorf("want %d, got %d", tt.want, total)
			}
		})
	}
}

func TestServiceJobPagination(t *testing.T) {
	svc := newTestService()
	for i := 0; i < 5; i++ {
		svc.CreateJob(model.Job{Title: "Job" + string(rune('A'+i)), Department: "Tech", Level: "P5", Headcount: 1})
	}
	items, total, err := svc.ListJobs(model.JobFilter{}, 1, 2)
	if err != nil {
		t.Fatalf("list jobs page1: %v", err)
	}
	if total != 5 || len(items) != 2 {
		t.Errorf("page1 mismatch: total=%d len=%d", total, len(items))
	}
	items, total, err = svc.ListJobs(model.JobFilter{}, 3, 2)
	if err != nil {
		t.Fatalf("list jobs page3: %v", err)
	}
	if total != 5 || len(items) != 1 {
		t.Errorf("page3 mismatch: total=%d len=%d", total, len(items))
	}
	items, total, err = svc.ListJobs(model.JobFilter{}, 10, 2)
	if err != nil {
		t.Fatalf("list jobs page10: %v", err)
	}
	if total != 5 || len(items) != 0 {
		t.Errorf("page10 mismatch: total=%d len=%d", total, len(items))
	}
}

func TestServiceBatchCloseJobs(t *testing.T) {
	svc := newTestService()
	j1, _ := svc.CreateJob(model.Job{Title: "Dev1", Department: "Tech", Level: "P4", Headcount: 1})
	j2, _ := svc.CreateJob(model.Job{Title: "Dev2", Department: "Tech", Level: "P5", Headcount: 1})
	j3, _ := svc.CreateJob(model.Job{Title: "Dev3", Department: "Tech", Level: "P5", Headcount: 1, Status: model.JobStatusClosed})
	count, err := svc.BatchCloseJobs([]string{j1.ID, j2.ID, j3.ID})
	if err != nil {
		t.Fatalf("batch close: %v", err)
	}
	if count != 2 {
		t.Errorf("expect 2 closed, got %d", count)
	}
	for _, id := range []string{j1.ID, j2.ID} {
		j, _ := svc.GetJob(id)
		if j.Status != model.JobStatusClosed {
			t.Errorf("expect closed, got %s", j.Status)
		}
	}
	if _, err := svc.BatchCloseJobs([]string{}); err == nil {
		t.Error("expect error for empty ids")
	}
}

func TestServiceCandidateCRUD(t *testing.T) {
	svc := newTestService()
	c, err := svc.CreateCandidate(model.Candidate{Name: "Alice", Email: "alice@example.com", Phone: "123", YearsOfExperience: 3})
	if err != nil {
		t.Fatalf("create candidate: %v", err)
	}
	if c.Status != model.CandidateStatusActive {
		t.Errorf("expect default active, got %s", c.Status)
	}
	if _, err := svc.CreateCandidate(model.Candidate{Name: "Alice2", Email: "alice@example.com", Phone: "456"}); err == nil {
		t.Error("expect error for duplicate email")
	}
	got, err := svc.GetCandidate(c.ID)
	if err != nil {
		t.Fatalf("get candidate: %v", err)
	}
	if got.Name != "Alice" {
		t.Error("name mismatch")
	}
	_, total, err := svc.ListCandidates(model.CandidateFilter{Keyword: "ali"}, 1, 10)
	if err != nil {
		t.Fatalf("list candidates: %v", err)
	}
	if total != 1 {
		t.Errorf("list candidates count mismatch: %d", total)
	}
	updated, err := svc.UpdateCandidate(c.ID, model.Candidate{Name: "Alice Updated"})
	if err != nil {
		t.Fatalf("update candidate: %v", err)
	}
	if updated.Name != "Alice Updated" {
		t.Error("update name mismatch")
	}
	if err := svc.DeleteCandidate(c.ID); err != nil {
		t.Fatalf("delete candidate: %v", err)
	}
	if _, err := svc.GetCandidate(c.ID); err == nil {
		t.Error("expect error after delete")
	}
}

func TestServiceCandidateValidation(t *testing.T) {
	svc := newTestService()
	tests := []struct {
		name string
		c    model.Candidate
	}{
		{"empty name", model.Candidate{Name: "", Email: "a@b.com", Phone: "123"}},
		{"empty email", model.Candidate{Name: "Alice", Email: "", Phone: "123"}},
		{"invalid email", model.Candidate{Name: "Alice", Email: "notemail", Phone: "123"}},
		{"empty phone", model.Candidate{Name: "Alice", Email: "a@b.com", Phone: ""}},
		{"negative yoe", model.Candidate{Name: "Alice", Email: "a@b.com", Phone: "123", YearsOfExperience: -1}},
		{"invalid status", model.Candidate{Name: "Alice", Email: "a@b.com", Phone: "123", Status: "invalid"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if _, err := svc.CreateCandidate(tt.c); err == nil {
				t.Error("expect validation error")
			}
		})
	}
}

func TestServiceCandidateFilter(t *testing.T) {
	svc := newTestService()
	svc.CreateCandidate(model.Candidate{Name: "Alice", Email: "alice@example.com", Phone: "123", Status: model.CandidateStatusActive})
	svc.CreateCandidate(model.Candidate{Name: "Bob", Email: "bob@example.com", Phone: "456", Status: model.CandidateStatusWithdrawn})
	tests := []struct {
		name   string
		filter model.CandidateFilter
		want   int
	}{
		{"keyword ali", model.CandidateFilter{Keyword: "ali"}, 1},
		{"keyword bob", model.CandidateFilter{Keyword: "bob"}, 1},
		{"status active", model.CandidateFilter{Status: model.CandidateStatusActive}, 1},
		{"status withdrawn", model.CandidateFilter{Status: model.CandidateStatusWithdrawn}, 1},
		{"no match", model.CandidateFilter{Keyword: "charlie"}, 0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, total, err := svc.ListCandidates(tt.filter, 1, 10)
			if err != nil {
				t.Fatalf("list candidates: %v", err)
			}
			if total != tt.want {
				t.Errorf("want %d, got %d", tt.want, total)
			}
		})
	}
}

func TestServiceResumeForeignKey(t *testing.T) {
	svc := newTestService()
	if _, err := svc.CreateResume(model.Resume{CandidateID: "notfound", Summary: "summary", Education: "BS"}); err == nil {
		t.Error("expect error for nonexistent candidate")
	}
	c, _ := svc.CreateCandidate(model.Candidate{Name: "Bob", Email: "bob@example.com", Phone: "789"})
	r, err := svc.CreateResume(model.Resume{CandidateID: c.ID, Summary: "summary", Education: "BS"})
	if err != nil {
		t.Fatalf("create resume: %v", err)
	}
	if _, err := svc.GetResume(r.ID); err != nil {
		t.Fatalf("get resume: %v", err)
	}
	_, total, err := svc.ListResumes(model.ResumeFilter{CandidateID: c.ID}, 1, 10)
	if err != nil {
		t.Fatalf("list resumes: %v", err)
	}
	if total != 1 {
		t.Errorf("expect 1 resume, got %d", total)
	}
	updated, err := svc.UpdateResume(r.ID, model.Resume{Summary: "updated"})
	if err != nil {
		t.Fatalf("update resume: %v", err)
	}
	if updated.Summary != "updated" {
		t.Error("update summary mismatch")
	}
	if err := svc.DeleteResume(r.ID); err != nil {
		t.Fatalf("delete resume: %v", err)
	}
	if _, err := svc.GetResume(r.ID); err == nil {
		t.Error("expect error after delete")
	}
}

func TestServiceResumeValidation(t *testing.T) {
	svc := newTestService()
	c, _ := svc.CreateCandidate(model.Candidate{Name: "Bob", Email: "bob@example.com", Phone: "789"})
	tests := []struct {
		name string
		r    model.Resume
	}{
		{"empty candidate_id", model.Resume{CandidateID: "", Summary: "s", Education: "BS"}},
		{"empty summary", model.Resume{CandidateID: c.ID, Summary: "", Education: "BS"}},
		{"empty education", model.Resume{CandidateID: c.ID, Summary: "s", Education: ""}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if _, err := svc.CreateResume(tt.r); err == nil {
				t.Error("expect validation error")
			}
		})
	}
}

func TestServiceInterviewStateMachine(t *testing.T) {
	svc := newTestService()
	j, _ := svc.CreateJob(model.Job{Title: "Dev", Department: "Tech", Level: "P4", Headcount: 1})
	c, _ := svc.CreateCandidate(model.Candidate{Name: "Tom", Email: "tom@example.com", Phone: "000"})
	iv, err := svc.CreateInterview(model.Interview{JobID: j.ID, CandidateID: c.ID, Round: 1, Interviewer: "M1", ScheduledAt: time.Now().Add(time.Hour)})
	if err != nil {
		t.Fatalf("create interview: %v", err)
	}
	if iv.Status != model.InterviewStatusScheduled {
		t.Errorf("expect default scheduled, got %s", iv.Status)
	}
	completed, err := svc.UpdateInterview(iv.ID, model.Interview{Status: model.InterviewStatusCompleted, Score: 85, Feedback: "good"})
	if err != nil {
		t.Fatalf("scheduled->completed: %v", err)
	}
	if completed.Status != model.InterviewStatusCompleted {
		t.Errorf("expect completed, got %s", completed.Status)
	}
	if completed.Score != 85 {
		t.Errorf("expect score 85, got %d", completed.Score)
	}

	iv2, _ := svc.CreateInterview(model.Interview{JobID: j.ID, CandidateID: c.ID, Round: 2, Interviewer: "M2", ScheduledAt: time.Now().Add(2 * time.Hour)})
	cancelled, err := svc.UpdateInterview(iv2.ID, model.Interview{Status: model.InterviewStatusCancelled})
	if err != nil {
		t.Fatalf("scheduled->cancelled: %v", err)
	}
	if cancelled.Status != model.InterviewStatusCancelled {
		t.Errorf("expect cancelled, got %s", cancelled.Status)
	}
	if _, err := svc.UpdateInterview(iv2.ID, model.Interview{Status: model.InterviewStatusCompleted}); err == nil {
		t.Error("expect error for cancelled->completed")
	}
	if _, err := svc.UpdateInterview(iv2.ID, model.Interview{Status: model.InterviewStatusScheduled}); err == nil {
		t.Error("expect error for cancelled->scheduled")
	}
}

func TestServiceInterviewValidation(t *testing.T) {
	svc := newTestService()
	j, _ := svc.CreateJob(model.Job{Title: "Dev", Department: "Tech", Level: "P4", Headcount: 1})
	c, _ := svc.CreateCandidate(model.Candidate{Name: "Tom", Email: "tom@example.com", Phone: "000"})
	tests := []struct {
		name string
		iv   model.Interview
	}{
		{"empty job_id", model.Interview{JobID: "", CandidateID: c.ID, Round: 1, Interviewer: "M", ScheduledAt: time.Now().Add(time.Hour)}},
		{"empty candidate_id", model.Interview{JobID: j.ID, CandidateID: "", Round: 1, Interviewer: "M", ScheduledAt: time.Now().Add(time.Hour)}},
		{"zero round", model.Interview{JobID: j.ID, CandidateID: c.ID, Round: 0, Interviewer: "M", ScheduledAt: time.Now().Add(time.Hour)}},
		{"empty interviewer", model.Interview{JobID: j.ID, CandidateID: c.ID, Round: 1, Interviewer: "", ScheduledAt: time.Now().Add(time.Hour)}},
		{"zero scheduled_at", model.Interview{JobID: j.ID, CandidateID: c.ID, Round: 1, Interviewer: "M", ScheduledAt: time.Time{}}},
		{"invalid status", model.Interview{JobID: j.ID, CandidateID: c.ID, Round: 1, Interviewer: "M", ScheduledAt: time.Now().Add(time.Hour), Status: "invalid"}},
		{"negative score", model.Interview{JobID: j.ID, CandidateID: c.ID, Round: 1, Interviewer: "M", ScheduledAt: time.Now().Add(time.Hour), Score: -1}},
		{"score over 100", model.Interview{JobID: j.ID, CandidateID: c.ID, Round: 1, Interviewer: "M", ScheduledAt: time.Now().Add(time.Hour), Score: 101}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if _, err := svc.CreateInterview(tt.iv); err == nil {
				t.Error("expect validation error")
			}
		})
	}
}

func TestServiceInterviewForeignKey(t *testing.T) {
	svc := newTestService()
	c, _ := svc.CreateCandidate(model.Candidate{Name: "Tom", Email: "tom@example.com", Phone: "000"})
	if _, err := svc.CreateInterview(model.Interview{JobID: "notfound", CandidateID: c.ID, Round: 1, Interviewer: "M", ScheduledAt: time.Now().Add(time.Hour)}); err == nil {
		t.Error("expect error for nonexistent job")
	}
	j, _ := svc.CreateJob(model.Job{Title: "Dev", Department: "Tech", Level: "P4", Headcount: 1})
	if _, err := svc.CreateInterview(model.Interview{JobID: j.ID, CandidateID: "notfound", Round: 1, Interviewer: "M", ScheduledAt: time.Now().Add(time.Hour)}); err == nil {
		t.Error("expect error for nonexistent candidate")
	}
}

func TestServiceInterviewPagination(t *testing.T) {
	svc := newTestService()
	j, _ := svc.CreateJob(model.Job{Title: "Dev", Department: "Tech", Level: "P4", Headcount: 1})
	c, _ := svc.CreateCandidate(model.Candidate{Name: "Tom", Email: "tom@example.com", Phone: "000"})
	for i := 0; i < 5; i++ {
		svc.CreateInterview(model.Interview{JobID: j.ID, CandidateID: c.ID, Round: i + 1, Interviewer: "M", ScheduledAt: time.Now().Add(time.Hour)})
	}
	items, total, err := svc.ListInterviews(model.InterviewFilter{}, 1, 2)
	if err != nil {
		t.Fatalf("list interviews: %v", err)
	}
	if total != 5 || len(items) != 2 {
		t.Errorf("page1 mismatch: total=%d len=%d", total, len(items))
	}
}

func TestServiceBatchUpdateInterviewStatus(t *testing.T) {
	svc := newTestService()
	j, _ := svc.CreateJob(model.Job{Title: "Dev", Department: "Tech", Level: "P4", Headcount: 1})
	c, _ := svc.CreateCandidate(model.Candidate{Name: "Tom", Email: "tom@example.com", Phone: "000"})
	iv1, _ := svc.CreateInterview(model.Interview{JobID: j.ID, CandidateID: c.ID, Round: 1, Interviewer: "M1", ScheduledAt: time.Now().Add(time.Hour)})
	iv2, _ := svc.CreateInterview(model.Interview{JobID: j.ID, CandidateID: c.ID, Round: 2, Interviewer: "M2", ScheduledAt: time.Now().Add(2 * time.Hour)})
	iv3, _ := svc.CreateInterview(model.Interview{JobID: j.ID, CandidateID: c.ID, Round: 3, Interviewer: "M3", ScheduledAt: time.Now().Add(3 * time.Hour)})
	_, _ = svc.UpdateInterview(iv3.ID, model.Interview{Status: model.InterviewStatusCompleted})

	count, err := svc.BatchUpdateInterviewStatus([]string{iv1.ID, iv2.ID, iv3.ID}, model.InterviewStatusCompleted)
	if err != nil {
		t.Fatalf("batch update: %v", err)
	}
	if count != 2 {
		t.Errorf("expect 2 updated, got %d", count)
	}
	if _, err := svc.BatchUpdateInterviewStatus([]string{}, model.InterviewStatusCompleted); err == nil {
		t.Error("expect error for empty ids")
	}
	if _, err := svc.BatchUpdateInterviewStatus([]string{iv1.ID}, "invalid"); err == nil {
		t.Error("expect error for invalid status")
	}
}

func TestServiceOfferStateMachine(t *testing.T) {
	svc := newTestService()
	j, _ := svc.CreateJob(model.Job{Title: "Dev", Department: "Tech", Level: "P4", Headcount: 1})
	c, _ := svc.CreateCandidate(model.Candidate{Name: "Tom", Email: "tom@example.com", Phone: "000"})
	o, err := svc.CreateOffer(model.Offer{JobID: j.ID, CandidateID: c.ID, Salary: 30000, ExpiresAt: time.Now().Add(7 * 24 * time.Hour)})
	if err != nil {
		t.Fatalf("create offer: %v", err)
	}
	if o.Status != model.OfferStatusPending {
		t.Errorf("expect default pending, got %s", o.Status)
	}
	accepted, err := svc.UpdateOffer(o.ID, model.Offer{Status: model.OfferStatusAccepted})
	if err != nil {
		t.Fatalf("pending->accepted: %v", err)
	}
	if accepted.Status != model.OfferStatusAccepted {
		t.Errorf("expect accepted, got %s", accepted.Status)
	}
	job, _ := svc.GetJob(j.ID)
	if job.Status != model.JobStatusFilled {
		t.Errorf("expect job filled after offer accepted, got %s", job.Status)
	}

	o2, _ := svc.CreateOffer(model.Offer{JobID: j.ID, CandidateID: c.ID, Salary: 35000, ExpiresAt: time.Now().Add(7 * 24 * time.Hour)})
	declined, err := svc.UpdateOffer(o2.ID, model.Offer{Status: model.OfferStatusDeclined})
	if err != nil {
		t.Fatalf("pending->declined: %v", err)
	}
	if declined.Status != model.OfferStatusDeclined {
		t.Errorf("expect declined, got %s", declined.Status)
	}
	if _, err := svc.UpdateOffer(o2.ID, model.Offer{Status: model.OfferStatusAccepted}); err == nil {
		t.Error("expect error for declined->accepted")
	}
	if _, err := svc.UpdateOffer(o2.ID, model.Offer{Status: model.OfferStatusPending}); err == nil {
		t.Error("expect error for declined->pending")
	}
}

func TestServiceOfferValidation(t *testing.T) {
	svc := newTestService()
	j, _ := svc.CreateJob(model.Job{Title: "Dev", Department: "Tech", Level: "P4", Headcount: 1})
	c, _ := svc.CreateCandidate(model.Candidate{Name: "Tom", Email: "tom@example.com", Phone: "000"})
	tests := []struct {
		name string
		o    model.Offer
	}{
		{"empty job_id", model.Offer{JobID: "", CandidateID: c.ID, Salary: 30000, ExpiresAt: time.Now().Add(7 * 24 * time.Hour)}},
		{"empty candidate_id", model.Offer{JobID: j.ID, CandidateID: "", Salary: 30000, ExpiresAt: time.Now().Add(7 * 24 * time.Hour)}},
		{"zero salary", model.Offer{JobID: j.ID, CandidateID: c.ID, Salary: 0, ExpiresAt: time.Now().Add(7 * 24 * time.Hour)}},
		{"negative salary", model.Offer{JobID: j.ID, CandidateID: c.ID, Salary: -1, ExpiresAt: time.Now().Add(7 * 24 * time.Hour)}},
		{"zero expires_at", model.Offer{JobID: j.ID, CandidateID: c.ID, Salary: 30000, ExpiresAt: time.Time{}}},
		{"invalid status", model.Offer{JobID: j.ID, CandidateID: c.ID, Salary: 30000, ExpiresAt: time.Now().Add(7 * 24 * time.Hour), Status: "invalid"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if _, err := svc.CreateOffer(tt.o); err == nil {
				t.Error("expect validation error")
			}
		})
	}
}

func TestServiceOfferForeignKey(t *testing.T) {
	svc := newTestService()
	c, _ := svc.CreateCandidate(model.Candidate{Name: "Tom", Email: "tom@example.com", Phone: "000"})
	if _, err := svc.CreateOffer(model.Offer{JobID: "notfound", CandidateID: c.ID, Salary: 30000, ExpiresAt: time.Now().Add(7 * 24 * time.Hour)}); err == nil {
		t.Error("expect error for nonexistent job")
	}
	j, _ := svc.CreateJob(model.Job{Title: "Dev", Department: "Tech", Level: "P4", Headcount: 1})
	if _, err := svc.CreateOffer(model.Offer{JobID: j.ID, CandidateID: "notfound", Salary: 30000, ExpiresAt: time.Now().Add(7 * 24 * time.Hour)}); err == nil {
		t.Error("expect error for nonexistent candidate")
	}
}

func TestServiceOfferPagination(t *testing.T) {
	svc := newTestService()
	j, _ := svc.CreateJob(model.Job{Title: "Dev", Department: "Tech", Level: "P4", Headcount: 1})
	c, _ := svc.CreateCandidate(model.Candidate{Name: "Tom", Email: "tom@example.com", Phone: "000"})
	for i := 0; i < 5; i++ {
		svc.CreateOffer(model.Offer{JobID: j.ID, CandidateID: c.ID, Salary: 30000 + i*1000, ExpiresAt: time.Now().Add(7 * 24 * time.Hour)})
	}
	items, total, err := svc.ListOffers(model.OfferFilter{}, 1, 2)
	if err != nil {
		t.Fatalf("list offers: %v", err)
	}
	if total != 5 || len(items) != 2 {
		t.Errorf("page1 mismatch: total=%d len=%d", total, len(items))
	}
}

func TestServiceOfferFilledJobMultipleAccepts(t *testing.T) {
	svc := newTestService()
	j, _ := svc.CreateJob(model.Job{Title: "Dev", Department: "Tech", Level: "P4", Headcount: 2})
	c1, _ := svc.CreateCandidate(model.Candidate{Name: "Tom", Email: "tom@example.com", Phone: "000"})
	c2, _ := svc.CreateCandidate(model.Candidate{Name: "Jerry", Email: "jerry@example.com", Phone: "111"})
	o1, _ := svc.CreateOffer(model.Offer{JobID: j.ID, CandidateID: c1.ID, Salary: 30000, ExpiresAt: time.Now().Add(7 * 24 * time.Hour)})
	o2, _ := svc.CreateOffer(model.Offer{JobID: j.ID, CandidateID: c2.ID, Salary: 30000, ExpiresAt: time.Now().Add(7 * 24 * time.Hour)})
	if _, err := svc.UpdateOffer(o1.ID, model.Offer{Status: model.OfferStatusAccepted}); err != nil {
		t.Fatalf("accept offer1: %v", err)
	}
	job, _ := svc.GetJob(j.ID)
	if job.Status != model.JobStatusOpen {
		t.Errorf("expect job still open after 1 accept, got %s", job.Status)
	}
	if _, err := svc.UpdateOffer(o2.ID, model.Offer{Status: model.OfferStatusAccepted}); err != nil {
		t.Fatalf("accept offer2: %v", err)
	}
	job, _ = svc.GetJob(j.ID)
	if job.Status != model.JobStatusFilled {
		t.Errorf("expect job filled after 2 accepts, got %s", job.Status)
	}
}

func TestServiceStats(t *testing.T) {
	svc := newTestService()
	svc.CreateJob(model.Job{Title: "Dev1", Department: "Tech", Level: "P4", Headcount: 2})
	svc.CreateJob(model.Job{Title: "Dev2", Department: "Tech", Level: "P5", Headcount: 1})
	svc.CreateJob(model.Job{Title: "PM", Department: "Product", Level: "P5", Headcount: 1})
	deptStats, err := svc.DepartmentStats()
	if err != nil {
		t.Fatalf("department stats: %v", err)
	}
	m := make(map[string]int)
	for _, st := range deptStats {
		m[st.Department] = st.JobCount
	}
	if m["Tech"] != 2 || m["Product"] != 1 {
		t.Errorf("department stats mismatch: %v", m)
	}
	funnel, err := svc.JobFunnelStats()
	if err != nil {
		t.Fatalf("funnel stats: %v", err)
	}
	if len(funnel) != 3 {
		t.Errorf("expect 3 funnel stats, got %d", len(funnel))
	}
}

func TestServiceStatsWithData(t *testing.T) {
	svc := newTestService()
	j, _ := svc.CreateJob(model.Job{Title: "Dev", Department: "Tech", Level: "P4", Headcount: 1})
	c1, _ := svc.CreateCandidate(model.Candidate{Name: "Tom", Email: "tom@example.com", Phone: "000"})
	c2, _ := svc.CreateCandidate(model.Candidate{Name: "Jerry", Email: "jerry@example.com", Phone: "111"})
	svc.CreateInterview(model.Interview{JobID: j.ID, CandidateID: c1.ID, Round: 1, Interviewer: "M1", ScheduledAt: time.Now().Add(time.Hour)})
	svc.CreateInterview(model.Interview{JobID: j.ID, CandidateID: c2.ID, Round: 1, Interviewer: "M2", ScheduledAt: time.Now().Add(2 * time.Hour)})
	svc.CreateOffer(model.Offer{JobID: j.ID, CandidateID: c1.ID, Salary: 30000, ExpiresAt: time.Now().Add(7 * 24 * time.Hour)})

	funnel, err := svc.JobFunnelStats()
	if err != nil {
		t.Fatalf("funnel stats: %v", err)
	}
	if len(funnel) != 1 {
		t.Fatalf("expect 1 funnel stat, got %d", len(funnel))
	}
	st := funnel[0]
	if st.CandidateCount != 2 {
		t.Errorf("expect candidate count 2, got %d", st.CandidateCount)
	}
	if st.InterviewCount != 2 {
		t.Errorf("expect interview count 2, got %d", st.InterviewCount)
	}
	if st.OfferCount != 1 {
		t.Errorf("expect offer count 1, got %d", st.OfferCount)
	}
}

func TestServiceEmptyStats(t *testing.T) {
	svc := newTestService()
	deptStats, err := svc.DepartmentStats()
	if err != nil {
		t.Fatalf("empty department stats: %v", err)
	}
	if len(deptStats) != 0 {
		t.Errorf("expect empty department stats, got %d", len(deptStats))
	}
	funnel, err := svc.JobFunnelStats()
	if err != nil {
		t.Fatalf("empty funnel stats: %v", err)
	}
	if len(funnel) != 0 {
		t.Errorf("expect empty funnel stats, got %d", len(funnel))
	}
}
