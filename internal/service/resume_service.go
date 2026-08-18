package service

import (
	"time"

	"recruit/internal/model"
	"recruit/pkg/idgen"
)

func (s *Service) CreateResume(input model.Resume) (*model.Resume, error) {
	input.ID = idgen.Hex()
	input.CreatedAt = time.Now()
	input.UpdatedAt = input.CreatedAt
	if err := input.Validate(); err != nil {
		return nil, err
	}
	if _, err := s.store.GetCandidate(input.CandidateID); err != nil {
		return nil, model.NewValidationError("candidate_id", "候选人不存在")
	}
	if err := s.store.CreateResume(&input); err != nil {
		return nil, err
	}
	return &input, nil
}

func (s *Service) GetResume(id string) (*model.Resume, error) {
	if id == "" {
		return nil, model.NewValidationError("id", "ID 不能为空")
	}
	return s.store.GetResume(id)
}

func (s *Service) ListResumes(filter model.ResumeFilter, page, size int) ([]*model.Resume, int, error) {
	all := s.store.ListResumes()
	matched := make([]*model.Resume, 0, len(all))
	for _, r := range all {
		if filter.Match(r) {
			matched = append(matched, r)
		}
	}
	total := len(matched)
	start, end := pageBounds(page, size, total)
	return matched[start:end], total, nil
}

func (s *Service) UpdateResume(id string, input model.Resume) (*model.Resume, error) {
	if id == "" {
		return nil, model.NewValidationError("id", "ID 不能为空")
	}
	r, err := s.store.GetResume(id)
	if err != nil {
		return nil, err
	}
	if input.Summary != "" {
		r.Summary = input.Summary
	}
	if input.Education != "" {
		r.Education = input.Education
	}
	if input.WorkExperience != "" {
		r.WorkExperience = input.WorkExperience
	}
	r.UpdatedAt = time.Now()
	if err := r.Validate(); err != nil {
		return nil, err
	}
	if err := s.store.UpdateResume(r); err != nil {
		return nil, err
	}
	return r, nil
}

func (s *Service) DeleteResume(id string) error {
	if id == "" {
		return model.NewValidationError("id", "ID 不能为空")
	}
	return s.store.DeleteResume(id)
}
