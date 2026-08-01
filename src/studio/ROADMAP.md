# ROADMAP

> 本 scope 规划对齐根 [ROADMAP.md](../../ROADMAP.md) 的迁移阶段。

## [v0.3] — LMS 页面迁移

> **目标**：承接 `qtclass` 与 `qtcloud-course` 的全部 LMS 页面，两仓库回归单一职责。

### 交付物

- [x] 登录页（飞书登录 · 本地模拟）
- [x] 我的课程 / 选课报名页
- [x] 学习进度页
- [x] 课表（自 `qtclass` schedule）与考勤（自 `qtclass` classroom）
- [x] 考核列表 / 详情 / 提交（自 `qtcloud-course` assessment）
- [x] 班级管理页 + 仪表盘（自 `qtcloud-course` class / dashboard）
- [x] 学习记录由 localStorage 改为服务端进度数据
- [x] 模型 / 服务全量迁移（student / teacher / class_teaching / assessment / submission / enums、data_service / app_state），明细见根 ROADMAP 的 `docs/lms-inventory.md`

### 测试

- [x] 各页面 widget 测试
- [ ] `qtclass` 移除 LMS 模块后播放器测试保持全绿（待 v0.5）
