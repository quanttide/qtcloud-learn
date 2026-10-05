package store

import "github.com/quanttide/qtcloud-learn-provider/internal/domain"

// TaskStore 是 Task 的内存存储（写后落盘）。
type TaskStore struct {
	*BaseStore[domain.Task]
}

func NewTaskStore() *TaskStore {
	return &TaskStore{BaseStore: NewBaseStore[domain.Task]("task")}
}

func (s *TaskStore) Create(t *domain.Task) *domain.Task {
	s.mu.Lock()
	defer s.mu.Unlock()
	clone := *t
	if clone.ID == "" {
		clone.ID = s.nextID()
	}
	s.data[clone.ID] = &clone
	s.persist()
	return &clone
}

func (s *TaskStore) Update(t *domain.Task) (*domain.Task, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	existing, ok := s.data[t.ID]
	if !ok {
		return nil, false
	}
	existing.Title = t.Title
	existing.Description = t.Description
	s.persist()
	return existing, true
}
