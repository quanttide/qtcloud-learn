# CHANGELOG

## [0.1.0] - 2026-08-01

### Added

- 统一领域模型：Student / Teacher / Class / Session / Enrollment / Progress / Assessment / Submission（`internal/domain`）
- 内存存储与 CRUD handler：class / student / teacher / session / assessment / submission / enrollment / progress
- LMS API 路由（`/api/v1/*`），自 `qtcloud-course/provider` 移植 `class.go`（domain / store / handler）与测试
- 领域 / 存储 / handler 层测试全覆盖（含 enrollment / progress / teacher / session 新增测试）

### Fixed

- `version.go` 由根目录 `package main` 修正为 `internal/version`，`go build ./...` 恢复可用

## [0.0.1] - 2026-08-01

### Added

- 服务入口 `cmd/server`，`/healthz` 健康检查
