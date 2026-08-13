# Items — 条目组件

可独立存在的列表条目，可点击，由原子组件组合而成。

---

## 1. CourseCard — 课程卡片

| 属性 | 值 |
|------|-----|
| 选择器 | `.course-card` |
| 布局 | 横向 flex，`gap: 20px`，`padding: 22px 28px` |

**子元素**：`.cc-num`（44px 编号圆）、`.cc-info`（`.cc-name` + `.cc-desc`）、`.badge`、箭头 `→`

**状态**：默认白底 → hover 蓝边框+右移 4px+阴影；`.active-card` 编号圆变蓝底白字

**使用位置**：课程列表页 `.card-grid`

---

## 2. LessonItem — 课时条目

| 属性 | 值 |
|------|-----|
| 选择器 | `.lesson-item` |
| 布局 | flex 横向，`gap: 11px`，`padding: 14px 16px`，底部分割线 |

**子元素**：`.ldot`（7px 圆点）、标题文本、`.duration`（右对齐 12px 灰色）

**状态**：hover 浅蓝背景+右移 4px、`.ldot.done` 圆点变绿

**使用位置**：课程内页 `.card` 内

---

## 3. NavItem — 侧边栏导航项

| 属性 | 值 |
|------|-----|
| 选择器 | `.nav-item` |
| 布局 | `display:block`，`padding: 9px 14px`，圆角 6px |

**变体**：`.nav-item.sub`（`padding-left: 32px`）、`.nav-item.placeholder`（`opacity:0.6`，不可点击）

**状态**：hover 浅蓝背景、`.active` 浅蓝背景+蓝色文字+粗体

**使用位置**：后台 `.sidenav` > `.nav-kids`

---

## 4. StatCard — 统计卡片

| 属性 | 值 |
|------|-----|
| 选择器 | `.stat-card` |
| 布局 | flex:1 均分，`padding: 20px` |

**子元素**：`.num`（30px 粗体）、`.label`（13px 灰色）

**状态**：hover 上浮 1px+阴影加深。容器：`.stat-row`（flex 横向）

**使用位置**：后台概览

---

## 5. TimelineItem — 时间线节点

| 属性 | 值 |
|------|-----|
| 选择器 | `.timeline-item` |
| 装饰 | `::before`：12px 蓝色圆点+白色边框+灰色外环 |

**子元素**：`.ts`（11px 灰色时间戳）、`.event`（14px 事件描述）

**容器**：`.timeline`（`padding-left:20px`，左侧 2px 边框）

**使用位置**：后台审批中心

---

## 6. PipelineStage — 流程阶段节点

| 属性 | 值 |
|------|-----|
| 选择器 | `.pipeline .stage` |
| 布局 | `flex:1`，`text-align:center`，节点间 `::after` 渲染 `→` |

**子元素**：`.icon`（44px 圆）、`.name`（13px 粗体）、`.hint`（11px 灰色）

**使用位置**：后台双创项目管理，5 阶段流程
