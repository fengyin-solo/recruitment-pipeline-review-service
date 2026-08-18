package service

import (
	"testing"
	"time"

	"recruit/internal/model"
)

func TestAcceptedOfferFillsSingleHeadcountJob(t *testing.T) {
	svc := newTestService()
	j, _ := svc.CreateJob(model.Job{Title: "SRE", Department: "Infra", Level: "P6", Headcount: 1})
	c, _ := svc.CreateCandidate(model.Candidate{Name: "Uma", Email: "uma@example.com", Phone: "10004"})
	o, err := svc.CreateOffer(model.Offer{JobID: j.ID, CandidateID: c.ID, Salary: 500, ExpiresAt: time.Now().Add(48 * time.Hour)})
	if err != nil {
		t.Fatalf("create offer: %v", err)
	}
	if _, err := svc.UpdateOffer(o.ID, model.Offer{Status: model.OfferStatusAccepted}); err != nil {
		t.Fatalf("accept offer: %v", err)
	}
	got, err := svc.GetJob(j.ID)
	if err != nil {
		t.Fatalf("get job: %v", err)
	}
	if got.Status != model.JobStatusFilled {
		t.Fatalf("single-headcount job should be filled after accepted offer, got %q", got.Status)
	}
}
