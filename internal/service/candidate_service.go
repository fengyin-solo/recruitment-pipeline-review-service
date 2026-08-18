package service

import (
	"sort"
	"time"

	"recruit/internal/model"
	"recruit/pkg/idgen"
)

func (s *Service) CreateCandidate(input model.Candidate) (*model.Candidate, error) {
	input.ID = idgen.Hex()
	input.CreatedAt = time.Now()
	input.UpdatedAt = input.CreatedAt
	if err := input.Validate(); err != nil {
		return nil, err
	}
	if _, err := s.store.GetCandidateByEmail(input.Email); err == nil {
		return nil, model.NewValidationError("email", "该邮箱已存在")
	}
	if err := s.store.CreateCandidate(&input); err != nil {
		return nil, err
	}
	return &input, nil
}

func (s *Service) GetCandidate(id string) (*model.Candidate, error) {
	if id == "" {
		return nil, model.NewValidationError("id", "ID 不能为空")
	}
	return s.store.GetCandidate(id)
}

func (s *Service) ListCandidates(filter model.CandidateFilter, page, size int) ([]*model.Candidate, int, error) {
	all := s.store.ListCandidates()
	matched := make([]*model.Candidate, 0, len(all))
	for _, c := range all {
		if filter.Match(c) {
			matched = append(matched, c)
		}
	}
	sort.Slice(matched, func(i, j int) bool {
		return matched[i].CreatedAt.After(matched[j].CreatedAt)
	})
	total := len(matched)
	start := (page - 1) * size
	if start >= total {
		return []*model.Candidate{}, total, nil
	}
	end := start + size
	if end > total {
		end = total
	}
	return matched[start:end], total, nil
}

func (s *Service) UpdateCandidate(id string, input model.Candidate) (*model.Candidate, error) {
	if id == "" {
		return nil, model.NewValidationError("id", "ID 不能为空")
	}
	c, err := s.store.GetCandidate(id)
	if err != nil {
		return nil, err
	}
	if input.Name != "" {
		c.Name = input.Name
	}
	if input.Email != "" {
		c.Email = input.Email
	}
	if input.Phone != "" {
		c.Phone = input.Phone
	}
	if input.YearsOfExperience >= 0 {
		c.YearsOfExperience = input.YearsOfExperience
	}
	if len(input.Skills) > 0 {
		c.Skills = input.Skills
	}
	if input.Status != "" {
		c.Status = input.Status
	}
	c.UpdatedAt = time.Now()
	if err := c.Validate(); err != nil {
		return nil, err
	}
	if err := s.store.UpdateCandidate(c); err != nil {
		return nil, err
	}
	return c, nil
}

func (s *Service) DeleteCandidate(id string) error {
	if id == "" {
		return model.NewValidationError("id", "ID 不能为空")
	}
	return s.store.DeleteCandidate(id)
}
