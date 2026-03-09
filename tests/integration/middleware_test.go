package integration_test

import (
	"net/http"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"

	"mqfm-backend/internal/infrastructure/middleware"
	"mqfm-backend/tests/testutil"
)

func TestMiddleware_JWTAuth_Valid(t *testing.T) {
	r := testutil.SetupRouter()
	r.GET("/protected", middleware.JWTAuth(), func(c *gin.Context) {
		userID, _ := c.Get("user_id")
		role, _ := c.Get("role")
		c.JSON(http.StatusOK, gin.H{"user_id": userID, "role": role})
	})

	req := testutil.MakeAuthRequest("GET", "/protected", nil, 1, "admin")
	w := testutil.PerformRequest(r, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestMiddleware_JWTAuth_Missing(t *testing.T) {
	r := testutil.SetupRouter()
	r.GET("/protected", middleware.JWTAuth(), func(c *gin.Context) {
		c.JSON(http.StatusOK, nil)
	})

	req := testutil.MakeRequest("GET", "/protected", nil)
	w := testutil.PerformRequest(r, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestMiddleware_JWTAuth_InvalidFormat(t *testing.T) {
	r := testutil.SetupRouter()
	r.GET("/protected", middleware.JWTAuth(), func(c *gin.Context) {
		c.JSON(http.StatusOK, nil)
	})

	req := testutil.MakeRequest("GET", "/protected", nil)
	req.Header.Set("Authorization", "InvalidFormat token123")
	w := testutil.PerformRequest(r, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestMiddleware_JWTAuth_InvalidToken(t *testing.T) {
	r := testutil.SetupRouter()
	r.GET("/protected", middleware.JWTAuth(), func(c *gin.Context) {
		c.JSON(http.StatusOK, nil)
	})

	req := testutil.MakeRequest("GET", "/protected", nil)
	req.Header.Set("Authorization", "Bearer invalid.token.here")
	w := testutil.PerformRequest(r, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestMiddleware_OptionalJWTAuth_WithToken(t *testing.T) {
	r := testutil.SetupRouter()
	r.GET("/optional", middleware.OptionalJWTAuth(), func(c *gin.Context) {
		userID, exists := c.Get("user_id")
		if exists {
			c.JSON(http.StatusOK, gin.H{"user_id": userID})
		} else {
			c.JSON(http.StatusOK, gin.H{"user_id": nil})
		}
	})

	req := testutil.MakeAuthRequest("GET", "/optional", nil, 5, "user")
	w := testutil.PerformRequest(r, req)

	assert.Equal(t, http.StatusOK, w.Code)
	resp := testutil.ParseResponse(w)
	assert.NotNil(t, resp["user_id"])
}

func TestMiddleware_OptionalJWTAuth_WithoutToken(t *testing.T) {
	r := testutil.SetupRouter()
	r.GET("/optional", middleware.OptionalJWTAuth(), func(c *gin.Context) {
		_, exists := c.Get("user_id")
		c.JSON(http.StatusOK, gin.H{"has_user": exists})
	})

	req := testutil.MakeRequest("GET", "/optional", nil)
	w := testutil.PerformRequest(r, req)

	assert.Equal(t, http.StatusOK, w.Code)
	resp := testutil.ParseResponse(w)
	assert.Equal(t, false, resp["has_user"])
}

func TestMiddleware_Security_Headers(t *testing.T) {
	r := testutil.SetupRouter()
	r.Use(middleware.Security())
	r.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, nil)
	})

	req := testutil.MakeRequest("GET", "/test", nil)
	w := testutil.PerformRequest(r, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "nosniff", w.Header().Get("X-Content-Type-Options"))
	assert.Equal(t, "DENY", w.Header().Get("X-Frame-Options"))
	assert.Equal(t, "1; mode=block", w.Header().Get("X-XSS-Protection"))
	assert.Contains(t, w.Header().Get("Strict-Transport-Security"), "max-age=31536000")
	assert.Equal(t, "strict-origin-when-cross-origin", w.Header().Get("Referrer-Policy"))
	assert.Equal(t, "default-src 'self'", w.Header().Get("Content-Security-Policy"))
	assert.Contains(t, w.Header().Get("Permissions-Policy"), "camera=()")
}

func TestMiddleware_RateLimit(t *testing.T) {
	r := testutil.SetupRouter()
	r.Use(middleware.RateLimit(2, 2))
	r.GET("/limited", func(c *gin.Context) {
		c.JSON(http.StatusOK, nil)
	})

	for i := 0; i < 2; i++ {
		req := testutil.MakeRequest("GET", "/limited", nil)
		w := testutil.PerformRequest(r, req)
		assert.Equal(t, http.StatusOK, w.Code)
	}

	req := testutil.MakeRequest("GET", "/limited", nil)
	w := testutil.PerformRequest(r, req)
	assert.Equal(t, http.StatusTooManyRequests, w.Code)
}
