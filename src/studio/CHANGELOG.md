# CHANGELOG

## [0.1.1] - 2026-08-17

### Changed

- LMS 后台网页标题更新（`web/index.html`）

## [0.1.0] - 2026-08-16

### Added

- LMS 管理后台首发（对应 tag `studio/v0.1.0`）：概览 / 学员 / 进度 / 立项 四页
- 独立入口 `lib/main_admin.dart`，与学员端（`main.dart`）分离部署，前后台各构建一个 HTML
- 学员管理（进度 / 立项）与立项管理（组队 / 删除 / 历史），对齐赵原型
- 部署 workflow：`studio/*` tag → Flutter Web 构建（`--target=lib/main_admin.dart`）→ OSS 桶（qtcloud-learn-studio）→ CDN（learn.cloud.quanttide.com）

## [0.0.1] - 2026-08-01

> 学员端原型开发期，未纳入 tag 发布（仅保留在 git 历史中）。

### Added

- Flutter 项目骨架（Linux / Android / iOS / Web / macOS）
- 学员端原型：登录门控（飞书登录 · 本地模拟）+ 主壳（仪表盘 / 我的课程 / 课表 / 班级）
- 选课报名、学习进度上报、课表 / 课堂考勤（自 `qtclass` 迁移）、班级 / 考核管理（自 `qtcloud-course` 迁移）
- 离线演示数据（assets/）与 widget 测试（登录 / 导航 / 各页面）
