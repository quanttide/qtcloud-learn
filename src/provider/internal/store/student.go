package store

import "github.com/quanttide/qtcloud-learn-provider/internal/domain"

// StudentStore 是 Student 的内存存储。
type StudentStore struct {
	*BaseStore[domain.Student]
}

func NewStudentStore() *StudentStore {
	return &StudentStore{BaseStore: NewBaseStore[domain.Student]("stu")}
}

func (s *StudentStore) Create(st *domain.Student) *domain.Student {
	s.mu.Lock()
	defer s.mu.Unlock()
	clone := *st
	clone.ID = s.nextID()
	s.data[clone.ID] = &clone
	return &clone
}

func (s *StudentStore) Update(st *domain.Student) (*domain.Student, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	existing, ok := s.data[st.ID]
	if !ok {
		return nil, false
	}
	existing.Name = st.Name
	existing.Email = st.Email
	existing.Avatar = st.Avatar
	existing.Plan = st.Plan
	return existing, true
}
