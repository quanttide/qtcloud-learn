package store

import "github.com/quanttide/qtcloud-learn-provider/internal/domain"

// SessionStore 是 Session 的内存存储。
type SessionStore struct {
	*BaseStore[domain.Session]
}

func NewSessionStore() *SessionStore {
	return &SessionStore{BaseStore: NewBaseStore[domain.Session]("sess")}
}

func (s *SessionStore) Create(sess *domain.Session) *domain.Session {
	s.mu.Lock()
	defer s.mu.Unlock()
	clone := *sess
	clone.ID = s.nextID()
	if clone.Attendances == nil {
		clone.Attendances = []domain.Attendance{}
	}
	s.data[clone.ID] = &clone
	return &clone
}

func (s *SessionStore) Update(sess *domain.Session) (*domain.Session, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	existing, ok := s.data[sess.ID]
	if !ok {
		return nil, false
	}
	existing.ClassID = sess.ClassID
	existing.LessonTitle = sess.LessonTitle
	existing.TeacherID = sess.TeacherID
	existing.StartTime = sess.StartTime
	existing.DurationMinutes = sess.DurationMinutes
	existing.Location = sess.Location
	existing.Status = sess.Status
	existing.Attendances = sess.Attendances
	return existing, true
}
