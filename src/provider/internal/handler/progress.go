package handler

import (
	"github.com/quanttide/qtcloud-learn-provider/internal/domain"
	"github.com/quanttide/qtcloud-learn-provider/internal/store"
)

// ProgressHandler 提供 Progress 的标准 CRUD。
type ProgressHandler = CRUDHandler[domain.Progress]

// NewProgressHandler 创建 Progress handler。
func NewProgressHandler(s *store.ProgressStore) *ProgressHandler {
	return NewCRUDHandler(
		s,
		func(p *domain.Progress) string {
			if p.StudentID == "" {
				return "studentId is required"
			}
			if p.ClassID == "" {
				return "classId is required"
			}
			return ""
		},
		func(p *domain.Progress, id string) { p.ID = id },
	)
}
