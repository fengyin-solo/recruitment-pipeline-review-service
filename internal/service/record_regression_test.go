package service

import (
	"testing"
	"time"

	"recruit/internal/model"
)

func TestFailedInterviewScoreUpdateKeepsStoredInterview(t *testing.T) {
	svc := newTestService()
	j, _ := svc.CreateJob(model.Job{Title: "Backend", Department: "Tech", Level: "P5", Headcount: 1})
	c, _ := svc.CreateCandidate(model.Candidate{Name: "Ren", Email: "ren@example.com", Phone: "10010"})
	iv, err := svc.CreateInterview(model.Interview{JobID: j.ID, CandidateID: c.ID, Round: 1, Interviewer: "Lee", ScheduledAt: time.Now().Add(time.Hour), Score: 80})
	if err != nil {
		t.Fatalf("create interview: %v", err)
	}
	if _, err := svc.UpdateInterview(iv.ID, model.Interview{Score: 130}); err == nil {
		t.Fatalf("expected score validation error")
	}
	got, err := svc.GetInterview(iv.ID)
	if err != nil {
		t.Fatalf("get interview: %v", err)
	}
	if got.Score != 80 {
		t.Fatalf("failed update changed stored score: %d", got.Score)
	}
}
