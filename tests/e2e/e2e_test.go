package e2e_test

import (
"errors"
"mime/multipart"
"net/http"
"strings"
"testing"
"time"

"github.com/gin-gonic/gin"
"github.com/stretchr/testify/assert"

admin "mqfm-backend/internal/adapter/handler/admin"
"mqfm-backend/internal/domain/entity"
"mqfm-backend/internal/infrastructure/middleware"
"mqfm-backend/internal/shared/dto/request"
audiomock "mqfm-backend/tests/mocks/audio"
authmock "mqfm-backend/tests/mocks/auth"
categorymock "mqfm-backend/tests/mocks/category"
historymock "mqfm-backend/tests/mocks/history"
"mqfm-backend/tests/testutil"
)

type e2eState struct {
categories map[uint]*entity.Category
audios     map[uint]*entity.Audio
nextCatID  uint
nextAudID  uint
}

func setupE2ERouter() (*gin.Engine, *e2eState) {
st := &e2eState{categories: make(map[uint]*entity.Category), audios: make(map[uint]*entity.Audio), nextCatID: 1, nextAudID: 1}

adminSvc := &authmock.MockAdminAuthService{
RegisterFn: func(req request.AdminRegisterRequest) (*entity.Admin, error) {
return &entity.Admin{ID: 1, Username: req.Username, Email: req.Email, Role: "admin", CreatedAt: time.Now(), UpdatedAt: time.Now()}, nil
},
LoginFn: func(req request.AdminLoginRequest) (string, *entity.Admin, error) {
return "e2e-jwt-token", &entity.Admin{ID: 1, Email: req.Email, Role: "admin", CreatedAt: time.Now(), UpdatedAt: time.Now()}, nil
},
GetByIDFn: func(id uint) (*entity.Admin, error) {
return &entity.Admin{ID: id, Username: "admin", Email: "admin@test.com", Role: "admin", CreatedAt: time.Now(), UpdatedAt: time.Now()}, nil
},
}
catSvc := &categorymock.MockCategoryService{
CreateFn: func(req request.CreateCategoryRequest, f *multipart.FileHeader) (*entity.Category, error) {
c := &entity.Category{ID: st.nextCatID, Name: req.Name, CreatedAt: time.Now(), UpdatedAt: time.Now()}
st.categories[st.nextCatID] = c
st.nextCatID++
return c, nil
},
FindAllFn: func() ([]entity.Category, error) {
r := make([]entity.Category, 0, len(st.categories))
for _, c := range st.categories { r = append(r, *c) }
return r, nil
},
FindByIDFn: func(id uint) (*entity.Category, error) {
c, ok := st.categories[id]; if !ok { return nil, errors.New("not found") }; return c, nil
},
UpdateFn: func(id uint, req request.UpdateCategoryRequest, f *multipart.FileHeader) (*entity.Category, error) {
c, ok := st.categories[id]; if !ok { return nil, errors.New("not found") }
if req.Name != "" { c.Name = req.Name }
return c, nil
},
DeleteFn: func(id uint) error { delete(st.categories, id); return nil },
SearchFn: func(q string) ([]entity.Category, error) {
var r []entity.Category
for _, c := range st.categories { if strings.Contains(c.Name, q) { r = append(r, *c) } }
return r, nil
},
}
audioSvc := &audiomock.MockAudioService{
FindAllFn:  func() ([]entity.Audio, error) { return nil, nil },
FindByIDFn: func(id uint) (*entity.Audio, error) { a, ok := st.audios[id]; if !ok { return nil, errors.New("not found") }; return a, nil },
DeleteFn:   func(id uint) error { delete(st.audios, id); return nil },
SearchFn:   func(q string) ([]entity.Audio, error) { return nil, nil },
}
historySvc := &historymock.MockHistoryService{RecordPlayFn: func(u uint, r request.HistoryRequest) error { return nil }}

r := testutil.SetupRouter()
r.Use(middleware.Security())
ah := admin.NewAuthHandler(adminSvc)
ch := admin.NewCategoryHandler(catSvc)
auh := admin.NewAudioHandler(audioSvc, historySvc)
r.POST("/admin/register", ah.Register)
r.POST("/admin/login", ah.Login)
r.GET("/admin/me", middleware.JWTAuth(), ah.Me)
r.POST("/admin/logout", middleware.JWTAuth(), ah.Logout)
cg := r.Group("/admin/categories", middleware.JWTAuth())
cg.POST("", ch.Create); cg.GET("", ch.FindAll); cg.GET("/:id", ch.FindByID)
cg.PUT("/:id", ch.Update); cg.DELETE("/:id", ch.Delete); cg.GET("/search", ch.Search)
ag := r.Group("/admin/audios", middleware.JWTAuth())
ag.GET("", auh.FindAll); ag.GET("/:id", auh.FindByID); ag.DELETE("/:id", auh.Delete)
return r, st
}

