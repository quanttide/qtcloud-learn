# 学习管理系统（LMS）代码盘点

> 迁移来源清单：从 `qtclass` 与 `qtcloud-course` 收拢到 `qtcloud-learn`。
> 本清单是 [ROADMAP.md](../ROADMAP.md) 的执行依据；条目完成后勾选。

## qtcloud-course（量潮课程云）

### provider（Go）

| 文件 | 内容 | 迁移目标 |
|------|------|---------|
| `internal/domain/class.go` | Class 领域模型 | `qtcloud-learn/provider/internal/domain` |
| `internal/store/class.go` | ClassStore | `qtcloud-learn/provider/internal/store` |
| `internal/handler/class.go` | ClassHandler CRUD | `qtcloud-learn/provider/internal/handler` |
| `internal/{domain,store,handler}/*_test.go` 中 class 相关 | 测试 | 随代码迁移 |
| `cmd/server/main.go` 中 class 路由 | 路由注册 | 随代码迁移 |

### studio（Flutter）

| 文件 | 内容 | 迁移目标 |
|------|------|---------|
| `lib/models/class_teaching.dart` | 班级教学模型 | `qtcloud-learn/studio/lib/models` |
| `lib/models/student.dart` / `teacher.dart` | 学员 / 讲师 | 同上 |
| `lib/models/assessment.dart` / `submission.dart` | 考核 / 提交 | 同上 |
| `lib/models/enums.dart`（LMS 部分） | 状态枚举 | 同上 |
| `lib/screens/class_screen.dart` | 班级管理页 | `qtcloud-learn/studio/lib/screens` |
| `lib/screens/assessment_list_screen.dart` / `assessment_detail_screen.dart` | 考核页 | 同上 |
| `lib/screens/dashboard_screen.dart` | 仪表盘 | 同上 |
| `lib/services/assessment_service.dart` | 考核服务 | `qtcloud-learn/studio/lib/services` |
| `lib/services/data_service.dart`（LMS 部分） | 数据服务 | 同上 |
| `assets/classes.json` / `students.json` / `teachers.json` / `assessments.json` / `submissions.json` | 演示数据 | 迁移 / 替换 |

### cli（Rust）

`ROADMAP.md` 中规划的 class / student 子命令（尚未实现）→ 在 `qtcloud-learn/cli` 实现。

## qtclass（量潮课堂）

### studio（Flutter）

| 文件 | 内容 | 迁移目标 |
|------|------|---------|
| `lib/models/session.dart` | Session / Teacher / Student / Attendance（课表 + 考勤） | `qtcloud-learn/studio/lib/models` |
| `lib/models/learning_record.dart` | 学习记录 | 同上（转为服务端进度数据） |
| `lib/screens/schedule_screen.dart` | 课表页 | `qtcloud-learn/studio/lib/screens` |
| `lib/screens/classroom_screen.dart` | 课堂 / 考勤页 | 同上 |
| `lib/screens/result_screen.dart` | 学习结果页 | 同上 |
| `lib/services/app_state.dart`（sessions 加载部分） | 状态加载 | `qtcloud-learn/studio/lib/services` |
| `lib/services/history_service.dart` | 学习记录 localStorage | 改为服务端进度 |
| `assets/sessions.json` | 课表演示数据 | 迁移 |

> 注：`lectures.json` 与 player / lecture 相关代码属播放器职责，**保留**在 `qtclass`。

## 统一领域模型（qtcloud-learn）

| 模型 | 来源 |
|------|------|
| Student（学员，含付费 / VIP / 免费） | `qtcloud-course` student.dart + `qtclass` session.dart Student |
| Teacher（讲师） | 两仓库均有 |
| Class（班级） | `qtcloud-course` Class / ClassTeaching |
| Session（课次，含考勤） | `qtclass` session.dart |
| Enrollment（选课 / 报名） | 新增 |
| Progress（学习进度） | `qtcloud-course` Class.progress + `qtclass` learning_record |
| Assessment / Submission（考核 / 提交） | `qtcloud-course` |
| LearningRecord（学习记录） | `qtclass` learning_record.dart（转为服务端进度数据） |

## 迁移状态

| 阶段 | 状态 | 说明 |
|------|------|------|
| v0.2 Provider（Go） | ✅ 已完成 | 领域模型 + 8 资源 CRUD + `/api/v1` 路由已落地；`class.go`（domain / store / handler）已移植 |
| v0.3 Studio（Flutter） | ✅ 已完成 | 登录（本地模拟）/ 我的课程 / 进度 / 课表考勤 / 班级 / 考核页面与服务已迁移 |
| v0.4 CLI（Rust） | ✅ 已完成 | student / class / enrollment / progress / assessment 子命令已实现 |
| v0.5 旧代码移除 | ✅ 已完成 | `qtcloud-course` 移除 class.go 与 LMS 页面/服务/assets；`qtclass` 移除 session/学习记录/课表考勤/result（播放器保留） |

## 统一领域模型字段映射（v0.1 三仓库对齐核对）

| qtcloud-learn 字段 | 来源 | 说明 |
|------|------|------|
| Student.ID / Name / Email / Avatar | `qtcloud-course` student.dart | 学员基本信息 |
| Student.Plan（free / paid / vip） | 新增 | 付费 / VIP 权益 |
| Student.attendance | `qtclass` session.dart Student | 不随学员存储，归入 Session.Attendances |
| Teacher.ID / Name / Email / Title | `qtcloud-course` teacher.dart（title）+ `qtclass` session.dart Teacher | 讲师 |
| Class.ID / Name / Slug / RefName / RefType / RefID / Status / StartDate / EndDate / StudentCount / Progress | `qtcloud-course` class.go + class_teaching.dart | 班级 |
| Class.TeacherIDs / StudentIDs | `qtcloud-course` class_teaching.dart | 成员引用 |
| Session.ClassID / LessonTitle / TeacherID / StartTime / DurationMinutes / Location / Status | `qtclass` session.dart | 课次；className / courseName 由 ClassID 关联派生 |
| Session.Attendances（StudentID / Status） | `qtclass` session.dart AttendanceStatus | 考勤 |
| Enrollment.ClassID / StudentID / Status / EnrolledAt | 新增 | 选课 / 报名 |
| Progress.StudentID / ClassID / Percent / Finished / UpdatedAt | `qtcloud-course` Class.progress + `qtclass` learning_record.finished | 学习记录转为服务端进度数据 |
| Assessment.ClassID / Type / Title / FullScore / PassScore / Deadline | `qtcloud-course` assessment.dart | 考核 |
| Submission.AssessmentID / StudentID / Status / Score / Comment / SubmittedAt | `qtcloud-course` submission.dart | 提交 |

> 说明：`qtclass` learning_record.dart 的 env / runState 为播放器环境细节，不纳入统一进度模型。
