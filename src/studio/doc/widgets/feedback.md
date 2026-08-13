# Feedback — 反馈组件

---

## 1. Toast (.toast)
轻提示：`position:fixed;bottom:28px`，黑底白字 12px，圆角 999px。`.visible` 时 `opacity:1`，1.8s 自动隐藏。调用 `showToast(message)`。

## 2. ProgressBar (.progress-bar)
进度条：5px 高浅蓝轨道 + `.fill` 蓝色填充。宽度动态更新（生产实习 20%/40%/60%/80%/100%，通用课程按比例）。位于 `.course-hero` 内。

## 3. Overlay（待实现）
覆盖层：参考 qtclass 播放器的互动/完成/历史/确认覆盖层。LMS 原型可用于提交确认、审批弹窗、删除确认。

## 4. Dialog（待实现）
对话框：当前用 `alert()` 临时替代。后续需要 Modal 组件支持审批通过/驳回、保存草稿提示。
