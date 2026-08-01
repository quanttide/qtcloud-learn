# 开发指南

量潮学习云的开发指引。

## 仓库结构

```
src/
├── provider/  # 服务端（Go）
├── studio/    # 客户端（Flutter）
└── cli/       # 命令行工具（Rust）
```

## 入门

| 子项目 | 技术栈 | 命令 |
|--------|--------|------|
| provider | Go | `go test ./...` / `go run ./cmd/server` |
| studio | Flutter | `flutter test` / `flutter run -d linux` |
| cli | Rust | `cargo test` / `cargo run -- version` |

## 相关文档

- [ROADMAP.md](../../ROADMAP.md) — 产品路线图
- `src/*/ROADMAP.md` — 各 scope 详细路线图
