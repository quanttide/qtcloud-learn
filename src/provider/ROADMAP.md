# ROADMAP

> 本 scope 规划对齐根 [ROADMAP.md](../../ROADMAP.md) 的迁移阶段。

## [v0.1] — LMS API 迁移（对齐根 ROADMAP v0.2）

> **目标**：承接 `qtcloud-course/provider` 的 LMS API，成为学习云唯一服务端。

### 交付物

- [ ] 领域模型：Student / Teacher / Class / Session / Enrollment / Progress / Assessment / Submission
- [ ] Store + Handler：class / student / assessment / submission / enrollment / progress CRUD
- [ ] 路由注册（`/api/v1/...`）与健康检查
- [ ] 从 `qtcloud-course/provider` 移植 class.go（domain / store / handler）

### 测试

- [ ] 移植的 CRUD 测试全部通过
- [ ] 新增 enrollment / progress API 测试
