package main

import (
	"net/http"
	"os"

	"github.com/quanttide/qtcloud-learn-provider/internal/handler"
	"github.com/quanttide/qtcloud-learn-provider/internal/store"
)

// newRouter 创建并配置所有路由，可单独测试。
// API 统一挂在 /api/v1 前缀下，资源对齐《量潮学习管理标准》（docs/specification）：
// Learner × Criterion → Completion。
// 持久化（Learner / Criterion / Completion 三个实体）：
//   - OSS_BUCKET 非空 → OSS 对象存储（生产 FC：实例盘不持久，跨实例/发版不丢）
//   - 否则 DATA_DIR 非空 → 本地 JSON 文件（dev/测试）
//   - 都为空 → 纯内存（测试默认）
func newRouter() *http.ServeMux {
	var persister store.Persister
	if bucket := os.Getenv("OSS_BUCKET"); bucket != "" {
		var err error
		persister, err = store.NewOSSPersister(
			bucket,
			os.Getenv("OSS_ENDPOINT"),
			os.Getenv("OSS_KEY_PREFIX"),
			os.Getenv("ALIYUN_ACCESS_KEY_ID"),
			os.Getenv("ALIYUN_ACCESS_KEY_SECRET"),
		)
		if err != nil {
			panic(err)
		}
	} else if dataDir := os.Getenv("DATA_DIR"); dataDir != "" {
		persister = store.NewFilePersister(dataDir)
	}

	learnerStore := store.NewLearnerStore()
	criterionStore := store.NewCriterionStore()
	completionStore := store.NewCompletionStore()

	if persister != nil {
		learnerStore.BaseStore.SetPersister(persister)
		_ = learnerStore.BaseStore.Load("lea.json")
		criterionStore.BaseStore.SetPersister(persister)
		_ = criterionStore.BaseStore.Load("cri.json")
		completionStore.BaseStore.SetPersister(persister)
		_ = completionStore.BaseStore.Load("com.json")
	}

	learnerh := handler.NewLearnerHandler(learnerStore)
	crith := handler.NewCriterionHandler(criterionStore)
	comh := handler.NewCompletionHandler(completionStore)

	mux := http.NewServeMux()

	// Learner（学习者）
	mux.HandleFunc("GET /api/v1/learners", learnerh.List)
	mux.HandleFunc("POST /api/v1/learners", learnerh.Create)
	mux.HandleFunc("GET /api/v1/learners/{id}", learnerh.Get)
	mux.HandleFunc("PUT /api/v1/learners/{id}", learnerh.Update)
	mux.HandleFunc("DELETE /api/v1/learners/{id}", learnerh.Delete)

	// Criterion（验收标准）
	mux.HandleFunc("GET /api/v1/criteria", crith.List)
	mux.HandleFunc("POST /api/v1/criteria", crith.Create)
	mux.HandleFunc("GET /api/v1/criteria/{id}", crith.Get)
	mux.HandleFunc("PUT /api/v1/criteria/{id}", crith.Update)
	mux.HandleFunc("DELETE /api/v1/criteria/{id}", crith.Delete)

	// Completion（完成记录）
	mux.HandleFunc("GET /api/v1/completions", comh.List)
	mux.HandleFunc("POST /api/v1/completions", comh.Create)
	mux.HandleFunc("GET /api/v1/completions/{id}", comh.Get)
	mux.HandleFunc("PUT /api/v1/completions/{id}", comh.Update)
	mux.HandleFunc("DELETE /api/v1/completions/{id}", comh.Delete)

	// Health
	mux.HandleFunc("GET /healthz", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok\n"))
	})

	return mux
}
