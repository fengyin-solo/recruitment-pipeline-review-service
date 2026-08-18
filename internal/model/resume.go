package model

import (
	"strings"
	"time"
)

type Resume struct {
	ID             string    `json:"id"`
	CandidateID    string    `json:"candidate_id"`
	Summary        string    `json:"summary"`
	Education      string    `json:"education"`
	WorkExperience string    `json:"work_experience"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

func (r *Resume) Validate() error {
	r.CandidateID = strings.TrimSpace(r.CandidateID)
	r.Summary = strings.TrimSpace(r.Summary)
	r.Education = strings.TrimSpace(r.Education)
	if r.CandidateID == "" {
		return NewValidationError("candidate_id", "候选人 ID 不能为空")
	}
	if r.Summary == "" {
		return NewValidationError("summary", "个人简介不能为空")
	}
	if r.Education == "" {
		return NewValidationError("education", "教育经历不能为空")
	}
	return nil
}

type ResumeFilter struct {
	CandidateID string
}

func (f ResumeFilter) Match(r *Resume) bool {
	if f.CandidateID != "" && r.CandidateID != f.CandidateID {
		return false
	}
	return true
}
