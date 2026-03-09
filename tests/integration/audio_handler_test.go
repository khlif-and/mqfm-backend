package integration_test

import (
	"errors"
	"net/http"
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

func setupAudioRouter(audioSvc *mocks.MockAudioService, historySvc *mocks.MockHistoryService) *gin.Engine {
	r := testutil.SetupRouter()
	handler := admin.NewAudioHandler(audioSvc, historySvc)

	r.GET("/audios", handler.FindAll)
	r.GET("/audios/:id", middleware.OptionalJWTAuth(), handler.FindByID)
	r.DELETE("/audios/:id", handler.Delete)
	r.GET("/audios/search", handler.Search)
	return r
}

func TestAudioHandler_FindAll_Success(t *testing.T) {
	audioSvc := &mocks.MockAudioService{
		FindAllFn: func() ([]entity.Audio, error) {
			return []entity.Audio{
				{ID: 1, Title: "Kajian 1", Artist: "Ustadz A", Status: "active",
					CreatedAt: time.Now(), UpdatedAt: time.Now()},
			}, nil
		},
	}

	r := setupAudioRouter(audioSvc, &mocks.MockHistoryService{})
	req := testutil.MakeRequest("GET", "/audios", nil)
	w := testutil.PerformRequest(r, req)

	assert.Equal(t, http.StatusOK, w.Code)
	resp := testutil.ParseResponse(w)
	data := resp["data"].([]interface{})
	assert.Len(t, data, 1)
}

func TestAudioHandler_FindAll_Error(t *testing.T) {
	audioSvc := &mocks.MockAudioService{
		FindAllFn: func() ([]entity.Audio, error) { return nil, errors.New("db error") },
	}

	r := setupAudioRouter(audioSvc, &mocks.MockHistoryService{})
	req := testutil.MakeRequest("GET", "/audios", nil)
	w := testutil.PerformRequest(r, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestAudioHandler_FindByID_Success(t *testing.T) {
	audioSvc := &mocks.MockAudioService{
		FindByIDFn: func(id uint) (*entity.Audio, error) {
			return &entity.Audio{
				ID: id, Title: "Kajian 1", Artist: "Ustadz A",
				CreatedAt: time.Now(), UpdatedAt: time.Now(),
			}, nil
		},
	}
	historySvc := &mocks.MockHistoryService{
		RecordPlayFn: func(userID uint, req request.HistoryRequest) error { return nil },
	}

	r := setupAudioRouter(audioSvc, historySvc)
	req := testutil.MakeAuthRequest("GET", "/audios/1", nil, 1, "user")
	w := testutil.PerformRequest(r, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestAudioHandler_FindByID_NotFound(t *testing.T) {
	audioSvc := &mocks.MockAudioService{
		FindByIDFn: func(id uint) (*entity.Audio, error) { return nil, errors.New("not found") },
	}

	r := setupAudioRouter(audioSvc, &mocks.MockHistoryService{})
	req := testutil.MakeRequest("GET", "/audios/999", nil)
	w := testutil.PerformRequest(r, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestAudioHandler_FindByID_InvalidID(t *testing.T) {
	audioSvc := &mocks.MockAudioService{}
	r := setupAudioRouter(audioSvc, &mocks.MockHistoryService{})

	req := testutil.MakeRequest("GET", "/audios/abc", nil)
	w := testutil.PerformRequest(r, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestAudioHandler_Delete_Success(t *testing.T) {
	audioSvc := &mocks.MockAudioService{
		DeleteFn: func(id uint) error { return nil },
	}

	r := setupAudioRouter(audioSvc, &mocks.MockHistoryService{})
	req := testutil.MakeRequest("DELETE", "/audios/1", nil)
	w := testutil.PerformRequest(r, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestAudioHandler_Delete_Error(t *testing.T) {
	audioSvc := &mocks.MockAudioService{
		DeleteFn: func(id uint) error { return errors.New("cannot delete") },
	}

	r := setupAudioRouter(audioSvc, &mocks.MockHistoryService{})
	req := testutil.MakeRequest("DELETE", "/audios/1", nil)
	w := testutil.PerformRequest(r, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestAudioHandler_Search_Success(t *testing.T) {
	audioSvc := &mocks.MockAudioService{
		SearchFn: func(query string) ([]entity.Audio, error) {
			return []entity.Audio{
				{ID: 1, Title: "Kajian Fiqih", CreatedAt: time.Now(), UpdatedAt: time.Now()},
			}, nil
		},
	}

	r := setupAudioRouter(audioSvc, &mocks.MockHistoryService{})
	req := testutil.MakeRequest("GET", "/audios/search?q=Fiqih", nil)
	w := testutil.PerformRequest(r, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestAudioHandler_Search_EmptyQuery(t *testing.T) {
	audioSvc := &mocks.MockAudioService{}
	r := setupAudioRouter(audioSvc, &mocks.MockHistoryService{})

	req := testutil.MakeRequest("GET", "/audios/search", nil)
	w := testutil.PerformRequest(r, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}
