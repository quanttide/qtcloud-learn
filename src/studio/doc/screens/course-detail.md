# 课程详情页 (view-front / view-generic)

课程详情页，包含**课程首页（Hero）**和**课时列表**两部分。从课程列表点击进入后先看 Hero，点击"继续学习"进入课时列表。生产实习（课程 5）为硬编码 HTML，课程 1-4 为 JS 动态渲染的通用模板。

## 页面定位

- **路由**: 从课程列表点击进入，无独立 URL（SPA 内部视图切换）
- **状态**:
  - `view-front`：生产实习专用，5 个模块硬编码 HTML，`currentCourse = 'prod'`
  - `view-generic`：课程 1-4 共用模板，JS 动态渲染，`genCourseId` 存储当前课程 ID
- **数据依赖**:
  - 生产实习：硬编码 HTML（模块 1-5 + 咨询表单），进度读取 `localStorage('qt-progress-prod')`
  - 课程 1-4：`COURSES` JS 对象 + `localStorage('qt-progress-{courseId}')`

## 页面结构总览

```
+-- AppBar (sticky, 56px) ----------------------------------------------+
|  [~ 量潮课堂 · LMS]          [🏠 学习前台] [⚙ LMS 管理后台]   v0.12  Z |
+-----------------------------------------------------------------------+

+-- .main (overflow-y:auto) --------------------------------------------+
|  +-- .main-inner (max-width:860px, 居中) ---------------------------+ |
|  |                                                                   | |
|  |  <- 返回课程列表  (.back-link)                                    | |
|  |                                                                   | |
|  |  +-- StepBar (.steps / .stepbar) ------------------------------+  | |
|  |  | <- 课程首页  --  ①量潮是谁 -- ②业务与市场 -- ③方法论 -- ...|  | |
|  |  |              ✓已完成      ●当前        ○未完成              |  | |
|  |  +-------------------------------------------------------------+  | |
|  |                                                                   | |
|  |  +-- 课程首页 (.course-hero) -----------------------------------+ | |
|  |  |  🏭 生产实习 · 微型创业                                       | | |
|  |  |  量潮课堂 · 实训基地                                          | | |
|  |  |  描述文字：量潮是"总部"，你是"分公司"...                       | | |
|  |  |  📚 5 个模块 · ⏱ 预计 2 周 · 👤 38 人在学 · 🎯 同一套卷子    | | |
|  |  |  ▓▓░░░░░░░░  20%                                              | | |
|  |  |  [▶ 继续学习]  [👥 组队广场]                                   | | |
|  |  +---------------------------------------------------------------+ | |
|  |                                                                   | |
|  |  +-- 模块内容 (.module-panel.show) -----------------------------+ | |
|  |  |  📖 一、量潮是谁                                              | | |
|  |  |  了解你所在的组织——它的起源、结构与核心业务                    | | |
|  |  |  · 1.1 量潮的创立故事                      阅读 10 min       | | |
|  |  |  · 1.2 组织架构与团队分工                   阅读 8 min        | | |
|  |  |  · 1.3 量潮云 + 量潮课堂：两大产品线         视频 15 min       | | |
|  |  |  · 1.4 我们服务谁：浙理工合作案例             阅读 10 min       | | |
|  |  |  [-> 下一模块]  [<- 返回课程首页]                              | | |
|  |  +---------------------------------------------------------------+ | |
|  +-------------------------------------------------------------------+ |
+-----------------------------------------------------------------------+
```

## 组件清单

### 顶层布局

| 组件 | 选择器/ID | 说明 |
|------|-----------|------|
| 视图容器 | `#view-front` / `#view-generic` | `display:none` 默认隐藏，`.show` 显示 |
| 主滚动区 | `.main` | `flex:1; overflow-y:auto`，内边距 `0 32px 60px` |
| 内容约束 | `.main-inner` | `max-width: 860px; margin: 0 auto` |

### 导航组件

