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

// Application 立项申请（生产实习第五步·微型创业）。
// 学员填写提交，服务端存储、后台可查（v0.1 不做审批流，status 固定 submitted）。
// teamMode: "personal"（个人独立，MemberNames=[个人姓名]）/ "partner"（搭档，[队长, 队员]）。
type Application struct {
	ID          string   `json:"id"`
	ProjectName string   `json:"projectName"`
	BlindSpot   string   `json:"blindSpot"`
	DemoPlan    string   `json:"demoPlan"`
	Direction   string   `json:"direction"`
	TeamMode    string   `json:"teamMode"`
	MemberNames []string `json:"memberNames"`
	StudentID   string   `json:"studentId"`
	StudentName string   `json:"studentName"`
	Status      string   `json:"status"` // "submitted"（不做审批）
	CreatedAt   string   `json:"createdAt,omitempty"`
}
