# ROADMAP

> 本 scope 规划对齐根 [ROADMAP.md](../../ROADMAP.md) 的迁移阶段。

## [v0.1] — LMS 页面迁移（对齐根 ROADMAP v0.3）

> **目标**：承接 `qtclass` 与 `qtcloud-course` 的全部 LMS 页面，两仓库回归单一职责。

### 交付物

- [ ] 登录页（飞书登录）
- [ ] 我的课程 / 选课报名页
- [ ] 学习进度页
- [ ] 课表（自 `qtclass` schedule）与考勤（自 `qtclass` classroom）
- [ ] 考核列表 / 详情 / 提交（自 `qtcloud-course` assessment）
- [ ] 班级管理页（自 `qtcloud-course` class）
- [ ] 学习记录由 localStorage 改为服务端进度数据

### 测试

- [ ] 各页面 widget 测试
- [ ] `qtclass` 移除 LMS 模块后播放器测试保持全绿
