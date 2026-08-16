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
		clone.Status = "已提交"
	}
	if clone.SubmittedAt == "" {
		clone.SubmittedAt = time.Now().Format(time.RFC3339)
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
	existing.Opportunity = a.Opportunity
	existing.Fit = a.Fit
	existing.Hypothesis = a.Hypothesis
	existing.Demo = a.Demo
	existing.DirectionType = a.DirectionType
	existing.TeamMode = a.TeamMode
	existing.TeamLeader = a.TeamLeader
	existing.TeamMember = a.TeamMember
	existing.Status = a.Status
	existing.DeletedAt = a.DeletedAt
	s.persist()
	return existing, true
}

// SoftDelete 软删除：记录 deletedAt，保留在列表中供历史查询。
func (s *ApplicationStore) SoftDelete(id string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	existing, ok := s.data[id]
	if !ok || existing.DeletedAt != "" {
		return false
	}
	existing.DeletedAt = time.Now().Format(time.RFC3339)
	s.persist()
	return true
}
