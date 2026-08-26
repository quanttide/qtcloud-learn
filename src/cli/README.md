# 学习云 CLI (`qtcloud-learn-cli`)

量潮学习云的 CLI 工具。Rust 编写。

领域模型对齐《量潮学习管理标准》（`docs/specification`）：核心模型 **Learner × Criterion → Completion**。

## 开发

```bash
# 构建
cargo build

# 运行（默认连接 http://localhost:8080 的 provider）
./target/debug/qtcloud-learn --base-url http://localhost:8080 learner list

# 测试
cargo test
```

## 子命令

| 子命令 | 操作 |
|--------|------|
| `learner` | create / list / get |
| `criterion` | create / list / get |
| `completion` | create / complete（标记完成）/ list / get |

连接的是 `qtcloud-learn-provider`（`/api/v1`）。

## 相关文档

- [ROADMAP.md](ROADMAP.md) — 路线图
- [CHANGELOG.md](CHANGELOG.md) — 变更记录
