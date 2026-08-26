package store

import (
	"testing"

	"github.com/quanttide/qtcloud-learn-provider/internal/domain"
)

func TestLearnerStore_CRUD(t *testing.T) {
	s := NewLearnerStore()

	if got := s.List(); len(got) != 0 {
		t.Fatalf("List() = %d", len(got))
	}
	if _, ok := s.Get("x"); ok {
		t.Fatal("Get() nonexistent ok = true")
	}

	l := s.Create(&domain.Learner{UserID: "user-123"})
	if l.ID == "" || l.UserID != "user-123" {
		t.Fatalf("Create() = %+v", l)
	}

	s.Create(&domain.Learner{})
	if got := s.List(); len(got) != 2 {
		t.Fatalf("List() = %d, want 2", len(got))
	}

	updated, ok := s.Update(&domain.Learner{ID: l.ID, UserID: "user-456"})
	if !ok || updated.UserID != "user-456" {
		t.Fatalf("Update() = %+v", updated)
	}
	if _, ok := s.Update(&domain.Learner{ID: "x"}); ok {
		t.Fatal("Update() nonexistent ok = true")
	}

	if ok := s.Delete(l.ID); !ok {
		t.Fatal("Delete() ok = false")
	}
	if ok := s.Delete(l.ID); ok {
		t.Fatal("Delete() again ok = true")
	}
	if ok := s.Delete("x"); ok {
		t.Fatal("Delete() nonexistent ok = true")
	}
}

func TestCriterionStore_CRUD(t *testing.T) {
	s := NewCriterionStore()

	c := s.Create(&domain.Criterion{Title: "会连接 Zed", Description: "成功建立 Zed 连接"})
	if c.ID == "" || c.Title != "会连接 Zed" || c.Description != "成功建立 Zed 连接" {
		t.Fatalf("Create() = %+v", c)
	}
	if got := s.List(); len(got) != 1 {
		t.Fatalf("List() = %d, want 1", len(got))
	}

	updated, ok := s.Update(&domain.Criterion{ID: c.ID, Title: "会用 Agent 执行任务", Description: "完成 Agent 任务"})
	if !ok || updated.Title != "会用 Agent 执行任务" || updated.Description != "完成 Agent 任务" {
		t.Fatalf("Update() = %+v", updated)
	}
	if _, ok := s.Update(&domain.Criterion{ID: "x"}); ok {
		t.Fatal("Update() nonexistent ok = true")
	}

	if ok := s.Delete(c.ID); !ok {
		t.Fatal("Delete() ok = false")
	}
}

func TestCompletionStore_CRUD(t *testing.T) {
	s := NewCompletionStore()

	c := s.Create(&domain.Completion{LearnerID: "lea-1", CriterionID: "cri-1"})
	if c.ID == "" || c.LearnerID != "lea-1" || c.CriterionID != "cri-1" {
		t.Fatalf("Create() = %+v", c)
	}
	if c.Status != "not_completed" {
		t.Fatalf("Create() default status = %q, want not_completed", c.Status)
	}
	if c.CreatedAt == "" || c.UpdatedAt == "" {
		t.Fatalf("Create() timestamps missing: %+v", c)
	}
	if got := s.List(); len(got) != 1 {
		t.Fatalf("List() = %d, want 1", len(got))
	}

	updated, ok := s.Update(&domain.Completion{ID: c.ID, LearnerID: "lea-1", CriterionID: "cri-1", Status: "completed"})
	if !ok || updated.Status != "completed" {
		t.Fatalf("Update() = %+v", updated)
	}
	if updated.CreatedAt == "" || updated.UpdatedAt == "" {
		t.Fatalf("Update() timestamps missing: %+v", updated)
	}
	if _, ok := s.Update(&domain.Completion{ID: "x"}); ok {
		t.Fatal("Update() nonexistent ok = true")
	}

	if ok := s.Delete(c.ID); !ok {
		t.Fatal("Delete() ok = false")
	}
}
