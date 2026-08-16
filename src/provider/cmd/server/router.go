package main

import (
	"net/http"
	"os"

	"github.com/quanttide/qtcloud-learn-provider/internal/handler"
	"github.com/quanttide/qtcloud-learn-provider/internal/store"
)

// newRouter 创建并配置所有路由，可单独测试。
// LMS API 统一挂在 /api/v1 前缀下。
// 持久化：DATA_DIR 环境变量非空时，学员/进度/立项三个后台核心实体
// 启用 JSON 文件存储（重启不丢；RDS 后续版本接入）。
func newRouter() *http.ServeMux {
	dataDir := os.Getenv("DATA_DIR")

	classStore := store.NewClassStore()
	studentStore := store.NewStudentStore()
	teacherStore := store.NewTeacherStore()
	sessionStore := store.NewSessionStore()
	assessmentStore := store.NewAssessmentStore()
	submissionStore := store.NewSubmissionStore()
	enrollmentStore := store.NewEnrollmentStore()
	progressStore := store.NewProgressStore()
	applicationStore := store.NewApplicationStore()

	if dataDir != "" {
		if err := os.MkdirAll(dataDir, 0o755); err != nil {
			panic(err)
		}
		applicationStore.BaseStore.SetPersistPath(dataDir + "/applications.json")
		_ = applicationStore.BaseStore.Load(dataDir + "/applications.json")
		studentStore.BaseStore.SetPersistPath(dataDir + "/students.json")
		_ = studentStore.BaseStore.Load(dataDir + "/students.json")
		progressStore.BaseStore.SetPersistPath(dataDir + "/progress.json")
		_ = progressStore.BaseStore.Load(dataDir + "/progress.json")
	}

	ch := handler.NewClassHandler(classStore)
	sh := handler.NewStudentHandler(studentStore)
	th := handler.NewTeacherHandler(teacherStore)
	sessh := handler.NewSessionHandler(sessionStore)
	ah := handler.NewAssessmentHandler(assessmentStore)
	subh := handler.NewSubmissionHandler(submissionStore)
	eh := handler.NewEnrollmentHandler(enrollmentStore)
	ph := handler.NewProgressHandler(progressStore)
	apph := handler.NewApplicationHandler(applicationStore)

	mux := http.NewServeMux()

	// Class
	mux.HandleFunc("GET /api/v1/classes", ch.List)
	mux.HandleFunc("POST /api/v1/classes", ch.Create)
	mux.HandleFunc("GET /api/v1/classes/{id}", ch.Get)
	mux.HandleFunc("PUT /api/v1/classes/{id}", ch.Update)
	mux.HandleFunc("DELETE /api/v1/classes/{id}", ch.Delete)

	// Student
	mux.HandleFunc("GET /api/v1/students", sh.List)
	mux.HandleFunc("POST /api/v1/students", sh.Create)
	mux.HandleFunc("GET /api/v1/students/{id}", sh.Get)
	mux.HandleFunc("PUT /api/v1/students/{id}", sh.Update)
	mux.HandleFunc("DELETE /api/v1/students/{id}", sh.Delete)

	// Teacher
	mux.HandleFunc("GET /api/v1/teachers", th.List)
	mux.HandleFunc("POST /api/v1/teachers", th.Create)
	mux.HandleFunc("GET /api/v1/teachers/{id}", th.Get)
	mux.HandleFunc("PUT /api/v1/teachers/{id}", th.Update)
	mux.HandleFunc("DELETE /api/v1/teachers/{id}", th.Delete)

	// Session（课次 / 考勤）
	mux.HandleFunc("GET /api/v1/sessions", sessh.List)
	mux.HandleFunc("POST /api/v1/sessions", sessh.Create)
	mux.HandleFunc("GET /api/v1/sessions/{id}", sessh.Get)
	mux.HandleFunc("PUT /api/v1/sessions/{id}", sessh.Update)
	mux.HandleFunc("DELETE /api/v1/sessions/{id}", sessh.Delete)

	// Assessment
	mux.HandleFunc("GET /api/v1/assessments", ah.List)
	mux.HandleFunc("POST /api/v1/assessments", ah.Create)
	mux.HandleFunc("GET /api/v1/assessments/{id}", ah.Get)
	mux.HandleFunc("PUT /api/v1/assessments/{id}", ah.Update)
	mux.HandleFunc("DELETE /api/v1/assessments/{id}", ah.Delete)

	// Submission
	mux.HandleFunc("GET /api/v1/submissions", subh.List)
	mux.HandleFunc("POST /api/v1/submissions", subh.Create)
	mux.HandleFunc("GET /api/v1/submissions/{id}", subh.Get)
	mux.HandleFunc("PUT /api/v1/submissions/{id}", subh.Update)
	mux.HandleFunc("DELETE /api/v1/submissions/{id}", subh.Delete)

	// Enrollment
	mux.HandleFunc("GET /api/v1/enrollments", eh.List)
	mux.HandleFunc("POST /api/v1/enrollments", eh.Create)
	mux.HandleFunc("GET /api/v1/enrollments/{id}", eh.Get)
	mux.HandleFunc("PUT /api/v1/enrollments/{id}", eh.Update)
	mux.HandleFunc("DELETE /api/v1/enrollments/{id}", eh.Delete)

	// Progress
	mux.HandleFunc("GET /api/v1/progress", ph.List)
	mux.HandleFunc("POST /api/v1/progress", ph.Create)
	mux.HandleFunc("GET /api/v1/progress/{id}", ph.Get)
	mux.HandleFunc("PUT /api/v1/progress/{id}", ph.Update)
	mux.HandleFunc("DELETE /api/v1/progress/{id}", ph.Delete)

	// Application（立项申请：学员提交 + 后台查看，v0.1 不做审批）
	mux.HandleFunc("GET /api/v1/applications", apph.List)
	mux.HandleFunc("POST /api/v1/applications", apph.Create)
	mux.HandleFunc("GET /api/v1/applications/{id}", apph.Get)

	// Health
	mux.HandleFunc("GET /healthz", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok\n"))
	})

	return mux
}
