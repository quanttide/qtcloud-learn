# CHANGELOG

## [Unreleased]

### Added

- Criterion（验收标准）与 Completion（完成记录）实体及 CRUD API：`/api/v1/criteria`、`/api/v1/completions`

### Changed

- 领域模型对齐《量潮学习管理标准》（Learner × Criterion → Completion），JSON 字段使用 spec 定义（`learner_id` / `criterion_id` / `created_at` / `updated_at`）
- Learner 精简为 spec 定义（`id` / `user_id`），移除进度上报 / 自动建档逻辑

### Removed

- 移除 spec 未定义的 LMS 实体与 API：Student / Teacher / Class / Session / Assessment / Submission / Enrollment / Progress / Application（含对应 store / handler / 路由 / 测试）
- 移除统一账号 JWT 鉴权中间件（原仅保护已删除的立项提交与进度上报接口）

## [0.1.0] - 2026-08-01

### Added

- 统一领域模型：Student / Teacher / Class / Session / Enrollment / Progress / Assessment / Submission（`internal/domain`）
- 内存存储与 CRUD handler：class / student / teacher / session / assessment / submission / enrollment / progress
- LMS API 路由（`/api/v1/*`），自 `qtcloud-course/provider` 移植 `class.go`（domain / store / handler）与测试
- 领域 / 存储 / handler 层测试全覆盖（含 enrollment / progress / teacher / session 新增测试）

### Fixed

- `version.go` 由根目录 `package main` 修正为 `internal/version`，`go build ./...` 恢复可用
- `PUT` 改为 JSON 合并语义：仅覆盖请求体中的字段（修复评分等局部更新抹掉 assessmentId / studentId 的问题）

## [0.0.1] - 2026-08-01

### Added

- 服务入口 `cmd/server`，`/healthz` 健康检查
