# CHANGELOG

## [Unreleased]

### Added

- 新增 Schedule / Task 领域模型与 CRUD API：`/schedules`、`/tasks`
- 领域模型改为引用 `quanttide-learn-toolkit/packages/go` type alias，避免跨仓复制模型
- 新增 `seeds/learning.json`，承载从学习档案/qtclass 学习中心清洗出的训练营与任务种子数据

### Changed

- Completion 从旧课程域 `criterion_id` 对齐为学习管理标准 `task_id`
- Learner 支持 `schedule_id` 字段，关联学习路径/训练营

## [0.1.0-alpha.5] - 2026-08-28

### Removed / Changed

- **破坏性**：移除 Criterion 本地模型与 `/criteria` 接口——验收标准的定义面在课程云，学习云不再注册副本；`completion.criterion_id` 字段名保留，改为直指课程域 Criterion.id（同源直连）

## [0.1.0-alpha.4] - 2026-08-26

### Fixed

- `DELETE` 操作落盘（`BaseStore.Delete` 调用 `persist()`），重启后已删除记录不再复活

## [0.1.0-alpha.3] - 2026-08-26

### Added

- Criterion（验收标准）与 Completion（完成记录）实体及 CRUD API：`/criteria`、`/completions`

### Changed

- 领域模型对齐《量潮学习管理标准》（Learner × Criterion → Completion），JSON 字段使用 spec 定义（`learner_id` / `criterion_id` / `created_at` / `updated_at`）
- Learner 精简为 spec 定义（`id` / `user_id`），移除进度上报 / 自动建档逻辑
- API 移除 `/api/v1` 前缀，资源直挂根路径

### Removed

- 移除 spec 未定义的 LMS 实体与 API：Student / Teacher / Class / Session / Assessment / Submission / Enrollment / Progress / Application（含对应 store / handler / 路由 / 测试）
- 移除统一账号 JWT 鉴权中间件（原仅保护已删除的立项提交与进度上报接口）

## [0.1.0-alpha.2] - 2026-08-23

### Added

- 统一账号 JWT 鉴权中间件（对应 tag `provider/v0.1.0-alpha.2`）：qtcloud-auth RSA 公钥验签，`JWT_PUBLIC_KEY` 未配置时本地 dev/测试可用
- 登录保护：`POST /api/proposals`（提交立项）、`POST /api/courses/prod/progress`（进度上报）需携带 Bearer Token

## [0.1.0-alpha.1] - 2026-08-16

### Added

- 统一领域模型：Student / Teacher / Class / Session / Enrollment / Progress / Assessment / Submission（`internal/domain`）
- 内存存储与 CRUD handler：class / student / teacher / session / assessment / submission / enrollment / progress
- LMS API 路由（`/api/v1/*`），自 `qtcloud-course/provider` 移植 `class.go`（domain / store / handler）与测试
- 立项申请 API（`/api/proposals`：5 问 + 方向类型 + 组队姓名，软删除 + 历史）+ 学员自动建档 + 进度上报
- 持久化升级：Persister 抽象（文件 / OSS 双后端，生产 FC 用 OSS 桶），立项与学员档案 JSON 落盘、重启恢复
- 部署配置：Dockerfile + terraform + CI workflow + OSS 持久化桶（对应 tag `provider/v0.1.0-alpha.1`）
- 领域 / 存储 / handler 层测试全覆盖（含 enrollment / progress / teacher / session 新增测试）

### Fixed

- `version.go` 由根目录 `package main` 修正为 `internal/version`，`go build ./...` 恢复可用
- `PUT` 改为 JSON 合并语义：仅覆盖请求体中的字段（修复评分等局部更新抹掉 assessmentId / studentId 的问题）
