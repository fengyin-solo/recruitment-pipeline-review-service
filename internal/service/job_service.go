package service

import (
	"sort"
	"time"

	"recruit/internal/model"
	"recruit/pkg/idgen"
)

func (s *Service) CreateJob(input model.Job) (*model.Job, error) {
	input.ID = idgen.Hex()
	input.CreatedAt = time.Now()
	input.UpdatedAt = input.CreatedAt
	if err := input.Validate(); err != nil {
		return nil, err
	}
	if err := s.store.CreateJob(&input); err != nil {
		return nil, err
	}
	return &input, nil
}

func (s *Service) GetJob(id string) (*model.Job, error) {
	if id == "" {
		return nil, model.NewValidationError("id", "ID 不能为空")
	}
	return s.store.GetJob(id)
}

func (s *Service) ListJobs(filter model.JobFilter, page, size int) ([]*model.Job, int, error) {
	page, size = normalizePagination(page, size)
	all := s.store.ListJobs()
	matched := make([]*model.Job, 0, len(all))
	for _, j := range all {
		if filter.Match(j) {
			matched = append(matched, j)
		}
	}
	sort.Slice(matched, func(i, j int) bool {
		return matched[i].CreatedAt.After(matched[j].CreatedAt)
	})
	total := len(matched)
	start := (page - 1) * size
	if start >= total {
		return []*model.Job{}, total, nil
	}
	end := start + size
	if end > total {
		end = total
	}
	return matched[start:end], total, nil
}

func (s *Service) UpdateJob(id string, input model.Job) (*model.Job, error) {
	if id == "" {
		return nil, model.NewValidationError("id", "ID 不能为空")
	}
	j, err := s.store.GetJob(id)
	if err != nil {
		return nil, err
	}
	if input.Title != "" {
		j.Title = input.Title
	}
	if input.Department != "" {
		j.Department = input.Department
	}
	if input.Level != "" {
		j.Level = input.Level
	}
	if input.Headcount > 0 {
		j.Headcount = input.Headcount
	}
	if input.SalaryMin >= 0 {
		j.SalaryMin = input.SalaryMin
	}
	if input.SalaryMax >= 0 {
		j.SalaryMax = input.SalaryMax
	}
	if input.Status != "" {
		j.Status = input.Status
	}
	if input.Description != "" {
		j.Description = input.Description
	}
	j.UpdatedAt = time.Now()
	if err := j.Validate(); err != nil {
		return nil, err
	}
	if err := s.store.UpdateJob(j); err != nil {
		return nil, err
	}
	return j, nil
}

func (s *Service) DeleteJob(id string) error {
	if id == "" {
		return model.NewValidationError("id", "ID 不能为空")
	}
	return s.store.DeleteJob(id)
}

func (s *Service) BatchCloseJobs(ids []string) (int, error) {
	if len(ids) == 0 {
		return 0, model.NewValidationError("ids", "ID 列表不能为空")
	}
	count := 0
	for _, id := range ids {
		j, err := s.store.GetJob(id)
		if err != nil {
			s.log.Warnf("批量关闭职位时获取职位 %s 失败: %v", id, err)
			continue
		}
		if j.Status == model.JobStatusClosed || j.Status == model.JobStatusFilled {
			continue
		}
		j.Status = model.JobStatusClosed
		j.UpdatedAt = time.Now()
		if err := s.store.UpdateJob(j); err != nil {
			s.log.Warnf("批量关闭职位时更新职位 %s 失败: %v", id, err)
			continue
		}
		count++
	}
	return count, nil
}
