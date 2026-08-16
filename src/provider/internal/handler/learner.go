package handler

// 学员档案 API（对齐原型 qt-students / 进度上报契约）：
//   GET  /api/learners                后台学员表
//   POST /api/courses/prod/progress   学员端进度上报（body {moduleId, name}）→ {max, last}
// 学员由上报进度 / 提交立项自动建档（按姓名 upsert，原型无登录态）。

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/quanttide/qtcloud-learn-provider/internal/domain"
)

// LearnerLister 后台学员表。
type LearnerLister interface {
	List() []*domain.Learner
}

// LearnerProgresser 进度上报所需。
type LearnerProgresser interface {
	UpsertByName(name, course string, progressMax, progressTotal int, projectName string) *domain.Learner
}

// LearnerHandler 学员档案。
type LearnerHandler struct {
	list LearnerLister
	upsert LearnerProgresser
}

func NewLearnerHandler(l LearnerLister, u LearnerProgresser) *LearnerHandler {
	return &LearnerHandler{list: l, upsert: u}
}

// List GET /api/learners 后台学员表（最近活跃在前）。
func (h *LearnerHandler) List(w http.ResponseWriter, r *http.Request) {
	learners := h.list.List()
	// 倒序：最近活跃在前
	for i, j := 0, len(learners)-1; i < j; i, j = i+1, j-1 {
		learners[i], learners[j] = learners[j], learners[i]
	}
	writeJSON(w, http.StatusOK, map[string]any{"learners": learners})
}

// ReportProgress POST /api/courses/prod/progress 上报进度。
// body: {moduleId: "m3", name: "张三"} → {max: 3, last: "m3"}
func (h *LearnerHandler) ReportProgress(w http.ResponseWriter, r *http.Request) {
	var req struct {
		ModuleID string `json:"moduleId"`
		Name     string `json:"name"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	req.ModuleID = strings.TrimSpace(req.ModuleID)
	req.Name = strings.TrimSpace(req.Name)
	if req.ModuleID == "" || req.Name == "" {
		writeError(w, http.StatusBadRequest, "moduleId and name are required")
		return
	}

	// moduleId → 模块序号（m1..m5 → 1..5）
	idx := moduleIndex(req.ModuleID)
	if idx == 0 {
		writeError(w, http.StatusBadRequest, "invalid moduleId (m1..m5)")
		return
	}
	// 上报进度：max 只增不减（upsert 语义）；5 模块课程
	h.upsert.UpsertByName(req.Name, "生产实习", idx, 5, "")
	writeJSON(w, http.StatusOK, map[string]any{"max": idx, "last": req.ModuleID})
}

func moduleIndex(moduleID string) int {
	if len(moduleID) != 2 || moduleID[0] != 'm' {
		return 0
	}
	n := moduleID[1] - '0'
	if n >= 1 && n <= 5 {
		return int(n)
	}
	return 0
}
