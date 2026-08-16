package store

// 持久化测试：写后落盘、新实例 Load 恢复（模拟重启不丢）。

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/quanttide/qtcloud-learn-provider/internal/domain"
)

func TestApplicationPersistence(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "applications.json")

	s1 := NewApplicationStore()
	s1.BaseStore.SetPersistPath(path)
	created := s1.Create(&domain.Application{
		ProjectName:  "持久化项目",
		BlindSpot:    "盲区",
		DemoPlan:     "方案",
		TeamMode:     "personal",
		MemberNames:  []string{"张三"},
		StudentID:    "stu-1",
		StudentName:  "张三",
		Status:       "submitted",
	})
	if created.ID == "" {
		t.Fatal("create failed")
	}

	// 新实例模拟重启：Load 后数据仍在
	s2 := NewApplicationStore()
	if err := s2.BaseStore.Load(path); err != nil {
		t.Fatalf("load: %v", err)
	}
	got, ok := s2.Get(created.ID)
	if !ok {
		t.Fatalf("restored app %s not found", created.ID)
	}
	if got.ProjectName != "持久化项目" || len(got.MemberNames) != 1 || got.MemberNames[0] != "张三" {
		t.Errorf("restored = %+v", got)
	}

	// 序号恢复：新创建的 ID 不与旧记录冲突
	next := s2.Create(&domain.Application{ProjectName: "第二个", TeamMode: "personal", MemberNames: []string{"李四"}})
	if next.ID == created.ID {
		t.Errorf("seq not restored: %s == %s", next.ID, created.ID)
	}

	// 文件不存在时 Load 静默跳过（首启）
	s3 := NewApplicationStore()
	if err := s3.BaseStore.Load(filepath.Join(dir, "missing.json")); err != nil {
		t.Errorf("missing file should be no-op, got %v", err)
	}
}

func TestPersistFileWritten(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "progress.json")
	s := NewProgressStore()
	s.BaseStore.SetPersistPath(path)
	s.Create(&domain.Progress{StudentID: "stu-1", Percent: 0.5, Finished: false})
	if _, err := os.Stat(path); err != nil {
		t.Fatalf("persist file not written: %v", err)
	}
}
