package service

import (
	"recruit/internal/model"
	"testing"
)

func TestCandidateSkillsSnapshotDoesNotLeak(t *testing.T) {
	svc := newTestService()
	c, err := svc.CreateCandidate(model.Candidate{Name: "Mina", Email: "mina@example.com", Phone: "10086", Skills: []string{"go", "sql"}})
	if err != nil {
		t.Fatalf("create candidate: %v", err)
	}
	c.Skills[0] = "php"
	got, err := svc.GetCandidate(c.ID)
	if err != nil {
		t.Fatalf("get candidate: %v", err)
	}
	if got.Skills[0] != "go" {
		t.Fatalf("candidate skills leaked through returned object: got %q", got.Skills[0])
	}
}
