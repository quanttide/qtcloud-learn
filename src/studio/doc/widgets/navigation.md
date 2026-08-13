# Navigation — 导航组件

---

## 1. AppBar (.appbar)
56px sticky 顶栏，毛玻璃背景。Logo + 前台/后台切换按钮 + 版本号 + Avatar。

## 2. SideNav (.sidenav)
220px 侧边栏，`.nav-group` > `.nav-label`（可折叠）+ `.nav-kids`（`.nav-item` 列表）。3 个分组：LMS管理、表单配置、关联系统。折叠动画 `max-height` transition。

## 3. StepBar (.steps / .stepbar)
横向步骤条，flex 居中，浅蓝背景卡片。`.bs`（回退链接）+ `.row`（`.step` 节点 + `.conn` 连接线）。两实例：`#stepbar`（生产实习 5 步）、`#gen-stepbar`（通用课程 4 步）。overview 时隐藏。

## 4. BackLink (.back-link)
"← 返回课程列表"，蓝色文字，hover 浅蓝背景，调用 `showCourseList()`。

## 5. Pipeline (.pipeline)
横向流程条，flex 均分，`.stage` 节点间 `::after →`。5 阶段：💡发现盲区 → 📝提交立项 → 🔍总部审批 → 🚀2周执行 → 🎤Sale给总部。
