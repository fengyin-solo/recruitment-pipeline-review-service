package model

import (
	"strings"
	"testing"
	"time"
)

func TestJobValidate(t *testing.T) {
	tests := []struct {
		name    string
		job     Job
		wantErr bool
		errMsg  string
	}{
		{
			name: "valid job",
			job:  Job{Title: "Dev", Department: "Tech", Level: "P5", Headcount: 1, SalaryMin: 10000, SalaryMax: 20000, Status: JobStatusOpen},
			wantErr: false,
		},
		{
			name:    "empty title",
			job:     Job{Title: "", Department: "Tech", Level: "P5", Headcount: 1},
			wantErr: true,
			errMsg:  "title",
		},
		{
			name:    "empty department",
			job:     Job{Title: "Dev", Department: "", Level: "P5", Headcount: 1},
			wantErr: true,
			errMsg:  "department",
		},
		{
			name:    "empty level",
			job:     Job{Title: "Dev", Department: "Tech", Level: "", Headcount: 1},
			wantErr: true,
			errMsg:  "level",
		},
		{
			name:    "zero headcount",
			job:     Job{Title: "Dev", Department: "Tech", Level: "P5", Headcount: 0},
			wantErr: true,
			errMsg:  "headcount",
		},
		{
			name:    "negative salary min",
			job:     Job{Title: "Dev", Department: "Tech", Level: "P5", Headcount: 1, SalaryMin: -1},
			wantErr: true,
			errMsg:  "salary",
		},
		{
			name:    "invalid salary range",
			job:     Job{Title: "Dev", Department: "Tech", Level: "P5", Headcount: 1, SalaryMin: 50000, SalaryMax: 30000},
			wantErr: true,
			errMsg:  "salary",
		},
		{
			name:    "invalid status",
			job:     Job{Title: "Dev", Department: "Tech", Level: "P5", Headcount: 1, Status: "invalid"},
			wantErr: true,
			errMsg:  "status",
		},
		{
			name:    "default status",
			job:     Job{Title: "Dev", Department: "Tech", Level: "P5", Headcount: 1},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.job.Validate()
			if tt.wantErr {
				if err == nil {
					t.Errorf("expect error")
					return
				}
				if !strings.Contains(err.Error(), tt.errMsg) {
					t.Errorf("expect error contains %s, got %s", tt.errMsg, err.Error())
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
				if tt.job.Status == "" {
					t.Errorf("expect default status set")
				}
			}
		})
	}
}

func TestJobFilter(t *testing.T) {
	jobs := []*Job{
		{Title: "Go Dev", Department: "Tech", Status: JobStatusOpen, Description: "backend"},
		{Title: "Java Dev", Department: "Tech", Status: JobStatusClosed, Description: "backend"},
		{Title: "PM", Department: "Product", Status: JobStatusOpen, Description: "product"},
	}
	tests := []struct {
		name   string
		filter JobFilter
		want   []bool
	}{
		{"match all", JobFilter{}, []bool{true, true, true}},
		{"department tech", JobFilter{Department: "Tech"}, []bool{true, true, false}},
		{"status open", JobFilter{Status: JobStatusOpen}, []bool{true, false, true}},
		{"keyword go", JobFilter{Keyword: "go"}, []bool{true, false, false}},
		{"keyword backend", JobFilter{Keyword: "backend"}, []bool{true, true, false}},
		{"combined", JobFilter{Department: "Tech", Status: JobStatusOpen}, []bool{true, false, false}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for i, j := range jobs {
				got := tt.filter.Match(j)
				if got != tt.want[i] {
					t.Errorf("job %d: want %v, got %v", i, tt.want[i], got)
				}
			}
		})
	}
}

func TestCandidateValidate(t *testing.T) {
	tests := []struct {
		name    string
		c       Candidate
		wantErr bool
		errMsg  string
	}{
		{
			name: "valid candidate",
			c:    Candidate{Name: "Alice", Email: "alice@example.com", Phone: "123", YearsOfExperience: 3, Status: CandidateStatusActive},
			wantErr: false,
		},
		{
			name:    "empty name",
			c:       Candidate{Name: "", Email: "a@b.com", Phone: "123"},
			wantErr: true,
			errMsg:  "name",
		},
		{
			name:    "empty email",
			c:       Candidate{Name: "Alice", Email: "", Phone: "123"},
			wantErr: true,
			errMsg:  "email",
		},
		{
			name:    "invalid email",
			c:       Candidate{Name: "Alice", Email: "notemail", Phone: "123"},
			wantErr: true,
			errMsg:  "email",
		},
		{
			name:    "empty phone",
			c:       Candidate{Name: "Alice", Email: "a@b.com", Phone: ""},
			wantErr: true,
			errMsg:  "phone",
		},
		{
			name:    "negative yoe",
			c:       Candidate{Name: "Alice", Email: "a@b.com", Phone: "123", YearsOfExperience: -1},
			wantErr: true,
			errMsg:  "years_of_experience",
		},
		{
			name:    "invalid status",
			c:       Candidate{Name: "Alice", Email: "a@b.com", Phone: "123", Status: "invalid"},
			wantErr: true,
			errMsg:  "status",
		},
		{
			name:    "default status",
			c:       Candidate{Name: "Alice", Email: "a@b.com", Phone: "123"},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.c.Validate()
			if tt.wantErr {
				if err == nil {
					t.Errorf("expect error")
					return
				}
				if !strings.Contains(err.Error(), tt.errMsg) {
					t.Errorf("expect error contains %s, got %s", tt.errMsg, err.Error())
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
				if tt.c.Status == "" {
					t.Errorf("expect default status set")
				}
			}
		})
	}
}