| 组件 | 选择器/ID | 说明 |
|------|-----------|------|
| 返回链接 | `.back-link` | "← 返回课程列表"，调用 `showCourseList()`，hover 浅蓝背景 |
| 步骤条容器 | `.steps` / `.stepbar` | flex 居中，浅蓝背景卡片，圆角 14px，`margin-bottom: 20px` |
| 回退链接 | `.steps .bs` | "← 课程首页"，灰色文字，hover 变蓝，调用 `showModule('overview')` 或 `genShowModule('overview')` |
| 步骤行 | `.steps .row` / `.stepbar .row` | flex 横向排列步骤节点 |

### 步骤条 (StepBar)

| 子组件 | 选择器 | 说明 |
|--------|--------|------|
| 步骤节点 | `.step` | 可点击，含 `.dot`（圆点 36px）和 `.lbl`（标签 13px 粗体） |
| 连接线 | `.conn` | 24px 宽 2px 高横线，`.ok` 变绿 |
| 完成态 | `.dot.ok` / `.lbl.ok` | 绿底白勾 + 绿色标签，`box-shadow: 0 0 0 4px rgba(34,197,94,.08)` |
| 当前态 | `.dot.on` / `.lbl.on` | 蓝底白字 + 蓝色标签，`box-shadow: 0 0 0 6px rgba(22,119,255,.10)` |
| 默认态 | `.dot` / `.lbl` | 灰底灰字，未激活 |

**两个实例：**
- `#stepbar` — 生产实习专用，5 个步骤（m1-m5）
- `#gen-stepbar` — 课程 1-4 通用，步骤数由 `COURSES[courseId].stages.length` 决定（4 个）

### 课程 Hero

| 子组件 | 选择器/ID | 说明 |
|--------|-----------|------|
| Hero 容器 | `.course-hero` | 蓝色渐变背景 `linear-gradient(135deg, #E3F2FD, #BBDEFB, #E8F5FF)`，圆角 18px，40px 48px 内边距 |
| 课程标签 | `.badge-course` | 半透明蓝色胶囊，显示 "🏭 生产实习 · 微型创业" 或 `{icon} {name} · {badge}` |
| 标题 | `h2` / `#gen-title` | 28px 粗体 |
| 描述 | `.hero-desc` / `#gen-desc` | 14px 二级文字，`max-width: 680px` |
| 元信息 | `.hero-meta` / `#gen-meta` | flex 横向排列，12px 三级文字 |
| 进度条 | `.progress-bar` | 5px 高浅蓝底，`.fill` 蓝色填充，`max-width: 400px` |
| CTA 按钮 | `.btn-filled.btn-lg` | "▶ 继续学习" 主按钮 + "👥 组队广场" 次按钮 |

### 模块面板 (ModulePanel)

| 子组件 | 选择器/ID | 说明 |
|--------|-----------|------|
| 面板容器 | `.module-panel` | `display:none` 默认隐藏，`.show` 变为 `display:block`，同时最多一个显示 |
| 课程首页面板 | `#mod-overview` / `#gen-mod-overview` | Hero 卡片，步骤条隐藏 |
| 模块内容面板 | `#mod-m1`~`#mod-m5` / `#gen-mod-s1`~`#gen-mod-s4` | 卡片 `.card` + 课时条目 `.lesson-item` + 操作按钮 `.acts` |
| 咨询面板 | `#mod-consult` | 特殊面板，步骤条隐藏，含咨询表单 |
| 卡片 | `.card` | 浅蓝灰背景 `#f5f9fc`，圆角 14px，36px 40px 内边距，hover 阴影加深 |
| 卡片标题 | `.card h2` | 20px 粗体 |
| 卡片副标题 | `.card .subtitle` | 13px 灰色，`margin-bottom: 28px` |

### 课时条目 (Lesson Item)

