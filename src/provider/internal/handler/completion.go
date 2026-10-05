package handler

import (
	"github.com/quanttide/qtcloud-learn-provider/internal/domain"
	"github.com/quanttide/qtcloud-learn-provider/internal/store"
)

// CompletionHandler 提供 Completion 的标准 CRUD。
type CompletionHandler = CRUDHandler[domain.Completion]

// NewCompletionHandler 创建 Completion handler。
// status 可选（缺省由 store 置为 "not_completed"），传入时须为合法枚举值。
func NewCompletionHandler(s *store.CompletionStore) *CompletionHandler {
	return NewCRUDHandler(
		s,
		func(c *domain.Completion) string {
			if c.LearnerID == "" {
				return "learner_id is required"
			}
			if c.TaskID == "" {
				return "task_id is required"
			}
			if c.Status != "" && c.Status != domain.CompletionStatusCompleted && c.Status != domain.CompletionStatusNotCompleted {
				return "status must be completed or not_completed"
			}
			return ""
		},
		func(c *domain.Completion, id string) { c.ID = id },
	)
}
