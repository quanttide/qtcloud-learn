package handler

// 立项申请 API 测试：提交（含姓名栏校验）/ 列表 / 详情。

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/quanttide/qtcloud-learn-provider/internal/store"
)

func setupApplicationMux() (*http.ServeMux, *store.ApplicationStore) {
	appStore := store.NewApplicationStore()
	apph := NewApplicationHandler(appStore)
	mux := http.NewServeMux()
	mux.HandleFunc("POST /api/v1/applications", apph.Create)
	mux.HandleFunc("GET /api/v1/applications", apph.List)
	mux.HandleFunc("GET /api/v1/applications/{id}", apph.Get)
	return mux, appStore
}

func TestApplicationCreate(t *testing.T) {
	mux, _ := setupApplicationMux()

	// 个人独立：memberNames=[个人姓名]
	body := `{"projectName":"校园选课助手","blindSpot":"选课信息分散","demoPlan":"做一个聚合查询页","direction":"内容","teamMode":"personal","memberNames":["张三"],"studentId":"stu-1","studentName":"张三"}`
	req := httptest.NewRequest("POST", "/api/v1/applications", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Fatalf("status = %d, body=%s", w.Code, w.Body.String())
	}
	var app struct {
		ID          string   `json:"id"`
		Status      string   `json:"status"`
		TeamMode    string   `json:"teamMode"`
		MemberNames []string `json:"memberNames"`
		CreatedAt   string   `json:"createdAt"`
	}
	json.NewDecoder(w.Body).Decode(&app)
	if app.ID == "" || app.Status != "submitted" || app.TeamMode != "personal" {
		t.Errorf("app = %+v", app)
	}
	if len(app.MemberNames) != 1 || app.MemberNames[0] != "张三" {
		t.Errorf("memberNames = %v", app.MemberNames)
	}
	if app.CreatedAt == "" {
		t.Error("createdAt empty")
	}
}

func TestApplicationCreatePartner(t *testing.T) {
	mux, _ := setupApplicationMux()

	// 搭档：memberNames=[队长, 队员]
	body := `{"projectName":"盲区地图","blindSpot":"业务盲区无记录","demoPlan":"共享表单","teamMode":"partner","memberNames":["李四","王五"],"studentId":"stu-2","studentName":"李四"}`
	req := httptest.NewRequest("POST", "/api/v1/applications", bytes.NewBufferString(body))
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)
	if w.Code != http.StatusCreated {
		t.Fatalf("status = %d", w.Code)
	}
	var app struct {
		MemberNames []string `json:"memberNames"`
	}
	json.NewDecoder(w.Body).Decode(&app)
	if len(app.MemberNames) != 2 || app.MemberNames[0] != "李四" || app.MemberNames[1] != "王五" {
		t.Errorf("memberNames = %v", app.MemberNames)
	}
}

func TestApplicationValidation(t *testing.T) {
	mux, _ := setupApplicationMux()
	cases := []struct {
		name string
		body string
	}{
		{"缺项目名称", `{"blindSpot":"x","demoPlan":"y","teamMode":"personal","memberNames":["张三"]}`},
		{"缺盲区描述", `{"projectName":"p","demoPlan":"y","teamMode":"personal","memberNames":["张三"]}`},
		{"组队方式非法", `{"projectName":"p","blindSpot":"x","demoPlan":"y","teamMode":"solo","memberNames":["张三"]}`},
		{"姓名栏为空", `{"projectName":"p","blindSpot":"x","demoPlan":"y","teamMode":"personal","memberNames":[]}`},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			req := httptest.NewRequest("POST", "/api/v1/applications", bytes.NewBufferString(c.body))
			w := httptest.NewRecorder()
			mux.ServeHTTP(w, req)
			if w.Code != http.StatusBadRequest {
				t.Errorf("status = %d, want 400", w.Code)
			}
		})
	}
}

func TestApplicationListAndGet(t *testing.T) {
	mux, _ := setupApplicationMux()
	// 提交两条
	for _, body := range []string{
		`{"projectName":"A项目","blindSpot":"x","demoPlan":"y","teamMode":"personal","memberNames":["张三"]}`,
		`{"projectName":"B项目","blindSpot":"x","demoPlan":"y","teamMode":"partner","memberNames":["李四","王五"]}`,
	} {
		req := httptest.NewRequest("POST", "/api/v1/applications", bytes.NewBufferString(body))
		w := httptest.NewRecorder()
		mux.ServeHTTP(w, req)
		if w.Code != http.StatusCreated {
			t.Fatalf("create status = %d", w.Code)
		}
	}

	// 列表：2 条，后提交在前
	req := httptest.NewRequest("GET", "/api/v1/applications", nil)
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("list status = %d", w.Code)
	}
	var body struct {
		Applications []struct {
			ID   string `json:"id"`
			Name string `json:"projectName"`
		} `json:"applications"`
	}
	json.NewDecoder(w.Body).Decode(&body)
	if len(body.Applications) != 2 {
		t.Fatalf("applications = %d, want 2", len(body.Applications))
	}
	if body.Applications[0].Name != "B项目" {
		t.Errorf("first = %s, want B项目（最近在前）", body.Applications[0].Name)
	}

	// 详情
	req2 := httptest.NewRequest("GET", "/api/v1/applications/"+body.Applications[0].ID, nil)
	w2 := httptest.NewRecorder()
	mux.ServeHTTP(w2, req2)
	if w2.Code != http.StatusOK {
		t.Errorf("get status = %d", w2.Code)
	}
	// 不存在 → 404
	req3 := httptest.NewRequest("GET", "/api/v1/applications/nope", nil)
	w3 := httptest.NewRecorder()
	mux.ServeHTTP(w3, req3)
	if w3.Code != http.StatusNotFound {
		t.Errorf("missing status = %d, want 404", w3.Code)
	}
}
