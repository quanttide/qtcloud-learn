# 量潮学习云 (`qtcloud-learn`)

AI 原生学员学习中心。

## 定位

承接学员侧学习能力，与 [量潮课程云](../qtcloud-course)（内容生产侧）、[量潮课堂](../qtclass)（客户端）构成闭环：

- 学员注册与认证（飞书登录）
- 选课 / 报名
- 学习进度追踪
- 考核提交与成绩
- 付费 / VIP 权益

## 仓库结构

```
src/
├── provider/  # 服务端（Go）
├── studio/    # 客户端（Flutter）
└── cli/       # 命令行工具（Rust）
```

## 开发

见 [docs/dev-guide/index.md](docs/dev-guide/index.md)。
