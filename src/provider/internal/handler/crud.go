package handler

import (
	"encoding/json"
	"net/http"
)

// CRUDStore 是泛型处理器所需的存储接口。
type CRUDStore[T any] interface {
	List() []*T
	Get(string) (*T, bool)
	Create(*T) *T
	Update(*T) (*T, bool)
	Delete(string) bool
}

// CRUDHandler 提供标准 CRUD HTTP handler，适合无特殊校验的资源。
type CRUDHandler[T any] struct {
	store       CRUDStore[T]
	validateFn  func(*T) string     // 返回空串表示合法，否则返回错误消息
	setIDFn     func(*T, string)    // 设置实体 ID（从路径取值）
	nameChecker func(string) string // 可选：返回空串表示可用，否则返回冲突消息
	getNameFn   func(*T) string     // 可选：从实体提取 name
}

// NewCRUDHandler 创建泛型 CRUD handler。
func NewCRUDHandler[T any](store CRUDStore[T], validateFn func(*T) string, setIDFn func(*T, string)) *CRUDHandler[T] {
	return &CRUDHandler[T]{store: store, validateFn: validateFn, setIDFn: setIDFn}
}

// WithNameCheck 添加 name 唯一性校验。
// checker 接收 name 返回 "" 表示允许或错误消息；getName 从实体提取 name。
func (h *CRUDHandler[T]) WithNameCheck(checker func(string) string, getName func(*T) string) *CRUDHandler[T] {
	h.nameChecker = checker
	h.getNameFn = getName
	return h
}

func (h *CRUDHandler[T]) List(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, h.store.List())
}

func (h *CRUDHandler[T]) Create(w http.ResponseWriter, r *http.Request) {
	var entity T
	if err := json.NewDecoder(r.Body).Decode(&entity); err != nil {
		http.Error(w, `{"error":"invalid request body"}`, http.StatusBadRequest)
		return
	}
	if msg := h.validateFn(&entity); msg != "" {
		http.Error(w, `{"error":"`+msg+`"}`, http.StatusBadRequest)
		return
	}
	if h.nameChecker != nil && h.getNameFn != nil {
		if msg := h.nameChecker(h.getNameFn(&entity)); msg != "" {
			http.Error(w, `{"error":"`+msg+`"}`, http.StatusConflict)
			return
		}
	}
	writeJSON(w, http.StatusCreated, h.store.Create(&entity))
}

func (h *CRUDHandler[T]) Get(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	entity, ok := h.store.Get(id)
	if !ok {
		http.Error(w, `{"error":"not found"}`, http.StatusNotFound)
		return
	}
	writeJSON(w, http.StatusOK, entity)
}

func (h *CRUDHandler[T]) Update(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	existing, ok := h.store.Get(id)
	if !ok {
		http.Error(w, `{"error":"not found"}`, http.StatusNotFound)
		return
	}

	// 合并语义：仅覆盖请求体中出现的字段，未提交字段保持不变。
	merged, err := mergeJSON(existing, r)
	if err != nil {
		http.Error(w, `{"error":"invalid request body"}`, http.StatusBadRequest)
		return
	}

	var entity T
	if err := json.Unmarshal(merged, &entity); err != nil {
		http.Error(w, `{"error":"invalid request body"}`, http.StatusBadRequest)
		return
	}
	h.setIDFn(&entity, id)
	updated, ok := h.store.Update(&entity)
	if !ok {
		http.Error(w, `{"error":"not found"}`, http.StatusNotFound)
		return
	}
	writeJSON(w, http.StatusOK, updated)
}

// mergeJSON 将请求体作为补丁合并到现有实体上，返回合并后的 JSON。
func mergeJSON(existing any, r *http.Request) ([]byte, error) {
	base, err := json.Marshal(existing)
	if err != nil {
		return nil, err
	}
	var merged map[string]json.RawMessage
	if err := json.Unmarshal(base, &merged); err != nil {
		return nil, err
	}
	var patch map[string]json.RawMessage
	if err := json.NewDecoder(r.Body).Decode(&patch); err != nil {
		return nil, err
	}
	for k, v := range patch {
		merged[k] = v
	}
	return json.Marshal(merged)
}

func (h *CRUDHandler[T]) Delete(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if !h.store.Delete(id) {
		http.Error(w, `{"error":"not found"}`, http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
