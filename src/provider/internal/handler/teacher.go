package handler

import (
	"github.com/quanttide/qtcloud-learn-provider/internal/domain"
	"github.com/quanttide/qtcloud-learn-provider/internal/store"
)

// TeacherHandler 提供 Teacher 的标准 CRUD。
type TeacherHandler = CRUDHandler[domain.Teacher]

// NewTeacherHandler 创建 Teacher handler。
func NewTeacherHandler(s *store.TeacherStore) *TeacherHandler {
	return NewCRUDHandler(
		s,
		func(t *domain.Teacher) string {
			if t.Name == "" {
				return "name is required"
			}
			return ""
		},
		func(t *domain.Teacher, id string) { t.ID = id },
	)
}
