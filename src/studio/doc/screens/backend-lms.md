# LMS 管理后台 (view-back)

学习管理系统的后台管理页面，统一管理课题、审批、双创项目、成果、成员、课程和学员。通过左侧 220px 侧边栏导航，右侧主内容区展示 13 个业务板块。

## 页面定位

- **路由**: 点击 AppBar "⚙ LMS 管理后台"按钮
- **状态**: `currentView === 'back'` — 独立视图，通过侧边栏导航切换内容区域
- **数据依赖**: 无（当前为静态 Demo 数据，所有表格/卡片内容硬编码在 HTML 中）

## 页面结构总览

```
+-- AppBar (sticky, 56px) ----------------------------------------------+
|  [~ 量潮课堂 · LMS]          [🏠 学习前台] [⚙ LMS 管理后台]   v0.12  Z |
+-----------------------------------------------------------------------+

+-- #view-back (flex) --------------------------------------------------+
|                                                                       |
|  +-- .sidenav (220px) ----+  +-- .main (overflow-y:auto) ----------+ |
|  |                        |  |                                      | |
|  | v LMS · 学习管理       |  |  +- .section · 概览 --------------+  | |
|  |   📋 课题管理          |  |  |  stat-row                       |  | |
|  |   🔍 审批中心          |  |  |  +------++------++------++----+ |  | |
|  |   🎓 双创项目管理      |  |  |  |5     |||12    |||8     |||3   | |  | |
|  |   📁 成果仓库          |  |  |  |待审批|||执行中|||已完成|||否决| |  | |
|  |   👥 成员管理          |  |  |  +------++------++------++----+ |  | |
|  |   📖 课程管理          |  |  +---------------------------------+  | |
|  |   👤 学员管理          |  |                                      | |
|  |                        |  |  +- .section · 课题管理 ----------+  | |
|  | v 表单配置             |  |  |  .card > .datatable             |  | |
|  |   📝 立项申请表        |  |  |  4 行课题数据 + 审批操作按钮     |  | |
|  |   📋 阶段报告模板      |  |  +---------------------------------+  | |
|  |   ✅ 验收评审表        |  |                                      | |
|  |                        |  |  +- .section · 审批中心 ----------+  | |
|  | v 关联系统             |  |  |  .twocol                        |  | |
|  |   📚 课程研发          |  |  |  |-- .timeline (4 节点)         |  | |
|  |   💬 咨询/工单         |  |  |  +-- .side (审批表单)           |  | |
|  +------------------------+  |  +---------------------------------+  | |
|                               |                                      | |
|                               |  +- .section · 双创项目管理 -------+  | |
|                               |  |  .card > .pipeline (5 阶段)     |  | |
|                               |  +---------------------------------+  | |
|                               |                                      | |
|                               |  ... (共 13 个 .section)             | |
|                               +--------------------------------------+ |
+-----------------------------------------------------------------------+
```

## 组件清单

### 侧边栏 (SideNav)

| 组件 | 选择器/ID | 说明 |
|------|-----------|------|
| 侧边栏容器 | `.sidenav` | 220px 宽，`#f5f9fc` 背景，右侧 1px 边框，`padding: 16px 0` |
| 导航分组 | `.nav-group` | `padding: 0 16px`，`margin-bottom: 16px` |
| 分组标题 | `.nav-label` | 11px 字母间距 1px，含 `.arr` 箭头（`▼` 展开 / ▶ 折叠），可点击折叠 |
| 箭头 | `.arr` | 10px 字号，`.folded` 时 `transform: rotate(-90deg)` |
| 子项容器 | `.nav-kids` | `overflow: hidden`，通过 `max-height` transition 折叠动画 |
| 折叠态 | `.nav-kids.folded` | `max-height: 0!important` |
| 导航项 | `.nav-item` | 13px，`padding: 9px 14px`，圆角 6px，hover 浅蓝背景，含 ripple |
| 子级导航项 | `.nav-item.sub` | `padding-left: 32px`（缩进 32px） |
| 占位项 | `.nav-item.placeholder` | `opacity: 0.6`，灰色不可点击 |
| 激活态 | `.nav-item.active` | 浅蓝背景 + 蓝色文字 + 粗体 |

### 分组结构

```
v LMS · 学习管理              (.nav-label, onclick="toggleNav('lms')")
  |-- 📋 课题管理              (.nav-item -> scrollToCard('topic-table'))
  |-- 🔍 审批中心              (.nav-item -> scrollToCard('approval-card'))
  |-- 🎓 双创项目管理          (.nav-item -> scrollToCard('sc-card'))
  |-- 📁 成果仓库              (.nav-item -> scrollToCard('archive-card'))
  |-- 👥 成员管理              (.nav-item -> scrollToCard('member-card'))
  |-- 📖 课程管理              (.nav-item -> scrollToCard('course-mgmt-card'))
  +-- 👤 学员管理              (.nav-item -> scrollToCard('student-card'))

v 表单配置                    (.nav-label, onclick="toggleNav('form')")
  |-- 📝 立项申请表            (.nav-item.sub -> scrollToCard('apply-card'))
  |-- 📋 阶段报告模板          (.nav-item.sub -> scrollToCard('stage-report-card'))
  +-- ✅ 验收评审表            (.nav-item.sub -> scrollToCard('review-card'))

v 关联系统                    (.nav-label, onclick="toggleNav('rel')")
  |-- 📚 课程研发              (.nav-item.sub -> scrollToCard('course-dev-card'))
  +-- 💬 咨询/工单             (.nav-item.sub -> scrollToCard('consult-card'))
```

