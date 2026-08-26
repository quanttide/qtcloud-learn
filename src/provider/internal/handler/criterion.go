package handler

import (
	"github.com/quanttide/qtcloud-learn-provider/internal/domain"
	"github.com/quanttide/qtcloud-learn-provider/internal/store"
)

// CriterionHandler 提供 Criterion 的标准 CRUD。
type CriterionHandler = CRUDHandler[domain.Criterion]

// NewCriterionHandler 创建 Criterion handler。
func NewCriterionHandler(s *store.CriterionStore) *CriterionHandler {
	return NewCRUDHandler(
		s,
		func(c *domain.Criterion) string {
			if c.Title == "" {
				return "title is required"
			}
			if c.Description == "" {
				return "description is required"
			}
			return ""
		},
		func(c *domain.Criterion, id string) { c.ID = id },
	)
}
