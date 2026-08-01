package store

import "github.com/quanttide/qtcloud-learn-provider/internal/domain"

// AssessmentStore 是 Assessment 的内存存储。
type AssessmentStore struct {
	*BaseStore[domain.Assessment]
}

func NewAssessmentStore() *AssessmentStore {
	return &AssessmentStore{BaseStore: NewBaseStore[domain.Assessment]("assess")}
}

func (s *AssessmentStore) Create(a *domain.Assessment) *domain.Assessment {
	s.mu.Lock()
	defer s.mu.Unlock()
	clone := *a
	clone.ID = s.nextID()
	s.data[clone.ID] = &clone
	return &clone
}

func (s *AssessmentStore) Update(a *domain.Assessment) (*domain.Assessment, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	existing, ok := s.data[a.ID]
	if !ok {
		return nil, false
	}
	existing.ClassID = a.ClassID
	existing.Type = a.Type
	existing.Title = a.Title
	existing.FullScore = a.FullScore
	existing.PassScore = a.PassScore
	existing.Deadline = a.Deadline
	return existing, true
}
