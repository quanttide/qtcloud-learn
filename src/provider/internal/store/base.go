package store

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"
)

// BaseStore 提供通用的内存存储骨架：List/Get/Delete 和 ID 生成。
// 具体类型嵌入后只需实现 Create/Update（因字段各异）。
// 持久化：SetPersistPath 后各 store 写操作尾部调 persist() 全量落盘，
// 启动时 Load 恢复——JSON 文件存储，重启不丢（v0.1 MVP；RDS 后续版本接入）。
type BaseStore[T any] struct {
	mu       sync.RWMutex
	data     map[string]*T
	seq      int
	idPrefix string

	persistPath string
}

// SetPersistPath 启用文件持久化（全量快照，原子写：临时文件 + rename）。
func (s *BaseStore[T]) SetPersistPath(path string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.persistPath = path
}

// Load 从 JSON 文件恢复数据（文件不存在时静默跳过——首启）。
func (s *BaseStore[T]) Load(path string) error {
	raw, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	var snapshot struct {
		Data map[string]*T `json:"data"`
		Seq  int           `json:"seq"`
	}
	if err := json.Unmarshal(raw, &snapshot); err != nil {
		return fmt.Errorf("restore %s: %w", path, err)
	}
	if snapshot.Data != nil {
		s.data = snapshot.Data
	}
	s.seq = snapshot.Seq
	return nil
}

// persist 全量快照落盘（写操作后调用；未启用持久化时为空操作）。
func (s *BaseStore[T]) persist() {
	if s.persistPath == "" {
		return
	}
	snapshot := struct {
		Data map[string]*T `json:"data"`
		Seq  int           `json:"seq"`
	}{Data: s.data, Seq: s.seq}
	raw, err := json.Marshal(snapshot)
	if err != nil {
		return
	}
	tmp := s.persistPath + ".tmp"
	if err := os.WriteFile(tmp, raw, 0o644); err != nil {
		return
	}
	_ = os.Rename(tmp, s.persistPath)
}

// NewBaseStore 创建泛型存储。
func NewBaseStore[T any](idPrefix string) *BaseStore[T] {
	return &BaseStore[T]{
		data:     make(map[string]*T),
		seq:      1,
		idPrefix: idPrefix,
	}
}

// nextID 生成自增 ID，如 "class-1"、"stu-2"。
func (s *BaseStore[T]) nextID() string {
	id := fmt.Sprintf("%s-%d", s.idPrefix, s.seq)
	s.seq++
	return id
}

// List 返回全部实体。
func (s *BaseStore[T]) List() []*T {
	s.mu.RLock()
	defer s.mu.RUnlock()
	result := make([]*T, 0, len(s.data))
	for _, v := range s.data {
		result = append(result, v)
	}
	return result
}

// Get 按 ID 查找实体。
func (s *BaseStore[T]) Get(id string) (*T, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	v, ok := s.data[id]
	return v, ok
}

// Delete 按 ID 删除实体。
func (s *BaseStore[T]) Delete(id string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	_, ok := s.data[id]
	if ok {
		delete(s.data, id)
	}
	return ok
}

// NameExists 检查 name 是否已被占用。
func (s *BaseStore[T]) NameExists(name string, getName func(*T) string) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	for _, v := range s.data {
		if getName(v) == name {
			return true
		}
	}
	return false
}
