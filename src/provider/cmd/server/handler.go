package main

import "net/http"

// handler 返回服务路由。后续在 internal/handler 中扩展业务 API。
func handler() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /healthz", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok\n"))
	})
	return mux
}
