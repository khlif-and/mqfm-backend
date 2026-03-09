package performance_test

import (
	"fmt"
	"net/http"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"

	admin "mqfm-backend/internal/adapter/handler/admin"
	"mqfm-backend/internal/domain/entity"
	"mqfm-backend/internal/infrastructure/middleware"
	"mqfm-backend/internal/shared/dto/request"
	"mqfm-backend/tests/mocks"
	"mqfm-backend/tests/testutil"
)

func setupPerfRouter() *gin.Engine {
	r := testutil.SetupRouter()
	audioSvc := &mocks.MockAudioService{
		FindAllFn:  func() ([]entity.Audio, error) { return []entity.Audio{{ID: 1, Title: "A1", Status: "active", CreatedAt: time.Now(), UpdatedAt: time.Now()}, {ID: 2, Title: "A2", Status: "active", CreatedAt: time.Now(), UpdatedAt: time.Now()}}, nil },
		FindByIDFn: func(id uint) (*entity.Audio, error) { return &entity.Audio{ID: id, Title: "Audio", Artist: "A", CreatedAt: time.Now(), UpdatedAt: time.Now()}, nil },
		SearchFn:   func(q string) ([]entity.Audio, error) { return []entity.Audio{{ID: 1, Title: q, CreatedAt: time.Now(), UpdatedAt: time.Now()}}, nil },
	}
	historySvc := &mocks.MockHistoryService{RecordPlayFn: func(u uint, r request.HistoryRequest) error { return nil }}
	adminSvc := &mocks.MockAdminAuthService{
		LoginFn:   func(r request.AdminLoginRequest) (string, *entity.Admin, error) { return "token", &entity.Admin{ID: 1, Email: r.Email, Role: "admin", CreatedAt: time.Now(), UpdatedAt: time.Now()}, nil },
		GetByIDFn: func(id uint) (*entity.Admin, error) { return &entity.Admin{ID: id, Username: "admin", CreatedAt: time.Now(), UpdatedAt: time.Now()}, nil },
	}
	ah := admin.NewAudioHandler(audioSvc, historySvc)
	auth := admin.NewAuthHandler(adminSvc)
	r.GET("/audios", ah.FindAll)
	r.GET("/audios/:id", middleware.OptionalJWTAuth(), ah.FindByID)
	r.GET("/audios/search", ah.Search)
	r.POST("/admin/login", auth.Login)
	r.GET("/admin/me", middleware.JWTAuth(), auth.Me)
	return r
}

func runConcurrent(r *gin.Engine, n int, mkReq func(int) *http.Request) (int64, time.Duration) {
	var wg sync.WaitGroup
	var ok int64
	start := time.Now()
	for i := 0; i < n; i++ {
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()
			if testutil.PerformRequest(r, mkReq(idx)).Code == http.StatusOK { atomic.AddInt64(&ok, 1) }
		}(i)
	}
	wg.Wait()
	return ok, time.Since(start)
}

func TestPerformance_500ConcurrentReads(t *testing.T) {
	r := setupPerfRouter()
	ok, elapsed := runConcurrent(r, 500, func(_ int) *http.Request { return testutil.MakeRequest("GET", "/audios", nil) })
	assert.Equal(t, int64(500), ok)
	assert.Less(t, elapsed, 10*time.Second)
	t.Logf("500 concurrent reads: %v (%.0f req/s)", elapsed, 500.0/elapsed.Seconds())
}

func TestPerformance_200ConcurrentAuth(t *testing.T) {
	r := setupPerfRouter()
	ok, elapsed := runConcurrent(r, 200, func(i int) *http.Request {
		return testutil.MakeRequest("POST", "/admin/login", map[string]string{"email": fmt.Sprintf("a%d@t.com", i), "password": "p"})
	})
	assert.Equal(t, int64(200), ok)
	t.Logf("200 concurrent auth: %v (%.0f req/s)", elapsed, 200.0/elapsed.Seconds())
}

