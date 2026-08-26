package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/quanttide/qtcloud-learn-provider/internal/store"
)

// setupMux 创建注册了全部资源路由的 mux，用于 handler 测试。
func setupMux() *http.ServeMux {
	learnerStore := store.NewLearnerStore()
	criterionStore := store.NewCriterionStore()
	completionStore := store.NewCompletionStore()

	learnerh := NewLearnerHandler(learnerStore)
	crith := NewCriterionHandler(criterionStore)
	comh := NewCompletionHandler(completionStore)

	mux := http.NewServeMux()
	mux.HandleFunc("GET /learners", learnerh.List)
	mux.HandleFunc("POST /learners", learnerh.Create)
	mux.HandleFunc("GET /learners/{id}", learnerh.Get)
	mux.HandleFunc("PUT /learners/{id}", learnerh.Update)
	mux.HandleFunc("DELETE /learners/{id}", learnerh.Delete)
	mux.HandleFunc("GET /criteria", crith.List)
	mux.HandleFunc("POST /criteria", crith.Create)
	mux.HandleFunc("GET /criteria/{id}", crith.Get)
	mux.HandleFunc("PUT /criteria/{id}", crith.Update)
	mux.HandleFunc("DELETE /criteria/{id}", crith.Delete)
	mux.HandleFunc("GET /completions", comh.List)
	mux.HandleFunc("POST /completions", comh.Create)
	mux.HandleFunc("GET /completions/{id}", comh.Get)
	mux.HandleFunc("PUT /completions/{id}", comh.Update)
	mux.HandleFunc("DELETE /completions/{id}", comh.Delete)
	return mux
}

func request(t *testing.T, mux *http.ServeMux, method, path, body string) *httptest.ResponseRecorder {
	t.Helper()
	r := httptest.NewRequest(method, path, strings.NewReader(body))
	if body != "" {
		r.Header.Set("Content-Type", "application/json")
	}
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, r)
	return w
}

func assertStatus(t *testing.T, w *httptest.ResponseRecorder, want int) {
	t.Helper()
	if w.Code != want {
		t.Errorf("status = %d, want %d; body = %s", w.Code, want, w.Body.String())
	}
}

func assertJSON(t *testing.T, w *httptest.ResponseRecorder) map[string]any {
	t.Helper()
	var data map[string]any
	if err := json.Unmarshal(w.Body.Bytes(), &data); err != nil {
		t.Fatalf("invalid JSON: %v; body=%s", err, w.Body.String())
	}
	return data
}

func assertJSONArray(t *testing.T, w *httptest.ResponseRecorder) []any {
	t.Helper()
	var data []any
	if err := json.Unmarshal(w.Body.Bytes(), &data); err != nil {
		t.Fatalf("invalid JSON array: %v; body=%s", err, w.Body.String())
	}
	return data
}

// --- Learner ---

func TestLearnerHandler_CRUD(t *testing.T) {
	mux := setupMux()

	w := request(t, mux, "GET", "/learners", "")
	assertStatus(t, w, 200)
	assertJSONArray(t, w)

	w = request(t, mux, "POST", "/learners", `{"user_id":"user-123"}`)
	assertStatus(t, w, 201)
	l := assertJSON(t, w)
	lid := l["id"].(string)
	if l["user_id"] != "user-123" {
		t.Fatalf("Create = %v", l)
	}

	w = request(t, mux, "POST", "/learners", `{invalid`)
	assertStatus(t, w, 400)

	w = request(t, mux, "GET", fmt.Sprintf("/learners/%s", lid), "")
	assertStatus(t, w, 200)

	w = request(t, mux, "GET", "/learners/nonexistent", "")
	assertStatus(t, w, 404)

	w = request(t, mux, "PUT", fmt.Sprintf("/learners/%s", lid), `{"user_id":"user-456"}`)
	assertStatus(t, w, 200)
	l = assertJSON(t, w)
	if l["user_id"] != "user-456" {
		t.Fatalf("Update = %v", l)
	}

	w = request(t, mux, "PUT", "/learners/nonexistent", `{"user_id":"x"}`)
	assertStatus(t, w, 404)

	w = request(t, mux, "DELETE", fmt.Sprintf("/learners/%s", lid), "")
	assertStatus(t, w, 204)
	w = request(t, mux, "DELETE", "/learners/nonexistent", "")
	assertStatus(t, w, 404)
}

// --- Criterion ---