### 主内容区布局组件

| 组件 | 选择器/ID | 说明 |
|------|-----------|------|
| 主滚动区 | `.main` | `flex:1; overflow-y:auto`，`padding: 0 32px 60px` |
| 内容分区 | `.section` | `margin-bottom: 24px` |
| 分区标题栏 | `.section-hd` | flex 横向，`h3` 标题 + 右侧说明文字 |
| 分区标题 | `.section-hd h3` | 16px 粗体 |
| 分区内容 | `.section-body` | 卡片容器的直接父容器 |

### 内容组件

| 组件 | 选择器/ID | 说明 |
|------|-----------|------|
| 统计卡片行 | `.stat-row` | flex 横向，`gap: 16px`，均分 |
| 统计卡片 | `.stat-card` | flex:1，白色背景，大数字（30px 粗体）+ 标签（13px 灰色），hover 上浮 1px |
| 数据表格 | `.datatable` | 全宽表格，13px 字号，表头灰底 2px 底线，行 hover 浅蓝背景 |
| 表格容器卡片 | `.card` (含 `.datatable`) | `padding: 0`，表格填满卡片 |
| 时间线 | `.timeline` | `padding-left: 20px`，左侧 2px 实线边框 |
| 时间线节点 | `.timeline-item` | `::before` 伪元素渲染 12px 圆点 + 蓝色边框 |
| 时间戳 | `.timeline-item .ts` | 11px 灰色 |
| 事件描述 | `.timeline-item .event` | 14px |
| 流程条 | `.pipeline` | flex 横向均分，阶段节点之间 `::after` 伪元素渲染 `→` |
| 阶段节点 | `.pipeline .stage` | flex:1，含 `.icon`（44px 圆）、`.name`（13px 粗体）、`.hint`（11px 灰色） |
| 双栏布局 | `.twocol` | flex 横向，`gap: 24px`，子元素 flex:1 |
| 固定侧栏 | `.twocol .side` | `flex: 0 0 280px` |
| 表单组 | `.form-group` | `margin-bottom: 16px`，含 `label` + `input/textarea/select` |
| 表单行 | `.form-row` | flex 横向，`gap: 16px`，子元素 flex:1 |
| 分割线 | `hr.divider` | `border-top: 1px solid var(--line-light)`，`margin: 20px 0` |
| 项目卡片 | 成果仓库内卡片 | flex:1 min-width:200px，浅灰背景，图标 + 标题 + 作者时间 + 状态标签 |

### 按钮

| 变体 | 选择器 | 用途 |
|------|--------|------|
| Filled | `.btn.btn-filled` | 主操作（审批、提交、通过） |
| Filled 小 | `.btn.btn-filled.btn-sm` | 表格内紧凑操作 |
| Outlined | `.btn.btn-outlined` | 次要操作（审批、新建） |
| Text | `.btn.btn-text` | 辅助操作（详情、查看、返回、驳回） |
| Text 小 | `.btn.btn-text.btn-sm` | 表格内文字操作 |

### Badge 状态标签

| 变体 | 选择器 | 用途 |
|------|--------|------|
| 进行中 | `.badge.active-status` | 浅绿背景，执行中/活跃状态 |
| 待审批 | `.badge.review` | 浅蓝背景，等待审批 |
| 已完成 | `.badge.done` | 浅紫背景，已验收/已结项 |
| 草稿 | `.badge.draft` | 浅橙背景，未发布 |
| 入门 | `.badge.beginner` | 浅蓝背景，课程难度 |
| 进阶 | `.badge.intermediate` | 浅靛背景 |
| 实战 | `.badge.advanced` | 浅橙背景 |
| 微型创业 | `.badge.capstone` | 浅红背景 |

## 后台板块清单（13 个）

