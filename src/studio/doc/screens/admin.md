# LMS 管理后台 · 框架页 (admin)

LMS 管理后台的页面框架（Shell）：AppBar + 侧边栏导航 + 主内容区。侧边栏菜单项跳转到各业务板块页面（本目录下各页面文档），是后台所有板块页的公共外壳。

## 页面定位

- **路由**: 后台入口 `/admin`，登录管理员身份进入
- **状态**: 侧边栏分组折叠状态（`sidebarFolded`）+ 当前激活菜单项
- **数据依赖**: 无（纯框架，内容由各板块页加载）

## 页面结构

```
+-- AppBar (sticky, 56px) ----------------------------------------------+
|  [~ 量潮学习云 · LMS]          [🏠 前台] [⚙ 管理后台]        v0.12  Z |
+-----------------------------------------------------------------------+

+-- .admin-layout (flex) -----------------------------------------------+
|                                                                       |
|  +-- .sidenav (220px) ------+  +-- .main (overflow-y:auto) --------+  |
|  |                          |  |                                    |  |
|  | v LMS · 学习管理         |  |  路由出口（Router Outlet）          |  |
|  |   📋 课题管理            |  |                                    |  |
|  |   🔍 审批中心            |  |  当前板块页渲染区域：               |  |
|  |   🎓 双创项目管理        |  |  dashboard / topic-management /     |  |
|  |   📁 成果仓库            |  |  approval-center / ...              |  |
|  |   👥 成员管理            |  |                                    |  |
|  |   📖 课程管理            |  |                                    |  |
|  |   👤 学员管理            |  |                                    |  |
|  |                          |  |                                    |  |
|  | v 表单配置               |  |                                    |  |
|  |   📝 立项申请表          |  |                                    |  |
|  |   📋 阶段报告模板        |  |                                    |  |
|  |   ✅ 验收评审表          |  |                                    |  |
|  |                          |  |                                    |  |
|  | v 关联系统               |  |                                    |  |
|  |   📚 课程研发            |  |                                    |  |
|  |   💬 咨询/工单           |  |                                    |  |
|  +--------------------------+  +------------------------------------+  |
+-----------------------------------------------------------------------+
```

## 组件清单

侧边栏组件详见 [widgets/navigation.md](../widgets/navigation.md)（SideNav）与 [widgets/atoms.md](../widgets/atoms.md)（AppBar 相关）。

| 组件 | 选择器/ID | 说明 |
|------|-----------|------|
| 侧边栏容器 | `.sidenav` | 220px 宽，`#f5f9fc` 背景，右侧 1px 边框 |
| 导航分组 | `.nav-group` | 3 组：LMS 学习管理 / 表单配置 / 关联系统 |
| 分组标题 | `.nav-label` | 可点击折叠，`.arr` 箭头（`▼` 展开 / `▶` 折叠） |
| 子项容器 | `.nav-kids` | `max-height` transition 折叠动画，`.folded` 完全折叠 |
| 导航项 | `.nav-item` | 13px，圆角 6px，hover 浅蓝背景；`.sub` 缩进 32px；`.placeholder` 灰显 |
| 激活态 | `.nav-item.active` | 浅蓝背景 + 蓝色文字 + 粗体 |
| 主内容区 | `.main` | `flex:1; overflow-y:auto`，路由出口 |

## 导航规则

| 交互 | 触发条件 | 行为 |
|------|---------|------|
| 点击侧边栏菜单项 | 点击 `.nav-item` | 路由跳转到对应板块页（原 `scrollToCard` 语义 → 页面导航），高亮当前项 |
| 点击分组标题 | 点击 `.nav-label` | 折叠/展开子项，箭头旋转 90° |
| 切换前后台 | AppBar 按钮 | 前台（qtcloud-learn 学员端）↔ 后台 |
| 深链进入 | 后台任意 URL | 加载框架 + 对应板块页，无默认高亮（停在各板块入口） |

## 分组与页面映射

| 分组 | 菜单项 | 页面文档 |
|------|--------|---------|
| LMS · 学习管理 | 📋 课题管理 | [topic-management.md](topic-management.md) |
| | 🔍 审批中心 | [approval-center.md](approval-center.md) |
| | 🎓 双创项目管理 | [innovation-project.md](innovation-project.md) |
| | 📁 成果仓库 | [archive.md](archive.md) |
| | 👥 成员管理 | [member-management.md](member-management.md) |
| | 📖 课程管理 | [course-management.md](course-management.md) |
| | 👤 学员管理 | [student-management.md](student-management.md) |
| 表单配置 | 📝 立项申请表 | [apply-form.md](apply-form.md) |
| | 📋 阶段报告模板 | [stage-report-form.md](stage-report-form.md) |
| | ✅ 验收评审表 | [review-form.md](review-form.md) |
| 关联系统 | 📚 课程研发 | [course-dev.md](course-dev.md) |
| | 💬 咨询/工单 | [consult.md](consult.md) |

> 概览（[dashboard.md](dashboard.md)）为后台默认首页，不占侧边栏菜单。

## 设计原则

1. **左右联动** — 侧边栏与板块页顺序一致，点击即跳转，无刷新
2. **分组即权限域** — 三个分组对应三类功能域（业务管理 / 表单模板 / 关联系统），便于按角色裁剪
3. **折叠态记忆** — 分组折叠状态持久化，刷新后保持

## 状态机

```
[登录] --> [后台框架]
             |
             |-- 默认 --> [dashboard 概览]
             |
             |-- 点击菜单 --> [对应板块页] <-- 高亮当前菜单项
             |
             +-- 点击"前台" --> [学员端]
```
