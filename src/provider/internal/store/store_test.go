package store

import (
	"testing"

	"github.com/quanttide/qtcloud-learn-provider/internal/domain"
)

// TestClassStore_CRUD 移植自 qtcloud-course/provider 的 TestClassStore_CRUD，扩展 TeacherIDs / StudentIDs。
func TestClassStore_CRUD(t *testing.T) {
	s := NewClassStore()

	if got := s.List(); len(got) != 0 {
		t.Fatalf("List() = %d", len(got))
	}
	if _, ok := s.Get("x"); ok {
		t.Fatal("Get() nonexistent ok = true")
	}

	c := s.Create(&domain.Class{
		Name: "浙理班级", RefName: "大数据微专业", RefType: "program", RefID: "prog-1",
		StartDate: "2026-09-01", EndDate: "2027-01-15", StudentCount: 30,
		TeacherIDs: []string{"tea-1"}, StudentIDs: []string{"stu-1"},
	})
	if c.ID == "" || c.Name != "浙理班级" || c.StudentCount != 30 || c.Slug == "" {
		t.Fatalf("Create() = %+v", c)
	}
	if len(c.TeacherIDs) != 1 || c.TeacherIDs[0] != "tea-1" {
		t.Fatalf("Create().TeacherIDs = %v", c.TeacherIDs)
	}
	if len(c.StudentIDs) != 1 || c.StudentIDs[0] != "stu-1" {
		t.Fatalf("Create().StudentIDs = %v", c.StudentIDs)
	}

	// nil slices → 初始化为空切片
	c2 := s.Create(&domain.Class{Name: "杭电班级", RefID: "prog-2"})
	if c2.TeacherIDs == nil || c2.StudentIDs == nil {
		t.Fatal("Create(): TeacherIDs/StudentIDs should not be nil")
	}

	if got := s.List(); len(got) != 2 {
		t.Fatalf("List() = %d, want 2", len(got))
	}

	updated, ok := s.Update(&domain.Class{ID: c.ID, Name: "浙理班级v2", RefName: "大数据微专业v2", RefType: "course", RefID: "cour-1", Status: "active", StartDate: "2026-09-15", EndDate: "2027-02-01", StudentCount: 35, Progress: 0.5, TeacherIDs: []string{"tea-2"}, StudentIDs: []string{"stu-1", "stu-2"}})
	if !ok || updated.Name != "浙理班级v2" || updated.RefType != "course" || updated.StudentCount != 35 || updated.Progress != 0.5 {
		t.Fatalf("Update() = %+v", updated)
	}
	if len(updated.TeacherIDs) != 1 || updated.TeacherIDs[0] != "tea-2" {
		t.Fatalf("Update().TeacherIDs = %v", updated.TeacherIDs)
	}
	if len(updated.StudentIDs) != 2 {
		t.Fatalf("Update().StudentIDs = %v", updated.StudentIDs)
	}
	if _, ok := s.Update(&domain.Class{ID: "x"}); ok {
		t.Fatal("Update() nonexistent ok = true")
	}

	if ok := s.Delete(c.ID); !ok {
		t.Fatal("Delete() ok = false")
	}
	if ok := s.Delete(c.ID); ok {
		t.Fatal("Delete() again ok = true")
	}
	if ok := s.Delete("x"); ok {
		t.Fatal("Delete() nonexistent ok = true")
	}
}

func TestStudentStore_CRUD(t *testing.T) {
	s := NewStudentStore()

	if got := s.List(); len(got) != 0 {
		t.Fatalf("List() = %d", len(got))
	}
	if _, ok := s.Get("x"); ok {
		t.Fatal("Get() nonexistent ok = true")
	}

	st := s.Create(&domain.Student{Name: "张三", Email: "zhangsan@example.com", Plan: "free"})
	if st.ID == "" || st.Name != "张三" || st.Email != "zhangsan@example.com" || st.Plan != "free" {
		t.Fatalf("Create() = %+v", st)
	}

	s.Create(&domain.Student{Name: "李四", Plan: "vip"})
	if got := s.List(); len(got) != 2 {
		t.Fatalf("List() = %d, want 2", len(got))
	}

	updated, ok := s.Update(&domain.Student{ID: st.ID, Name: "张三丰", Email: "new@example.com", Avatar: "a.png", Plan: "paid"})
	if !ok || updated.Name != "张三丰" || updated.Avatar != "a.png" || updated.Plan != "paid" {
		t.Fatalf("Update() = %+v", updated)
	}
	if _, ok := s.Update(&domain.Student{ID: "x"}); ok {
		t.Fatal("Update() nonexistent ok = true")
	}

	if ok := s.Delete(st.ID); !ok {
		t.Fatal("Delete() ok = false")
	}
	if ok := s.Delete("x"); ok {
		t.Fatal("Delete() nonexistent ok = true")
	}
}

