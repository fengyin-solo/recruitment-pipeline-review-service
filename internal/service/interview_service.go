package service

import (
	"sort"
	"time"

	"recruit/internal/model"
	"recruit/pkg/idgen"
)

func (s *Service) CreateInterview(input model.Interview) (*model.Interview, error) {
	input.ID = idgen.Hex()
	input.CreatedAt = time.Now()
	input.UpdatedAt = input.CreatedAt
	if err := input.Validate(); err != nil {
		return nil, err
	}
	if _, err := s.store.GetJob(input.JobID); err != nil {
		return nil, model.NewValidationError("job_id", "职位不存在")
	}
	if _, err := s.store.GetCandidate(input.CandidateID); err != nil {
		return nil, model.NewValidationError("candidate_id", "候选人不存在")
	}
	if err := s.store.CreateInterview(&input); err != nil {
		return nil, err
	}
	return &input, nil
}

func (s *Service) GetInterview(id string) (*model.Interview, error) {
	if id == "" {
		return nil, model.NewValidationError("id", "ID 不能为空")
	}
	return s.store.GetInterview(id)
}

func (s *Service) ListInterviews(filter model.InterviewFilter, page, size int) ([]*model.Interview, int, error) {
	all := s.store.ListInterviews()
	matched := make([]*model.Interview, 0, len(all))
	for _, i := range all {
		if filter.Match(i) {
			matched = append(matched, i)
		}
	}
	sort.Slice(matched, func(i, j int) bool {
		return matched[i].CreatedAt.After(matched[j].CreatedAt)
	})
	total := len(matched)
	start := (page - 1) * size
	if start >= total {
		return []*model.Interview{}, total, nil
	}
	end := start + size
	if end > total {
		end = total
	}
	return matched[start:end], total, nil
}

func (s *Service) UpdateInterview(id string, input model.Interview) (*model.Interview, error) {
	if id == "" {
		return nil, model.NewValidationError("id", "ID 不能为空")
	}
	iv, err := s.store.GetInterview(id)
	if err != nil {
		return nil, err
	}
	if input.Round > 0 {
		iv.Round = input.Round
	}
	if input.Interviewer != "" {
		iv.Interviewer = input.Interviewer
	}
	if !input.ScheduledAt.IsZero() {
		iv.ScheduledAt = input.ScheduledAt
	}
	if input.Status != "" && input.Status != iv.Status {
		if !model.CanTransitionInterview(iv.Status, input.Status) {
			return nil, model.NewValidationError("status", "非法的状态流转: "+iv.Status+" -> "+input.Status)
		}
		iv.Status = input.Status
	}
	if input.Feedback != "" {
		iv.Feedback = input.Feedback
	}
	if input.Score >= 0 {
		iv.Score = input.Score
	}
	iv.UpdatedAt = time.Now()
	if err := iv.Validate(); err != nil {
		return nil, err
	}
	if err := s.store.UpdateInterview(iv); err != nil {
		return nil, err
	}
	return iv, nil
}

func (s *Service) DeleteInterview(id string) error {
	if id == "" {
		return model.NewValidationError("id", "ID 不能为空")
	}
	return s.store.DeleteInterview(id)
}

func (s *Service) BatchUpdateInterviewStatus(ids []string, status string) (int, error) {
	if len(ids) == 0 {
		return 0, model.NewValidationError("ids", "ID 列表不能为空")
	}
	if status != model.InterviewStatusCompleted && status != model.InterviewStatusCancelled {
		return 0, model.NewValidationError("status", "不支持的批量更新状态")
	}
	count := 0
	for _, id := range ids {
		iv, err := s.store.GetInterview(id)
		if err != nil {
			s.log.Warnf("批量更新面试状态时获取面试 %s 失败: %v", id, err)
			continue
		}
		if iv.Status == status {
			continue
		}
		if !model.CanTransitionInterview(iv.Status, status) {
			s.log.Warnf("批量更新面试状态时面试 %s 状态无法从 %s 转为 %s", id, iv.Status, status)
			continue
		}
		iv.Status = status
		iv.UpdatedAt = time.Now()
		if err := s.store.UpdateInterview(iv); err != nil {
			s.log.Warnf("批量更新面试状态时更新面试 %s 失败: %v", id, err)
			continue
		}
		count++
	}
	return count, nil
}
