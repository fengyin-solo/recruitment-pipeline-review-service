package service

import (
	"recruit/internal/model"
	"testing"
)

func TestListJobsZeroPageDoesNotPanic(t *testing.T) {
	svc := newTestService()
	for _, title := range []string{"A", "B", "C"} {
		if _, err := svc.CreateJob(model.Job{Title: title, Department: "Tech", Level: "P4", Headcount: 1}); err != nil {
			t.Fatalf("create job: %v", err)
		}
	}
	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("list jobs panicked for page zero: %v", r)
		}
	}()
	items, total, err := svc.ListJobs(model.JobFilter{}, 0, 2)
	if err != nil {
		t.Fatalf("list jobs: %v", err)
	}
	if total != 3 || len(items) != 2 {
		t.Fatalf("unexpected page result: total=%d len=%d", total, len(items))
	}
}
