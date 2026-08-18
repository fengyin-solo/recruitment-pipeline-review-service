package store

import (
	"sync"

	"recruit/internal/model"
)

type MemoryStore struct {
	mu          sync.RWMutex
	jobs        map[string]*model.Job
	candidates  map[string]*model.Candidate
	resumes     map[string]*model.Resume
	interviews  map[string]*model.Interview
	offers      map[string]*model.Offer
}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		jobs:        make(map[string]*model.Job),
		candidates:  make(map[string]*model.Candidate),
		resumes:     make(map[string]*model.Resume),
		interviews:  make(map[string]*model.Interview),
		offers:      make(map[string]*model.Offer),
	}
}

var _ Store = (*MemoryStore)(nil)