func TestCandidateFilter(t *testing.T) {
	candidates := []*Candidate{
		{Name: "Alice", Email: "alice@example.com", Status: CandidateStatusActive},
		{Name: "Bob", Email: "bob@example.com", Status: CandidateStatusWithdrawn},
	}
	tests := []struct {
		name   string
		filter CandidateFilter
		want   []bool
	}{
		{"match all", CandidateFilter{}, []bool{true, true}},
		{"keyword ali", CandidateFilter{Keyword: "ali"}, []bool{true, false}},
		{"keyword bob", CandidateFilter{Keyword: "bob"}, []bool{false, true}},
		{"keyword email", CandidateFilter{Keyword: "example"}, []bool{true, true}},
		{"status active", CandidateFilter{Status: CandidateStatusActive}, []bool{true, false}},
		{"no match", CandidateFilter{Keyword: "charlie"}, []bool{false, false}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for i, c := range candidates {
				got := tt.filter.Match(c)
				if got != tt.want[i] {
					t.Errorf("candidate %d: want %v, got %v", i, tt.want[i], got)
				}
			}
		})
	}
}

func TestResumeValidate(t *testing.T) {
	tests := []struct {
		name    string
		r       Resume
		wantErr bool
		errMsg  string
	}{
		{
			name:    "valid resume",
			r:       Resume{CandidateID: "c1", Summary: "good", Education: "BS"},
			wantErr: false,
		},
		{
			name:    "empty candidate_id",
			r:       Resume{CandidateID: "", Summary: "good", Education: "BS"},
			wantErr: true,
			errMsg:  "candidate_id",
		},
		{
			name:    "empty summary",
			r:       Resume{CandidateID: "c1", Summary: "", Education: "BS"},
			wantErr: true,
			errMsg:  "summary",
		},
		{
			name:    "empty education",
			r:       Resume{CandidateID: "c1", Summary: "good", Education: ""},
			wantErr: true,
			errMsg:  "education",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.r.Validate()
			if tt.wantErr {
				if err == nil {
					t.Errorf("expect error")
					return
				}
				if !strings.Contains(err.Error(), tt.errMsg) {
					t.Errorf("expect error contains %s, got %s", tt.errMsg, err.Error())
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
			}
		})
	}
}

func TestResumeFilter(t *testing.T) {
	resumes := []*Resume{
		{CandidateID: "c1"},
		{CandidateID: "c2"},
	}
	f := ResumeFilter{CandidateID: "c1"}
	if !f.Match(resumes[0]) {
		t.Error("expect match c1")
	}
	if f.Match(resumes[1]) {
		t.Error("expect not match c2")
	}
}

