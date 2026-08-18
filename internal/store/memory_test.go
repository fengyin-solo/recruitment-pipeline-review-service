package store

import (
	"testing"
	"time"

	"recruit/internal/model"
)

func TestMemoryStoreJobCRUD(t *testing.T) {
	s := NewMemoryStore()
	j := &model.Job{ID: "j1", Title: "Go Engineer", Department: "Tech", Level: "P5", Headcount: 2, SalaryMin: 20000, SalaryMax: 40000, Status: model.JobStatusOpen, CreatedAt: time.Now(), UpdatedAt: time.Now()}
	if err := s.CreateJob(j); err != nil {
		t.Fatalf("create job: %v", err)
	}
	got, err := s.GetJob("j1")
	if err != nil {
		t.Fatalf("get job: %v", err)
	}
	if got.Title != "Go Engineer" {
		t.Errorf("title mismatch")
	}
	if _, err := s.GetJob("notfound"); err != ErrNotFound {
		t.Errorf("expect ErrNotFound, got %v", err)
	}
	if len(s.ListJobs()) != 1 {
		t.Errorf("list jobs count mismatch")
	}
	j.Title = "Senior Go Engineer"
	if err := s.UpdateJob(j); err != nil {
		t.Fatalf("update job: %v", err)
	}
	got, _ = s.GetJob("j1")
	if got.Title != "Senior Go Engineer" {
		t.Errorf("update title mismatch")
	}
	if err := s.DeleteJob("j1"); err != nil {
		t.Fatalf("delete job: %v", err)
	}
	if _, err := s.GetJob("j1"); err != ErrNotFound {
		t.Errorf("expect ErrNotFound after delete, got %v", err)
	}
	if err := s.DeleteJob("j1"); err != ErrNotFound {
		t.Errorf("expect ErrNotFound on double delete, got %v", err)
	}
	if err := s.UpdateJob(j); err != ErrNotFound {
		t.Errorf("expect ErrNotFound on update deleted, got %v", err)
	}
}

func TestMemoryStoreJobListMultiple(t *testing.T) {
	s := NewMemoryStore()
	for i := 0; i < 5; i++ {
		j := &model.Job{ID: "j" + string(rune('0'+i)), Title: "Job" + string(rune('0'+i)), Department: "Tech", Level: "P5", Headcount: 1, Status: model.JobStatusOpen, CreatedAt: time.Now(), UpdatedAt: time.Now()}
		_ = s.CreateJob(j)
	}
	if len(s.ListJobs()) != 5 {
		t.Errorf("expect 5 jobs, got %d", len(s.ListJobs()))
	}
}

func TestMemoryStoreCandidateCRUD(t *testing.T) {
	s := NewMemoryStore()
	c := &model.Candidate{ID: "c1", Name: "Alice", Email: "alice@example.com", Phone: "123", Status: model.CandidateStatusActive, CreatedAt: time.Now(), UpdatedAt: time.Now()}
	if err := s.CreateCandidate(c); err != nil {
		t.Fatalf("create candidate: %v", err)
	}
	if _, err := s.GetCandidate("c1"); err != nil {
		t.Fatalf("get candidate: %v", err)
	}
	if _, err := s.GetCandidateByEmail("alice@example.com"); err != nil {
		t.Fatalf("get candidate by email: %v", err)
	}
	if _, err := s.GetCandidateByEmail("notfound"); err != ErrNotFound {
		t.Errorf("expect ErrNotFound, got %v", err)
	}
	if _, err := s.GetCandidate("notfound"); err != ErrNotFound {
		t.Errorf("expect ErrNotFound by id, got %v", err)
	}
	cDup := &model.Candidate{ID: "cdup", Name: "Dup", Email: "alice@example.com", Phone: "999", Status: model.CandidateStatusActive, CreatedAt: time.Now(), UpdatedAt: time.Now()}
	if err := s.CreateCandidate(cDup); err != ErrConflict {
		t.Errorf("expect ErrConflict for duplicate email create, got %v", err)
	}
	c2 := &model.Candidate{ID: "c2", Name: "Bob", Email: "bob@example.com", Phone: "456", Status: model.CandidateStatusActive, CreatedAt: time.Now(), UpdatedAt: time.Now()}
	if err := s.CreateCandidate(c2); err != nil {
		t.Fatalf("create candidate c2: %v", err)
	}
	c.Email = "alice2@example.com"
	if err := s.UpdateCandidate(c); err != nil {
		t.Fatalf("update candidate: %v", err)
	}
	c2.Email = "alice2@example.com"
	if err := s.UpdateCandidate(c2); err != ErrConflict {
		t.Errorf("expect ErrConflict for duplicate email on update, got %v", err)
	}
	if err := s.DeleteCandidate("c1"); err != nil {
		t.Fatalf("delete candidate: %v", err)
	}
	if _, err := s.GetCandidate("c1"); err != ErrNotFound {
		t.Errorf("expect ErrNotFound after delete, got %v", err)
	}
	if err := s.DeleteCandidate("c1"); err != ErrNotFound {
		t.Errorf("expect ErrNotFound on double delete, got %v", err)
	}
	if err := s.UpdateCandidate(c); err != ErrNotFound {
		t.Errorf("expect ErrNotFound on update deleted, got %v", err)
	}
}

