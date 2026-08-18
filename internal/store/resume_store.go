package store

import "recruit/internal/model"

func (s *MemoryStore) CreateResume(r *model.Resume) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.resumes[r.ID] = r
	return nil
}

func (s *MemoryStore) GetResume(id string) (*model.Resume, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	r, ok := s.resumes[id]
	if !ok {
		return nil, ErrNotFound
	}
	return r, nil
}

func (s *MemoryStore) ListResumes() []*model.Resume {
	s.mu.RLock()
	defer s.mu.RUnlock()
	list := make([]*model.Resume, 0, len(s.resumes))
	for _, r := range s.resumes {
		list = append(list, r)
	}
	return list
}

func (s *MemoryStore) UpdateResume(r *model.Resume) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.resumes[r.ID]; !ok {
		return ErrNotFound
	}
	s.resumes[r.ID] = r
	return nil
}

func (s *MemoryStore) DeleteResume(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.resumes[id]; !ok {
		return ErrNotFound
	}
	delete(s.resumes, id)
	return nil
}
