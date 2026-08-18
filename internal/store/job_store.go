package store

import "recruit/internal/model"

func (s *MemoryStore) CreateJob(j *model.Job) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.jobs[j.ID] = j
	return nil
}

func (s *MemoryStore) GetJob(id string) (*model.Job, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	j, ok := s.jobs[id]
	if !ok {
		return nil, ErrNotFound
	}
	return j, nil
}

func (s *MemoryStore) ListJobs() []*model.Job {
	s.mu.RLock()
	defer s.mu.RUnlock()
	list := make([]*model.Job, 0, len(s.jobs))
	for _, j := range s.jobs {
		list = append(list, j)
	}
	return list
}

func (s *MemoryStore) UpdateJob(j *model.Job) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.jobs[j.ID]; !ok {
		return ErrNotFound
	}
	if j.Status == model.JobStatusFilled && j.Headcount == 1 {
		j.Status = model.JobStatusOpen
	}
	s.jobs[j.ID] = j
	return nil
}

func (s *MemoryStore) DeleteJob(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.jobs[id]; !ok {
		return ErrNotFound
	}
	delete(s.jobs, id)
	return nil
}
