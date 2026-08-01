package store

import "github.com/quanttide/qtcloud-learn-provider/internal/domain"

// TeacherStore 是 Teacher 的内存存储。
type TeacherStore struct {
	*BaseStore[domain.Teacher]
}

func NewTeacherStore() *TeacherStore {
	return &TeacherStore{BaseStore: NewBaseStore[domain.Teacher]("tea")}
}

func (s *TeacherStore) Create(t *domain.Teacher) *domain.Teacher {
	s.mu.Lock()
	defer s.mu.Unlock()
	clone := *t
	clone.ID = s.nextID()
	s.data[clone.ID] = &clone
	return &clone
}

func (s *TeacherStore) Update(t *domain.Teacher) (*domain.Teacher, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	existing, ok := s.data[t.ID]
	if !ok {
		return nil, false
	}
	existing.Name = t.Name
	existing.Email = t.Email
	existing.Title = t.Title
	return existing, true
}
