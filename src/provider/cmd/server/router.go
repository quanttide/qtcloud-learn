package main

import (
	"net/http"
	"os"

	"github.com/quanttide/qtcloud-learn-provider/internal/handler"
	"github.com/quanttide/qtcloud-learn-provider/internal/store"
)

// newRouter 创建并配置所有路由，可单独测试。
// API 无版本前缀（资源直挂根路径），对齐《量潮学习管理标准》（docs/specification）：
// Learner × Task → Completion，Schedule 为 Task 的有序集合。
// 持久化（Learner / Completion / Schedule / Task）：
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
	completionStore := store.NewCompletionStore()
	scheduleStore := store.NewScheduleStore()
	taskStore := store.NewTaskStore()

	if persister != nil {
		learnerStore.BaseStore.SetPersister(persister)
		_ = learnerStore.BaseStore.Load("lea.json")
		completionStore.BaseStore.SetPersister(persister)
		_ = completionStore.BaseStore.Load("com.json")
		scheduleStore.BaseStore.SetPersister(persister)
		_ = scheduleStore.BaseStore.Load("schedule.json")
		taskStore.BaseStore.SetPersister(persister)
		_ = taskStore.BaseStore.Load("task.json")
	}

	learnerh := handler.NewLearnerHandler(learnerStore)
	comh := handler.NewCompletionHandler(completionStore)
	scheduleh := handler.NewScheduleHandler(scheduleStore)
	taskh := handler.NewTaskHandler(taskStore)

	mux := http.NewServeMux()

	// Learner（学习者）
	mux.HandleFunc("GET /learners", learnerh.List)
	mux.HandleFunc("POST /learners", learnerh.Create)
	mux.HandleFunc("GET /learners/{id}", learnerh.Get)
	mux.HandleFunc("PUT /learners/{id}", learnerh.Update)
	mux.HandleFunc("DELETE /learners/{id}", learnerh.Delete)

	// Completion（完成记录）
	mux.HandleFunc("GET /completions", comh.List)
	mux.HandleFunc("POST /completions", comh.Create)
	mux.HandleFunc("GET /completions/{id}", comh.Get)
	mux.HandleFunc("PUT /completions/{id}", comh.Update)
	mux.HandleFunc("DELETE /completions/{id}", comh.Delete)

	// Schedule（训练营/学习路径）
	mux.HandleFunc("GET /schedules", scheduleh.List)
	mux.HandleFunc("POST /schedules", scheduleh.Create)
	mux.HandleFunc("GET /schedules/{id}", scheduleh.Get)
	mux.HandleFunc("PUT /schedules/{id}", scheduleh.Update)
	mux.HandleFunc("DELETE /schedules/{id}", scheduleh.Delete)

	// Task（任务）
	mux.HandleFunc("GET /tasks", taskh.List)
	mux.HandleFunc("POST /tasks", taskh.Create)
	mux.HandleFunc("GET /tasks/{id}", taskh.Get)
	mux.HandleFunc("PUT /tasks/{id}", taskh.Update)
	mux.HandleFunc("DELETE /tasks/{id}", taskh.Delete)

	// Health
	mux.HandleFunc("GET /healthz", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok\n"))
	})

	return mux
}
