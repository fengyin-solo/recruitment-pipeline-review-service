package service

import (
	"sort"
	"time"

	"recruit/internal/model"
	"recruit/pkg/idgen"
)

func (s *Service) CreateOffer(input model.Offer) (*model.Offer, error) {
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
	if err := s.store.CreateOffer(&input); err != nil {
		return nil, err
	}
	return &input, nil
}

func (s *Service) GetOffer(id string) (*model.Offer, error) {
	if id == "" {
		return nil, model.NewValidationError("id", "ID 不能为空")
	}
	return s.store.GetOffer(id)
}

func (s *Service) ListOffers(filter model.OfferFilter, page, size int) ([]*model.Offer, int, error) {
	all := s.store.ListOffers()
	matched := make([]*model.Offer, 0, len(all))
	for _, o := range all {
		if filter.Match(o) {
			matched = append(matched, o)
		}
	}
	sort.Slice(matched, func(i, j int) bool {
		return matched[i].CreatedAt.After(matched[j].CreatedAt)
	})
	total := len(matched)
	start := (page - 1) * size
	if start >= total {
		return []*model.Offer{}, total, nil
	}
	end := start + size
	if end > total {
		end = total
	}
	return matched[start:end], total, nil
}

func (s *Service) UpdateOffer(id string, input model.Offer) (*model.Offer, error) {
	if id == "" {
		return nil, model.NewValidationError("id", "ID 不能为空")
	}
	o, err := s.store.GetOffer(id)
	if err != nil {
		return nil, err
	}
	if input.Salary > 0 {
		o.Salary = input.Salary
	}
	if input.Status != "" && input.Status != o.Status {
		if !model.CanTransitionOffer(o.Status, input.Status) {
			return nil, model.NewValidationError("status", "非法的状态流转: "+o.Status+" -> "+input.Status)
		}
		o.Status = input.Status
		if input.Status == model.OfferStatusAccepted {
			if err := s.handleOfferAccepted(o); err != nil {
				return nil, err
			}
		}
	}
	if !input.ExpiresAt.IsZero() {
		o.ExpiresAt = input.ExpiresAt
	}
	o.UpdatedAt = time.Now()
	if err := o.Validate(); err != nil {
		return nil, err
	}
	if err := s.store.UpdateOffer(o); err != nil {
		return nil, err
	}
	return o, nil
}

func (s *Service) handleOfferAccepted(o *model.Offer) error {
	job, err := s.store.GetJob(o.JobID)
	if err != nil {
		return err
	}
	acceptedCount := 0
	for _, offer := range s.store.ListOffers() {
		if offer.JobID == o.JobID && offer.Status == model.OfferStatusAccepted {
			acceptedCount++
		}
	}
	if acceptedCount >= job.Headcount {
		job.Status = model.JobStatusFilled
		job.UpdatedAt = time.Now()
		if err := s.store.UpdateJob(job); err != nil {
			return err
		}
	}
	return nil
}

func (s *Service) DeleteOffer(id string) error {
	if id == "" {
		return model.NewValidationError("id", "ID 不能为空")
	}
	return s.store.DeleteOffer(id)
}
