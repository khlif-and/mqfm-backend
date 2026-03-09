package testutil

import (
	"bytes"
	"encoding/json"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"strings"

	"github.com/gin-gonic/gin"

	"mqfm-backend/internal/shared/security"
)

func init() {
	gin.SetMode(gin.TestMode)
}

func SetupRouter() *gin.Engine {
	return gin.New()
}

func MakeRequest(method, url string, body interface{}) *http.Request {
	var reader io.Reader
	if body != nil {
		jsonBytes, _ := json.Marshal(body)
		reader = bytes.NewBuffer(jsonBytes)
	}
	req, _ := http.NewRequest(method, url, reader)
	req.Header.Set("Content-Type", "application/json")
	return req
}

func MakeAuthRequest(method, url string, body interface{}, userID uint, role string) *http.Request {
	req := MakeRequest(method, url, body)
	token, _ := security.GenerateToken(userID, role)
	req.Header.Set("Authorization", "Bearer "+token)
	return req
}

func PerformRequest(r *gin.Engine, req *http.Request) *httptest.ResponseRecorder {
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}

func ParseResponse(w *httptest.ResponseRecorder) map[string]interface{} {
	var response map[string]interface{}
	_ = json.Unmarshal(w.Body.Bytes(), &response)
	return response
}

func MakeMultipartRequest(method, url string, fields map[string]string) (*http.Request, error) {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	for key, val := range fields {
		_ = writer.WriteField(key, val)
	}
	_ = writer.Close()
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())
	return req, nil
}

func MakeAuthMultipartRequest(method, url string, fields map[string]string, userID uint, role string) (*http.Request, error) {
	req, err := MakeMultipartRequest(method, url, fields)
	if err != nil {
		return nil, err
	}
	token, _ := security.GenerateToken(userID, role)
	req.Header.Set("Authorization", "Bearer "+token)
	return req, nil
}

func MakeFormRequest(method, url string, fields map[string]string) *http.Request {
	form := make([]string, 0)
	for k, v := range fields {
		form = append(form, k+"="+v)
	}
	body := strings.Join(form, "&")
	req, _ := http.NewRequest(method, url, strings.NewReader(body))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	return req
}
