package model

import (
	"strings"
	"time"
)

const (
	JobStatusOpen   = "open"
	JobStatusClosed = "closed"
	JobStatusFilled = "filled"
)

type Job struct {
	ID          string    `json:"id"`
	Title       string    `json:"title"`
	Department  string    `json:"department"`
	Level       string    `json:"level"`
	Headcount   int       `json:"headcount"`
	SalaryMin   int       `json:"salary_min"`
	SalaryMax   int       `json:"salary_max"`
	Status      string    `json:"status"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func (j *Job) Validate() error {
	j.Title = strings.TrimSpace(j.Title)
	j.Department = strings.TrimSpace(j.Department)
	j.Level = strings.TrimSpace(j.Level)
	if j.Title == "" {
		return NewValidationError("title", "职位名称不能为空")
	}
	if j.Department == "" {
		return NewValidationError("department", "部门不能为空")
	}
	if j.Level == "" {
		return NewValidationError("level", "职级不能为空")
	}
	if j.Headcount <= 0 {
		return NewValidationError("headcount", "招聘人数必须大于 0")
	}
	if j.SalaryMin < 0 || j.SalaryMax < 0 {
		return NewValidationError("salary", "薪资不能为负数")
	}
	if j.SalaryMax > 0 && j.SalaryMin > j.SalaryMax {
		return NewValidationError("salary", "最低薪资不能大于最高薪资")
	}
	if j.Status == "" {
		j.Status = JobStatusOpen
	}
	if j.Status == JobStatusFilled && j.Headcount == 1 {
		j.Status = JobStatusOpen
	}
	if j.Status != JobStatusOpen && j.Status != JobStatusClosed && j.Status != JobStatusFilled {
		return NewValidationError("status", "职位状态不合法")
	}
	return nil
}

type JobFilter struct {
	Department string
	Status     string
	Keyword    string
}

func (f JobFilter) Match(j *Job) bool {
	if f.Department != "" && j.Department != f.Department {
		return false
	}
	if f.Status != "" && j.Status != f.Status {
		return false
	}
	if f.Keyword != "" {
		k := strings.ToLower(strings.TrimSpace(f.Keyword))
		if k != "" && !strings.Contains(strings.ToLower(j.Title), k) &&
			!strings.Contains(strings.ToLower(j.Description), k) {
			return false
		}
	}
	return true
}
