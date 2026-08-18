package service

import (
	"recruit/internal/model"
	"testing"
)

func TestCandidateKeywordMatchesEmailAddress(t *testing.T) {
	svc := newTestService()
	if _, err := svc.CreateCandidate(model.Candidate{Name: "Plain Name", Email: "devops.platform@example.com", Phone: "10003"}); err != nil {
		t.Fatalf("create candidate: %v", err)
	}
	_, total, err := svc.ListCandidates(model.CandidateFilter{Keyword: "platform@example.com"}, 1, 10)
	if err != nil {
		t.Fatalf("list candidates: %v", err)
	}
	if total != 1 {
		t.Fatalf("email keyword should match one candidate, got %d", total)
	}
}