func TestMemoryStoreCandidateListMultiple(t *testing.T) {
	s := NewMemoryStore()
	for i := 0; i < 3; i++ {
		c := &model.Candidate{ID: "c" + string(rune('0'+i)), Name: "User" + string(rune('0'+i)), Email: "user" + string(rune('0'+i)) + "@example.com", Phone: "000", Status: model.CandidateStatusActive, CreatedAt: time.Now(), UpdatedAt: time.Now()}
		_ = s.CreateCandidate(c)
	}
	if len(s.ListCandidates()) != 3 {
		t.Errorf("expect 3 candidates, got %d", len(s.ListCandidates()))
	}
}

func TestMemoryStoreResumeCRUD(t *testing.T) {
	s := NewMemoryStore()
	r := &model.Resume{ID: "r1", CandidateID: "c1", Summary: "summary", Education: "BS", CreatedAt: time.Now(), UpdatedAt: time.Now()}
	if err := s.CreateResume(r); err != nil {
		t.Fatalf("create resume: %v", err)
	}
	if _, err := s.GetResume("r1"); err != nil {
		t.Fatalf("get resume: %v", err)
	}
	if len(s.ListResumes()) != 1 {
		t.Errorf("list resumes count mismatch")
	}
	r.Summary = "updated"
	if err := s.UpdateResume(r); err != nil {
		t.Fatalf("update resume: %v", err)
	}
	got, _ := s.GetResume("r1")
	if got.Summary != "updated" {
		t.Errorf("update summary mismatch")
	}
	if err := s.DeleteResume("r1"); err != nil {
		t.Fatalf("delete resume: %v", err)
	}
	if _, err := s.GetResume("r1"); err != ErrNotFound {
		t.Errorf("expect ErrNotFound after delete, got %v", err)
	}
	if err := s.DeleteResume("r1"); err != ErrNotFound {
		t.Errorf("expect ErrNotFound on double delete, got %v", err)
	}
	if err := s.UpdateResume(r); err != ErrNotFound {
		t.Errorf("expect ErrNotFound on update deleted, got %v", err)
	}
}

func TestMemoryStoreResumeListMultiple(t *testing.T) {
	s := NewMemoryStore()
	for i := 0; i < 4; i++ {
		r := &model.Resume{ID: "r" + string(rune('0'+i)), CandidateID: "c1", Summary: "s", Education: "BS", CreatedAt: time.Now(), UpdatedAt: time.Now()}
		_ = s.CreateResume(r)
	}
	if len(s.ListResumes()) != 4 {
		t.Errorf("expect 4 resumes, got %d", len(s.ListResumes()))
	}
}

func TestMemoryStoreInterviewCRUD(t *testing.T) {
	s := NewMemoryStore()
	iv := &model.Interview{ID: "i1", JobID: "j1", CandidateID: "c1", Round: 1, Interviewer: "Bob", ScheduledAt: time.Now(), Status: model.InterviewStatusScheduled, CreatedAt: time.Now(), UpdatedAt: time.Now()}
	if err := s.CreateInterview(iv); err != nil {
		t.Fatalf("create interview: %v", err)
	}
	if _, err := s.GetInterview("i1"); err != nil {
		t.Fatalf("get interview: %v", err)
	}
	if len(s.ListInterviews()) != 1 {
		t.Errorf("list interviews count mismatch")
	}
	iv.Round = 2
	if err := s.UpdateInterview(iv); err != nil {
		t.Fatalf("update interview: %v", err)
	}
	got, _ := s.GetInterview("i1")
	if got.Round != 2 {
		t.Errorf("update round mismatch")
	}
	if err := s.DeleteInterview("i1"); err != nil {
		t.Fatalf("delete interview: %v", err)
	}
	if _, err := s.GetInterview("i1"); err != ErrNotFound {
		t.Errorf("expect ErrNotFound after delete, got %v", err)
	}
	if err := s.DeleteInterview("i1"); err != ErrNotFound {
		t.Errorf("expect ErrNotFound on double delete, got %v", err)
	}
	if err := s.UpdateInterview(iv); err != ErrNotFound {
		t.Errorf("expect ErrNotFound on update deleted, got %v", err)
	}
}

