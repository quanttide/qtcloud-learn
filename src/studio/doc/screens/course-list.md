# 课程列表页 (view-courses)

系统默认入口页面，展示五门课程的阶梯式导航。所有用户首先看到的页面。

## 页面定位

- **路由**: `/` 或点击 AppBar "学习前台"
- **状态**: `currentView === 'courses'` — 默认视图
- **数据依赖**: 无（课程 1-4 数据定义在 `COURSES` JS 对象中，课程 5 硬编码在 HTML）；生产实习卡片可根据 `localStorage('qt-progress-prod')` 动态显示进度

## 页面结构

```
+-- AppBar (sticky, 56px) ----------------------------------------------+
|  [~ 量潮课堂 · LMS]          [🏠 学习前台] [⚙ LMS 管理后台]   v0.12  Z |
+-----------------------------------------------------------------------+

                              +- page-header -+
                              |                |
                              |  量潮课堂·实训基地  |
                              |  从基础到实战，五门  |
                              |  课程形成一条完整的  |
                              |  成长阶梯           |
                              |       💬 一对一咨询  |
                              +----------------+

                    +- card-grid (max-width: 700px) -+
                    |                                 |
                    |  +- course-card --------------+  |
                    |  | ① | 知识工作    [入门] ->  |  |
                    |  |   从整理文档开始，建立...   |  |
                    |  +----------------------------+  |
                    |                                 |
                    |  +- course-card --------------+  |
                    |  | ② | 氛围编程    [入门] ->  |  |
                    |  |   面向非技术岗位...         |  |
                    |  +----------------------------+  |
                    |                                 |
                    |  +- course-card --------------+  |
                    |  | ③ | 大数据导论  [进阶] ->  |  |
                    |  |   培养数据思维...           |  |
                    |  +----------------------------+  |
                    |                                 |
                    |  +- course-card --------------+  |
                    |  | ④ | 数据工程    [实战] ->  |  |
                    |  |   从数据到工程...           |  |
                    |  +----------------------------+  |
                    |                                 |
                    |  +- course-card (.active-card) + |
                    |  | ⑤ | 生产实习  [微型创业] ->  |  |
                    |  |   扮演分公司，自主选题...   |  |
                    |  +----------------------------+  |
                    +---------------------------------+
```

## 组件清单

### 顶层布局

| 组件 | 选择器/ID | 说明 |
|------|-----------|------|
| 应用导航栏 | `.appbar` | sticky top-0，56px，毛玻璃背景，品牌 Logo + 前台/后台切换 + 版本号 + 头像 |
| logo | `.appbar .logo` | "~ 量潮课堂 · LMS"，点击回到课程列表 |
| 视图切换组 | `#switch-btns` | 两个 button：学习前台 / LMS 管理后台，`.active` 高亮 |
| 元信息 | `.appbar .meta` | 版本号 `v0.12 · MD` + 头像圆 `Z` |
| 海浪背景 | `.bg-atmo` | 固定定位，径向发光 + SVG 波浪，纯装饰 |
| 工作区 | `#workspace` | flex 容器，`height: calc(100vh - 56px)` |

### 课程列表页

| 组件 | 选择器/ID | 说明 |
|------|-----------|------|
| 视图容器 | `#view-courses` | 默认 `.show`，flex 布局承载课程列表 |
| 页面容器 | `.courses-page` | `max-width: 900px`，居中，上下 60px 内边距 |
| 页面标题 | `.page-header` | 居中对齐，`h1` 标题 + `.desc` 描述文字 |
| 咨询入口 | `a[href]` | 蓝色文字链接 `💬 一对一咨询`，跳转到生产实习咨询表单 |
| 卡片网格 | `.card-grid` | 纵向排列，`max-width: 700px` 居中，`gap: 12px` |
| 课程卡片 | `.course-card` | 横向 flex，编号圆 + 课程信息 + 难度标签 + 入口箭头 |
| 编号圆 | `.cc-num` | 44px 圆，默认灰色，`.active-card` 下蓝色 |
| 课程信息 | `.cc-info` | flex:1，含 `.cc-name`（课程名）和 `.cc-desc`（描述） |
| 难度标签 | `.badge` | 胶囊形，4 种变体：`.beginner` / `.intermediate` / `.advanced` / `.capstone` |

