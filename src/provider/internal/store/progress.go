package store

import "github.com/quanttide/qtcloud-learn-provider/internal/domain"

// ProgressStore 是 Progress 的内存存储。
type ProgressStore struct {
	*BaseStore[domain.Progress]
}

func NewProgressStore() *ProgressStore {
	return &ProgressStore{BaseStore: NewBaseStore[domain.Progress]("prog")}
}

func (s *ProgressStore) Create(p *domain.Progress) *domain.Progress {
	s.mu.Lock()
	defer s.mu.Unlock()
	clone := *p
	clone.ID = s.nextID()
	s.data[clone.ID] = &clone
	return &clone
}

func (s *ProgressStore) Update(p *domain.Progress) (*domain.Progress, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	existing, ok := s.data[p.ID]
	if !ok {
		return nil, false
	}
	existing.StudentID = p.StudentID
	existing.ClassID = p.ClassID
	existing.Percent = p.Percent
	existing.Finished = p.Finished
	existing.UpdatedAt = p.UpdatedAt
	return existing, true
}
