# 学习云 Provider (`qtcloud-learn-provider`)

量潮学习云的服务端。Go 编写。

## 开发

```bash
# 运行测试
go test ./...

# 启动服务（默认 :8080）
go run ./cmd/server
```

## 目录

```
cmd/server/          # 服务入口与路由（/api/v1/*，/healthz）
internal/domain/     # 统一领域模型（Student / Teacher / Class / Session / Enrollment / Progress / Assessment / Submission）
internal/store/      # 内存存储（BaseStore + 各实体 Store）
internal/handler/    # CRUD handler（泛型 CRUDHandler + 各实体 Handler）
internal/version/    # 版本信息
```

## API

LMS API 统一挂在 `/api/v1/` 前缀下，每个资源提供标准 CRUD：

| 资源 | 路径 |
|------|------|
| 班级 | `/api/v1/classes` |
| 学员 | `/api/v1/students` |
| 考核 | `/api/v1/assessments` |
| 提交 | `/api/v1/submissions` |
| 选课 | `/api/v1/enrollments` |
| 进度 | `/api/v1/progress` |

健康检查：`GET /healthz`。

## 相关文档

- [ROADMAP.md](ROADMAP.md) — 路线图
- [CHANGELOG.md](CHANGELOG.md) — 变更记录
