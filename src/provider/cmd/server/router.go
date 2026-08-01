package main

import (
	"net/http"

	"github.com/quanttide/qtcloud-learn-provider/internal/handler"
	"github.com/quanttide/qtcloud-learn-provider/internal/store"
)

// newRouter 创建并配置所有路由，可单独测试。
// LMS API 统一挂在 /api/v1 前缀下。
func newRouter() *http.ServeMux {
	classStore := store.NewClassStore()
	studentStore := store.NewStudentStore()
	assessmentStore := store.NewAssessmentStore()
	submissionStore := store.NewSubmissionStore()
	enrollmentStore := store.NewEnrollmentStore()
	progressStore := store.NewProgressStore()

	ch := handler.NewClassHandler(classStore)
	sh := handler.NewStudentHandler(studentStore)
	ah := handler.NewAssessmentHandler(assessmentStore)
	subh := handler.NewSubmissionHandler(submissionStore)
	eh := handler.NewEnrollmentHandler(enrollmentStore)
	ph := handler.NewProgressHandler(progressStore)

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

	// Health
	mux.HandleFunc("GET /healthz", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok\n"))
	})

	return mux
}
