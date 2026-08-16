package store

import (
	"encoding/json"
	"fmt"
	"sync"
)

// BaseStore 提供通用的内存存储骨架：List/Get/Delete 和 ID 生成。
// 具体类型嵌入后只需实现 Create/Update（因字段各异）。
// 持久化：SetPersister 后各 store 写操作尾部调 persist() 全量快照，
// 启动时 Load 恢复——本地文件 / OSS 均可（v0.1 MVP；RDS 后续版本接入）。
type BaseStore[T any] struct {
	mu       sync.RWMutex
	data     map[string]*T
	seq      int
	idPrefix string

	persister Persister
}

// Persister 快照读写后端（文件 / OSS）。key 为 store 专属（如 applications.json）。
type Persister interface {
	Load(key string) ([]byte, error)
	Store(key string, data []byte) error
}

// SetPersister 启用持久化（原子写由各后端保证；未设置时 persist 为空操作）。
func (s *BaseStore[T]) SetPersister(p Persister) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.persister = p
}

// Load 从持久化后端恢复数据（后端无数据时静默跳过——首启）。
func (s *BaseStore[T]) Load(key string) error {
	if s.persister == nil {
		return nil
	}
	raw, err := s.persister.Load(key)
	if err != nil {
		return err
	}
	if len(raw) == 0 {
		return nil
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	var snapshot struct {
		Data map[string]*T `json:"data"`
		Seq  int           `json:"seq"`
	}
	if err := json.Unmarshal(raw, &snapshot); err != nil {
		return fmt.Errorf("restore %s: %w", key, err)
	}
	if snapshot.Data != nil {
		s.data = snapshot.Data
	}
	s.seq = snapshot.Seq
	return nil
}

// persist 全量快照（写操作后调用；未启用持久化时为空操作）。
func (s *BaseStore[T]) persist() {
	if s.persister == nil {
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
	_ = s.persister.Store(s.idPrefix+".json", raw)
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
