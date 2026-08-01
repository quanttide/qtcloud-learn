package handler

import (
	"github.com/quanttide/qtcloud-learn-provider/internal/domain"
	"github.com/quanttide/qtcloud-learn-provider/internal/store"
)

// AssessmentHandler 提供 Assessment 的标准 CRUD。
type AssessmentHandler = CRUDHandler[domain.Assessment]

// NewAssessmentHandler 创建 Assessment handler。
func NewAssessmentHandler(s *store.AssessmentStore) *AssessmentHandler {
	return NewCRUDHandler(
		s,
		func(a *domain.Assessment) string {
			if a.ClassID == "" {
				return "classId is required"
			}
			if a.Title == "" {
				return "title is required"
			}
			return ""
		},
		func(a *domain.Assessment, id string) { a.ID = id },
	)
}
