# Atoms — 原子组件

最小不可拆的视觉单元，构成所有复杂组件的基础。

---

## 1. Button — 按钮

| 选择器 | 说明 |
|--------|------|
| `.btn` | 基础：`padding: 10px 28px`，圆角 8px，`font-weight: 600` |
| `.btn-filled` | 主操作：蓝色渐变，白色文字，`box-shadow: 0 2px 8px` |
| `.btn-outlined` | 次要：透明背景，1.5px 边框，蓝色文字 |
| `.btn-text` | 文字按钮：无背景无边框，hover 浅蓝背景 |
| `.btn-sm` | 小尺寸：`padding: 6px 14px; font-size: 12px` |
| `.btn-lg` | 大尺寸：`padding: 14px 32px; font-size: 16px` |

**状态**：hover 上浮 1px、focus-visible 蓝色 outline

---

## 2. Badge — 标签

| 选择器 | 背景 | 文字色 | 用途 |
|--------|------|--------|------|
| `.badge.beginner` | `#e8f4ff` | `#1677ff` | 入门 |
| `.badge.intermediate` | `#ede9fe` | `#5b21b6` | 进阶 |
| `.badge.advanced` | `#fff3e0` | `#e65100` | 实战 |
| `.badge.capstone` | `#ffeaea` | `#c62828` | 微型创业 |
| `.badge.review` | `#e8f4ff` | `#1677ff` | 待审批 |
| `.badge.active-status` | `#f0fdf4` | `#22c55e` | 进行中 |
| `.badge.done` | `#f5f0ff` | `#6a1b9a` | 已完成 |
| `.badge.draft` | `#fff3e0` | `#e65100` | 草稿 |

---

## 3. Dot — 圆点

| 选择器 | 尺寸 | 说明 |
|--------|------|------|
| `.lesson-item .ldot` | 7px | 课时条目前置圆点，蓝色半透明 |
| `.step .dot` | 36px | 步骤条圆点：默认灰 / `.ok` 绿 / `.on` 蓝+发光 |
| `.step .dot.ok` | 36px | 绿色背景 + 白色✓ + `box-shadow: 0 0 0 4px` |
| `.step .dot.on` | 36px | 蓝色背景 + 白色数字 + `box-shadow: 0 0 0 6px` |

---

## 4. Icon — 图标

当前统一使用 emoji，未引入图标库。
常用：`📋🔍🎓📁👥📖👤📝📋✅📚💬📖🏢🔧💡📝🏭💻📊⚙️`

---

## 5. Avatar — 头像

| 选择器 | 说明 |
|--------|------|
| `.appbar .avatar` | 32px 圆，浅蓝背景 `#e8f4ff` + 蓝色文字，`font-weight: 600` |

---

## 6. Connector — 连接线

| 选择器 | 说明 |
|--------|------|
| `.step .conn` | 步骤条节点间 24px×2px 横线，默认灰，`.ok` 变绿 `#bbf7d0` |

---

## 7. Divider — 分割线

| 选择器 | 说明 |
|--------|------|
| `hr.divider` | `border-top: 1px solid var(--line-light)`，`margin: 20px 0` |
