package service

import (
	"testing"
	"time"

	"recruit/internal/model"
)

func TestOfferReturnedObjectCannotMutateStore(t *testing.T) {
	svc := newTestService()
	j, _ := svc.CreateJob(model.Job{Title: "Data", Department: "Tech", Level: "P6", Headcount: 1})
	c, _ := svc.CreateCandidate(model.Candidate{Name: "Nia", Email: "nia@example.com", Phone: "10001"})
	o, err := svc.CreateOffer(model.Offer{JobID: j.ID, CandidateID: c.ID, Salary: 300, ExpiresAt: time.Now().Add(24 * time.Hour)})
	if err != nil {
		t.Fatalf("create offer: %v", err)
	}
	got, _ := svc.GetOffer(o.ID)
	got.Status = model.OfferStatusAccepted
	again, err := svc.GetOffer(o.ID)
	if err != nil {
		t.Fatalf("get offer again: %v", err)
	}
	if again.Status != model.OfferStatusPending {
		t.Fatalf("external offer pointer mutation changed store: %q", again.Status)
	}
}
