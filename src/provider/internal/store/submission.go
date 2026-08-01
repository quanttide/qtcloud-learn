package store

import "github.com/quanttide/qtcloud-learn-provider/internal/domain"

// SubmissionStore 是 Submission 的内存存储。
type SubmissionStore struct {
	*BaseStore[domain.Submission]
}

func NewSubmissionStore() *SubmissionStore {
	return &SubmissionStore{BaseStore: NewBaseStore[domain.Submission]("sub")}
}

func (s *SubmissionStore) Create(sub *domain.Submission) *domain.Submission {
	s.mu.Lock()
	defer s.mu.Unlock()
	clone := *sub
	clone.ID = s.nextID()
	s.data[clone.ID] = &clone
	return &clone
}

func (s *SubmissionStore) Update(sub *domain.Submission) (*domain.Submission, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	existing, ok := s.data[sub.ID]
	if !ok {
		return nil, false
	}
	existing.AssessmentID = sub.AssessmentID
	existing.StudentID = sub.StudentID
	existing.Status = sub.Status
	existing.Score = sub.Score
	existing.Comment = sub.Comment
	existing.SubmittedAt = sub.SubmittedAt
	return existing, true
}
