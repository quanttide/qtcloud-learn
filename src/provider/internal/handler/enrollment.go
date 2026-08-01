package handler

import (
	"github.com/quanttide/qtcloud-learn-provider/internal/domain"
	"github.com/quanttide/qtcloud-learn-provider/internal/store"
)

// EnrollmentHandler 提供 Enrollment 的标准 CRUD。
type EnrollmentHandler = CRUDHandler[domain.Enrollment]

// NewEnrollmentHandler 创建 Enrollment handler。
func NewEnrollmentHandler(s *store.EnrollmentStore) *EnrollmentHandler {
	return NewCRUDHandler(
		s,
		func(e *domain.Enrollment) string {
			if e.ClassID == "" {
				return "classId is required"
			}
			if e.StudentID == "" {
				return "studentId is required"
			}
			return ""
		},
		func(e *domain.Enrollment, id string) { e.ID = id },
	)
}
