package model

import (
	"strings"
	"time"
)

const (
	OfferStatusPending  = "pending"
	OfferStatusAccepted = "accepted"
	OfferStatusDeclined = "declined"
)

var offerTransitions = map[string]map[string]bool{
	OfferStatusPending: {OfferStatusAccepted: true, OfferStatusDeclined: true},
}

func CanTransitionOffer(from, to string) bool {
	if m, ok := offerTransitions[from]; ok {
		return m[to]
	}
	return false
}

type Offer struct {
	ID          string    `json:"id"`
	JobID       string    `json:"job_id"`
	CandidateID string    `json:"candidate_id"`
	Salary      int       `json:"salary"`
	Status      string    `json:"status"`
	ExpiresAt   time.Time `json:"expires_at"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func (o *Offer) Validate() error {
	o.JobID = strings.TrimSpace(o.JobID)
	o.CandidateID = strings.TrimSpace(o.CandidateID)
	if o.JobID == "" {
		return NewValidationError("job_id", "职位 ID 不能为空")
	}
	if o.CandidateID == "" {
		return NewValidationError("candidate_id", "候选人 ID 不能为空")
	}
	if o.Salary <= 0 {
		return NewValidationError("salary", "薪资必须大于 0")
	}
	if o.Status == "" {
		o.Status = OfferStatusPending
	}
	if o.Status != OfferStatusPending && o.Status != OfferStatusAccepted && o.Status != OfferStatusDeclined {
		return NewValidationError("status", "Offer 状态不合法")
	}
	if o.ExpiresAt.IsZero() {
		return NewValidationError("expires_at", "过期时间不能为空")
	}
	return nil
}

type OfferFilter struct {
	JobID       string
	CandidateID string
	Status      string
}

func (f OfferFilter) Match(o *Offer) bool {
	if f.JobID != "" && o.JobID != f.JobID {
		return false
	}
	if f.CandidateID != "" && o.CandidateID != f.CandidateID {
		return false
	}
	if f.Status != "" && o.Status != f.Status {
		return false
	}
	return true
}
