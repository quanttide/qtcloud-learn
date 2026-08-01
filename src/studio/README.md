# 学习云 Studio (`qtcloud_learn_studio`)

量潮学习云的客户端。Flutter 编写。

## 开发

```bash
# 运行测试
flutter test

# 静态检查
flutter analyze

# 运行（Linux）
flutter run -d linux
```

## 目录

```
lib/models/     # 统一领域模型（Student / Teacher / Class / Session / Enrollment / Progress / Assessment / Submission）
lib/services/   # 数据服务（LearnDataService / AssessmentService，支持 API 与离线 assets）
lib/screens/    # 页面（登录 / 仪表盘 / 我的课程 / 进度 / 课表 / 课堂考勤 / 班级 / 考核）
lib/widgets/    # 通用组件（StatusChip / MetricCard 等）
assets/         # 离线演示数据（本地模式默认读取）
```

## 数据模式

- **本地模式（默认）**：从 `assets/*.json` 读取演示数据。
- **API 模式**：`LearnDataService(baseUrl: ...)` / `AssessmentService(baseUrl: ...)`
  对接 `qtcloud-learn-provider`（`/api/v1`），失败自动回退本地 assets。
- **登录**：认证 API 未排期，当前为本地模拟登录（默认第一个学员）。

## 相关文档

- [ROADMAP.md](ROADMAP.md) — 路线图
- [CHANGELOG.md](CHANGELOG.md) — 变更记录
