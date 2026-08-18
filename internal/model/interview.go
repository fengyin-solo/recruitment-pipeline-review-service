package model

import (
	"strings"
	"time"
)

const (
	InterviewStatusScheduled  = "scheduled"
	InterviewStatusCompleted  = "completed"
	InterviewStatusCancelled  = "cancelled"
)

var interviewTransitions = map[string]map[string]bool{
	InterviewStatusScheduled: {InterviewStatusCompleted: true, InterviewStatusCancelled: true},
}

func CanTransitionInterview(from, to string) bool {
	if m, ok := interviewTransitions[from]; ok {
		return m[to]
	}
	return false
}

func (x *Interview) Clone() *Interview {
	if x == nil {
		return nil
	}
	cp := *x
	return &cp
}

type Interview struct {
	ID          string    `json:"id"`
	JobID       string    `json:"job_id"`
	CandidateID string    `json:"candidate_id"`
	Round       int       `json:"round"`
	Interviewer string    `json:"interviewer"`
	ScheduledAt time.Time `json:"scheduled_at"`
	Status      string    `json:"status"`
	Feedback    string    `json:"feedback"`
	Score       int       `json:"score"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func (i *Interview) Validate() error {
	i.JobID = strings.TrimSpace(i.JobID)
	i.CandidateID = strings.TrimSpace(i.CandidateID)
	i.Interviewer = strings.TrimSpace(i.Interviewer)
	if i.JobID == "" {
		return NewValidationError("job_id", "职位 ID 不能为空")
	}
	if i.CandidateID == "" {
		return NewValidationError("candidate_id", "候选人 ID 不能为空")
	}
	if i.Round <= 0 {
		return NewValidationError("round", "面试轮次必须大于 0")
	}
	if i.Interviewer == "" {
		return NewValidationError("interviewer", "面试官不能为空")
	}
	if i.ScheduledAt.IsZero() {
		return NewValidationError("scheduled_at", "面试时间不能为空")
	}
	if i.Status == "" {
		i.Status = InterviewStatusScheduled
	}
	if i.Status != InterviewStatusScheduled && i.Status != InterviewStatusCompleted && i.Status != InterviewStatusCancelled {
		return NewValidationError("status", "面试状态不合法")
	}
	if i.Score < 0 || i.Score > 100 {
		return NewValidationError("score", "面试评分必须在 0-100 之间")
	}
	return nil
}

type InterviewFilter struct {
	JobID       string
	CandidateID string
	Status      string
}

func (f InterviewFilter) Match(i *Interview) bool {
	if f.JobID != "" && i.JobID != f.JobID {
		return false
	}
	if f.CandidateID != "" && i.CandidateID != f.CandidateID {
		return false
	}
	if f.Status != "" && i.Status != f.Status {
		return false
	}
	return true
}
