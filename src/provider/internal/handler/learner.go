package handler

import (
	"github.com/quanttide/qtcloud-learn-provider/internal/domain"
	"github.com/quanttide/qtcloud-learn-provider/internal/store"
)

// LearnerHandler 提供 Learner 的标准 CRUD。
type LearnerHandler = CRUDHandler[domain.Learner]

// NewLearnerHandler 创建 Learner handler。
func NewLearnerHandler(s *store.LearnerStore) *LearnerHandler {
	return NewCRUDHandler(
		s,
		func(l *domain.Learner) string { return "" }, // 无必填字段（id 自动生成，user_id 预留）
		func(l *domain.Learner, id string) { l.ID = id },
	)
}
