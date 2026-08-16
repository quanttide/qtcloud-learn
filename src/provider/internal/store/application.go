package store

import (
	"time"

	"github.com/quanttide/qtcloud-learn-provider/internal/domain"
)

// ApplicationStore 是立项申请的内存存储（写后落盘，见 BaseStore.persist）。
type ApplicationStore struct {
	*BaseStore[domain.Application]
}

func NewApplicationStore() *ApplicationStore {
	return &ApplicationStore{BaseStore: NewBaseStore[domain.Application]("appl")}
}

func (s *ApplicationStore) Create(a *domain.Application) *domain.Application {
	s.mu.Lock()
	defer s.mu.Unlock()
	clone := *a
	clone.ID = s.nextID()
	if clone.Status == "" {
		clone.Status = "submitted"
	}
	if clone.CreatedAt == "" {
		clone.CreatedAt = time.Now().Format(time.RFC3339)
	}
	s.data[clone.ID] = &clone
	s.persist()
	return &clone
}

func (s *ApplicationStore) Update(a *domain.Application) (*domain.Application, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	existing, ok := s.data[a.ID]
	if !ok {
		return nil, false
	}
	existing.ProjectName = a.ProjectName
	existing.BlindSpot = a.BlindSpot
	existing.DemoPlan = a.DemoPlan
	existing.Direction = a.Direction
	existing.TeamMode = a.TeamMode
	existing.MemberNames = a.MemberNames
	existing.Status = a.Status
	s.persist()
	return existing, true
}
