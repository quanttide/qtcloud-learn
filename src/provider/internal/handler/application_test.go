package handler

// 立项申请 API 测试（对齐原型契约 /api/proposals）：
// 提交（5 问 + 姓名栏）/ 列表 / 软删除 + 历史 / 校验 / 学员自动建档。

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/quanttide/qtcloud-learn-provider/internal/store"
)

func setupApplicationMux() (*http.ServeMux, *store.ApplicationStore, *store.LearnerStore) {
	appStore := store.NewApplicationStore()
	learnerStore := store.NewLearnerStore()
	apph := NewApplicationHandler(appStore, learnerStore)
	learnerh := NewLearnerHandler(learnerStore, learnerStore)
	mux := http.NewServeMux()
	mux.HandleFunc("POST /api/proposals", apph.Create)
	mux.HandleFunc("GET /api/proposals", apph.List)
	mux.HandleFunc("GET /api/proposals/history", apph.History)
	mux.HandleFunc("DELETE /api/proposals/{id}", apph.Delete)
	mux.HandleFunc("GET /api/learners", learnerh.List)
	mux.HandleFunc("POST /api/courses/prod/progress", learnerh.ReportProgress)
	return mux, appStore, learnerStore
}

func TestProposalCreatePersonal(t *testing.T) {
	mux, _, learnerStore := setupApplicationMux()

	// 个人独立：teamLeader=个人姓名
	body := `{"projectName":"校园选课助手","opportunity":"选课信息分散","fit":"量潮有数据能力","hypothesis":"聚合查询能提升效率","demo":"查询页原型","directionType":"内容","teamMode":"personal","teamLeader":"张三","studentName":"张三"}`
	req := httptest.NewRequest("POST", "/api/proposals", bytes.NewBufferString(body))
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Fatalf("status = %d, body=%s", w.Code, w.Body.String())
	}
	var app struct {
		ID          string `json:"id"`
		Status      string `json:"status"`
		TeamLeader  string `json:"teamLeader"`
		Opportunity string `json:"opportunity"`
		SubmittedAt string `json:"submittedAt"`
	}
	json.NewDecoder(w.Body).Decode(&app)
	if app.ID == "" || app.Status != "已提交" || app.TeamLeader != "张三" {
		t.Errorf("app = %+v", app)
	}
	if app.Opportunity != "选课信息分散" || app.SubmittedAt == "" {
		t.Errorf("app fields = %+v", app)
	}
	// 学员自动建档（提交立项 → 进度 100%）
	learner := learnerStore.GetByName("张三")
	if learner == nil {
		t.Fatal("learner not auto-created")
	}
	if learner.ProgressMax != 5 || learner.ProgressTotal != 5 || learner.Status != "已完成" {
		t.Errorf("learner = %+v", learner)
	}
	if learner.ProjectName != "校园选课助手" {
		t.Errorf("learner project = %s", learner.ProjectName)
	}
}

func TestProposalCreatePartner(t *testing.T) {
	mux, _, _ := setupApplicationMux()

	// 搭档：队长+队员姓名
	body := `{"projectName":"盲区地图","teamMode":"partner","teamLeader":"李四","teamMember":"王五、赵六","studentName":"李四"}`
	req := httptest.NewRequest("POST", "/api/proposals", bytes.NewBufferString(body))
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)
	if w.Code != http.StatusCreated {
		t.Fatalf("status = %d", w.Code)
	}
	var app struct {
		TeamLeader string `json:"teamLeader"`
		TeamMember string `json:"teamMember"`
		StudentName string `json:"studentName"`
	}
	json.NewDecoder(w.Body).Decode(&app)
	if app.TeamLeader != "李四" || app.TeamMember != "王五、赵六" {
		t.Errorf("team = %+v", app)
	}
	if app.StudentName != "李四" {
		t.Errorf("studentName = %s, want 李四（队长）", app.StudentName)
	}
}

func TestProposalValidation(t *testing.T) {
	mux, _, _ := setupApplicationMux()
	cases := []struct {
		name string
		body string
	}{
		{"缺项目名称", `{"teamMode":"personal","teamLeader":"张三"}`},
		{"组队方式非法", `{"projectName":"p","teamMode":"solo","teamLeader":"张三"}`},
		{"姓名栏为空", `{"projectName":"p","teamMode":"personal","teamLeader":"  "}`},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			req := httptest.NewRequest("POST", "/api/proposals", bytes.NewBufferString(c.body))
			w := httptest.NewRecorder()
			mux.ServeHTTP(w, req)
			if w.Code != http.StatusBadRequest {
				t.Errorf("status = %d, want 400", w.Code)
			}
		})
	}
}

