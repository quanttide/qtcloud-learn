# ROADMAP

> 产品级版本规划，侧重测试与文档。各 scope 详细路线图见 `src/*/ROADMAP.md`。

## 背景

`qtclass`（量潮课堂）与 `qtcloud-course`（量潮课程云）中存在大量早期"学习管理系统"（LMS）代码，
与各自核心职责（播放器客户端 / 课程制作）混杂，影响迭代。本路线图规划将这些 LMS 能力
统一收拢到 `qtcloud-learn`（量潮学习云），迁移完成后从旧仓库移除。

**迁移原则：**

1. **先建新、后删旧**：LMS 能力先在 `qtcloud-learn` 落地并通过测试，再从旧仓库移除。
2. **冻结旧代码**：自 v0.1 起，`qtclass` / `qtcloud-course` 不再新增 LMS 功能，仅做缺陷修复（冻结声明落点：`qtcloud-course/ROADMAP.md`、`qtclass/STATUS.md`）。
3. **模型唯一**：以 `qtcloud-learn` 为唯一事实来源；旧仓库的重复模型（如 `Class`、`Student`）只映射、不复用。

**职责边界：**

- 认证（飞书登录）与学员数据属 `qtcloud-learn/provider`；`qtcloud-course` 不再规划认证 API。
- 播放器（player / lecture）保留在 `qtclass`；学员侧入口统一由 `qtcloud-learn/studio` 承接，`qtcloud-course` 不再规划学生端播放器。

**版本依赖：**

- v0.2（provider）与 v0.3（studio）可并行；v0.3 学习进度页依赖 v0.2 的 progress API。
- v0.3 登录页（飞书登录）依赖 provider 认证能力，认证 API 未在本路线图排期，需另行规划或先本地模拟。
- v0.5 依赖 v0.1–v0.4 全部完成。

代码盘点清单见 [docs/lms-inventory.md](docs/lms-inventory.md)。

## [v0.1] — 盘点与冻结

> **目标**：确认迁移边界，冻结旧仓库 LMS 代码演化。

### 交付物

- [x] 盘点两仓库 LMS 代码清单（见 `docs/lms-inventory.md`）
- [x] `qtclass` / `qtcloud-course` 声明冻结 LMS 功能迭代（落点：`qtcloud-course/ROADMAP.md`、`qtclass/STATUS.md`）
- [x] 定义 `qtcloud-learn` 统一领域模型：Student / Teacher / Class / Session / Enrollment / Progress / Assessment / Submission；LearningRecord 不独立建模，并入 Progress 作为服务端进度数据

### 测试

- [x] 三仓库模型字段对齐核对（qtclass ↔ qtcloud-course ↔ qtcloud-learn）

## [v0.2] — Provider 迁移（Go 服务端）

> **目标**：`qtcloud-learn/provider` 承接全部 LMS API，替代 `qtcloud-course/provider` 的 class 相关接口。

### 交付物

- [x] 领域模型：Student / Teacher / Class / Session / Enrollment / Progress / Assessment / Submission
- [x] Store + Handler：class / student / assessment / submission / enrollment / progress 的 CRUD
- [x] 路由注册（`/api/v1/...`）与健康检查
- [x] 从 `qtcloud-course/provider` 移植 `class.go`（domain / store / handler）与测试

### 测试

- [x] 移植的 CRUD 测试全部通过
- [x] 新增 enrollment / progress API 测试
- [x] 健康检查 `/healthz` 测试

## [v0.3] — Studio 迁移（Flutter 客户端）

> **目标**：`qtcloud-learn/studio` 承接全部 LMS 页面，`qtclass` / `qtcloud-course` 移除 LMS 页面。

### 交付物

- [x] 登录页（飞书登录 · 本地模拟）
- [x] 我的课程 / 选课报名页
- [x] 学习进度页（依赖 v0.2 progress API）
- [x] 课表（自 `qtclass` schedule）与考勤（自 `qtclass` classroom）
- [x] 考核列表 / 详情 / 提交（自 `qtcloud-course` assessment）
- [x] 班级管理页 + 仪表盘（自 `qtcloud-course` class / dashboard）
- [x] 学习记录由 localStorage（`qtclass` history_service）改为服务端进度数据
- [x] 模型 / 服务全量迁移：student / teacher / class_teaching / assessment / submission / enums（LMS 部分）与 data_service / app_state（sessions 加载），文件级明细见 `docs/lms-inventory.md`

### 测试

- [x] 各页面 widget 测试
- [ ] `qtclass` 移除 LMS 模块后，播放器测试保持全绿（待 v0.5）

## [v0.4] — CLI 迁移（Rust）

> **目标**：`qtcloud-learn/cli` 提供学员侧子命令（承接 `qtcloud-course` cli 中规划的 class / student 能力）。

### 交付物

- [ ] student / class / enrollment / progress / assessment 子命令

### 测试

- [ ] 子命令 CLI 测试

## [v0.5] — 旧代码移除

> **目标**：`qtcloud-course` 与 `qtclass` 中 LMS 代码全部移除，仓库回归单一职责。

### 交付物

- [ ] `qtcloud-course/provider`：移除 `class.go`（domain / store / handler）与相关测试
- [ ] `qtcloud-course/studio`：移除 assessment / class / student / teacher / submission 模型、页面、服务，`enums.dart`（LMS 部分）、`dashboard_screen.dart`、`data_service.dart`（LMS 部分）与对应 assets JSON
- [ ] `qtclass/studio`：移除 session（课表/考勤）、learning_record、schedule / classroom / result 页面与 `sessions.json`（`lectures.json` 与播放器代码保留）
- [ ] 数据迁移：旧 assets JSON 与 localStorage 数据导入 `qtcloud-learn`（localStorage 在客户端本地，需一次性导入脚本写入服务端进度数据）

### 测试

- [ ] 两仓库移除后 build / test 全绿
- [ ] 按 `docs/lms-inventory.md` 文件清单逐项核对移除；`grep` 校验注意 `class` 是 Dart 关键字，需限定文件名 / 标识符避免误报

## 后续

> 本路线图仅覆盖 LMS 迁移阶段（v0.1–v0.5）。迁移完成后，`qtcloud-learn` 的产品能力演进
> （认证完善、付费 / VIP 权益、学员侧播放等）另行规划。
