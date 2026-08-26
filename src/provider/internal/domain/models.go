// Package domain 定义量潮学习云的领域模型。
//
// 模型对齐《量潮学习管理标准》（quanttide-learn/docs/specification/index.md），
// 核心模型：Learner × criterion_id（→ 课程域 Criterion.id，同源直连）→ Completion，以本包为唯一事实来源。
// 本领域不设验收标准实体——Criterion 由课程云一等实体承载，学习云只保存跨域引用并负责记录与聚合。
package domain

// Learner 学习者：本领域的学习者主体。
// 对齐 spec：id（本领域的学习者 ID）、user_id（预留，关联 auth 领域的用户 ID）。
type Learner struct {
	ID     string `json:"id"`
	UserID string `json:"user_id,omitempty"` // 预留：关联 auth 领域的用户 ID
}

// Completion 完成记录：Learner 与课程域验收标准的交叉记录，记录通过状态。
// criterion_id 同源直连课程域 Criterion.id；created_at / updated_at 为记录时间戳。
type Completion struct {
	ID          string `json:"id"`
	LearnerID   string `json:"learner_id"`   // → Learner.id
	CriterionID string `json:"criterion_id"` // → 课程域 Criterion.id（同源直连）
	Status      string `json:"status"`       // "completed" / "not_completed"
	CreatedAt   string `json:"created_at,omitempty"`
	UpdatedAt   string `json:"updated_at,omitempty"`
}
