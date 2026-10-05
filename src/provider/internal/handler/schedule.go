package handler

import (
	"github.com/quanttide/qtcloud-learn-provider/internal/domain"
	"github.com/quanttide/qtcloud-learn-provider/internal/store"
)

// ScheduleHandler 提供 Schedule 的标准 CRUD。
type ScheduleHandler = CRUDHandler[domain.Schedule]

// NewScheduleHandler 创建 Schedule handler。
func NewScheduleHandler(s *store.ScheduleStore) *ScheduleHandler {
	return NewCRUDHandler(
		s,
		func(schedule *domain.Schedule) string {
			if schedule.Title == "" {
				return "title is required"
			}
			if schedule.Tasks == nil {
				return "tasks is required"
			}
			for _, task := range schedule.Tasks {
				if task.ID == "" {
					return "tasks.id is required"
				}
				if task.Title == "" {
					return "tasks.title is required"
				}
				if task.Description == "" {
					return "tasks.description is required"
				}
			}
			return ""
		},
		func(schedule *domain.Schedule, id string) { schedule.ID = id },
	)
}
