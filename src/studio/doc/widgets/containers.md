# Containers — 容器组件

承载和排列 items 的布局容器。

---

## 1. Card (.card)
白色卡片：`background:#f5f9fc`，圆角 14px，`padding:36px 40px`，1px 边框，hover 阴影加深。含表格时 `padding:0`。

## 2. Section (.section)
内容分区：`margin-bottom:24px`。子元素 `.section-hd`（h3+说明文字）+ `.section-body`。

## 3. TwoCol (.twocol)
双栏：`display:flex;gap:24px`，子元素 `flex:1`。变体 `.twocol .side`（`flex:0 0 280px`）。

## 4. CourseHero (.course-hero)
课程首页 Hero：蓝色渐变背景，圆角 18px，`padding:40px 48px`。子元素：`.badge-course`、`h2`、`.hero-desc`、`.hero-meta`、`.progress-bar`、CTA 按钮。

## 5. ModulePanel (.module-panel)
课时面板：默认 `display:none`，`.show` 显示。面板 ID：`#mod-overview`（首页）、`#mod-m1`~`#mod-m5`（生产实习）、`#gen-mod-s1`~`#gen-mod-s4`（通用课程）、`#mod-consult`（咨询）。

## 6. PageHeader (.page-header)
页面标题：`text-align:center`，`margin-bottom:56px`。`h1`（30px）+ `.desc`（16px 灰色，`max-width:560px`）。

## 7. Workspace (#workspace)
全局工作区：`display:flex`，`height:calc(100vh - 56px)`。前台全宽，后台侧边栏 220px + 主区 flex:1。

## 8. FormRow (.form-row)
表单行：`display:flex;gap:16px`，子元素 `flex:1` 均分。
