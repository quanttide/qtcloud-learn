package handler

// 立项申请 API（对齐原型 qt-proposals 契约）：
//   POST   /api/proposals            提交立项（学员端，5 问 + 方向类型 + 组队姓名）
//   GET    /api/proposals            后台列表（不含已删除）
//   GET    /api/proposals/history    历史记录（软删除留痕）
//   DELETE /api/proposals/{id}       软删除
// v0.1 不做审批流，status 固定"已提交"。

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/quanttide/qtcloud-learn-provider/internal/domain"
)

// ApplicationStorer 是 ApplicationHandler 所需存储。
type ApplicationStorer interface {
	List() []*domain.Application
	Get(id string) (*domain.Application, bool)
	Create(a *domain.Application) *domain.Application
	Update(a *domain.Application) (*domain.Application, bool)
	SoftDelete(id string) bool
}

// LearnerUpserter 立项/进度上报时自动建档学员。
type LearnerUpserter interface {
	UpsertByName(name, course string, progressMax, progressTotal int, projectName string) *domain.Learner
}

// ApplicationHandler 立项申请。
type ApplicationHandler struct {
	store  ApplicationStorer
	learner LearnerUpserter
}

func NewApplicationHandler(s ApplicationStorer, l LearnerUpserter) *ApplicationHandler {
	return &ApplicationHandler{store: s, learner: l}
}

// Create POST /api/proposals 提交立项。
// 姓名栏（赵追加需求）：teamMode=personal → TeamLeader=个人姓名；partner → TeamLeader=队长 + TeamMember=队员。
func (h *ApplicationHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req struct {
		ProjectName   string `json:"projectName"`
		Opportunity   string `json:"opportunity"`
		Fit           string `json:"fit"`
		Hypothesis    string `json:"hypothesis"`
		Demo          string `json:"demo"`
		DirectionType string `json:"directionType"`
		TeamMode      string `json:"teamMode"`
		TeamLeader    string `json:"teamLeader"`
		TeamMember    string `json:"teamMember"`
		StudentName   string `json:"studentName"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	req.ProjectName = strings.TrimSpace(req.ProjectName)
	req.TeamMode = strings.TrimSpace(req.TeamMode)
	req.TeamLeader = strings.TrimSpace(req.TeamLeader)

	switch {
	case req.ProjectName == "":
		writeError(w, http.StatusBadRequest, "projectName is required")
		return
	case req.TeamMode != "personal" && req.TeamMode != "partner":
		writeError(w, http.StatusBadRequest, "teamMode must be personal or partner")
		return
	case req.TeamLeader == "":
		writeError(w, http.StatusBadRequest, "teamLeader is required")
		return
	}
	// 当前学员身份：优先表单姓名，兜底队长姓名（对齐原型 qt-learner 语义）
	learnerName := strings.TrimSpace(req.StudentName)
	if learnerName == "" {
		learnerName = req.TeamLeader
	}

	app := h.store.Create(&domain.Application{
		ProjectName:   req.ProjectName,
		Opportunity:   req.Opportunity,
		Fit:           req.Fit,
		Hypothesis:    req.Hypothesis,
		Demo:          req.Demo,
		DirectionType: req.DirectionType,
		TeamMode:      req.TeamMode,
		TeamLeader:    req.TeamLeader,
		TeamMember:    strings.TrimSpace(req.TeamMember),
		StudentName:   learnerName,
		Status:        "已提交",
	})
	// 自动建档学员：进度置满（提交立项 = 进度 100%）
	if h.learner != nil {
		h.learner.UpsertByName(learnerName, "生产实习", 5, 5, app.ProjectName)
	}
	writeJSON(w, http.StatusCreated, app)
}

// List GET /api/proposals 后台列表（最近提交在前，不含已删除）。
func (h *ApplicationHandler) List(w http.ResponseWriter, r *http.Request) {
	apps := h.store.List()
	out := make([]*domain.Application, 0, len(apps))
	for _, a := range apps {
		if a.DeletedAt == "" {
			out = append(out, a)
		}
	}
	// 倒序：最近提交在前
	for i, j := 0, len(out)-1; i < j; i, j = i+1, j-1 {
		out[i], out[j] = out[j], out[i]
	}
	writeJSON(w, http.StatusOK, map[string]any{"proposals": out})
}

// History GET /api/proposals/history 软删除历史（含删除时间）。
func (h *ApplicationHandler) History(w http.ResponseWriter, r *http.Request) {
	apps := h.store.List()
	out := make([]*domain.Application, 0, len(apps))
	for _, a := range apps {
		if a.DeletedAt != "" {
			out = append(out, a)
		}
	}
	// 倒序：最近删除在前
	for i, j := 0, len(out)-1; i < j; i, j = i+1, j-1 {
		out[i], out[j] = out[j], out[i]
	}
	writeJSON(w, http.StatusOK, map[string]any{"history": out})
}

// Delete DELETE /api/proposals/{id} 软删除。
func (h *ApplicationHandler) Delete(w http.ResponseWriter, r *http.Request) {
	if !h.store.SoftDelete(r.PathValue("id")) {
		writeError(w, http.StatusNotFound, "proposal not found")
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"success": true})
}
