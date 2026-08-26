# CHANGELOG

## [Unreleased]

### Added

- 子命令对齐《量潮学习管理标准》：learner / criterion / completion（含 `completion complete` 标记完成）

### Removed

- 移除已从 provider 砍掉的子命令：student / class / enrollment / progress / assessment

## [0.1.0] - 2026-08-01

### Added

- 学员侧子命令：student / class / enrollment / progress / assessment（承接 `qtcloud-course` cli 规划能力）
- Provider API 客户端（ureq，`/api/v1` 前缀，`--base-url` 可配置）
- 集成测试：进程内迷你 provider（内存 CRUD）端到端验证各子命令

## [0.0.1] - 2026-08-01

### Added

- CLI 入口，`version` 子命令
