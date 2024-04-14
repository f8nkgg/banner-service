package tests

import (
	v1 "banner/internal/controller/http/v1"
	"banner/internal/entity"
	"context"
	"errors"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
)

func (s *APITestSuite) TestCreateBanner_Success() {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	v1.RegisterRoutes(router, s.handler)
	r := s.Require()
	requestBody := `{
		"tag_ids": [4, 5, 6],
		"feature_id": 123,
		"content": {
			"text": "some_text",
			"title": "some_title",
			"url": "some_url"
		},
		"is_active": true
	}`
	req, _ := http.NewRequest("POST", "/banner", strings.NewReader(requestBody))
	req.Header.Set("Content-type", "application/json")
	req.Header.Set("token", "admin_token")
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)
	r.Equal(http.StatusCreated, resp.Result().StatusCode)
	responseBody, err := io.ReadAll(resp.Body)
	s.NoError(err)
	r.Equal("{\"banner_id\":1}", string(responseBody))
	defer func() {
		_, err := s.db.Pool.Exec(context.Background(), "DELETE FROM banners WHERE id = $1", 1)
		s.NoError(err)
	}()
}
func (s *APITestSuite) TestCreateBanner_Unauthorized() {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	v1.RegisterRoutes(router, s.handler)
	r := s.Require()
	requestBody := `{
		"tag_ids": [4, 5, 6],
		"feature_id": 123,
		"content": {
			"text": "some_text",
			"title": "some_title",
			"url": "some_url"
		},
		"is_active": true
	}`
	req, _ := http.NewRequest("POST", "/banner", strings.NewReader(requestBody))
	req.Header.Set("Content-type", "application/json")
	req.Header.Set("token", "")
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)
	r.Equal(http.StatusUnauthorized, resp.Result().StatusCode)
}
func (s *APITestSuite) TestCreateBanner_Forbidden() {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	v1.RegisterRoutes(router, s.handler)
	r := s.Require()
	requestBody := `{
		"tag_ids": [4, 5, 6],
		"feature_id": 123,
		"content": {
			"text": "some_text",
			"title": "some_title",
			"url": "some_url"
		},
		"is_active": true
	}`
	req, _ := http.NewRequest("POST", "/banner", strings.NewReader(requestBody))
	req.Header.Set("Content-type", "application/json")
	req.Header.Set("token", "user_token")
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)
	r.Equal(http.StatusForbidden, resp.Result().StatusCode)
}
func (s *APITestSuite) TestCreateBanner_InternalServerError() {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	mockService := &MockBannerService{
		SaveFunc: func(ctx context.Context, banner *entity.Banner) (int32, error) {
			return -1, errors.New("internal server error")
		},
	}
	v1.RegisterRoutes(router, v1.NewBannerController(mockService, s.logger))
	r := s.Require()
	requestBody := `{
		"tag_ids": [4, 5, 6],
		"feature_id": 123,
		"content": {
			"text": "some_text",
			"title": "some_title",
			"url": "some_url"
		},
		"is_active": true
	}`
	req, _ := http.NewRequest("POST", "/banner", strings.NewReader(requestBody))
	req.Header.Set("Content-type", "application/json")
	req.Header.Set("token", "admin_token")
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)
	r.Equal(http.StatusInternalServerError, resp.Result().StatusCode)
	responseBody, err := io.ReadAll(resp.Body)
	s.NoError(err)
	r.Equal("{\"error\":\"Внутренняя ошибка сервера\"}", string(responseBody))
}