func TestCriterionHandler_CRUD(t *testing.T) {
	mux := setupMux()

	w := request(t, mux, "GET", "/criteria", "")
	assertStatus(t, w, 200)
	assertJSONArray(t, w)

	w = request(t, mux, "POST", "/criteria", `{"title":"会连接 Zed","description":"成功建立 Zed 连接"}`)
	assertStatus(t, w, 201)
	c := assertJSON(t, w)
	cid := c["id"].(string)
	if c["title"] != "会连接 Zed" || c["description"] != "成功建立 Zed 连接" {
		t.Fatalf("Create = %v", c)
	}

	w = request(t, mux, "POST", "/criteria", `{invalid`)
	assertStatus(t, w, 400)

	w = request(t, mux, "POST", "/criteria", `{"description":"缺 title"}`)
	assertStatus(t, w, 400)

	w = request(t, mux, "POST", "/criteria", `{"title":"缺 description"}`)
	assertStatus(t, w, 400)

	w = request(t, mux, "GET", fmt.Sprintf("/criteria/%s", cid), "")
	assertStatus(t, w, 200)

	w = request(t, mux, "GET", "/criteria/nonexistent", "")
	assertStatus(t, w, 404)

	w = request(t, mux, "PUT", fmt.Sprintf("/criteria/%s", cid), `{"title":"会用 Agent 执行任务","description":"完成 Agent 任务"}`)
	assertStatus(t, w, 200)
	c = assertJSON(t, w)
	if c["title"] != "会用 Agent 执行任务" || c["description"] != "完成 Agent 任务" {
		t.Fatalf("Update = %v", c)
	}

	w = request(t, mux, "PUT", "/criteria/nonexistent", `{"title":"x","description":"x"}`)
	assertStatus(t, w, 404)

	w = request(t, mux, "DELETE", fmt.Sprintf("/criteria/%s", cid), "")
	assertStatus(t, w, 204)
	w = request(t, mux, "DELETE", "/criteria/nonexistent", "")
	assertStatus(t, w, 404)
}

// --- Completion ---

func TestCompletionHandler_CRUD(t *testing.T) {
	mux := setupMux()

	w := request(t, mux, "GET", "/completions", "")
	assertStatus(t, w, 200)
	assertJSONArray(t, w)

	w = request(t, mux, "POST", "/completions", `{"learner_id":"lea-1","criterion_id":"cri-1"}`)
	assertStatus(t, w, 201)
	com := assertJSON(t, w)
	cid := com["id"].(string)
	if com["learner_id"] != "lea-1" || com["criterion_id"] != "cri-1" {
		t.Fatalf("Create = %v", com)
	}
	// 缺省 status → not_completed
	if com["status"] != "not_completed" {
		t.Fatalf("Create default status = %v", com["status"])
	}
	if com["created_at"] == "" || com["updated_at"] == "" {
		t.Fatalf("Create timestamps missing = %v", com)
	}

	w = request(t, mux, "POST", "/completions", `{invalid`)
	assertStatus(t, w, 400)

	w = request(t, mux, "POST", "/completions", `{"criterion_id":"cri-1"}`)
	assertStatus(t, w, 400)

	w = request(t, mux, "POST", "/completions", `{"learner_id":"lea-1"}`)
	assertStatus(t, w, 400)

	w = request(t, mux, "POST", "/completions", `{"learner_id":"lea-1","criterion_id":"cri-1","status":"invalid"}`)
	assertStatus(t, w, 400)

	w = request(t, mux, "GET", fmt.Sprintf("/completions/%s", cid), "")
	assertStatus(t, w, 200)

	w = request(t, mux, "GET", "/completions/nonexistent", "")
	assertStatus(t, w, 404)

	w = request(t, mux, "PUT", fmt.Sprintf("/completions/%s", cid), `{"learner_id":"lea-1","criterion_id":"cri-1","status":"completed"}`)
	assertStatus(t, w, 200)
	com = assertJSON(t, w)
	if com["status"] != "completed" {
		t.Fatalf("Update = %v", com)
	}

	// 局部更新（合并语义）：仅改 status 不抹掉 learner_id / criterion_id
	w = request(t, mux, "PUT", fmt.Sprintf("/completions/%s", cid), `{"status":"not_completed"}`)
	assertStatus(t, w, 200)
	com = assertJSON(t, w)
	if com["status"] != "not_completed" || com["learner_id"] != "lea-1" || com["criterion_id"] != "cri-1" {
		t.Fatalf("Partial Update = %v", com)
	}

	w = request(t, mux, "PUT", "/completions/nonexistent", `{"status":"completed"}`)
	assertStatus(t, w, 404)

	w = request(t, mux, "DELETE", fmt.Sprintf("/completions/%s", cid), "")
	assertStatus(t, w, 204)
	w = request(t, mux, "DELETE", "/completions/nonexistent", "")
	assertStatus(t, w, 404)
}
