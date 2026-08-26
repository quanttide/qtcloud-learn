package store

import (
	"github.com/quanttide/qtcloud-learn-provider/internal/domain"
)

// CriterionStore 是 Criterion 的内存存储（写后落盘）。
type CriterionStore struct {
	*BaseStore[domain.Criterion]
}

func NewCriterionStore() *CriterionStore {
	return &CriterionStore{BaseStore: NewBaseStore[domain.Criterion]("cri")}
}

func (s *CriterionStore) Create(c *domain.Criterion) *domain.Criterion {
	s.mu.Lock()
	defer s.mu.Unlock()
	clone := *c
	clone.ID = s.nextID()
	s.data[clone.ID] = &clone
	s.persist()
	return &clone
}

func (s *CriterionStore) Update(c *domain.Criterion) (*domain.Criterion, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	existing, ok := s.data[c.ID]
	if !ok {
		return nil, false
	}
	existing.Title = c.Title
	existing.Description = c.Description
	s.persist()
	return existing, true
}