func TestE2E_FullAdminWorkflow(t *testing.T) {
r, _ := setupE2ERouter()

w := testutil.PerformRequest(r, testutil.MakeRequest("POST", "/admin/register", map[string]string{"username": "a", "email": "a@t.com", "password": "password123"}))
assert.Equal(t, http.StatusCreated, w.Code)

w = testutil.PerformRequest(r, testutil.MakeRequest("POST", "/admin/login", map[string]string{"email": "a@t.com", "password": "password123"}))
assert.Equal(t, http.StatusOK, w.Code)
assert.NotEmpty(t, testutil.ParseResponse(w)["data"].(map[string]interface{})["token"])

w = testutil.PerformRequest(r, testutil.MakeAuthRequest("GET", "/admin/me", nil, 1, "admin"))
assert.Equal(t, http.StatusOK, w.Code)

w = testutil.PerformRequest(r, testutil.MakeAuthRequest("POST", "/admin/logout", nil, 1, "admin"))
assert.Equal(t, http.StatusOK, w.Code)
}

func TestE2E_CategoryCRUD(t *testing.T) {
r, _ := setupE2ERouter()
mk := func(m, u string, b interface{}) *http.Request { return testutil.MakeAuthRequest(m, u, b, 1, "admin") }

req, _ := testutil.MakeAuthMultipartRequest("POST", "/admin/categories", map[string]string{"name": "Fiqih"}, 1, "admin")
w := testutil.PerformRequest(r, req)
assert.Equal(t, http.StatusCreated, w.Code)
assert.Equal(t, "Fiqih", testutil.ParseResponse(w)["data"].(map[string]interface{})["name"])

w = testutil.PerformRequest(r, mk("GET", "/admin/categories", nil))
assert.Equal(t, http.StatusOK, w.Code)

w = testutil.PerformRequest(r, mk("GET", "/admin/categories/1", nil))
assert.Equal(t, http.StatusOK, w.Code)

req, _ = testutil.MakeAuthMultipartRequest("PUT", "/admin/categories/1", map[string]string{"name": "Updated"}, 1, "admin")
w = testutil.PerformRequest(r, req)
assert.Equal(t, http.StatusOK, w.Code)
assert.Equal(t, "Updated", testutil.ParseResponse(w)["data"].(map[string]interface{})["name"])

w = testutil.PerformRequest(r, mk("DELETE", "/admin/categories/1", nil))
assert.Equal(t, http.StatusOK, w.Code)

w = testutil.PerformRequest(r, mk("GET", "/admin/categories/1", nil))
assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestE2E_UnauthorizedAccess(t *testing.T) {
r, _ := setupE2ERouter()
for _, ep := range []struct{ m, p string }{{"GET", "/admin/me"}, {"GET", "/admin/categories"}, {"POST", "/admin/categories"}, {"GET", "/admin/audios"}} {
w := testutil.PerformRequest(r, testutil.MakeRequest(ep.m, ep.p, nil))
assert.Equal(t, http.StatusUnauthorized, w.Code, "%s %s should require auth", ep.m, ep.p)
}
}

func TestE2E_SecurityHeaders(t *testing.T) {
r, _ := setupE2ERouter()
w := testutil.PerformRequest(r, testutil.MakeRequest("POST", "/admin/login", map[string]string{"email": "a@t.com", "password": "p"}))
assert.Equal(t, "nosniff", w.Header().Get("X-Content-Type-Options"))
assert.Equal(t, "DENY", w.Header().Get("X-Frame-Options"))
assert.Contains(t, w.Header().Get("Strict-Transport-Security"), "max-age=31536000")
}

func TestE2E_InputValidation(t *testing.T) {
r, _ := setupE2ERouter()
cases := []struct{ body map[string]string; code int }{
{map[string]string{"username": "a", "email": "not-email", "password": "password123"}, http.StatusBadRequest},
{map[string]string{"username": "a", "email": "a@t.com", "password": "ab"}, http.StatusBadRequest},
{map[string]string{}, http.StatusBadRequest},
}
for _, tc := range cases {
w := testutil.PerformRequest(r, testutil.MakeRequest("POST", "/admin/register", tc.body))
assert.Equal(t, tc.code, w.Code)
}
}

func TestE2E_ResponseFormat(t *testing.T) {
r, _ := setupE2ERouter()
w := testutil.PerformRequest(r, testutil.MakeRequest("POST", "/admin/login", map[string]string{"email": "a@t.com", "password": "p"}))
resp := testutil.ParseResponse(w)
_, s := resp["status"]; _, m := resp["message"]; _, d := resp["data"]
assert.True(t, s && m && d)
}
