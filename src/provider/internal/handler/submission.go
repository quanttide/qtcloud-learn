package handler

import (
	"github.com/quanttide/qtcloud-learn-provider/internal/domain"
	"github.com/quanttide/qtcloud-learn-provider/internal/store"
)

// SubmissionHandler 提供 Submission 的标准 CRUD。
type SubmissionHandler = CRUDHandler[domain.Submission]

// NewSubmissionHandler 创建 Submission handler。
func NewSubmissionHandler(s *store.SubmissionStore) *SubmissionHandler {
	return NewCRUDHandler(
		s,
		func(sub *domain.Submission) string {
			if sub.AssessmentID == "" {
				return "assessmentId is required"
			}
			if sub.StudentID == "" {
				return "studentId is required"
			}
			return ""
		},
		func(sub *domain.Submission, id string) { sub.ID = id },
	)
}
