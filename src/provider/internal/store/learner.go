package store

import (
	"time"

	"github.com/quanttide/qtcloud-learn-provider/internal/domain"
)

// LearnerStore 是学员档案的内存存储（写后落盘）。
type LearnerStore struct {
	*BaseStore[domain.Learner]
}

func NewLearnerStore() *LearnerStore {
	return &LearnerStore{BaseStore: NewBaseStore[domain.Learner]("lea")}
}

// UpsertByName 按账号/姓名建档或更新（上报进度或提交立项时调用）。
// userID 非空时优先按账号匹配（统一账号体系下同名学员互不干扰），其次按姓名（历史数据兼容）。
// progressMax 只增不减；返回更新后的档案。
func (s *LearnerStore) UpsertByName(name, course string, progressMax, progressTotal int, projectName, userID string) *domain.Learner {
	s.mu.Lock()
	defer s.mu.Unlock()
	var existing *domain.Learner
	for _, l := range s.data {
		if userID != "" && l.UserID == userID {
			existing = l
			break
		}
		if userID == "" && l.Name == name {
			existing = l
			break
		}
	}
	now := time.Now().Format(time.RFC3339)
	if existing == nil {
		status := "在读"
		if progressMax >= progressTotal {
			status = "已完成"
		}
		learner := &domain.Learner{
			ID:            s.nextID(),
			UserID:        userID,
			Name:          name,
			Course:        course,
			ProgressMax:   progressMax,
			ProgressTotal: progressTotal,
			ActiveAt:      now,
			Status:        status,
			ProjectName:   projectName,
		}
		s.data[learner.ID] = learner
		s.persist()
		return learner
	}
	// 更新：进度只增不减、活跃时间刷新、立项关联、账号绑定
	if progressMax > existing.ProgressMax {
		existing.ProgressMax = progressMax
	}
	existing.ProgressTotal = progressTotal
	existing.ActiveAt = now
	if projectName != "" {
		existing.ProjectName = projectName
	}
	if userID != "" {
		existing.UserID = userID
	}
	if existing.ProgressMax >= existing.ProgressTotal {
		existing.Status = "已完成"
	}
	s.persist()
	return existing
}

// GetByName 按姓名查找档案。
func (s *LearnerStore) GetByName(name string) *domain.Learner {
	s.mu.RLock()
	defer s.mu.RUnlock()
	for _, l := range s.data {
		if l.Name == name {
			return l
		}
	}
	return nil
}
