package security_test

import (
"encoding/base64"
"fmt"
"net/http"
"strings"
"sync"
"sync/atomic"
"testing"

"github.com/gin-gonic/gin"
"github.com/stretchr/testify/assert"

"mqfm-backend/internal/infrastructure/middleware"
"mqfm-backend/internal/shared/security"
"mqfm-backend/tests/testutil"
)

func TestSecurity_PasswordHashingAndUniqueness(t *testing.T) {
password := "secure_password_123"
h1, _ := security.HashPassword(password)
h2, _ := security.HashPassword(password)
assert.NotEqual(t, password, h1)
assert.NotEqual(t, h1, h2)
assert.True(t, security.CheckPassword(password, h1))
assert.True(t, security.CheckPassword(password, h2))
assert.False(t, security.CheckPassword("wrong_password", h1))
}

func TestSecurity_JWTFullLifecycle(t *testing.T) {
token, err := security.GenerateToken(1, "admin")
assert.NoError(t, err)
assert.True(t, len(strings.Split(token, ".")) == 3)
parsed, err := security.ValidateToken(token)
assert.NoError(t, err)
assert.True(t, parsed.Valid)
}

func TestSecurity_JWTTamperedAndInvalidTokens(t *testing.T) {
token, _ := security.GenerateToken(1, "admin")
bad := []string{"invalid.token.string", token[:len(token)-5] + "XXXXX", "", "Bearer", base64.StdEncoding.EncodeToString([]byte("fake")), strings.Repeat("a", 1000), "eyJ0eXAiOiJKV1QiLCJhbGciOiJub25lIn0.eyJzdWIiOiIxIn0."}
for _, tc := range bad {
_, err := security.ValidateToken(tc)
assert.Error(t, err)
}
}

func TestSecurity_JWTGetUserID(t *testing.T) {
r := testutil.SetupRouter()
r.GET("/test", middleware.JWTAuth(), func(c *gin.Context) { c.JSON(http.StatusOK, gin.H{"user_id": security.GetUserID(c)}) })
w := testutil.PerformRequest(r, testutil.MakeAuthRequest("GET", "/test", nil, 42, "user"))
assert.Equal(t, http.StatusOK, w.Code)
assert.Equal(t, float64(42), testutil.ParseResponse(w)["user_id"])
}

func TestSecurity_AllSecurityHeaders(t *testing.T) {
r := testutil.SetupRouter()
r.Use(middleware.Security())
r.GET("/test", func(c *gin.Context) { c.Status(http.StatusOK) })
w := testutil.PerformRequest(r, testutil.MakeRequest("GET", "/test", nil))
for k, v := range map[string]string{"X-Content-Type-Options": "nosniff", "X-Frame-Options": "DENY", "X-XSS-Protection": "1; mode=block", "Referrer-Policy": "strict-origin-when-cross-origin", "Content-Security-Policy": "default-src 'self'"} {
assert.Equal(t, v, w.Header().Get(k), "Header %s", k)
}
assert.Contains(t, w.Header().Get("Strict-Transport-Security"), "max-age=31536000")
assert.Contains(t, w.Header().Get("Strict-Transport-Security"), "includeSubDomains")
assert.Contains(t, w.Header().Get("Permissions-Policy"), "camera=()")
assert.Contains(t, w.Header().Get("Permissions-Policy"), "microphone=()")
}

func TestSecurity_SQLInjection_HeavyPayloads(t *testing.T) {
r := testutil.SetupRouter()
r.GET("/test", func(c *gin.Context) {
q := c.Query("q")
if strings.ContainsAny(q, "';\\\"") {
c.JSON(http.StatusBadRequest, gin.H{"error": "invalid"})
return
}
c.JSON(http.StatusOK, gin.H{"q": q})
})
payloads := []string{"'; DROP TABLE users; --", "1' OR '1'='1", "admin'--", "1; SELECT * FROM admins", "' UNION SELECT * FROM users --", "'; WAITFOR DELAY '0:0:5'; --", "1' AND 1=CONVERT(int,(SELECT TOP 1 table_name FROM information_schema.tables))--", "' OR 1=1 LIMIT 1; --", "1'; EXEC xp_cmdshell('whoami'); --", "' UNION SELECT username,password FROM users WHERE '1'='1"}
for _, p := range payloads {
w := testutil.PerformRequest(r, testutil.MakeRequest("GET", "/test?q="+p, nil))
assert.NotEqual(t, http.StatusInternalServerError, w.Code, "SQLi should not cause 500: %s", p)
}
}

