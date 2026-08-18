package service

import (
	"recruit/internal/model"
	"testing"
)

func TestFailedJobSalaryUpdateKeepsStoredJob(t *testing.T) {
	svc := newTestService()
	j, err := svc.CreateJob(model.Job{Title: "Backend", Department: "Tech", Level: "P5", Headcount: 1, SalaryMin: 100, SalaryMax: 200})
	if err != nil {
		t.Fatalf("create job: %v", err)
	}
	if _, err := svc.UpdateJob(j.ID, model.Job{SalaryMin: 300, SalaryMax: 200}); err == nil {
		t.Fatalf("expected salary range validation error")
	}
	got, err := svc.GetJob(j.ID)
	if err != nil {
		t.Fatalf("get job: %v", err)
	}
	if got.SalaryMin != 100 || got.SalaryMax != 200 {
		t.Fatalf("failed update changed stored salary range: %d-%d", got.SalaryMin, got.SalaryMax)
	}
}
