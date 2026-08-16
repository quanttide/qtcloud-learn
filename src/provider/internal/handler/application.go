package handler

// 立项申请 API：提交（学员端 POST）+ 列表/详情（后台查看）。
// v0.1 不做审批流，提交即 status=submitted 存档。

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/quanttide/qtcloud-learn-provider/internal/domain"
)

// ApplicationStorer 是 ApplicationHandler 所需存储。
type ApplicationStorer interface {
	List() []*domain.Application
	Get(id string) (*domain.Application, bool)
	Create(a *domain.Application) *domain.Application
	Update(a *domain.Application) (*domain.Application, bool)
}

// ApplicationHandler 立项申请。
type ApplicationHandler struct {
	store ApplicationStorer
}

func NewApplicationHandler(s ApplicationStorer) *ApplicationHandler {
	return &ApplicationHandler{store: s}
}

// Create POST /api/v1/applications 提交立项（微型创业姓名栏：
// teamMode=personal 时 MemberNames 为个人姓名；partner 时为队长+队员）。
func (h *ApplicationHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req struct {
		ProjectName  string   `json:"projectName"`
		BlindSpot    string   `json:"blindSpot"`
		DemoPlan     string   `json:"demoPlan"`
		Direction    string   `json:"direction"`
		TeamMode     string   `json:"teamMode"`
		MemberNames  []string `json:"memberNames"`
		StudentID    string   `json:"studentId"`
		StudentName  string   `json:"studentName"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	req.ProjectName = strings.TrimSpace(req.ProjectName)
	req.BlindSpot = strings.TrimSpace(req.BlindSpot)
	req.DemoPlan = strings.TrimSpace(req.DemoPlan)
	req.TeamMode = strings.TrimSpace(req.TeamMode)

	switch {
	case req.ProjectName == "":
		writeError(w, http.StatusBadRequest, "projectName is required")
		return
	case req.BlindSpot == "":
		writeError(w, http.StatusBadRequest, "blindSpot is required")
		return
	case req.DemoPlan == "":
		writeError(w, http.StatusBadRequest, "demoPlan is required")
		return
	case req.TeamMode != "personal" && req.TeamMode != "partner":
		writeError(w, http.StatusBadRequest, "teamMode must be personal or partner")
		return
	case len(req.MemberNames) == 0 || req.MemberNames[0] == "":
		writeError(w, http.StatusBadRequest, "memberNames is required")
		return
	}

	app := h.store.Create(&domain.Application{
		ProjectName: req.ProjectName,
		BlindSpot:   req.BlindSpot,
		DemoPlan:    req.DemoPlan,
		Direction:   req.Direction,
		TeamMode:    req.TeamMode,
		MemberNames: req.MemberNames,
		StudentID:   req.StudentID,
		StudentName: req.StudentName,
		Status:      "submitted",
		CreatedAt:   time.Now().Format(time.RFC3339),
	})
	writeJSON(w, http.StatusCreated, app)
}

// List GET /api/v1/applications 后台立项列表（新到旧）。
func (h *ApplicationHandler) List(w http.ResponseWriter, r *http.Request) {
	apps := h.store.List()
	// 倒序：最近提交在前
	for i, j := 0, len(apps)-1; i < j; i, j = i+1, j-1 {
		apps[i], apps[j] = apps[j], apps[i]
	}
	writeJSON(w, http.StatusOK, map[string]any{"applications": apps})
}

// Get GET /api/v1/applications/{id} 单条详情。
func (h *ApplicationHandler) Get(w http.ResponseWriter, r *http.Request) {
	app, ok := h.store.Get(r.PathValue("id"))
	if !ok {
		writeError(w, http.StatusNotFound, "application not found")
		return
	}
	writeJSON(w, http.StatusOK, app)
}