| 子组件 | 选择器 | 说明 |
|--------|--------|------|
| 条目行 | `.lesson-item` | flex 横向，14px 16px 内边距，底部分割线，hover 浅蓝背景 + 右移 4px |
| 圆点 | `.lesson-item .ldot` | 7px 蓝色半透明圆点 |
| 时长标签 | `.lesson-item .duration` | 右对齐，12px 灰色，显示 "阅读 10 min" / "视频 15 min" 等 |
| 完成态圆点 | `.ldot.done` | 绿色实心（当前全部为默认态，内容待建设） |

### 操作按钮组

| 子组件 | 选择器 | 说明 |
|--------|--------|------|
| 按钮容器 | `.acts` | flex 横向，`gap: 16px`，`margin-top: 30px` |
| 主按钮 | `.btn.btn-filled` | 蓝色渐变，白色文字，"→ 下一模块" 或 "→ 进入下一课程" |
| 次按钮 | `.btn.btn-text` | 蓝色文字无背景，"← 返回课程首页" |

### 生产实习模块详情

| 模块 | ID | 标题 | 课时数 |
|------|-----|------|--------|
| 模块 1 | `#mod-m1` | 📖 一、量潮是谁 | 4 课时（创立故事/组织架构/两大产品线/浙理工案例） |
| 模块 2 | `#mod-m2` | 🏢 二、业务与市场 | 4 课时（产教融合/创新成果收集/用户画像/规模化路径） |
| 模块 3 | `#mod-m3` | 🔧 三、工作方法论 | 4 课时（快速开发/知识工程/Web Coding/工作日志） |
| 模块 4 | `#mod-m4` | 💡 四、最新探索方向 | 4 课时（Build in Public/多业务平衡/核心竞争力/制度创新） |
| 模块 5 | `#mod-m5` | 📝 五、微型创业 | 立项表单（项目名称/盲区描述/Demo 方案/方向类型/组队方式） |

### 课程 1-4 阶段详情 (COURSES 对象)

| 课程 ID | 课程名 | 阶段 1 | 阶段 2 | 阶段 3 | 阶段 4 |
|---------|--------|--------|--------|--------|--------|
| `knowledge-work` | 知识工作 | 信息检索基础 | 文档整理规范 | 知识管理工具 | AI 辅助写作 |
| `vibe-coding` | 氛围编程 | 安装开发环境 | 第一个页面 | AI 辅助编程 | 独立小作品 |
| `big-data` | 大数据导论 | 数据意识培养 | 基本分析方法 | 可视化入门 | 数据驱动决策 |
| `data-engineering` | 数据工程 | 数据管道概念 | Python 数据处理 | 工程化实践 | 综合项目 |

每个阶段 3 个课时，共 12 课时/课程。

## 交互与导航

| 交互 | 触发条件 | 行为 |
|------|---------|------|
| 点击步骤节点 | 用户点击步骤条上的 `.step` | `showModule(id)` / `genShowModule(id)` — 切换到对应模块面板 |
| 点击"继续学习" | 课程首页 Hero 中 | `continueProd()` / `continueGen()` — 有进度跳到最后完成的下一步，无进度跳到模块 2/阶段 1 |
| 点击"→ 下一模块" | 模块内容底部 `.btn-filled` | 切换到下一个模块面板，同时 `updateStepBar()` + `saveProgress()` |
| 点击"← 返回课程首页" | 模块底部 `.btn-text` | `showModule('overview')` — 回到 Hero 首页，隐藏步骤条 |
| 点击"← 返回课程列表" | 顶部 `.back-link` | `showCourseList()` — 回到 `#view-courses` |
| 点击"→ 进入下一课程" | 最后模块底部（课程 1-4） | `goNextCourse(id)` — 调用 `renderGenericCourse(nextId)` 跳到下一门课首页 |
| 点击"👥 组队广场" | 生产实习 Hero | 预留按钮，当前无绑定事件 |
| 点击"✓ 提交申报" | 生产实习模块 5 底部 | 提交立项表单（当前为 Demo，无实际提交逻辑） |
| 点击"保存草稿" | 模块 5 / 咨询表单 | 预留按钮（当前无实际持久化） |
| 点击"💬 一对一咨询" | 课程列表页链接 | 跳转到生产实习 `#mod-consult` 面板 |
| 课时条目点击 | `.lesson-item` 的 `onclick` | 弹出 alert 提示内容待上线（部分课时有详细说明） |

