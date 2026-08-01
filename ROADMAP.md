# ROADMAP

> 产品级版本规划，侧重测试与文档。各 scope 详细路线图见 `src/*/ROADMAP.md`。

## 背景

LMS 迁移（旧 v0.1–v0.5，见 Git 历史）已完成：`qtclass` / `qtcloud-course` 中的学习管理代码
已统一收拢到 `qtcloud-learn`，旧仓库回归单一职责。本路线图规划迁移完成后的产品能力演进。

**现状基线（迁移完成时）：**

- `provider`（Go）：`/api/v1` CRUD 齐全（class / student / teacher / session / assessment / submission / enrollment / progress），**内存存储，重启即失**，无认证
- `studio`（Flutter）：登录（本地模拟）/ 仪表盘 / 我的课程 / 课表 / 班级 / 考核 / 进度页面齐全，数据**默认本地 assets + 离线回退**，未真实对接 provider
- `cli`（Rust）：student / class / enrollment / progress / assessment 五个子命令，对接 `/api/v1`

**迭代原则：**

1. 端到端优先：先打通 provider ↔ studio 真实链路，再叠加认证与权益。
2. 每版本可运行：每个版本结束时本地一键启动、全链路可用。

## [v1.0] — 可运行闭环（alpha）

> **目标**：本地完整可运行——数据持久化、端到端打通。

### 交付物

- [ ] `provider` 持久化：内存 store → JSON 文件存储（`data/` 目录），重启不丢
- [ ] `studio` 端到端：默认对接 provider `/api/v1`，本地 assets 仅作离线回退
- [ ] `cli` 端到端：子命令直连 provider 冒烟验证
- [ ] 本地一键启动脚本（provider + studio）

### 测试

- [ ] provider 持久化测试（重启后数据仍在）
- [ ] 集成测试：studio → provider 全链路（登录 → 选课 → 上报进度）
- [ ] cli 冒烟脚本

## [v1.1] — 认证与权限（beta）

> **目标**：真实飞书登录与角色权限。

### 交付物

- [ ] 飞书登录（OAuth）落地，替换 studio 本地模拟登录
- [ ] 角色权限：管理员 / 讲师 / 学生
- [ ] 登录态与服务端 token

### 测试

- [ ] 认证 API 测试
- [ ] 权限中间件测试（越权访问拒绝）

## [v1.2] — 权益与运营（beta）

> **目标**：付费 / VIP 权益落地。

### 交付物

- [ ] 学员套餐：free / paid / vip（`Student.plan` 字段已就位）
- [ ] 按套餐限制选课与考核
- [ ] 学员 / 班级运营报表

### 测试

- [ ] 权益校验测试

## [v2.0] — 完整学习体验（RC）

> **目标**：学员端完整闭环。

### 交付物

- [ ] 学员侧播放器（与 `qtclass` 播放器集成，进度回写服务端）
- [ ] 考核增强：题型 / 自动评分 / 计时 / 成绩统计
- [ ] 考勤与成绩报表

### 测试

- [ ] 播放器 ↔ provider 进度同步测试
- [ ] 考核全流程集成测试
