package service

import (
	"recruit/internal/model"
	"testing"
)

func TestJobReturnedObjectCannotMutateStore(t *testing.T) {
	svc := newTestService()
	j, err := svc.CreateJob(model.Job{Title: "Frontend", Department: "Product", Level: "P4", Headcount: 1})
	if err != nil {
		t.Fatalf("create job: %v", err)
	}
	got, _ := svc.GetJob(j.ID)
	got.Title = "Mutated outside"
	again, err := svc.GetJob(j.ID)
	if err != nil {
		t.Fatalf("get job again: %v", err)
	}
	if again.Title != "Frontend" {
		t.Fatalf("external job pointer mutation changed store: %q", again.Title)
	}
}
