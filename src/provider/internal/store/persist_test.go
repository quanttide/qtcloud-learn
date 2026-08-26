package store

// 持久化测试：写后落盘、新实例 Load 恢复（模拟重启不丢）。

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/quanttide/qtcloud-learn-provider/internal/domain"
)

func TestCompletionPersistence(t *testing.T) {
	dir := t.TempDir()

	s1 := NewCompletionStore()
	s1.BaseStore.SetPersister(NewFilePersister(dir))
	created := s1.Create(&domain.Completion{LearnerID: "lea-1", CriterionID: "cri-1", Status: "completed"})
	if created.ID == "" {
		t.Fatal("create failed")
	}

	// 新实例模拟重启：Load 后数据仍在
	s2 := NewCompletionStore()
	s2.BaseStore.SetPersister(NewFilePersister(dir))
	if err := s2.BaseStore.Load("com.json"); err != nil {
		t.Fatalf("load: %v", err)
	}
	got, ok := s2.Get(created.ID)
	if !ok {
		t.Fatalf("restored completion %s not found", created.ID)
	}
	if got.LearnerID != "lea-1" || got.CriterionID != "cri-1" || got.Status != "completed" {
		t.Errorf("restored = %+v", got)
	}

	// 序号恢复：新创建的 ID 不与旧记录冲突
	next := s2.Create(&domain.Completion{LearnerID: "lea-1", CriterionID: "cri-2"})
	if next.ID == created.ID {
		t.Errorf("seq not restored: %s == %s", next.ID, created.ID)
	}

	// 文件不存在时 Load 静默跳过（首启）
	s3 := NewCompletionStore()
	s3.BaseStore.SetPersister(NewFilePersister(dir))
	if err := s3.BaseStore.Load("missing.json"); err != nil {
		t.Errorf("missing file should be no-op, got %v", err)
	}
}

func TestPersistFileWritten(t *testing.T) {
	dir := t.TempDir()
	s := NewCriterionStore()
	s.BaseStore.SetPersister(NewFilePersister(dir))
	s.Create(&domain.Criterion{Title: "会连接 Zed", Description: "成功建立 Zed 连接"})
	path := filepath.Join(dir, "cri.json")
	if _, err := os.Stat(path); err != nil {
		t.Fatalf("persist file not written: %v", err)
	}
}

func TestDeletePersisted(t *testing.T) {
	dir := t.TempDir()

	s1 := NewCriterionStore()
	s1.BaseStore.SetPersister(NewFilePersister(dir))
	created := s1.Create(&domain.Criterion{Title: "会连接 Zed", Description: "成功建立 Zed 连接"})
	if !s1.Delete(created.ID) {
		t.Fatal("delete failed")
	}

	// 新实例模拟重启：删除结果已落盘，记录不应复活
	s2 := NewCriterionStore()
	s2.BaseStore.SetPersister(NewFilePersister(dir))
	if err := s2.BaseStore.Load("cri.json"); err != nil {
		t.Fatalf("load: %v", err)
	}
	if _, ok := s2.Get(created.ID); ok {
		t.Fatal("deleted record resurrected after restart")
	}
	// 删除不存在的记录返回 false 且不应影响快照
	if s2.Delete("x") {
		t.Fatal("delete nonexistent ok = true")
	}
}