### 卡片状态

| 状态 | 触发条件 | 视觉 |
|------|---------|------|
| 默认 | 所有课程卡片 | 白色卡片，灰色编号圆，`box-shadow: 0 1px 3px` |
| `.active-card` | 当前可在学课程（课程 1-4） | 编号圆变蓝底白字 |
| `:hover` | 鼠标悬停 | 边框变蓝，`translateX(4px)`，`box-shadow` 加深 |

## 交互与导航

| 交互 | 触发条件 | 行为 |
|------|---------|------|
| 点击课程 1-4 | 用户点击知识工作/氛围编程/大数据导论/数据工程 | `enterCourse(id)` → 切换至 `#view-generic`，调用 `renderGenericCourse(id)` 渲染课程内页 |
| 点击课程 5 | 用户点击生产实习 | `enterCourse('prod')` → 切换至 `#view-front`，调用 `showModule()` 恢复进度 |
| 点击"学习前台" | AppBar 按钮 | `showCourseList()` → 回到此页面 |
| 点击 Logo | AppBar 品牌区 | `showCourseList()` → 回到此页面 |
| 点击"LMS 管理后台" | AppBar 按钮 | `switchView('back')` → 进入后台 |
| 点击"一对一咨询" | 页面标题下方链接 | 跳转到生产实习的 `#mod-consult` 咨询表单 |

## 数据模型

### COURSES 对象结构

课程 1-4 的数据定义在全局 `COURSES` 对象中（详见实验室 [course.md](https://github.com/quanttide/quanttide-laboratory-of-course-development/blob/main/innovation-prototype/models/course.md)）：

```javascript
COURSES = {
  'knowledge-work': { name, icon, badge, badgeClass, next, nextName, desc, meta, stages[] },
  'vibe-coding':     { ... },  // next → 'big-data'
  'big-data':        { ... },  // next → 'data-engineering'
  'data-engineering': { ... },  // next → 'prod'
}
```

课程 5（生产实习）的数据硬编码在 HTML 中，不经过 `COURSES` 对象。

### 课程卡片渲染

课程卡片在 HTML 中静态编写（非 JS 动态渲染），每张卡片在 `.card-grid` 中直接硬编码，通过 `onclick="enterCourse(id)"` 绑定事件。

## 设计原则

1. **阶梯感** — 编号 ①-⑤ 暗示从基础到进阶的学习路径，每门课是下一门的前置
2. **信息密度适中** — 每张卡片展示课程名、一句话描述、难度标签，不做信息过载
3. **低门槛** — 五门课全部可点击进入，课程 1-4 有完整的阶段框架即使内容待建设
4. **进度感知** — 生产实习卡片蓝色高亮（`.active-card`），其余课程 hover 有反馈
5. **悬浮装饰** — 海浪 SVG + 径向发光背景，营造"海岸"品牌氛围但不干扰内容

## 状态机

```
[页面加载]
     |
     v
[课程列表] <-- showCourseList() -------------------------+
     |                                                    |
     |-- 点击课程 1-4 --> [view-generic · 通用课程内页]    |
     |                        |                            |
     |                        +-- 返回 --------------------+
     |                                                    |
     |-- 点击课程 5 -----> [view-front · 生产实习内页]     |
     |                        |                            |
     |                        +-- 返回 --------------------+
     |                                                    |
     +-- 点击"LMS管理后台" -> [view-back · 后台管理]       |
                                  |                        |
                                  +-- 点击"学习前台" ------+
```
