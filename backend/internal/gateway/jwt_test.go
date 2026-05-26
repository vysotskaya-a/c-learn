package gateway

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/c-learn/internal/auth"
	"github.com/gin-gonic/gin"
)

func TestIsPublicPath(t *testing.T) {
	public := []string{
		"/api/v1/auth/register",
		"/api/v1/auth/login",
		"/api/v1/auth/refresh",
	}
	for _, p := range public {
		if !isPublicPath(p) {
			t.Fatalf("%q should be public", p)
		}
	}
	if isPublicPath("/api/v1/courses/tree") {
		t.Fatal("protected path must not be public")
	}
}

func TestJWTMiddleware_PublicPath(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mgr := auth.NewJWTManager("secret")
	r := gin.New()
	r.Use(JWTMiddleware(mgr))
	r.GET("/api/v1/auth/login", func(c *gin.Context) { c.Status(http.StatusOK) })

	req := httptest.NewRequest(http.MethodGet, "/api/v1/auth/login", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200", w.Code)
	}
}

func TestJWTMiddleware_MissingAuth(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mgr := auth.NewJWTManager("secret")
	r := gin.New()
	r.Use(JWTMiddleware(mgr))
	r.GET("/api/v1/courses/tree", func(c *gin.Context) { c.Status(http.StatusOK) })

	req := httptest.NewRequest(http.MethodGet, "/api/v1/courses/tree", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Fatalf("status = %d, want 401", w.Code)
	}
}

func TestJWTMiddleware_ValidToken(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mgr := auth.NewJWTManager("secret")
	token, err := mgr.GenerateAccessToken("user-42", "admin")
	if err != nil {
		t.Fatalf("token: %v", err)
	}

	var gotUserID, gotRole string
	r := gin.New()
	r.Use(JWTMiddleware(mgr))
	r.GET("/protected", func(c *gin.Context) {
		gotUserID = c.GetHeader("X-User-ID")
		gotRole = c.GetHeader("X-User-Role")
		c.Status(http.StatusOK)
	})

	req := httptest.NewRequest(http.MethodGet, "/protected", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200", w.Code)
	}
	if gotUserID != "user-42" || gotRole != "admin" {
		t.Fatalf("headers: user=%q role=%q", gotUserID, gotRole)
	}
}

func TestJWTMiddleware_InvalidToken(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mgr := auth.NewJWTManager("secret")
	r := gin.New()
	r.Use(JWTMiddleware(mgr))
	r.GET("/protected", func(c *gin.Context) { c.Status(http.StatusOK) })

	req := httptest.NewRequest(http.MethodGet, "/protected", nil)
	req.Header.Set("Authorization", "Bearer invalid")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Fatalf("status = %d, want 401", w.Code)
	}
}
