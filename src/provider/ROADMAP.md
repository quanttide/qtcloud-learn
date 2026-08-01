# ROADMAP

> 本 scope 规划对齐根 [ROADMAP.md](../../ROADMAP.md) 的迁移阶段。

## [v0.2] — LMS API 迁移

> **目标**：承接 `qtcloud-course/provider` 的 LMS API，成为学习云唯一服务端。

### 交付物

- [x] 领域模型：Student / Teacher / Class / Session / Enrollment / Progress / Assessment / Submission
- [x] Store + Handler：class / student / assessment / submission / enrollment / progress CRUD
- [x] 路由注册（`/api/v1/...`）与健康检查
- [x] 从 `qtcloud-course/provider` 移植 class.go（domain / store / handler）

### 测试

- [x] 移植的 CRUD 测试全部通过
- [x] 新增 enrollment / progress API 测试
- [x] 健康检查 `/healthz` 测试
