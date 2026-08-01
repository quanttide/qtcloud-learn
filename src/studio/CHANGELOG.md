# CHANGELOG

## [0.1.0] - 2026-08-01

### Added

- 统一领域模型（models/）：Student / Teacher / Class / Session / Enrollment / Progress / Assessment / Submission / enums
- 数据服务（services/）：LearnDataService + AssessmentService，支持 API（`/api/v1`）与离线 assets 兜底
- 登录页（飞书登录 · 本地模拟）与主壳（底部导航：仪表盘 / 我的课程 / 课表 / 班级）
- 仪表盘（学员侧指标）、我的课程 / 选课报名、学习进度上报页
- 课表 / 课堂考勤页（自 `qtclass` schedule / classroom 迁移）
- 班级管理 / 考核管理（自 `qtcloud-course` class / assessment 迁移）
- 离线演示数据（assets/）与 widget 测试（登录 / 导航 / 各页面）

### Changed

- 首页替换为登录门控 + 主壳（原骨架首页移除）
- 学习记录由 localStorage 改为服务端进度数据（Progress 上报）

## [0.0.1] - 2026-08-01

### Added

- Flutter 项目骨架（Linux / Android / iOS / Web / macOS）
- 首页：量潮学习云 · AI 原生学员学习中心
