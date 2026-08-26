// Package domain 定义量潮学习云的领域模型。
//
// 模型对齐《量潮学习管理标准》（quanttide-learn/docs/specification/index.md），
// 核心模型：Learner × Criterion → Completion，以本包为唯一事实来源。
package domain

// Learner 学习者：本领域的学习者主体。
// 对齐 spec：id（本领域的学习者 ID）、user_id（预留，关联 auth 领域的用户 ID）。
type Learner struct {
	ID     string `json:"id"`
	UserID string `json:"user_id,omitempty"` // 预留：关联 auth 领域的用户 ID
}

// Criterion 验收标准：学习的原子单元，由课程档案定义。
// 对齐 spec：id（标准标识）、title（标准名称，人类可读）、
// description（具体规则描述）。
type Criterion struct {
	ID          string `json:"id"`
	Title       string `json:"title"`       // 标准名称（人类可读，用于展示与检索）
	Description string `json:"description"` // 具体规则描述
}

// Completion 完成记录：Learner 与 Criterion 的交叉记录，记录通过状态。
// 对齐 spec：learner_id → Learner.id、criterion_id → Criterion.id、
// status（completed | not_completed）、created_at / updated_at。
type Completion struct {
	ID          string `json:"id"`
	LearnerID   string `json:"learner_id"`   // → Learner.id
	CriterionID string `json:"criterion_id"` // → Criterion.id
	Status      string `json:"status"`       // "completed" / "not_completed"
	CreatedAt   string `json:"created_at,omitempty"`
	UpdatedAt   string `json:"updated_at,omitempty"`
}