func TestMemoryStoreInterviewListMultiple(t *testing.T) {
	s := NewMemoryStore()
	for i := 0; i < 3; i++ {
		iv := &model.Interview{ID: "i" + string(rune('0'+i)), JobID: "j1", CandidateID: "c1", Round: i + 1, Interviewer: "M", ScheduledAt: time.Now(), Status: model.InterviewStatusScheduled, CreatedAt: time.Now(), UpdatedAt: time.Now()}
		_ = s.CreateInterview(iv)
	}
	if len(s.ListInterviews()) != 3 {
		t.Errorf("expect 3 interviews, got %d", len(s.ListInterviews()))
	}
}

func TestMemoryStoreOfferCRUD(t *testing.T) {
	s := NewMemoryStore()
	o := &model.Offer{ID: "o1", JobID: "j1", CandidateID: "c1", Salary: 30000, Status: model.OfferStatusPending, ExpiresAt: time.Now().Add(7 * 24 * time.Hour), CreatedAt: time.Now(), UpdatedAt: time.Now()}
	if err := s.CreateOffer(o); err != nil {
		t.Fatalf("create offer: %v", err)
	}
	if _, err := s.GetOffer("o1"); err != nil {
		t.Fatalf("get offer: %v", err)
	}
	if len(s.ListOffers()) != 1 {
		t.Errorf("list offers count mismatch")
	}
	o.Salary = 35000
	if err := s.UpdateOffer(o); err != nil {
		t.Fatalf("update offer: %v", err)
	}
	got, _ := s.GetOffer("o1")
	if got.Salary != 35000 {
		t.Errorf("update salary mismatch")
	}
	if err := s.DeleteOffer("o1"); err != nil {
		t.Fatalf("delete offer: %v", err)
	}
	if _, err := s.GetOffer("o1"); err != ErrNotFound {
		t.Errorf("expect ErrNotFound after delete, got %v", err)
	}
	if err := s.DeleteOffer("o1"); err != ErrNotFound {
		t.Errorf("expect ErrNotFound on double delete, got %v", err)
	}
	if err := s.UpdateOffer(o); err != ErrNotFound {
		t.Errorf("expect ErrNotFound on update deleted, got %v", err)
	}
}

func TestMemoryStoreOfferListMultiple(t *testing.T) {
	s := NewMemoryStore()
	for i := 0; i < 3; i++ {
		o := &model.Offer{ID: "o" + string(rune('0'+i)), JobID: "j1", CandidateID: "c1", Salary: 30000, Status: model.OfferStatusPending, ExpiresAt: time.Now().Add(7 * 24 * time.Hour), CreatedAt: time.Now(), UpdatedAt: time.Now()}
		_ = s.CreateOffer(o)
	}
	if len(s.ListOffers()) != 3 {
		t.Errorf("expect 3 offers, got %d", len(s.ListOffers()))
	}
}

func TestMemoryStoreConcurrentAccess(t *testing.T) {
	s := NewMemoryStore()
	done := make(chan bool, 10)
	for i := 0; i < 10; i++ {
		go func(idx int) {
			j := &model.Job{ID: "jc" + string(rune('0'+idx)), Title: "Job", Department: "Tech", Level: "P5", Headcount: 1, Status: model.JobStatusOpen, CreatedAt: time.Now(), UpdatedAt: time.Now()}
			_ = s.CreateJob(j)
			_, _ = s.GetJob(j.ID)
			_ = s.ListJobs()
			j.Title = "Updated"
			_ = s.UpdateJob(j)
			done <- true
		}(i)
	}
	for i := 0; i < 10; i++ {
		<-done
	}
	if len(s.ListJobs()) != 10 {
		t.Errorf("expect 10 jobs after concurrent create, got %d", len(s.ListJobs()))
	}
}
