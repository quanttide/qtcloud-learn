# ROADMAP

> 本 scope 规划对齐根 [ROADMAP.md](../../ROADMAP.md) 的能力演进阶段。

## [v1.0] — 直连 provider（对齐根 v1.0）

> **目标**：子命令真实对接 provider，端到端可用。

### 交付物

- [ ] 子命令冒烟脚本（provider 启动 → student / class / enrollment / progress / assessment 全链路）
- [ ] 错误处理：连接失败 / 404 / 400 可读提示

### 测试

- [ ] 冒烟脚本 CI 化
- [ ] 错误路径测试

## [v1.1] — 认证与权益（对齐根 v1.1 / v1.2）

### 交付物

- [ ] 子命令支持 token 参数（登录态）
- [ ] plan 字段展示与权益查询

### 测试

- [ ] token 参数测试

## [v2.0] — 报表（对齐根 v2.0）

### 交付物

- [ ] assessment stats 报表（成绩分布 / 及格率）
- [ ] 考勤统计报表

### 测试

- [ ] 报表输出测试