func TestInterviewValidate(t *testing.T) {
	tests := []struct {
		name    string
		iv      Interview
		wantErr bool
		errMsg  string
	}{
		{
			name:    "valid interview",
			iv:      Interview{JobID: "j1", CandidateID: "c1", Round: 1, Interviewer: "M1", ScheduledAt: time.Now().Add(time.Hour), Status: InterviewStatusScheduled},
			wantErr: false,
		},
		{
			name:    "empty job_id",
			iv:      Interview{JobID: "", CandidateID: "c1", Round: 1, Interviewer: "M1", ScheduledAt: time.Now().Add(time.Hour)},
			wantErr: true,
			errMsg:  "job_id",
		},
		{
			name:    "empty candidate_id",
			iv:      Interview{JobID: "j1", CandidateID: "", Round: 1, Interviewer: "M1", ScheduledAt: time.Now().Add(time.Hour)},
			wantErr: true,
			errMsg:  "candidate_id",
		},
		{
			name:    "zero round",
			iv:      Interview{JobID: "j1", CandidateID: "c1", Round: 0, Interviewer: "M1", ScheduledAt: time.Now().Add(time.Hour)},
			wantErr: true,
			errMsg:  "round",
		},
		{
			name:    "empty interviewer",
			iv:      Interview{JobID: "j1", CandidateID: "c1", Round: 1, Interviewer: "", ScheduledAt: time.Now().Add(time.Hour)},
			wantErr: true,
			errMsg:  "interviewer",
		},
		{
			name:    "zero scheduled_at",
			iv:      Interview{JobID: "j1", CandidateID: "c1", Round: 1, Interviewer: "M1", ScheduledAt: time.Time{}},
			wantErr: true,
			errMsg:  "scheduled_at",
		},
		{
			name:    "invalid status",
			iv:      Interview{JobID: "j1", CandidateID: "c1", Round: 1, Interviewer: "M1", ScheduledAt: time.Now().Add(time.Hour), Status: "invalid"},
			wantErr: true,
			errMsg:  "status",
		},
		{
			name:    "negative score",
			iv:      Interview{JobID: "j1", CandidateID: "c1", Round: 1, Interviewer: "M1", ScheduledAt: time.Now().Add(time.Hour), Score: -1},
			wantErr: true,
			errMsg:  "score",
		},
		{
			name:    "score over 100",
			iv:      Interview{JobID: "j1", CandidateID: "c1", Round: 1, Interviewer: "M1", ScheduledAt: time.Now().Add(time.Hour), Score: 101},
			wantErr: true,
			errMsg:  "score",
		},
		{
			name:    "default status",
			iv:      Interview{JobID: "j1", CandidateID: "c1", Round: 1, Interviewer: "M1", ScheduledAt: time.Now().Add(time.Hour)},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.iv.Validate()
			if tt.wantErr {
				if err == nil {
					t.Errorf("expect error")
					return
				}
				if !strings.Contains(err.Error(), tt.errMsg) {
					t.Errorf("expect error contains %s, got %s", tt.errMsg, err.Error())
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
				if tt.iv.Status == "" {
					t.Errorf("expect default status set")
				}
			}
		})
	}
}

func TestInterviewCanTransition(t *testing.T) {
	tests := []struct {
		from string
		to   string
		want bool
	}{
		{InterviewStatusScheduled, InterviewStatusCompleted, true},
		{InterviewStatusScheduled, InterviewStatusCancelled, true},
		{InterviewStatusScheduled, InterviewStatusScheduled, false},
		{InterviewStatusCompleted, InterviewStatusScheduled, false},
		{InterviewStatusCompleted, InterviewStatusCancelled, false},
		{InterviewStatusCancelled, InterviewStatusCompleted, false},
		{InterviewStatusCancelled, InterviewStatusScheduled, false},
		{"invalid", InterviewStatusCompleted, false},
	}
	for _, tt := range tests {
		t.Run(tt.from+"->"+tt.to, func(t *testing.T) {
			got := CanTransitionInterview(tt.from, tt.to)
			if got != tt.want {
				t.Errorf("CanTransitionInterview(%s, %s) = %v, want %v", tt.from, tt.to, got, tt.want)
			}
		})
	}
}

func TestInterviewFilter(t *testing.T) {
	interviews := []*Interview{
		{JobID: "j1", CandidateID: "c1", Status: InterviewStatusScheduled},
		{JobID: "j1", CandidateID: "c2", Status: InterviewStatusCompleted},
		{JobID: "j2", CandidateID: "c1", Status: InterviewStatusCancelled},
	}
	tests := []struct {
		name   string
		filter InterviewFilter
		want   []bool
	}{
		{"match all", InterviewFilter{}, []bool{true, true, true}},
		{"job_id j1", InterviewFilter{JobID: "j1"}, []bool{true, true, false}},
		{"candidate_id c1", InterviewFilter{CandidateID: "c1"}, []bool{true, false, true}},
		{"status scheduled", InterviewFilter{Status: InterviewStatusScheduled}, []bool{true, false, false}},
		{"combined", InterviewFilter{JobID: "j1", Status: InterviewStatusScheduled}, []bool{true, false, false}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for i, iv := range interviews {
				got := tt.filter.Match(iv)
				if got != tt.want[i] {
					t.Errorf("interview %d: want %v, got %v", i, tt.want[i], got)
				}
			}
		})
	}
}