func TestSecurity_XSS_HeavyPayloads(t *testing.T) {
r := testutil.SetupRouter()
r.Use(middleware.Security())
r.GET("/test", func(c *gin.Context) { c.JSON(http.StatusOK, gin.H{"d": c.Query("i")}) })
payloads := []string{"<script>alert('xss')</script>", "<img src=x onerror=alert(1)>", "javascript:alert(1)", "<svg onload=alert(1)>", "<body onload=alert(1)>", "<iframe src='javascript:alert(1)'>"}
for _, p := range payloads {
w := testutil.PerformRequest(r, testutil.MakeRequest("GET", "/test?i="+p, nil))
assert.Equal(t, "nosniff", w.Header().Get("X-Content-Type-Options"))
assert.Equal(t, "1; mode=block", w.Header().Get("X-XSS-Protection"))
}
}

func TestSecurity_AuthTokenVariations(t *testing.T) {
r := testutil.SetupRouter()
r.GET("/p", middleware.JWTAuth(), func(c *gin.Context) { c.JSON(http.StatusOK, nil) })
for _, h := range []string{"", "Bearer ", "Token abc123", "Bearer invalid.token.here", "Bearer", "Basic " + base64.StdEncoding.EncodeToString([]byte("admin:pass")), "Bearer " + strings.Repeat("A", 500)} {
req := testutil.MakeRequest("GET", "/p", nil)
if h != "" {
req.Header.Set("Authorization", h)
}
w := testutil.PerformRequest(r, req)
assert.Equal(t, http.StatusUnauthorized, w.Code)
}
}

func TestSecurity_RateLimitAndBruteForce(t *testing.T) {
r := testutil.SetupRouter()
r.Use(middleware.RateLimit(5, 5))
r.POST("/login", func(c *gin.Context) { c.JSON(http.StatusOK, nil) })
var blocked int64
var wg sync.WaitGroup
for i := 0; i < 20; i++ {
wg.Add(1)
go func() {
defer wg.Done()
w := testutil.PerformRequest(r, testutil.MakeRequest("POST", "/login", nil))
if w.Code == http.StatusTooManyRequests {
atomic.AddInt64(&blocked, 1)
}
}()
}
wg.Wait()
assert.Greater(t, blocked, int64(0), "rate limiter must block some concurrent requests")

r2 := testutil.SetupRouter()
r2.Use(middleware.RateLimit(3, 3))
r2.POST("/login", func(c *gin.Context) { c.JSON(http.StatusOK, nil) })
for i := 0; i < 3; i++ {
w := testutil.PerformRequest(r2, testutil.MakeRequest("POST", "/login", map[string]string{"password": fmt.Sprintf("wrong%d", i)}))
assert.Equal(t, http.StatusOK, w.Code)
}
w := testutil.PerformRequest(r2, testutil.MakeRequest("POST", "/login", map[string]string{"password": "wrong99"}))
assert.Equal(t, http.StatusTooManyRequests, w.Code)
}

func TestSecurity_PasswordValidation(t *testing.T) {
r := testutil.SetupRouter()
r.POST("/reg", func(c *gin.Context) {
var in struct{ Password string `json:"password" binding:"required,min=6"` }
if err := c.ShouldBindJSON(&in); err != nil {
c.JSON(http.StatusBadRequest, nil)
return
}
c.JSON(http.StatusOK, nil)
})
assert.Equal(t, http.StatusBadRequest, testutil.PerformRequest(r, testutil.MakeRequest("POST", "/reg", map[string]string{"password": "abc"})).Code)
assert.Equal(t, http.StatusOK, testutil.PerformRequest(r, testutil.MakeRequest("POST", "/reg", map[string]string{"password": "password123"})).Code)
}

func TestSecurity_LargePayloadHandling(t *testing.T) {
r := testutil.SetupRouter()
r.POST("/test", func(c *gin.Context) {
var body map[string]string
if err := c.ShouldBindJSON(&body); err != nil {
c.JSON(http.StatusBadRequest, nil)
return
}
c.JSON(http.StatusOK, nil)
})
w := testutil.PerformRequest(r, testutil.MakeRequest("POST", "/test", map[string]string{"data": strings.Repeat("A", 100_000)}))
assert.NotEqual(t, http.StatusInternalServerError, w.Code)
}
