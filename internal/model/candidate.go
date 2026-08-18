package model

import (
	"strings"
	"time"
)

const (
	CandidateStatusActive    = "active"
	CandidateStatusWithdrawn = "withdrawn"
)

func (x *Candidate) Clone() *Candidate {
	if x == nil {
		return nil
	}
	cp := *x
	cp.Skills = append([]string(nil), x.Skills...)
	return &cp
}

type Candidate struct {
	ID                string    `json:"id"`
	Name              string    `json:"name"`
	Email             string    `json:"email"`
	Phone             string    `json:"phone"`
	YearsOfExperience int       `json:"years_of_experience"`
	Skills            []string  `json:"skills"`
	Status            string    `json:"status"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
}

func (c *Candidate) Validate() error {
	c.Name = strings.TrimSpace(c.Name)
	c.Email = strings.TrimSpace(c.Email)
	c.Phone = strings.TrimSpace(c.Phone)
	if c.Name == "" {
		return NewValidationError("name", "候选人姓名不能为空")
	}
	if c.Email == "" {
		return NewValidationError("email", "邮箱不能为空")
	}
	if !strings.Contains(c.Email, "@") {
		return NewValidationError("email", "邮箱格式不正确")
	}
	if c.Phone == "" {
		return NewValidationError("phone", "电话不能为空")
	}
	if c.YearsOfExperience < 0 {
		return NewValidationError("years_of_experience", "工作年限不能为负数")
	}
	if c.Status == "" {
		c.Status = CandidateStatusActive
	}
	if c.Status != CandidateStatusActive && c.Status != CandidateStatusWithdrawn {
		return NewValidationError("status", "候选人状态不合法")
	}
	return nil
}

type CandidateFilter struct {
	Keyword string
	Status  string
}

func (f CandidateFilter) Match(c *Candidate) bool {
	if f.Status != "" && c.Status != f.Status {
		return false
	}
	if f.Keyword != "" {
		k := strings.ToLower(strings.TrimSpace(f.Keyword))
		if k != "" && !strings.Contains(strings.ToLower(c.Name), k) &&
			!strings.Contains(strings.ToLower(c.Email), k) {
			return false
		}
	}
	return true
}
