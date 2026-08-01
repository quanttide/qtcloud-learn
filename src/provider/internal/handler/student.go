package handler

import (
	"github.com/quanttide/qtcloud-learn-provider/internal/domain"
	"github.com/quanttide/qtcloud-learn-provider/internal/store"
)

// StudentHandler 提供 Student 的标准 CRUD。
type StudentHandler = CRUDHandler[domain.Student]

// NewStudentHandler 创建 Student handler。
func NewStudentHandler(s *store.StudentStore) *StudentHandler {
	return NewCRUDHandler(
		s,
		func(st *domain.Student) string {
			if st.Name == "" {
				return "name is required"
			}
			return ""
		},
		func(st *domain.Student, id string) { st.ID = id },
	)
}
