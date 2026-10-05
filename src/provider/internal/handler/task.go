package handler

import (
	"github.com/quanttide/qtcloud-learn-provider/internal/domain"
	"github.com/quanttide/qtcloud-learn-provider/internal/store"
)

// TaskHandler 提供 Task 的标准 CRUD。
type TaskHandler = CRUDHandler[domain.Task]

// NewTaskHandler 创建 Task handler。
func NewTaskHandler(s *store.TaskStore) *TaskHandler {
	return NewCRUDHandler(
		s,
		func(t *domain.Task) string {
			if t.Title == "" {
				return "title is required"
			}
			if t.Description == "" {
				return "description is required"
			}
			return ""
		},
		func(t *domain.Task, id string) { t.ID = id },
	)
}
