// Package domain 定义量潮学习云的统一领域模型。
//
// 模型对齐 `qtcloud-course` / `qtclass` 两仓库的 LMS 模型（字段映射见
// qtcloud-learn/docs/lms-inventory.md），以本包为唯一事实来源。
package domain

// Student 学员。
// 对齐：qtcloud-course student.dart + qtclass session.dart Student。
type Student struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	Email  string `json:"email"`
	Avatar string `json:"avatar,omitempty"`
	Plan   string `json:"plan,omitempty"` // "free" / "paid" / "vip"
}

// Teacher 讲师。
// 对齐：qtcloud-course teacher.dart（含 title）+ qtclass session.dart Teacher。
type Teacher struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
	Title string `json:"title,omitempty"`
}

// Class 班级。
// 对齐：qtcloud-course Class（provider class.go）与 ClassTeaching（studio class_teaching.dart），
// 合并 TeacherIDs / StudentIDs。
type Class struct {
	ID           string   `json:"id"`
	Name         string   `json:"name"`
	Slug         string   `json:"slug"`
	RefName      string   `json:"refName"`            // 引用的专业/课程名称（展示用）
	RefType      string   `json:"refType,omitempty"`  // 引用类型："program" / "course"
	RefID        string   `json:"refId"`              // 引用的 Program/Course ID
	Status       string   `json:"status,omitempty"`   // "preparing" / "active" / "ended"
	StartDate    string   `json:"startDate"`          // ISO 日期
	EndDate      string   `json:"endDate"`            // ISO 日期
	StudentCount int      `json:"studentCount,omitempty"`
	Progress     float64  `json:"progress,omitempty"` // 教学进度（0.0 ~ 1.0）
	TeacherIDs   []string `json:"teacherIds"`
	StudentIDs   []string `json:"studentIds"`
}

// Session 课次（含考勤）。
// 对齐：qtclass session.dart，className/courseName 由 ClassID 关联 Class 派生。
type Session struct {
	ID              string       `json:"id"`
	ClassID         string       `json:"classId"`
	LessonTitle     string       `json:"lessonTitle"`
	TeacherID       string       `json:"teacherId"`
	StartTime       string       `json:"startTime"`       // ISO 时间
	DurationMinutes int          `json:"durationMinutes"`
	Location        string       `json:"location"`
	Status          string       `json:"status"`          // "upcoming" / "inProgress" / "completed"
	Attendances     []Attendance `json:"attendances"`
}

// Attendance 考勤记录。
// 对齐：qtclass session.dart AttendanceStatus（unknown / present / late / absent）。
type Attendance struct {
	StudentID string `json:"studentId"`
	Status    string `json:"status"`
}

// Enrollment 选课 / 报名。
type Enrollment struct {
	ID         string `json:"id"`
	ClassID    string `json:"classId"`
	StudentID  string `json:"studentId"`
	Status     string `json:"status,omitempty"` // "enrolled" / "withdrawn"
	EnrolledAt string `json:"enrolledAt,omitempty"`
}

// Progress 学习进度。
// 对齐：qtcloud-course Class.progress + qtclass learning_record.dart（转为服务端进度数据）。
type Progress struct {
	ID        string  `json:"id"`
	StudentID string  `json:"studentId"`
	ClassID   string  `json:"classId"`
	Percent   float64 `json:"percent"` // 学习进度（0.0 ~ 1.0）
	Finished  bool    `json:"finished"`
	UpdatedAt string  `json:"updatedAt,omitempty"`
}

// Assessment 考核。
// 对齐：qtcloud-course assessment.dart。
type Assessment struct {
	ID        string `json:"id"`
	ClassID   string `json:"classId"`
	Type      string `json:"type"` // "homework" / "exam"
	Title     string `json:"title"`
	FullScore int    `json:"fullScore"`
	PassScore int    `json:"passScore"`
	Deadline  string `json:"deadline"` // ISO 日期
}

// Submission 考核提交。
// 对齐：qtcloud-course submission.dart。
type Submission struct {
	ID           string  `json:"id"`
	AssessmentID string  `json:"assessmentId"`
	StudentID    string  `json:"studentId"`
	Status       string  `json:"status"` // "submitted" / "late" / "resubmitted"
	Score        float64 `json:"score,omitempty"`
	Comment      string  `json:"comment,omitempty"`
	SubmittedAt  string  `json:"submittedAt"`
}

// Application 立项申请（生产实习第五步·微型创业 Sell Your Demo）。
// 字段对齐原型 qt-proposals：5 问（机会/匹配/假设/Demo）+ 方向类型 + 组队姓名。
// teamMode: "personal"（个人独立，TeamLeader=个人姓名）/ "partner"（搭档，队长+队员）。
// v0.1 不做审批流，status 固定 "已提交"；删除为软删除（DeletedAt 非空进历史）。
type Application struct {
	ID            string `json:"id"`
	ProjectName   string `json:"projectName"`
	Opportunity   string `json:"opportunity"`   // 发现的机会（谁遇到什么问题、现有方案哪里不够好）
	Fit           string `json:"fit"`           // 为什么适合量潮（与已有业务/能力/资产的关系）
	Hypothesis    string `json:"hypothesis"`    // 核心假设（两周后最想证明什么）
	Demo          string `json:"demo"`          // 准备 Sale 什么 Demo
	DirectionType string `json:"directionType"` // 方向类型（内容/数据/渠道/方法…）
	TeamMode      string `json:"teamMode"`      // "personal" / "partner"
	TeamLeader    string `json:"teamLeader"`    // 队长姓名（个人独立时=本人姓名）
	TeamMember    string `json:"teamMember"`    // 队员姓名（多个顿号分隔；个人独立时空）
	StudentName   string `json:"studentName"`   // 当前学员身份（组队时=队长姓名）
	Status        string `json:"status"`        // "已提交"（不做审批）
	SubmittedAt   string `json:"submittedAt,omitempty"`
	DeletedAt     string `json:"deletedAt,omitempty"` // 软删除时间（非空=已删除）
}

// Learner 学员档案（原型 qt-students）：上报进度/提交立项时按姓名自动建档。
type Learner struct {
	ID           string `json:"id"`
	Name         string `json:"name"`
	Course       string `json:"course"`                 // 当前课程："生产实习"
	ProgressMax  int    `json:"progressMax"`            // 已到模块数 X
	ProgressTotal int   `json:"progressTotal"`          // 总模块数（5）
	ActiveAt     string `json:"activeAt,omitempty"`     // 最近活跃
	Status       string `json:"status"`                 // "在读" / "已完成"
	ProjectName  string `json:"projectName,omitempty"`  // 最近立项项目名（有立项时显示 ✓）
}
