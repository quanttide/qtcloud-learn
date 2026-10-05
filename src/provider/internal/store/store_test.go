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

func TestCompletionStore_CRUD(t *testing.T) {
	s := NewCompletionStore()

	c := s.Create(&domain.Completion{LearnerID: "lea-1", TaskID: "task-data-second-brain"})
	if c.ID == "" || c.LearnerID != "lea-1" || c.TaskID != "task-data-second-brain" {
		t.Fatalf("Create() = %+v", c)
	}
	if c.Status != domain.CompletionStatusNotCompleted {
		t.Fatalf("Create() default status = %q, want not_completed", c.Status)
	}
	if c.CreatedAt == "" || c.UpdatedAt == "" {
		t.Fatalf("Create() timestamps missing: %+v", c)
	}
	if got := s.List(); len(got) != 1 {
		t.Fatalf("List() = %d, want 1", len(got))
	}

	updated, ok := s.Update(&domain.Completion{ID: c.ID, LearnerID: "lea-1", TaskID: "task-data-second-brain", Status: domain.CompletionStatusCompleted})
	if !ok || updated.Status != domain.CompletionStatusCompleted {
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

func TestTaskStore_CRUD(t *testing.T) {
	s := NewTaskStore()

	task := s.Create(&domain.Task{ID: "task-data-second-brain", Title: "熟悉数据工程第二大脑", Description: "完成一条改进建议。"})
	if task.ID != "task-data-second-brain" || task.Title == "" || task.Description == "" {
		t.Fatalf("Create() = %+v", task)
	}

	created := s.Create(&domain.Task{Title: "整理数据工程意图", Description: "写清楚建设意图。"})
	if created.ID == "" || created.ID == "task-data-second-brain" {
		t.Fatalf("Create() generated id = %+v", created)
	}

	updated, ok := s.Update(&domain.Task{ID: task.ID, Title: "熟悉第二大脑", Description: "更新描述。"})
	if !ok || updated.Title != "熟悉第二大脑" || updated.Description != "更新描述。" {
		t.Fatalf("Update() = %+v", updated)
	}
	if _, ok := s.Update(&domain.Task{ID: "x"}); ok {
		t.Fatal("Update() nonexistent ok = true")
	}
}

func TestScheduleStore_CRUD(t *testing.T) {
	s := NewScheduleStore()
	tasks := []domain.Task{{ID: "task-data-second-brain", Title: "熟悉数据工程第二大脑", Description: "完成一条改进建议。"}}

	schedule := s.Create(&domain.Schedule{ID: "schedule-agent-engineer", Title: "智能体工程师训练营", Tasks: tasks})
	if schedule.ID != "schedule-agent-engineer" || len(schedule.Tasks) != 1 {
		t.Fatalf("Create() = %+v", schedule)
	}

	updated, ok := s.Update(&domain.Schedule{ID: schedule.ID, Title: "智能体工程师训练营 v2", Tasks: tasks})
	if !ok || updated.Title != "智能体工程师训练营 v2" || len(updated.Tasks) != 1 {
		t.Fatalf("Update() = %+v", updated)
	}
	if _, ok := s.Update(&domain.Schedule{ID: "x"}); ok {
		t.Fatal("Update() nonexistent ok = true")
	}
}
