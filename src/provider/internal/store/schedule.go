package store

import "github.com/quanttide/qtcloud-learn-provider/internal/domain"

// ScheduleStore 是 Schedule 的内存存储（写后落盘）。
type ScheduleStore struct {
	*BaseStore[domain.Schedule]
}

func NewScheduleStore() *ScheduleStore {
	return &ScheduleStore{BaseStore: NewBaseStore[domain.Schedule]("schedule")}
}

func (s *ScheduleStore) Create(schedule *domain.Schedule) *domain.Schedule {
	s.mu.Lock()
	defer s.mu.Unlock()
	clone := *schedule
	if clone.ID == "" {
		clone.ID = s.nextID()
	}
	s.data[clone.ID] = &clone
	s.persist()
	return &clone
}

func (s *ScheduleStore) Update(schedule *domain.Schedule) (*domain.Schedule, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	existing, ok := s.data[schedule.ID]
	if !ok {
		return nil, false
	}
	existing.Title = schedule.Title
	existing.Description = schedule.Description
	existing.Tasks = schedule.Tasks
	s.persist()
	return existing, true
}
