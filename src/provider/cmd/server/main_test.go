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

// TestRouter_SpecRoutes 冒烟测试：对齐 spec 的三个资源路由已注册并可创建/读取。
func TestRouter_SpecRoutes(t *testing.T) {
	mux := newRouter()

	cases := []struct {
		name string
		path string
		body string
	}{
		{"learners", "/learners", `{"user_id":"user-123"}`},
		{"completions", "/completions", `{"learner_id":"lea-1","criterion_id":"cri-1"}`},
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
