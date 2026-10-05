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
	completionStore := store.NewCompletionStore()
	scheduleStore := store.NewScheduleStore()
	taskStore := store.NewTaskStore()

	learnerh := NewLearnerHandler(learnerStore)
	comh := NewCompletionHandler(completionStore)
	scheduleh := NewScheduleHandler(scheduleStore)
	taskh := NewTaskHandler(taskStore)

	mux := http.NewServeMux()
	mux.HandleFunc("GET /learners", learnerh.List)
	mux.HandleFunc("POST /learners", learnerh.Create)
	mux.HandleFunc("GET /learners/{id}", learnerh.Get)
	mux.HandleFunc("PUT /learners/{id}", learnerh.Update)
	mux.HandleFunc("DELETE /learners/{id}", learnerh.Delete)
	mux.HandleFunc("GET /completions", comh.List)
	mux.HandleFunc("POST /completions", comh.Create)
	mux.HandleFunc("GET /completions/{id}", comh.Get)
	mux.HandleFunc("PUT /completions/{id}", comh.Update)
	mux.HandleFunc("DELETE /completions/{id}", comh.Delete)
	mux.HandleFunc("GET /schedules", scheduleh.List)
	mux.HandleFunc("POST /schedules", scheduleh.Create)
	mux.HandleFunc("GET /schedules/{id}", scheduleh.Get)
	mux.HandleFunc("PUT /schedules/{id}", scheduleh.Update)
	mux.HandleFunc("DELETE /schedules/{id}", scheduleh.Delete)
	mux.HandleFunc("GET /tasks", taskh.List)
	mux.HandleFunc("POST /tasks", taskh.Create)
	mux.HandleFunc("GET /tasks/{id}", taskh.Get)
	mux.HandleFunc("PUT /tasks/{id}", taskh.Update)
	mux.HandleFunc("DELETE /tasks/{id}", taskh.Delete)
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

	w = request(t, mux, "POST", "/learners", `{"user_id":"user-123","schedule_id":"schedule-agent-engineer"}`)
	assertStatus(t, w, 201)
	l := assertJSON(t, w)
	lid := l["id"].(string)
	if l["user_id"] != "user-123" || l["schedule_id"] != "schedule-agent-engineer" {
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

// --- Completion ---

func TestCompletionHandler_CRUD(t *testing.T) {
	mux := setupMux()

	w := request(t, mux, "GET", "/completions", "")
	assertStatus(t, w, 200)
	assertJSONArray(t, w)

	w = request(t, mux, "POST", "/completions", `{"learner_id":"lea-1","task_id":"task-data-second-brain"}`)
	assertStatus(t, w, 201)
	com := assertJSON(t, w)
	cid := com["id"].(string)
	if com["learner_id"] != "lea-1" || com["task_id"] != "task-data-second-brain" {
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

	w = request(t, mux, "POST", "/completions", `{"task_id":"task-data-second-brain"}`)
	assertStatus(t, w, 400)

	w = request(t, mux, "POST", "/completions", `{"learner_id":"lea-1"}`)
	assertStatus(t, w, 400)

	w = request(t, mux, "POST", "/completions", `{"learner_id":"lea-1","task_id":"task-data-second-brain","status":"invalid"}`)
	assertStatus(t, w, 400)

	w = request(t, mux, "GET", fmt.Sprintf("/completions/%s", cid), "")
	assertStatus(t, w, 200)

	w = request(t, mux, "GET", "/completions/nonexistent", "")
	assertStatus(t, w, 404)

	w = request(t, mux, "PUT", fmt.Sprintf("/completions/%s", cid), `{"learner_id":"lea-1","task_id":"task-data-second-brain","status":"completed"}`)
	assertStatus(t, w, 200)
	com = assertJSON(t, w)
	if com["status"] != "completed" {
		t.Fatalf("Update = %v", com)
	}

	// 局部更新（合并语义）：仅改 status 不抹掉 learner_id / task_id
	w = request(t, mux, "PUT", fmt.Sprintf("/completions/%s", cid), `{"status":"not_completed"}`)
	assertStatus(t, w, 200)
	com = assertJSON(t, w)
	if com["status"] != "not_completed" || com["learner_id"] != "lea-1" || com["task_id"] != "task-data-second-brain" {
		t.Fatalf("Partial Update = %v", com)
	}

	w = request(t, mux, "PUT", "/completions/nonexistent", `{"status":"completed"}`)
	assertStatus(t, w, 404)

	w = request(t, mux, "DELETE", fmt.Sprintf("/completions/%s", cid), "")
	assertStatus(t, w, 204)
	w = request(t, mux, "DELETE", "/completions/nonexistent", "")
	assertStatus(t, w, 404)
}

// --- Task ---

func TestTaskHandler_CRUD(t *testing.T) {
	mux := setupMux()

	w := request(t, mux, "GET", "/tasks", "")
	assertStatus(t, w, 200)
	assertJSONArray(t, w)

	w = request(t, mux, "POST", "/tasks", `{"id":"task-data-second-brain","title":"熟悉数据工程第二大脑","description":"完成一条改进建议。"}`)
	assertStatus(t, w, 201)
	task := assertJSON(t, w)
	if task["id"] != "task-data-second-brain" || task["title"] == "" || task["description"] == "" {
		t.Fatalf("Create = %v", task)
	}

	w = request(t, mux, "POST", "/tasks", `{"title":"缺描述"}`)
	assertStatus(t, w, 400)

	w = request(t, mux, "GET", "/tasks/task-data-second-brain", "")
	assertStatus(t, w, 200)

	w = request(t, mux, "PUT", "/tasks/task-data-second-brain", `{"title":"熟悉第二大脑"}`)
	assertStatus(t, w, 200)
	task = assertJSON(t, w)
	if task["title"] != "熟悉第二大脑" || task["description"] == "" {
		t.Fatalf("Partial Update = %v", task)
	}

	w = request(t, mux, "DELETE", "/tasks/task-data-second-brain", "")
	assertStatus(t, w, 204)
}

// --- Schedule ---

func TestScheduleHandler_CRUD(t *testing.T) {
	mux := setupMux()
	body := `{"id":"schedule-agent-engineer","title":"智能体工程师训练营","description":"按学习路径推进的训练计划。","tasks":[{"id":"task-data-second-brain","title":"熟悉数据工程第二大脑","description":"完成一条改进建议。"}]}`

	w := request(t, mux, "GET", "/schedules", "")
	assertStatus(t, w, 200)
	assertJSONArray(t, w)

	w = request(t, mux, "POST", "/schedules", body)
	assertStatus(t, w, 201)
	schedule := assertJSON(t, w)
	if schedule["id"] != "schedule-agent-engineer" || schedule["title"] == "" {
		t.Fatalf("Create = %v", schedule)
	}

	w = request(t, mux, "POST", "/schedules", `{"title":"缺任务"}`)
	assertStatus(t, w, 400)

	w = request(t, mux, "GET", "/schedules/schedule-agent-engineer", "")
	assertStatus(t, w, 200)

	w = request(t, mux, "PUT", "/schedules/schedule-agent-engineer", `{"title":"智能体工程师训练营 v2"}`)
	assertStatus(t, w, 200)
	schedule = assertJSON(t, w)
	tasks := schedule["tasks"].([]any)
	if schedule["title"] != "智能体工程师训练营 v2" || len(tasks) != 1 {
		t.Fatalf("Partial Update = %v", schedule)
	}

	w = request(t, mux, "DELETE", "/schedules/schedule-agent-engineer", "")
	assertStatus(t, w, 204)
}