func TestPerformance_300AuthenticatedEndpoint(t *testing.T) {
	r := setupPerfRouter()
	ok, elapsed := runConcurrent(r, 300, func(i int) *http.Request {
		return testutil.MakeAuthRequest("GET", "/admin/me", nil, uint(i+1), "admin")
	})
	assert.Equal(t, int64(300), ok)
	t.Logf("300 concurrent auth endpoint: %v (%.0f req/s)", elapsed, 300.0/elapsed.Seconds())
}

func TestPerformance_MixedWorkload(t *testing.T) {
	r := setupPerfRouter()
	var wg sync.WaitGroup
	var reads, auths, searches int64
	total := 300
	start := time.Now()
	for i := 0; i < total; i++ {
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()
			switch idx % 3 {
			case 0:
				if testutil.PerformRequest(r, testutil.MakeRequest("GET", "/audios", nil)).Code == http.StatusOK { atomic.AddInt64(&reads, 1) }
			case 1:
				if testutil.PerformRequest(r, testutil.MakeAuthRequest("GET", "/admin/me", nil, uint(idx+1), "admin")).Code == http.StatusOK { atomic.AddInt64(&auths, 1) }
			case 2:
				if testutil.PerformRequest(r, testutil.MakeRequest("GET", "/audios/search?q=fiqih", nil)).Code == http.StatusOK { atomic.AddInt64(&searches, 1) }
			}
		}(i)
	}
	wg.Wait()
	elapsed := time.Since(start)
	assert.Equal(t, int64(100), reads)
	assert.Equal(t, int64(100), auths)
	assert.Equal(t, int64(100), searches)
	t.Logf("Mixed %d reqs: %v (%.0f req/s) [reads=%d auths=%d searches=%d]", total, elapsed, float64(total)/elapsed.Seconds(), reads, auths, searches)
}

func TestPerformance_BurstTraffic(t *testing.T) {
	r := setupPerfRouter()
	for burst := 0; burst < 3; burst++ {
		ok, elapsed := runConcurrent(r, 100, func(_ int) *http.Request { return testutil.MakeRequest("GET", "/audios", nil) })
		assert.Equal(t, int64(100), ok)
		t.Logf("Burst %d: 100 reqs in %v", burst+1, elapsed)
		time.Sleep(50 * time.Millisecond)
	}
}

func TestPerformance_SequentialVsConcurrent(t *testing.T) {
	r := setupPerfRouter()
	n := 100
	start := time.Now()
	for i := 0; i < n; i++ {
		testutil.PerformRequest(r, testutil.MakeRequest("GET", "/audios", nil))
	}
	seqTime := time.Since(start)

	ok, concTime := runConcurrent(r, n, func(_ int) *http.Request { return testutil.MakeRequest("GET", "/audios", nil) })
	assert.Equal(t, int64(n), ok)
	t.Logf("Sequential: %v | Concurrent: %v | Speedup: %.1fx", seqTime, concTime, float64(seqTime)/float64(concTime))
}

func BenchmarkAudioFindAll(b *testing.B) {
	r := setupPerfRouter()
	for i := 0; i < b.N; i++ { testutil.PerformRequest(r, testutil.MakeRequest("GET", "/audios", nil)) }
}

func BenchmarkAudioFindByID(b *testing.B) {
	r := setupPerfRouter()
	for i := 0; i < b.N; i++ { testutil.PerformRequest(r, testutil.MakeRequest("GET", "/audios/1", nil)) }
}

func BenchmarkAdminLogin(b *testing.B) {
	r := setupPerfRouter()
	for i := 0; i < b.N; i++ { testutil.PerformRequest(r, testutil.MakeRequest("POST", "/admin/login", map[string]string{"email": "a@t.com", "password": "p"})) }
}

func BenchmarkAuthenticatedMe(b *testing.B) {
	r := setupPerfRouter()
	for i := 0; i < b.N; i++ { testutil.PerformRequest(r, testutil.MakeAuthRequest("GET", "/admin/me", nil, 1, "admin")) }
}