## 进度持久化

每次进入模块时调用 `saveProgress(courseId, moduleId)`，将当前模块 ID 和已完成的最大模块序号存入 `localStorage`。

| 课程 | 存储键 | 模块 ID 格式 | 示例值 |
|------|--------|-------------|--------|
| 生产实习 | `qt-progress-prod` | `m1`-`m5` | `{"max":3,"last":"m3"}` |
| 知识工作 | `qt-progress-knowledge-work` | `s1`-`s4` | `{"max":2,"last":"s2"}` |
| 氛围编程 | `qt-progress-vibe-coding` | `s1`-`s4` | `{"max":1,"last":"s1"}` |
| 大数据导论 | `qt-progress-big-data` | `s1`-`s4` | `{"max":0}` |
| 数据工程 | `qt-progress-data-engineering` | `s1`-`s4` | `{"max":4,"last":"s4"}` |

**关键规则**：
- `last` 为 `undefined` 时用户首次访问 → 显示课程首页（Hero）
- `max` 只增不减 — 用户回到之前的模块不会降低 max
- `saveProgress()` 在 `showModule()` / `genShowModule()` 中自动调用
- 每门课独立存储，互不干扰

## 数据模型

详见实验室 [course.md](https://github.com/quanttide/quanttide-laboratory-of-course-development/blob/main/innovation-prototype/models/course.md)，核心结构：

```javascript
// Course
{ name, icon, badge, badgeClass, next, nextName, desc, meta,
  stages: Stage[] }

// Stage
{ name, desc, lessons: Lesson[] }

// Lesson
{ t: string,    // 课时标题
  d: string,    // 资源类型 + 时长
  done: boolean // 完成状态（当前全部 false） }

// Progress (localStorage)
{ max: number,   // 已完成的最大模块/阶段序号
  last: string } // 最后访问的模块 ID
```

## 设计原则

1. **一屏一模块** — 课程首页 Hero 控制在一屏内，模块内容不溢出，每个面板独立切换
2. **进度可视化** — 步骤条（圆点颜色）+ 进度条（蓝色填充）双重反馈当前位置
3. **操作连贯** — "下一模块"按钮贯穿所有步骤，最后一步自动变为"→ 进入下一课程"
4. **状态恢复** — 关闭浏览器后重开，自动回到上次学习的模块（通过 `continueProd()` / `continueGen()`）
5. **双轨渲染** — 生产实习硬编码（内容定制化强），课程 1-4 共用模板（`renderGenericCourse` 动态生成），两种模式共享同一套 CSS 组件体系

## 状态机

```
[课程列表] --点击课程 1-4--> [view-generic]
                                |
                    renderGenericCourse(courseId)
                                |
                                v
                         [课程首页 · overview]
                           |              |
                  "继续学习"              "开始学习"
                           |              |
                           v              v
                    [阶段 N+1]        [阶段 1]
                           |
                    点击"下一阶段"
                           |
                           v
                    [阶段 N+1] --> ... --> [最后阶段]
                                              |
                                    点击"-> 进入下一课程"
                                              |
                                              v
                                    [下一门课首页]


[课程列表] --点击课程 5--> [view-front]
                                |
                    showModule(progress.last || 'overview')
                                |
                      +---------+---------+
                      |                   |
                      v                   v
               [课程首页]            [模块 N]
                 |                      |
            "继续学习"           点击"下一模块"
                 |                      |
                 v                      v
            [模块 N+1] --> ... --> [模块 5 · 微型创业]
                                        |
                                 "✓ 提交申报"
                                        |
                                        v
                                   [完成]
```
