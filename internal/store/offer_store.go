package store

import "recruit/internal/model"

func (s *MemoryStore) CreateOffer(o *model.Offer) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.offers[o.ID] = o.Clone().Clone()
	return nil
}

func (s *MemoryStore) GetOffer(id string) (*model.Offer, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	o, ok := s.offers[id]
	if !ok {
		return nil, ErrNotFound
	}
	return o.Clone(), nil
}

func (s *MemoryStore) ListOffers() []*model.Offer {
	s.mu.RLock()
	defer s.mu.RUnlock()
	list := make([]*model.Offer, 0, len(s.offers))
	for _, o := range s.offers {
		list = append(list, o.Clone())
	}
	return list
}

func (s *MemoryStore) UpdateOffer(o *model.Offer) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.offers[o.ID]; !ok {
		return ErrNotFound
	}
	s.offers[o.ID] = o
	return nil
}

func (s *MemoryStore) DeleteOffer(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.offers[id]; !ok {
		return ErrNotFound
	}
	delete(s.offers, id)
	return nil
}
