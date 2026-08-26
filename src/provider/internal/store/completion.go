package store

import (
	"time"

	"github.com/quanttide/qtcloud-learn-provider/internal/domain"
)

// CompletionStore 是 Completion 的内存存储（写后落盘）。
type CompletionStore struct {
	*BaseStore[domain.Completion]
}

func NewCompletionStore() *CompletionStore {
	return &CompletionStore{BaseStore: NewBaseStore[domain.Completion]("com")}
}

func (s *CompletionStore) Create(c *domain.Completion) *domain.Completion {
	s.mu.Lock()
	defer s.mu.Unlock()
	clone := *c
	clone.ID = s.nextID()
	if clone.Status == "" {
		clone.Status = "not_completed"
	}
	now := time.Now().Format(time.RFC3339)
	if clone.CreatedAt == "" {
		clone.CreatedAt = now
	}
	clone.UpdatedAt = now
	s.data[clone.ID] = &clone
	s.persist()
	return &clone
}

func (s *CompletionStore) Update(c *domain.Completion) (*domain.Completion, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	existing, ok := s.data[c.ID]
	if !ok {
		return nil, false
	}
	existing.LearnerID = c.LearnerID
	existing.CriterionID = c.CriterionID
	existing.Status = c.Status
	existing.UpdatedAt = time.Now().Format(time.RFC3339)
	s.persist()
	return existing, true
}