func TestProposalListDeleteHistory(t *testing.T) {
	mux, _, _ := setupApplicationMux()
	for _, body := range []string{
		`{"projectName":"A项目","teamMode":"personal","teamLeader":"张三"}`,
		`{"projectName":"B项目","teamMode":"partner","teamLeader":"李四","teamMember":"王五"}`,
	} {
		req := httptest.NewRequest("POST", "/api/proposals", bytes.NewBufferString(body))
		w := httptest.NewRecorder()
		mux.ServeHTTP(w, req)
		if w.Code != http.StatusCreated {
			t.Fatalf("create status = %d", w.Code)
		}
	}

	// 列表：2 条，最近在前，不含删除
	req := httptest.NewRequest("GET", "/api/proposals", nil)
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)
	var body struct {
		Proposals []struct {
			ID   string `json:"id"`
			Name string `json:"projectName"`
		} `json:"proposals"`
	}
	json.NewDecoder(w.Body).Decode(&body)
	if len(body.Proposals) != 2 {
		t.Fatalf("proposals = %d, want 2", len(body.Proposals))
	}
	if body.Proposals[0].Name != "B项目" {
		t.Errorf("first = %s, want B项目（最近在前）", body.Proposals[0].Name)
	}

	// 软删除 A 项目（后提交在前 → A 在第二位）
	delID := body.Proposals[1].ID
	reqDel := httptest.NewRequest("DELETE", "/api/proposals/"+delID, nil)
	wDel := httptest.NewRecorder()
	mux.ServeHTTP(wDel, reqDel)
	if wDel.Code != http.StatusOK {
		t.Fatalf("delete status = %d", wDel.Code)
	}

	// 列表只剩 1 条
	req2 := httptest.NewRequest("GET", "/api/proposals", nil)
	w2 := httptest.NewRecorder()
	mux.ServeHTTP(w2, req2)
	json.NewDecoder(w2.Body).Decode(&body)
	if len(body.Proposals) != 1 {
		t.Errorf("after delete proposals = %d, want 1", len(body.Proposals))
	}

	// 历史：1 条（含删除记录）
	req3 := httptest.NewRequest("GET", "/api/proposals/history", nil)
	w3 := httptest.NewRecorder()
	mux.ServeHTTP(w3, req3)
	var hist struct {
		History []struct {
			ID        string `json:"id"`
			DeletedAt string `json:"deletedAt"`
		} `json:"history"`
	}
	json.NewDecoder(w3.Body).Decode(&hist)
	if len(hist.History) != 1 || hist.History[0].ID != delID || hist.History[0].DeletedAt == "" {
		t.Errorf("history = %+v", hist.History)
	}

	// 重复删除 → 404
	req4 := httptest.NewRequest("DELETE", "/api/proposals/"+delID, nil)
	w4 := httptest.NewRecorder()
	mux.ServeHTTP(w4, req4)
	if w4.Code != http.StatusNotFound {
		t.Errorf("re-delete status = %d, want 404", w4.Code)
	}
}

func TestReportProgress(t *testing.T) {
	mux, _, learnerStore := setupApplicationMux()

	// m3 上报
	body := `{"moduleId":"m3","name":"演示学员"}`
	req := httptest.NewRequest("POST", "/api/courses/prod/progress", bytes.NewBufferString(body))
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("status = %d, body=%s", w.Code, w.Body.String())
	}
	var resp struct {
		Max  int    `json:"max"`
		Last string `json:"last"`
	}
	json.NewDecoder(w.Body).Decode(&resp)
	if resp.Max != 3 || resp.Last != "m3" {
		t.Errorf("resp = %+v", resp)
	}
	learner := learnerStore.GetByName("演示学员")
	if learner == nil || learner.ProgressMax != 3 || learner.Status != "在读" {
		t.Errorf("learner = %+v", learner)
	}

	// 回退上报（m1）不降低 max（只增不减）
	body2 := `{"moduleId":"m1","name":"演示学员"}`
	req2 := httptest.NewRequest("POST", "/api/courses/prod/progress", bytes.NewBufferString(body2))
	w2 := httptest.NewRecorder()
	mux.ServeHTTP(w2, req2)
	if w2.Code != http.StatusOK {
		t.Fatalf("status = %d", w2.Code)
	}
	if got := learnerStore.GetByName("演示学员").ProgressMax; got != 3 {
		t.Errorf("max = %d, want 3（只增不减）", got)
	}

	// 非法 moduleId → 400
	bad := `{"moduleId":"m9","name":"演示学员"}`
	req3 := httptest.NewRequest("POST", "/api/courses/prod/progress", bytes.NewBufferString(bad))
	w3 := httptest.NewRecorder()
	mux.ServeHTTP(w3, req3)
	if w3.Code != http.StatusBadRequest {
		t.Errorf("bad moduleId status = %d, want 400", w3.Code)
	}

	// 后台学员列表
	req4 := httptest.NewRequest("GET", "/api/learners", nil)
	w4 := httptest.NewRecorder()
	mux.ServeHTTP(w4, req4)
	if w4.Code != http.StatusOK {
		t.Errorf("learners status = %d", w4.Code)
	}
	var lb struct {
		Learners []struct {
			Name string `json:"name"`
		} `json:"learners"`
	}
	json.NewDecoder(w4.Body).Decode(&lb)
	if len(lb.Learners) != 1 || lb.Learners[0].Name != "演示学员" {
		t.Errorf("learners = %+v", lb.Learners)
	}
}