func TestAssessmentStore_CRUD(t *testing.T) {
	s := NewAssessmentStore()

	a := s.Create(&domain.Assessment{ClassID: "class-1", Type: "homework", Title: "作业1", FullScore: 100, PassScore: 60, Deadline: "2026-10-01"})
	if a.ID == "" || a.Title != "作业1" || a.FullScore != 100 || a.Deadline != "2026-10-01" {
		t.Fatalf("Create() = %+v", a)
	}
	if got := s.List(); len(got) != 1 {
		t.Fatalf("List() = %d, want 1", len(got))
	}

	updated, ok := s.Update(&domain.Assessment{ID: a.ID, ClassID: "class-1", Type: "exam", Title: "期中考试", FullScore: 150, PassScore: 90, Deadline: "2026-11-01"})
	if !ok || updated.Type != "exam" || updated.Title != "期中考试" || updated.FullScore != 150 {
		t.Fatalf("Update() = %+v", updated)
	}
	if _, ok := s.Update(&domain.Assessment{ID: "x"}); ok {
		t.Fatal("Update() nonexistent ok = true")
	}

	if ok := s.Delete(a.ID); !ok {
		t.Fatal("Delete() ok = false")
	}
}

func TestSubmissionStore_CRUD(t *testing.T) {
	s := NewSubmissionStore()

	sub := s.Create(&domain.Submission{AssessmentID: "assess-1", StudentID: "stu-1", Status: "submitted", SubmittedAt: "2026-09-30T10:00:00Z"})
	if sub.ID == "" || sub.AssessmentID != "assess-1" || sub.StudentID != "stu-1" || sub.Status != "submitted" {
		t.Fatalf("Create() = %+v", sub)
	}
	if got := s.List(); len(got) != 1 {
		t.Fatalf("List() = %d, want 1", len(got))
	}

	updated, ok := s.Update(&domain.Submission{ID: sub.ID, AssessmentID: "assess-1", StudentID: "stu-1", Status: "late", Score: 88.5, Comment: "不错", SubmittedAt: "2026-09-30T10:00:00Z"})
	if !ok || updated.Status != "late" || updated.Score != 88.5 || updated.Comment != "不错" {
		t.Fatalf("Update() = %+v", updated)
	}
	if _, ok := s.Update(&domain.Submission{ID: "x"}); ok {
		t.Fatal("Update() nonexistent ok = true")
	}

	if ok := s.Delete(sub.ID); !ok {
		t.Fatal("Delete() ok = false")
	}
}

func TestEnrollmentStore_CRUD(t *testing.T) {
	s := NewEnrollmentStore()

	e := s.Create(&domain.Enrollment{ClassID: "class-1", StudentID: "stu-1", Status: "enrolled", EnrolledAt: "2026-08-01"})
	if e.ID == "" || e.ClassID != "class-1" || e.StudentID != "stu-1" || e.Status != "enrolled" {
		t.Fatalf("Create() = %+v", e)
	}
	if got := s.List(); len(got) != 1 {
		t.Fatalf("List() = %d, want 1", len(got))
	}

	updated, ok := s.Update(&domain.Enrollment{ID: e.ID, ClassID: "class-1", StudentID: "stu-1", Status: "withdrawn", EnrolledAt: "2026-08-01"})
	if !ok || updated.Status != "withdrawn" {
		t.Fatalf("Update() = %+v", updated)
	}
	if _, ok := s.Update(&domain.Enrollment{ID: "x"}); ok {
		t.Fatal("Update() nonexistent ok = true")
	}

	if ok := s.Delete(e.ID); !ok {
		t.Fatal("Delete() ok = false")
	}
}

func TestProgressStore_CRUD(t *testing.T) {
	s := NewProgressStore()

	p := s.Create(&domain.Progress{StudentID: "stu-1", ClassID: "class-1", Percent: 0.0, Finished: false})
	if p.ID == "" || p.StudentID != "stu-1" || p.ClassID != "class-1" || p.Percent != 0.0 {
		t.Fatalf("Create() = %+v", p)
	}
	if got := s.List(); len(got) != 1 {
		t.Fatalf("List() = %d, want 1", len(got))
	}

	updated, ok := s.Update(&domain.Progress{ID: p.ID, StudentID: "stu-1", ClassID: "class-1", Percent: 1.0, Finished: true, UpdatedAt: "2026-08-02"})
	if !ok || updated.Percent != 1.0 || !updated.Finished || updated.UpdatedAt != "2026-08-02" {
		t.Fatalf("Update() = %+v", updated)
	}
	if _, ok := s.Update(&domain.Progress{ID: "x"}); ok {
		t.Fatal("Update() nonexistent ok = true")
	}

	if ok := s.Delete(p.ID); !ok {
		t.Fatal("Delete() ok = false")
	}
}
