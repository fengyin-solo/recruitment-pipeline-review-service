package service

import (
	"recruit/internal/model"
	"testing"
)

func TestResumeReturnedObjectCannotMutateStore(t *testing.T) {
	svc := newTestService()
	c, _ := svc.CreateCandidate(model.Candidate{Name: "Kai", Email: "kai@example.com", Phone: "10000"})
	r, err := svc.CreateResume(model.Resume{CandidateID: c.ID, Summary: "backend engineer", Education: "BS", WorkExperience: "two years"})
	if err != nil {
		t.Fatalf("create resume: %v", err)
	}
	got, _ := svc.GetResume(r.ID)
	got.Summary = "changed by caller"
	again, err := svc.GetResume(r.ID)
	if err != nil {
		t.Fatalf("get resume again: %v", err)
	}
	if again.Summary != "backend engineer" {
		t.Fatalf("external resume pointer mutation changed store: %q", again.Summary)
	}
}
