package domain

import (
	"encoding/json"
	"testing"
)

func TestLearner_JSON(t *testing.T) {
	l := Learner{ID: "lea-1", UserID: "user-123", ScheduleID: "schedule-agent-engineer"}
	b, _ := json.Marshal(l)
	var got Learner
	json.Unmarshal(b, &got)
	if got.ID != "lea-1" || got.UserID != "user-123" || got.ScheduleID != "schedule-agent-engineer" {
		t.Fatalf("roundtrip = %+v", got)
	}
}

func TestLearner_UserIDOptional(t *testing.T) {
	// user_id 预留字段：为空时不序列化
	b, _ := json.Marshal(Learner{ID: "lea-1"})
	if got := string(b); got != `{"id":"lea-1"}` {
		t.Fatalf("marshal = %s", got)
	}
}

func TestCompletion_JSON(t *testing.T) {
	c := Completion{
		ID: "com-1", LearnerID: "lea-1", TaskID: "task-data-second-brain",
		Status: CompletionStatusCompleted, CreatedAt: "2026-08-26T10:00:00Z", UpdatedAt: "2026-08-26T10:00:00Z",
	}
	b, _ := json.Marshal(c)
	var got Completion
	json.Unmarshal(b, &got)
	if got.LearnerID != "lea-1" || got.TaskID != "task-data-second-brain" || got.Status != CompletionStatusCompleted {
		t.Fatalf("roundtrip = %+v", got)
	}
	if got.CreatedAt == "" || got.UpdatedAt == "" {
		t.Fatalf("timestamps missing: %+v", got)
	}
}

func TestTask_JSON(t *testing.T) {
	task := Task{ID: "task-data-second-brain", Title: "熟悉数据工程第二大脑", Description: "完成一条改进建议。"}
	b, _ := json.Marshal(task)
	var got Task
	json.Unmarshal(b, &got)
	if got.ID != task.ID || got.Title != task.Title || got.Description != task.Description {
		t.Fatalf("roundtrip = %+v", got)
	}
}

func TestSchedule_JSON(t *testing.T) {
	schedule := Schedule{
		ID:          "schedule-agent-engineer",
		Title:       "智能体工程师训练营",
		Description: "按学习路径推进的训练计划。",
		Tasks: []Task{{
			ID:          "task-data-second-brain",
			Title:       "熟悉数据工程第二大脑",
			Description: "完成一条改进建议。",
		}},
	}
	b, _ := json.Marshal(schedule)
	var got Schedule
	json.Unmarshal(b, &got)
	if got.ID != schedule.ID || got.Title != schedule.Title || len(got.Tasks) != 1 || got.Tasks[0].ID != "task-data-second-brain" {
		t.Fatalf("roundtrip = %+v", got)
	}
}
