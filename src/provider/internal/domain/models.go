// Package domain 定义量潮学习云的领域模型。
//
// 模型对齐《量潮学习管理标准》（quanttide-learn/docs/specification/index.md），
// 单一事实源为 quanttide-learn-toolkit，本包只做 type alias，避免跨仓复制模型。
package domain

import toolkit "github.com/quanttide/quanttide-learn-toolkit/packages/go/pkg"

type CompletionStatus = toolkit.CompletionStatus

const (
	CompletionStatusCompleted    = toolkit.CompletionStatusCompleted
	CompletionStatusNotCompleted = toolkit.CompletionStatusNotCompleted
)

type Learner = toolkit.Learner
type Completion = toolkit.Completion
type Task = toolkit.Task
type Schedule = toolkit.Schedule
