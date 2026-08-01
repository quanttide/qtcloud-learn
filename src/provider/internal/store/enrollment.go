package store

import "github.com/quanttide/qtcloud-learn-provider/internal/domain"

// EnrollmentStore 是 Enrollment 的内存存储。
type EnrollmentStore struct {
	*BaseStore[domain.Enrollment]
}

func NewEnrollmentStore() *EnrollmentStore {
	return &EnrollmentStore{BaseStore: NewBaseStore[domain.Enrollment]("enr")}
}

func (s *EnrollmentStore) Create(e *domain.Enrollment) *domain.Enrollment {
	s.mu.Lock()
	defer s.mu.Unlock()
	clone := *e
	clone.ID = s.nextID()
	s.data[clone.ID] = &clone
	return &clone
}

func (s *EnrollmentStore) Update(e *domain.Enrollment) (*domain.Enrollment, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	existing, ok := s.data[e.ID]
	if !ok {
		return nil, false
	}
	existing.ClassID = e.ClassID
	existing.StudentID = e.StudentID
	existing.Status = e.Status
	existing.EnrolledAt = e.EnrolledAt
	return existing, true
}
