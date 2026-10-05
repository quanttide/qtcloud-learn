package store

import (
	"github.com/quanttide/qtcloud-learn-provider/internal/domain"
)

// LearnerStore 是 Learner 的内存存储（写后落盘）。
type LearnerStore struct {
	*BaseStore[domain.Learner]
}

func NewLearnerStore() *LearnerStore {
	return &LearnerStore{BaseStore: NewBaseStore[domain.Learner]("lea")}
}

func (s *LearnerStore) Create(l *domain.Learner) *domain.Learner {
	s.mu.Lock()
	defer s.mu.Unlock()
	clone := *l
	clone.ID = s.nextID()
	s.data[clone.ID] = &clone
	s.persist()
	return &clone
}

func (s *LearnerStore) Update(l *domain.Learner) (*domain.Learner, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	existing, ok := s.data[l.ID]
	if !ok {
		return nil, false
	}
	existing.UserID = l.UserID
	existing.ScheduleID = l.ScheduleID
	s.persist()
	return existing, true
}
