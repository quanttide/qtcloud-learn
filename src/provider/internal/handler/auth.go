package handler

// 统一账号接入（果总 2026-08-19：所有系统公用 qtcloud-auth 账号系统）。
// 本服务用 qtcloud-auth 的 RSA 公钥验签（env JWT_PUBLIC_KEY），
// 各应用侧做权限决策（零信任：登录 ≠ 授权，见 quanttide-report-of-authorization
// tech-decisions/unified-account-authorization-principles.md）。
// JWT_PUBLIC_KEY 未配置时中间件为 nil（本地 dev/测试可用），生产必须配置。

import (
	"context"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

type contextKey string

const ClaimsKey contextKey = "auth_claims"

// UserIDFrom 从请求上下文提取当前用户 ID（JWT sub）。
func UserIDFrom(r *http.Request) string {
	if claims, ok := r.Context().Value(ClaimsKey).(jwt.MapClaims); ok {
		if sub, ok := claims["sub"].(string); ok {
			return sub
		}
	}
	return ""
}

// AuthMiddleware 解析并验证 Bearer JWT（RS256），注入 claims 到请求上下文。
type AuthMiddleware struct {
	pubKey *rsa.PublicKey
}

// NewAuthMiddleware 从 env JWT_PUBLIC_KEY（RSA 公钥 PEM）构建中间件；
// 未配置时返回 nil（鉴权不启用，仅限本地 dev/测试）。
func NewAuthMiddleware() (*AuthMiddleware, error) {
	pemStr := strings.TrimSpace(os.Getenv("JWT_PUBLIC_KEY"))
	if pemStr == "" {
		return nil, nil
	}
	block, _ := pem.Decode([]byte(pemStr))
	if block == nil {
		return nil, fmt.Errorf("auth: invalid JWT public key PEM")
	}
	key, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("auth: parse JWT public key: %w", err)
	}
	pub, ok := key.(*rsa.PublicKey)
	if !ok {
		return nil, fmt.Errorf("auth: JWT key must be RSA public key")
	}
	return &AuthMiddleware{pubKey: pub}, nil
}

// Wrap 包裹需要登录的处理器。
func (m *AuthMiddleware) Wrap(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tokenStr := extractBearerToken(r)
		if tokenStr == "" {
			writeError(w, http.StatusUnauthorized, "missing bearer token")
			return
		}
		token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (any, error) {
			if t.Method != jwt.SigningMethodRS256 {
				return nil, jwt.ErrSignatureInvalid
			}
			return m.pubKey, nil
		})
		if err != nil || !token.Valid {
			writeError(w, http.StatusUnauthorized, "invalid or expired token")
			return
		}
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			writeError(w, http.StatusUnauthorized, "invalid token claims")
			return
		}
		ctx := context.WithValue(r.Context(), ClaimsKey, claims)
		next(w, r.WithContext(ctx))
	}
}

func extractBearerToken(r *http.Request) string {
	auth := r.Header.Get("Authorization")
	const prefix = "Bearer "
	if len(auth) > len(prefix) && strings.EqualFold(auth[:len(prefix)], prefix) {
		return strings.TrimSpace(auth[len(prefix):])
	}
	return ""
}
