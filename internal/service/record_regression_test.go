package service

import (
	"testing"
	"time"

	"recruit/internal/model"
)

func TestInterviewListObjectsCannotMutateStore(t *testing.T) {
	svc := newTestService()
	j, _ := svc.CreateJob(model.Job{Title: "QA", Department: "Tech", Level: "P4", Headcount: 1})
	c, _ := svc.CreateCandidate(model.Candidate{Name: "Owen", Email: "owen@example.com", Phone: "10002"})
	iv, err := svc.CreateInterview(model.Interview{JobID: j.ID, CandidateID: c.ID, Round: 1, Interviewer: "Mo", ScheduledAt: time.Now().Add(time.Hour), Score: 70})
	if err != nil {
		t.Fatalf("create interview: %v", err)
	}
	items, _, err := svc.ListInterviews(model.InterviewFilter{}, 1, 10)
	if err != nil {
		t.Fatalf("list interviews: %v", err)
	}
	items[0].Score = 5
	again, err := svc.GetInterview(iv.ID)
	if err != nil {
		t.Fatalf("get interview again: %v", err)
	}
	if again.Score != 70 {
		t.Fatalf("list result mutation changed stored interview: %d", again.Score)
	}
}
