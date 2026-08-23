package handler

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// 生成测试密钥对（RS256）：签发用私钥、验签用公钥。
func testKeyPair(t *testing.T) (*rsa.PrivateKey, string) {
	t.Helper()
	key, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatalf("generate key: %v", err)
	}
	pubDER, err := x509.MarshalPKIXPublicKey(&key.PublicKey)
	if err != nil {
		t.Fatalf("marshal pubkey: %v", err)
	}
	pubPEM := pem.EncodeToMemory(&pem.Block{Type: "PUBLIC KEY", Bytes: pubDER})
	return key, string(pubPEM)
}

// 用测试私钥签发 JWT。
func signToken(t *testing.T, key *rsa.PrivateKey, sub string) string {
	t.Helper()
	claims := jwt.MapClaims{
		"sub": sub,
		"exp": time.Now().Add(time.Hour).Unix(),
	}
	tok, err := jwt.NewWithClaims(jwt.SigningMethodRS256, claims).SignedString(key)
	if err != nil {
		t.Fatalf("sign token: %v", err)
	}
	return tok
}

func TestNewAuthMiddleware_NoEnvReturnsNil(t *testing.T) {
	os.Unsetenv("JWT_PUBLIC_KEY")
	mw, err := NewAuthMiddleware()
	if err != nil {
		t.Fatalf("NewAuthMiddleware with empty env: %v", err)
	}
	if mw != nil {
		t.Fatal("expected nil middleware when JWT_PUBLIC_KEY unset")
	}
}

func TestAuthMiddleware_RejectsMissingToken(t *testing.T) {
	_, pub := testKeyPair(t)
	t.Setenv("JWT_PUBLIC_KEY", pub)
	mw, err := NewAuthMiddleware()
	if err != nil {
		t.Fatalf("NewAuthMiddleware: %v", err)
	}
	handler := mw.Wrap(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	req := httptest.NewRequest(http.MethodPost, "/api/proposals", nil)
	rec := httptest.NewRecorder()
	handler(rec, req)
	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("status = %d, want 401", rec.Code)
	}
}

func TestAuthMiddleware_RejectsInvalidToken(t *testing.T) {
	_, pub := testKeyPair(t)
	t.Setenv("JWT_PUBLIC_KEY", pub)
	mw, err := NewAuthMiddleware()
	if err != nil {
		t.Fatalf("NewAuthMiddleware: %v", err)
	}
	handler := mw.Wrap(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	req := httptest.NewRequest(http.MethodPost, "/api/proposals", nil)
	req.Header.Set("Authorization", "Bearer not-a-jwt")
	rec := httptest.NewRecorder()
	handler(rec, req)
	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("status = %d, want 401", rec.Code)
	}
}

func TestAuthMiddleware_RejectsNonRS256Token(t *testing.T) {
	priv, pub := testKeyPair(t)
	t.Setenv("JWT_PUBLIC_KEY", pub)
	mw, err := NewAuthMiddleware()
	if err != nil {
		t.Fatalf("NewAuthMiddleware: %v", err)
	}
	handler := mw.Wrap(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	claims := jwt.MapClaims{
		"sub": "user-123",
		"exp": time.Now().Add(time.Hour).Unix(),
	}
	token, err := jwt.NewWithClaims(jwt.SigningMethodRS512, claims).SignedString(priv)
	if err != nil {
		t.Fatalf("sign token: %v", err)
	}
	req := httptest.NewRequest(http.MethodPost, "/api/proposals", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	rec := httptest.NewRecorder()
	handler(rec, req)
	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("status = %d, want 401", rec.Code)
	}
}

func TestAuthMiddleware_PassesValidTokenWithSub(t *testing.T) {
	priv, pub := testKeyPair(t)
	t.Setenv("JWT_PUBLIC_KEY", pub)
	mw, err := NewAuthMiddleware()
	if err != nil {
		t.Fatalf("NewAuthMiddleware: %v", err)
	}
	var gotSub string
	handler := mw.Wrap(func(w http.ResponseWriter, r *http.Request) {
		gotSub = UserIDFrom(r)
		w.WriteHeader(http.StatusOK)
	})
	req := httptest.NewRequest(http.MethodPost, "/api/proposals", nil)
	req.Header.Set("Authorization", "Bearer "+signToken(t, priv, "user-123"))
	rec := httptest.NewRecorder()
	handler(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200", rec.Code)
	}
	if gotSub != "user-123" {
		t.Fatalf("sub = %q, want user-123", gotSub)
	}
}

func TestAuthMiddleware_RejectsForeignKeyToken(t *testing.T) {
	// 配置公钥来自密钥对 A；token 用密钥对 B 签发 → 验签必须失败
	_, pubA := testKeyPair(t)
	t.Setenv("JWT_PUBLIC_KEY", pubA)
	mw, err := NewAuthMiddleware()
	if err != nil {
		t.Fatalf("NewAuthMiddleware: %v", err)
	}
	handler := mw.Wrap(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	privB, _ := testKeyPair(t)
	req := httptest.NewRequest(http.MethodPost, "/api/proposals", nil)
	req.Header.Set("Authorization", "Bearer "+signToken(t, privB, "evil-user"))
	rec := httptest.NewRecorder()
	handler(rec, req)
	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("status = %d, want 401", rec.Code)
	}
}
