package domain

import (
	"encoding/json"
	"testing"
)

func TestLearner_JSON(t *testing.T) {
	l := Learner{ID: "lea-1", UserID: "user-123"}
	b, _ := json.Marshal(l)
	var got Learner
	json.Unmarshal(b, &got)
	if got.ID != "lea-1" || got.UserID != "user-123" {
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
		ID: "com-1", LearnerID: "lea-1", CriterionID: "cri-1",
		Status: "completed", CreatedAt: "2026-08-26T10:00:00Z", UpdatedAt: "2026-08-26T10:00:00Z",
	}
	b, _ := json.Marshal(c)
	var got Completion
	json.Unmarshal(b, &got)
	if got.LearnerID != "lea-1" || got.CriterionID != "cri-1" || got.Status != "completed" {
		t.Fatalf("roundtrip = %+v", got)
	}
	if got.CreatedAt == "" || got.UpdatedAt == "" {
		t.Fatalf("timestamps missing: %+v", got)
	}
}
