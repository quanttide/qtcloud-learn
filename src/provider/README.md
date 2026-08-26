# 学习云 Provider (`qtcloud-learn-provider`)

量潮学习云的服务端。Go 编写。

领域模型对齐《量潮学习管理标准》（`docs/specification`）：核心模型 **Learner × Criterion → Completion**。

## 开发

```bash
# 运行测试
go test ./...

# 启动服务（默认 :8080）
go run ./cmd/server
```

## 目录

```
cmd/server/          # 服务入口与路由（资源 CRUD 与 /healthz）
internal/domain/     # 领域模型（Learner / Criterion / Completion）
internal/store/      # 内存存储（BaseStore + 各实体 Store）
internal/handler/    # CRUD handler（泛型 CRUDHandler + 各实体 Handler）
internal/version/    # 版本信息
```

## API

API 无版本前缀，每个资源在根路径提供标准 CRUD：

| 资源 | 路径 |
|------|------|
| 学习者 | `/learners` |
| 验收标准 | `/criteria` |
| 完成记录 | `/completions` |

健康检查：`GET /healthz`。

## 相关文档

- [ROADMAP.md](ROADMAP.md) — 路线图
- [CHANGELOG.md](CHANGELOG.md) — 变更记录
