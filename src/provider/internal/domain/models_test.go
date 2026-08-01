package domain

import (
	"encoding/json"
	"strings"
	"testing"
)

func TestClass_JSON(t *testing.T) {
	c := Class{
		ID: "class-1", Name: "浙理班级", Slug: "slug-class-1", RefName: "大数据微专业",
		RefType: "program", RefID: "prog-1", StartDate: "2026-09-01", EndDate: "2027-01-15",
		StudentCount: 30, Progress: 0.5, TeacherIDs: []string{"tea-1"}, StudentIDs: []string{"stu-1"},
	}
	b, _ := json.Marshal(c)
	var got Class
	json.Unmarshal(b, &got)
	if got.Name != "浙理班级" || got.Slug != "slug-class-1" || got.StudentCount != 30 || got.Progress != 0.5 {
		t.Fatalf("roundtrip = %+v", got)
	}
	if len(got.TeacherIDs) != 1 || got.TeacherIDs[0] != "tea-1" {
		t.Fatalf("TeacherIDs roundtrip = %v", got.TeacherIDs)
	}
	if len(got.StudentIDs) != 1 || got.StudentIDs[0] != "stu-1" {
		t.Fatalf("StudentIDs roundtrip = %v", got.StudentIDs)
	}
}

func TestStudent_JSON(t *testing.T) {
	st := Student{ID: "stu-1", Name: "张三", Email: "zhangsan@example.com", Plan: "vip"}
	b, _ := json.Marshal(st)
	var got Student
	json.Unmarshal(b, &got)
	if got.Name != "张三" || got.Email != "zhangsan@example.com" || got.Plan != "vip" {
		t.Fatalf("roundtrip = %+v", got)
	}
}

func TestTeacher_JSON(t *testing.T) {
	th := Teacher{ID: "tea-1", Name: "李老师", Email: "li@example.com", Title: "副教授"}
	b, _ := json.Marshal(th)
	var got Teacher
	json.Unmarshal(b, &got)
	if got.Name != "李老师" || got.Title != "副教授" {
		t.Fatalf("roundtrip = %+v", got)
	}
}

func TestSession_JSON(t *testing.T) {
	s := Session{
		ID: "sess-1", ClassID: "class-1", LessonTitle: "Git 入门", TeacherID: "tea-1",
		StartTime: "2026-09-02T09:00:00Z", DurationMinutes: 45, Location: "A-101",
		Status: "upcoming", Attendances: []Attendance{{StudentID: "stu-1", Status: "present"}},
	}
	b, _ := json.Marshal(s)
	var got Session
	json.Unmarshal(b, &got)
	if got.ClassID != "class-1" || got.LessonTitle != "Git 入门" || len(got.Attendances) != 1 {
		t.Fatalf("roundtrip = %+v", got)
	}
	if got.Attendances[0].Status != "present" {
		t.Fatalf("attendance roundtrip = %+v", got.Attendances)
	}
}

func TestEnrollment_JSON(t *testing.T) {
	e := Enrollment{ID: "enr-1", ClassID: "class-1", StudentID: "stu-1", Status: "enrolled"}
	b, _ := json.Marshal(e)
	var got Enrollment
	json.Unmarshal(b, &got)
	if got.ClassID != "class-1" || got.StudentID != "stu-1" || got.Status != "enrolled" {
		t.Fatalf("roundtrip = %+v", got)
	}
}

func TestProgress_JSON(t *testing.T) {
	p := Progress{ID: "prog-1", StudentID: "stu-1", ClassID: "class-1", Percent: 0.6, Finished: false}
	b, _ := json.Marshal(p)
	var got Progress
	json.Unmarshal(b, &got)
	if got.Percent != 0.6 || got.Finished {
		t.Fatalf("roundtrip = %+v", got)
	}
}

func TestAssessment_JSON(t *testing.T) {
	a := Assessment{ID: "assess-1", ClassID: "class-1", Type: "homework", Title: "作业1", FullScore: 100, PassScore: 60, Deadline: "2026-10-01"}
	b, _ := json.Marshal(a)
	var got Assessment
	json.Unmarshal(b, &got)
	if got.Title != "作业1" || got.FullScore != 100 || got.Deadline != "2026-10-01" {
		t.Fatalf("roundtrip = %+v", got)
	}
}

func TestSubmission_JSON(t *testing.T) {
	sub := Submission{ID: "sub-1", AssessmentID: "assess-1", StudentID: "stu-1", Status: "submitted", Score: 88.5, SubmittedAt: "2026-09-30T10:00:00Z"}
	b, _ := json.Marshal(sub)
	var got Submission
	json.Unmarshal(b, &got)
	if got.Score != 88.5 || got.Status != "submitted" || got.AssessmentID != "assess-1" {
		t.Fatalf("roundtrip = %+v", got)
	}
}

func TestMakeSlug_ASCII(t *testing.T) {
	tests := []struct {
		name   string
		input  string
		prefix string
		want   string
	}{
		{"simple", "Hello World", "h-1", "hello-world"},
		{"mixed case", "Go Programming 101", "gp-1", "go-programming-101"},
		{"dash preserved", "user-guide-v2", "ug-2", "user-guide-v2"},
		{"underscore", "my_var_name", "mvn-1", "my-var-name"},
		{"special chars", "data@science#101", "ds-1", "datascience101"},
		{"trailing dash", "hello- ", "h-2", "hello"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := MakeSlug(tt.input, tt.prefix)
			if got != tt.want {
				t.Errorf("MakeSlug(%q, %q) = %q, want %q", tt.input, tt.prefix, got, tt.want)
			}
		})
	}
}

func TestMakeSlug_Chinese(t *testing.T) {
	got := MakeSlug("浙理班级", "class-1")
	if got != "slug-class-1" {
		t.Errorf("MakeSlug(Chinese) = %q, want %q", got, "slug-class-1")
	}
}

func TestMakeSlug_Empty(t *testing.T) {
	got := MakeSlug("", "e-1")
	if got != "slug-e-1" {
		t.Errorf("MakeSlug(empty) = %q, want %q", got, "slug-e-1")
	}
}

func TestClass_EmptySlices(t *testing.T) {
	c := Class{ID: "class-1", TeacherIDs: []string{}, StudentIDs: []string{}}
	b, _ := json.Marshal(c)
	if !strings.Contains(string(b), `"teacherIds":[]`) {
		t.Fatalf("empty teacherIds should serialize as [], got %s", string(b))
	}
	if !strings.Contains(string(b), `"studentIds":[]`) {
		t.Fatalf("empty studentIds should serialize as [], got %s", string(b))
	}
}
