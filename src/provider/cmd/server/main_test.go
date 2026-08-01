package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestHealthz(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/healthz", nil)
	rec := httptest.NewRecorder()

	newRouter().ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rec.Code)
	}
	if got := rec.Body.String(); got != "ok\n" {
		t.Fatalf("expected body %q, got %q", "ok\n", got)
	}
}

// TestRouter_LMSRoutes 冒烟测试：/api/v1 下各资源路由已注册并可创建/读取。
func TestRouter_LMSRoutes(t *testing.T) {
	mux := newRouter()

	cases := []struct {
		name   string
		path   string
		body   string
	}{
		{"classes", "/api/v1/classes", `{"name":"浙理班级","refId":"prog-1"}`},
		{"students", "/api/v1/students", `{"name":"张三"}`},
		{"teachers", "/api/v1/teachers", `{"name":"王老师"}`},
		{"sessions", "/api/v1/sessions", `{"classId":"class-1","lessonTitle":"Git 入门"}`},
		{"assessments", "/api/v1/assessments", `{"classId":"class-1","title":"作业1"}`},
		{"submissions", "/api/v1/submissions", `{"assessmentId":"assess-1","studentId":"stu-1"}`},
		{"enrollments", "/api/v1/enrollments", `{"classId":"class-1","studentId":"stu-1"}`},
		{"progress", "/api/v1/progress", `{"studentId":"stu-1","classId":"class-1"}`},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			// GET 列表应返回 200
			req := httptest.NewRequest(http.MethodGet, tc.path, nil)
			rec := httptest.NewRecorder()
			mux.ServeHTTP(rec, req)
			if rec.Code != http.StatusOK {
				t.Fatalf("GET %s = %d, want 200; body=%s", tc.path, rec.Code, rec.Body.String())
			}

			// POST 创建应返回 201
			req = httptest.NewRequest(http.MethodPost, tc.path, strings.NewReader(tc.body))
			req.Header.Set("Content-Type", "application/json")
			rec = httptest.NewRecorder()
			mux.ServeHTTP(rec, req)
			if rec.Code != http.StatusCreated {
				t.Fatalf("POST %s = %d, want 201; body=%s", tc.path, rec.Code, rec.Body.String())
			}
		})
	}
}
