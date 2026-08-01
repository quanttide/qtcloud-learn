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

// setupMux 创建注册了全部 /api/v1 路由的 mux，用于 handler 测试。
func setupMux() *http.ServeMux {
	classStore := store.NewClassStore()
	studentStore := store.NewStudentStore()
	teacherStore := store.NewTeacherStore()
	sessionStore := store.NewSessionStore()
	assessmentStore := store.NewAssessmentStore()
	submissionStore := store.NewSubmissionStore()
	enrollmentStore := store.NewEnrollmentStore()
	progressStore := store.NewProgressStore()

	ch := NewClassHandler(classStore)
	sh := NewStudentHandler(studentStore)
	th := NewTeacherHandler(teacherStore)
	sessh := NewSessionHandler(sessionStore)
	ah := NewAssessmentHandler(assessmentStore)
	subh := NewSubmissionHandler(submissionStore)
	eh := NewEnrollmentHandler(enrollmentStore)
	ph := NewProgressHandler(progressStore)

	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/v1/classes", ch.List)
	mux.HandleFunc("POST /api/v1/classes", ch.Create)
	mux.HandleFunc("GET /api/v1/classes/{id}", ch.Get)
	mux.HandleFunc("PUT /api/v1/classes/{id}", ch.Update)
	mux.HandleFunc("DELETE /api/v1/classes/{id}", ch.Delete)
	mux.HandleFunc("GET /api/v1/students", sh.List)
	mux.HandleFunc("POST /api/v1/students", sh.Create)
	mux.HandleFunc("GET /api/v1/students/{id}", sh.Get)
	mux.HandleFunc("PUT /api/v1/students/{id}", sh.Update)
	mux.HandleFunc("DELETE /api/v1/students/{id}", sh.Delete)
	mux.HandleFunc("GET /api/v1/teachers", th.List)
	mux.HandleFunc("POST /api/v1/teachers", th.Create)
	mux.HandleFunc("GET /api/v1/teachers/{id}", th.Get)
	mux.HandleFunc("PUT /api/v1/teachers/{id}", th.Update)
	mux.HandleFunc("DELETE /api/v1/teachers/{id}", th.Delete)
	mux.HandleFunc("GET /api/v1/sessions", sessh.List)
	mux.HandleFunc("POST /api/v1/sessions", sessh.Create)
	mux.HandleFunc("GET /api/v1/sessions/{id}", sessh.Get)
	mux.HandleFunc("PUT /api/v1/sessions/{id}", sessh.Update)
	mux.HandleFunc("DELETE /api/v1/sessions/{id}", sessh.Delete)
	mux.HandleFunc("GET /api/v1/assessments", ah.List)
	mux.HandleFunc("POST /api/v1/assessments", ah.Create)
	mux.HandleFunc("GET /api/v1/assessments/{id}", ah.Get)
	mux.HandleFunc("PUT /api/v1/assessments/{id}", ah.Update)
	mux.HandleFunc("DELETE /api/v1/assessments/{id}", ah.Delete)
	mux.HandleFunc("GET /api/v1/submissions", subh.List)
	mux.HandleFunc("POST /api/v1/submissions", subh.Create)
	mux.HandleFunc("GET /api/v1/submissions/{id}", subh.Get)
	mux.HandleFunc("PUT /api/v1/submissions/{id}", subh.Update)
	mux.HandleFunc("DELETE /api/v1/submissions/{id}", subh.Delete)
	mux.HandleFunc("GET /api/v1/enrollments", eh.List)
	mux.HandleFunc("POST /api/v1/enrollments", eh.Create)
	mux.HandleFunc("GET /api/v1/enrollments/{id}", eh.Get)
	mux.HandleFunc("PUT /api/v1/enrollments/{id}", eh.Update)
	mux.HandleFunc("DELETE /api/v1/enrollments/{id}", eh.Delete)
	mux.HandleFunc("GET /api/v1/progress", ph.List)
	mux.HandleFunc("POST /api/v1/progress", ph.Create)
	mux.HandleFunc("GET /api/v1/progress/{id}", ph.Get)
	mux.HandleFunc("PUT /api/v1/progress/{id}", ph.Update)
	mux.HandleFunc("DELETE /api/v1/progress/{id}", ph.Delete)
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

// --- Class（移植自 qtcloud-course/provider 的 TestClassHandler_CRUD） ---

func TestClassHandler_CRUD(t *testing.T) {
	mux := setupMux()

	w := request(t, mux, "GET", "/api/v1/classes", "")
	assertStatus(t, w, 200)
	assertJSONArray(t, w)

	w = request(t, mux, "POST", "/api/v1/classes", `{"name":"浙理班级","refName":"大数据微专业","refType":"program","refId":"prog-1","startDate":"2026-09-01","endDate":"2027-01-15","studentCount":30}`)
	assertStatus(t, w, 201)
	c := assertJSON(t, w)
	cid := c["id"].(string)
	if c["slug"] == "" {
		t.Fatal("slug is empty")
	}

	w = request(t, mux, "POST", "/api/v1/classes", `{invalid`)
	assertStatus(t, w, 400)

	w = request(t, mux, "POST", "/api/v1/classes", `{"name":"x"}`)
	assertStatus(t, w, 400)

	w = request(t, mux, "POST", "/api/v1/classes", `{"refId":"prog-1"}`)
	assertStatus(t, w, 400)

	// Create duplicate name
	w = request(t, mux, "POST", "/api/v1/classes", `{"name":"浙理班级","refId":"prog-2"}`)
	assertStatus(t, w, 409)

	w = request(t, mux, "GET", fmt.Sprintf("/api/v1/classes/%s", cid), "")
	assertStatus(t, w, 200)

	w = request(t, mux, "GET", "/api/v1/classes/nonexistent", "")
	assertStatus(t, w, 404)

	w = request(t, mux, "PUT", fmt.Sprintf("/api/v1/classes/%s", cid), `{"name":"v2","refName":"v2","refType":"course","refId":"cour-1","status":"active","startDate":"2026-09-15","endDate":"2027-02-01","studentCount":35,"progress":0.5,"teacherIds":["tea-1"]}`)
	assertStatus(t, w, 200)
	c = assertJSON(t, w)
	if c["name"] != "v2" || c["studentCount"] != float64(35) || c["progress"] != 0.5 {
		t.Fatalf("Update = %v", c)
	}
	if ids, ok := c["teacherIds"].([]any); !ok || len(ids) != 1 {
		t.Fatalf("Update teacherIds = %v", c["teacherIds"])
	}

	w = request(t, mux, "PUT", fmt.Sprintf("/api/v1/classes/%s", cid), `{`)
	assertStatus(t, w, 400)

	w = request(t, mux, "PUT", "/api/v1/classes/nonexistent", `{"name":"x","refId":"prog-1"}`)
	assertStatus(t, w, 404)

	w = request(t, mux, "DELETE", fmt.Sprintf("/api/v1/classes/%s", cid), "")
	assertStatus(t, w, 204)
	w = request(t, mux, "DELETE", fmt.Sprintf("/api/v1/classes/%s", cid), "")
	assertStatus(t, w, 404)
	w = request(t, mux, "DELETE", "/api/v1/classes/nonexistent", "")
	assertStatus(t, w, 404)
}

// --- Student ---

func TestStudentHandler_CRUD(t *testing.T) {
	mux := setupMux()

	w := request(t, mux, "GET", "/api/v1/students", "")
	assertStatus(t, w, 200)
	assertJSONArray(t, w)

	w = request(t, mux, "POST", "/api/v1/students", `{"name":"张三","email":"z@example.com","plan":"vip"}`)
	assertStatus(t, w, 201)
	st := assertJSON(t, w)
	sid := st["id"].(string)
	if st["name"] != "张三" || st["plan"] != "vip" {
		t.Fatalf("Create = %v", st)
	}

	w = request(t, mux, "POST", "/api/v1/students", `{invalid`)
	assertStatus(t, w, 400)

	w = request(t, mux, "POST", "/api/v1/students", `{"name":""}`)
	assertStatus(t, w, 400)

	w = request(t, mux, "GET", fmt.Sprintf("/api/v1/students/%s", sid), "")
	assertStatus(t, w, 200)

	w = request(t, mux, "GET", "/api/v1/students/nonexistent", "")
	assertStatus(t, w, 404)

	w = request(t, mux, "PUT", fmt.Sprintf("/api/v1/students/%s", sid), `{"name":"张三丰","email":"z2@example.com","avatar":"a.png","plan":"paid"}`)
	assertStatus(t, w, 200)
	st = assertJSON(t, w)
	if st["name"] != "张三丰" || st["plan"] != "paid" {
		t.Fatalf("Update = %v", st)
	}

	w = request(t, mux, "PUT", "/api/v1/students/nonexistent", `{"name":"x"}`)
	assertStatus(t, w, 404)

	w = request(t, mux, "DELETE", fmt.Sprintf("/api/v1/students/%s", sid), "")
	assertStatus(t, w, 204)
	w = request(t, mux, "DELETE", "/api/v1/students/nonexistent", "")
	assertStatus(t, w, 404)
}

// --- Assessment ---

func TestAssessmentHandler_CRUD(t *testing.T) {
	mux := setupMux()

	w := request(t, mux, "POST", "/api/v1/assessments", `{"classId":"class-1","type":"homework","title":"作业1","fullScore":100,"passScore":60,"deadline":"2026-10-01"}`)
	assertStatus(t, w, 201)
	a := assertJSON(t, w)
	aid := a["id"].(string)
	if a["title"] != "作业1" || a["fullScore"] != float64(100) {
		t.Fatalf("Create = %v", a)
	}

	w = request(t, mux, "POST", "/api/v1/assessments", `{"title":"缺classId"}`)
	assertStatus(t, w, 400)

	w = request(t, mux, "POST", "/api/v1/assessments", `{"classId":"class-1"}`)
	assertStatus(t, w, 400)

	w = request(t, mux, "PUT", fmt.Sprintf("/api/v1/assessments/%s", aid), `{"classId":"class-1","type":"exam","title":"期中考试","fullScore":150,"passScore":90,"deadline":"2026-11-01"}`)
	assertStatus(t, w, 200)
	a = assertJSON(t, w)
	if a["type"] != "exam" || a["fullScore"] != float64(150) {
		t.Fatalf("Update = %v", a)
	}

	w = request(t, mux, "DELETE", fmt.Sprintf("/api/v1/assessments/%s", aid), "")
	assertStatus(t, w, 204)
}

// --- Submission ---

func TestSubmissionHandler_CRUD(t *testing.T) {
	mux := setupMux()

	w := request(t, mux, "POST", "/api/v1/submissions", `{"assessmentId":"assess-1","studentId":"stu-1","status":"submitted","submittedAt":"2026-09-30T10:00:00Z"}`)
	assertStatus(t, w, 201)
	sub := assertJSON(t, w)
	sid := sub["id"].(string)
	if sub["assessmentId"] != "assess-1" || sub["status"] != "submitted" {
		t.Fatalf("Create = %v", sub)
	}

	w = request(t, mux, "POST", "/api/v1/submissions", `{"studentId":"stu-1"}`)
	assertStatus(t, w, 400)

	w = request(t, mux, "POST", "/api/v1/submissions", `{"assessmentId":"assess-1"}`)
	assertStatus(t, w, 400)

	w = request(t, mux, "PUT", fmt.Sprintf("/api/v1/submissions/%s", sid), `{"assessmentId":"assess-1","studentId":"stu-1","status":"late","score":88.5,"comment":"不错","submittedAt":"2026-09-30T10:00:00Z"}`)
	assertStatus(t, w, 200)
	sub = assertJSON(t, w)
	if sub["status"] != "late" || sub["score"] != 88.5 {
		t.Fatalf("Update = %v", sub)
	}

	w = request(t, mux, "DELETE", fmt.Sprintf("/api/v1/submissions/%s", sid), "")
	assertStatus(t, w, 204)
}

func TestTeacherHandler_CRUD(t *testing.T) {
	mux := setupMux()

	w := request(t, mux, "POST", "/api/v1/teachers", `{"name":"王老师","email":"wang@example.com","title":"教授"}`)
	assertStatus(t, w, 201)
	tr := assertJSON(t, w)
	tid := tr["id"].(string)
	if tr["name"] != "王老师" || tr["title"] != "教授" {
		t.Fatalf("Create = %v", tr)
	}

	w = request(t, mux, "POST", "/api/v1/teachers", `{"name":""}`)
	assertStatus(t, w, 400)

	w = request(t, mux, "PUT", fmt.Sprintf("/api/v1/teachers/%s", tid), `{"name":"王老师v2","email":"w2@example.com","title":"副教授"}`)
	assertStatus(t, w, 200)
	tr = assertJSON(t, w)
	if tr["name"] != "王老师v2" || tr["title"] != "副教授" {
		t.Fatalf("Update = %v", tr)
	}

	w = request(t, mux, "DELETE", fmt.Sprintf("/api/v1/teachers/%s", tid), "")
	assertStatus(t, w, 204)
}

func TestSessionHandler_CRUD(t *testing.T) {
	mux := setupMux()

	w := request(t, mux, "POST", "/api/v1/sessions", `{"classId":"class-1","lessonTitle":"Git 入门","teacherId":"tea-1","startTime":"2026-09-02T09:00:00Z","durationMinutes":45,"location":"A-101","status":"upcoming","attendances":[{"studentId":"stu-1","status":"present"}]}`)
	assertStatus(t, w, 201)
	sess := assertJSON(t, w)
	sid := sess["id"].(string)
	if sess["lessonTitle"] != "Git 入门" || sess["durationMinutes"] != float64(45) {
		t.Fatalf("Create = %v", sess)
	}

	w = request(t, mux, "POST", "/api/v1/sessions", `{"lessonTitle":"缺classId"}`)
	assertStatus(t, w, 400)

	w = request(t, mux, "POST", "/api/v1/sessions", `{"classId":"class-1"}`)
	assertStatus(t, w, 400)

	w = request(t, mux, "PUT", fmt.Sprintf("/api/v1/sessions/%s", sid), `{"classId":"class-1","lessonTitle":"Git 进阶","teacherId":"tea-2","startTime":"2026-09-03T09:00:00Z","durationMinutes":60,"location":"B-202","status":"completed","attendances":[{"studentId":"stu-1","status":"absent"}]}`)
	assertStatus(t, w, 200)
	sess = assertJSON(t, w)
	if sess["lessonTitle"] != "Git 进阶" || sess["status"] != "completed" {
		t.Fatalf("Update = %v", sess)
	}

	w = request(t, mux, "DELETE", fmt.Sprintf("/api/v1/sessions/%s", sid), "")
	assertStatus(t, w, 204)
}

// --- Enrollment（新增） ---

func TestEnrollmentHandler_CRUD(t *testing.T) {
	mux := setupMux()

	w := request(t, mux, "GET", "/api/v1/enrollments", "")
	assertStatus(t, w, 200)
	assertJSONArray(t, w)

	w = request(t, mux, "POST", "/api/v1/enrollments", `{"classId":"class-1","studentId":"stu-1","status":"enrolled"}`)
	assertStatus(t, w, 201)
	e := assertJSON(t, w)
	eid := e["id"].(string)
	if e["classId"] != "class-1" || e["studentId"] != "stu-1" {
		t.Fatalf("Create = %v", e)
	}

	w = request(t, mux, "POST", "/api/v1/enrollments", `{"classId":"class-1"}`)
	assertStatus(t, w, 400)

	w = request(t, mux, "POST", "/api/v1/enrollments", `{"studentId":"stu-1"}`)
	assertStatus(t, w, 400)

	w = request(t, mux, "PUT", fmt.Sprintf("/api/v1/enrollments/%s", eid), `{"classId":"class-1","studentId":"stu-1","status":"withdrawn"}`)
	assertStatus(t, w, 200)
	e = assertJSON(t, w)
	if e["status"] != "withdrawn" {
		t.Fatalf("Update = %v", e)
	}

	w = request(t, mux, "GET", fmt.Sprintf("/api/v1/enrollments/%s", eid), "")
	assertStatus(t, w, 200)

	w = request(t, mux, "DELETE", fmt.Sprintf("/api/v1/enrollments/%s", eid), "")
	assertStatus(t, w, 204)
	w = request(t, mux, "DELETE", "/api/v1/enrollments/nonexistent", "")
	assertStatus(t, w, 404)
}

// --- Progress（新增） ---

func TestProgressHandler_CRUD(t *testing.T) {
	mux := setupMux()

	w := request(t, mux, "GET", "/api/v1/progress", "")
	assertStatus(t, w, 200)
	assertJSONArray(t, w)

	w = request(t, mux, "POST", "/api/v1/progress", `{"studentId":"stu-1","classId":"class-1","percent":0.0,"finished":false}`)
	assertStatus(t, w, 201)
	p := assertJSON(t, w)
	pid := p["id"].(string)
	if p["studentId"] != "stu-1" || p["classId"] != "class-1" || p["percent"] != float64(0) {
		t.Fatalf("Create = %v", p)
	}

	w = request(t, mux, "POST", "/api/v1/progress", `{"classId":"class-1"}`)
	assertStatus(t, w, 400)

	w = request(t, mux, "POST", "/api/v1/progress", `{"studentId":"stu-1"}`)
	assertStatus(t, w, 400)

	w = request(t, mux, "PUT", fmt.Sprintf("/api/v1/progress/%s", pid), `{"studentId":"stu-1","classId":"class-1","percent":1.0,"finished":true,"updatedAt":"2026-08-02"}`)
	assertStatus(t, w, 200)
	p = assertJSON(t, w)
	if p["percent"] != 1.0 || p["finished"] != true {
		t.Fatalf("Update = %v", p)
	}

	w = request(t, mux, "GET", fmt.Sprintf("/api/v1/progress/%s", pid), "")
	assertStatus(t, w, 200)

	w = request(t, mux, "DELETE", fmt.Sprintf("/api/v1/progress/%s", pid), "")
	assertStatus(t, w, 204)
	w = request(t, mux, "DELETE", "/api/v1/progress/nonexistent", "")
	assertStatus(t, w, 404)
}