| 顺序 | 板块 | ID / 锚点 | 类型 | 说明 |
|------|------|-----------|------|------|
| 1 | 概览 | - | 统计卡片行 | 4 张 stat-card：待审批(5)/执行中(12)/已完成(8)/已否决(3) |
| 2 | 课题管理 | `topic-table` | 数据表格 | 4 行课题：AI 内容流水线/数据看板/校企渠道/左移标准，含审批操作 |
| 3 | 审批中心 | `approval-card` | 时间线 + 表单 | 左栏 4 节点时间线（提交→分配→等待→待触发），右栏 280px 审批表单 |
| 4 | 双创项目管理 | `sc-card` | 流程条 + 说明 | Pipeline 五阶段 + 系统说明文字 |
| 5 | 成果仓库 | `archive-card` | 卡片网格 | 3 个项目卡片：数据看板/左移标准/AI 面试模拟器 |
| 6 | 成员管理 | `member-card` | 数据表格 | 4 人：李明/王芳/陈强/赵子奕，含角色和状态 |
| 7 | 课程管理 | `course-mgmt-card` | 数据表格 + 按钮 | 5 门课程配置信息 + "+ 新建课程"按钮 |
| 8 | 学员管理 | `student-card` | 数据表格 | 5 名学员：进度、活跃时间、状态 |
| 9 | 立项申请表 | `apply-card` | 内联表单 | 项目名称/盲区描述/Demo 方案/方向类型/组队方式 + 提交按钮 |
| 10 | 阶段报告模板 | `stage-report-card` | 表单 | 报告阶段/本周完成/困难/下周计划/帮助需求/进度评估 |
| 11 | 验收评审表 | `review-card` | 表单 + 评分 | 完成度/创新性/可复用性三维评分 + 验收结论下拉 |
| 12 | 课程研发 | `course-dev-card` | 卡片网格 | 2 张卡片：内容"左移"标准 + 知识工程教程 |
| 13 | 咨询/工单 | `consult-card` | 数据表格 | 3 条咨询记录，含紧急程度和状态 |

## 交互与导航

| 交互 | 触发条件 | 行为 |
|------|---------|------|
| 点击侧边栏菜单项 | 用户点击 `.nav-item` | `scrollToCard(id)` → 切换到 `view-back` + 取消所有菜单高亮 + 高亮当前项 + `scrollIntoView` 平滑滚动到目标卡片 |
| 点击导航分组标题 | 用户点击 `.nav-label` | `toggleNav(id)` → 折叠/展开子项，箭头旋转 90° |
| 点击"立项申请表" | 表单配置子项 | 直接滚动到 `#apply-card`（内联在后台页面中，非独立页） |
| 点击"LMS 管理后台" | AppBar 按钮 | `switchView('back')` → 进入后台，清除所有菜单高亮，页面从顶部开始 |
| 点击"学习前台" | AppBar 按钮 | `showCourseList()` → 回到课程列表页 |
| 点击"详情/查看" | 表格内 `.btn-text` 按钮 | 预留操作（当前无绑定事件） |
| 点击"审批/通过" | 表格/表单内 `.btn-filled` 按钮 | 预留操作（当前无实际审批逻辑） |
| 点击"驳回" | 审批表单 | 预留操作，红色文字 `#F25F5C` |
| 点击"+ 新建课程" | 课程管理卡片底部 | 预留按钮（当前无表单弹出） |

## 数据模型

详见实验室 [course.md](https://github.com/quanttide/quanttide-laboratory-of-course-development/blob/main/innovation-prototype/models/course.md)，与后台相关的核心结构：

### AppState.sidebarFolded

```javascript
sidebarFolded: {
  lms: false,   // LMS 分组折叠状态
  form: false,  // 表单配置分组折叠状态
  rel: false    // 关联系统分组折叠状态
}
```

### 侧边栏折叠逻辑 (toggleNav)

```
toggleNav(id):
  kids = document.getElementById('kids-' + id)
  arr  = document.getElementById('arr-' + id)

  if kids.classList.contains('folded'):
    → 展开：移除 .folded，恢复 max-height，箭头恢复 ▼
  else:
    → 折叠：设 max-height = scrollHeight → requestAnimationFrame → 添加 .folded，箭头旋转 ▶
```

### 菜单跳转逻辑 (scrollToCard)

```
scrollToCard(id):
  1. 切换视图到 view-back
  2. 更新 AppBar 按钮高亮
  3. 清除所有 .nav-item.active
  4. 找到 onclick 包含 id 的菜单项 → 添加 .active
  5. target.scrollIntoView({ behavior: 'smooth', block: 'start' })
```

## 设计原则

1. **左右联动** — 侧边栏和主内容区顺序严格一致，点击即跳转，无页面刷新
2. **信息密度适中** — 每个板块独立分区（`.section`），可独立滚动定位
3. **表单内联** — 立项申请表、阶段报告模板、验收评审表均嵌入后台页面，无需跳转外部
4. **状态可见** — 当前激活的菜单项蓝色高亮，折叠态箭头旋转指示
5. **操作就近** — 表格行内嵌入操作按钮（详情/审批），审批中心左侧时间线右侧即审批表单
6. **视觉层级** — 概览→数据表格→表单→详细信息，从上到下信息密度递增

## 状态机

```
[任意前台页面]
     |
     | 点击"⚙ LMS 管理后台"
     v
[view-back · 后台概览]
     |  (无默认高亮菜单项，页面停在顶部)
     |
     |-- 点击侧边栏菜单 --> scrollToCard(id) --> 滚动到目标板块 + 高亮菜单
     |
     |-- 点击分组标题 --> toggleNav(id) --> 折叠/展开
     |
     +-- 点击"🏠 学习前台" --> showCourseList() --> [课程列表页]
```