func TestOfferValidate(t *testing.T) {
	tests := []struct {
		name    string
		o       Offer
		wantErr bool
		errMsg  string
	}{
		{
			name:    "valid offer",
			o:       Offer{JobID: "j1", CandidateID: "c1", Salary: 30000, Status: OfferStatusPending, ExpiresAt: time.Now().Add(7 * 24 * time.Hour)},
			wantErr: false,
		},
		{
			name:    "empty job_id",
			o:       Offer{JobID: "", CandidateID: "c1", Salary: 30000, ExpiresAt: time.Now().Add(7 * 24 * time.Hour)},
			wantErr: true,
			errMsg:  "job_id",
		},
		{
			name:    "empty candidate_id",
			o:       Offer{JobID: "j1", CandidateID: "", Salary: 30000, ExpiresAt: time.Now().Add(7 * 24 * time.Hour)},
			wantErr: true,
			errMsg:  "candidate_id",
		},
		{
			name:    "zero salary",
			o:       Offer{JobID: "j1", CandidateID: "c1", Salary: 0, ExpiresAt: time.Now().Add(7 * 24 * time.Hour)},
			wantErr: true,
			errMsg:  "salary",
		},
		{
			name:    "negative salary",
			o:       Offer{JobID: "j1", CandidateID: "c1", Salary: -1, ExpiresAt: time.Now().Add(7 * 24 * time.Hour)},
			wantErr: true,
			errMsg:  "salary",
		},
		{
			name:    "zero expires_at",
			o:       Offer{JobID: "j1", CandidateID: "c1", Salary: 30000, ExpiresAt: time.Time{}},
			wantErr: true,
			errMsg:  "expires_at",
		},
		{
			name:    "invalid status",
			o:       Offer{JobID: "j1", CandidateID: "c1", Salary: 30000, ExpiresAt: time.Now().Add(7 * 24 * time.Hour), Status: "invalid"},
			wantErr: true,
			errMsg:  "status",
		},
		{
			name:    "default status",
			o:       Offer{JobID: "j1", CandidateID: "c1", Salary: 30000, ExpiresAt: time.Now().Add(7 * 24 * time.Hour)},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.o.Validate()
			if tt.wantErr {
				if err == nil {
					t.Errorf("expect error")
					return
				}
				if !strings.Contains(err.Error(), tt.errMsg) {
					t.Errorf("expect error contains %s, got %s", tt.errMsg, err.Error())
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
				if tt.o.Status == "" {
					t.Errorf("expect default status set")
				}
			}
		})
	}
}

func TestOfferCanTransition(t *testing.T) {
	tests := []struct {
		from string
		to   string
		want bool
	}{
		{OfferStatusPending, OfferStatusAccepted, true},
		{OfferStatusPending, OfferStatusDeclined, true},
		{OfferStatusPending, OfferStatusPending, false},
		{OfferStatusAccepted, OfferStatusPending, false},
		{OfferStatusAccepted, OfferStatusDeclined, false},
		{OfferStatusDeclined, OfferStatusAccepted, false},
		{OfferStatusDeclined, OfferStatusPending, false},
		{"invalid", OfferStatusAccepted, false},
	}
	for _, tt := range tests {
		t.Run(tt.from+"->"+tt.to, func(t *testing.T) {
			got := CanTransitionOffer(tt.from, tt.to)
			if got != tt.want {
				t.Errorf("CanTransitionOffer(%s, %s) = %v, want %v", tt.from, tt.to, got, tt.want)
			}
		})
	}
}

func TestOfferFilter(t *testing.T) {
	offers := []*Offer{
		{JobID: "j1", CandidateID: "c1", Status: OfferStatusPending},
		{JobID: "j1", CandidateID: "c2", Status: OfferStatusAccepted},
		{JobID: "j2", CandidateID: "c1", Status: OfferStatusDeclined},
	}
	tests := []struct {
		name   string
		filter OfferFilter
		want   []bool
	}{
		{"match all", OfferFilter{}, []bool{true, true, true}},
		{"job_id j1", OfferFilter{JobID: "j1"}, []bool{true, true, false}},
		{"candidate_id c1", OfferFilter{CandidateID: "c1"}, []bool{true, false, true}},
		{"status pending", OfferFilter{Status: OfferStatusPending}, []bool{true, false, false}},
		{"combined", OfferFilter{JobID: "j1", Status: OfferStatusPending}, []bool{true, false, false}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for i, o := range offers {
				got := tt.filter.Match(o)
				if got != tt.want[i] {
					t.Errorf("offer %d: want %v, got %v", i, tt.want[i], got)
				}
			}
		})
	}
}

func TestValidationError(t *testing.T) {
	e := NewValidationError("field", "message")
	if e.Error() != "field: message" {
		t.Errorf("expect 'field: message', got %s", e.Error())
	}
	e2 := NewValidationError("", "message")
	if e2.Error() != "message" {
		t.Errorf("expect 'message', got %s", e2.Error())
	}
	if !IsValidationError(e) {
		t.Error("expect IsValidationError true")
	}
	if IsValidationError(nil) {
		t.Error("expect IsValidationError false for nil")
	}
}
