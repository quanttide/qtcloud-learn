package handler

import (
	"github.com/quanttide/qtcloud-learn-provider/internal/domain"
	"github.com/quanttide/qtcloud-learn-provider/internal/store"
)

// SessionHandler 提供 Session 的标准 CRUD。
type SessionHandler = CRUDHandler[domain.Session]

// NewSessionHandler 创建 Session handler。
func NewSessionHandler(s *store.SessionStore) *SessionHandler {
	return NewCRUDHandler(
		s,
		func(sess *domain.Session) string {
			if sess.ClassID == "" {
				return "classId is required"
			}
			if sess.LessonTitle == "" {
				return "lessonTitle is required"
			}
			return ""
		},
		func(sess *domain.Session, id string) { sess.ID = id },
	)
}
