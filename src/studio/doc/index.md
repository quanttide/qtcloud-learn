# 量潮学习云 · Studio 文档索引

> LMS 管理后台页面与通用组件规格。
> 分解自实验室创新系统原型 v0.12（2026-08-07）。

---

## 内容层级

```
Course（课程）
  └── Lesson（课时）
       └── Scene（场景，嵌套在课时页面内）
```

| 层级 | 说明 | 示例 |
|------|------|------|
| Course | 一门完整课程，含课时列表和元信息 | 知识工作 / 氛围编程 / 大数据导论 / 数据工程 / 生产实习 |
| Lesson | 课程中的直接课时单元 | "搜索引擎高级技巧"（阅读 10 min） |
| Scene | 课时内的互动场景，通过分支选项串联 | intro / lecture / quiz / summary |

---

## 屏幕清单（后台 14 页）

### 框架与概览

| 屏幕 | 文档 |
|------|------|
| LMS 管理后台 · 框架页 | [admin.md](screens/admin.md) |
| 概览（后台首页） | [dashboard.md](screens/dashboard.md) |

### 学习管理

| 屏幕 | 文档 |
|------|------|
| 课题管理 | [topic-management.md](screens/topic-management.md) |
| 审批中心 | [approval-center.md](screens/approval-center.md) |
| 双创项目管理 | [innovation-project.md](screens/innovation-project.md) |
| 成果仓库 | [archive.md](screens/archive.md) |
| 成员管理 | [member-management.md](screens/member-management.md) |
| 课程管理 | [course-management.md](screens/course-management.md) |
| 学员管理 | [student-management.md](screens/student-management.md) |

### 表单配置

| 屏幕 | 文档 |
|------|------|
| 立项申请表 | [apply-form.md](screens/apply-form.md) |
| 阶段报告模板 | [stage-report-form.md](screens/stage-report-form.md) |
| 验收评审表 | [review-form.md](screens/review-form.md) |

### 关联系统

| 屏幕 | 文档 |
|------|------|
| 课程研发 | [course-dev.md](screens/course-dev.md) |
| 咨询/工单 | [consult.md](screens/consult.md) |

> 原型中 LMS 管理后台为单页 13 板块（`view-back`），此处按板块拆分为独立页面文档，由框架页（admin.md）侧边栏导航串联。

---

## 组件规格（Widgets）

| 组件 | 文件 | 说明 |
|------|------|------|
| Atoms | [atoms.md](widgets/atoms.md) | 原子组件：按钮 / 标签 / 圆点 / 图标 / 头像 / 连接线 / 分割线 |
| Containers | [containers.md](widgets/containers.md) | 容器组件：Card / Section / TwoCol / CourseHero / ModulePanel / PageHeader / Workspace / FormRow |
| Items | [items.md](widgets/items.md) | 条目组件：CourseCard / LessonItem / NavItem / StatCard / TimelineItem / PipelineStage |
| Navigation | [navigation.md](widgets/navigation.md) | 导航组件：AppBar / SideNav / StepBar / BackLink / Pipeline |
| Feedback | [feedback.md](widgets/feedback.md) | 反馈组件：Toast / ProgressBar / Overlay（待实现）/ Dialog（待实现） |
| Forms | [forms.md](widgets/forms.md) | 表单组件：FormGroup / Input / Textarea / Select / Radio |

---

## 文档格式说明

每屏一文件 `.md`，结构化文档规格：页面定位 → ASCII 结构图 → 组件清单（含选择器）→ 交互表 → 数据模型 → 状态机 → 设计原则。

数据模型（Course / Lesson / Scene / Progress / AppState 等）与 HTML 视觉原型保留在实验室仓库，未随页面分解迁移：

- 数据模型：[innovation-prototype/models/course.md](https://github.com/quanttide/quanttide-laboratory-of-course-development/blob/main/innovation-prototype/models/course.md)
- 视觉原型：[innovation-prototype/prototype.html](https://quanttide.github.io/quanttide-laboratory-of-course-development/innovation-prototype/prototype.html)

---

## 相关链接

- 实验室仓库：[quanttide-laboratory-of-course-development](https://github.com/quanttide/quanttide-laboratory-of-course-development)
- 设计系统：[DESIGN.md · 海岸设计系统](https://github.com/quanttide/quanttide-laboratory-of-course-development/blob/main/DESIGN.md)
- 文档体系参考：[qtclass/src/studio/doc](https://github.com/quanttide/qtclass/tree/main/src/studio/doc)
