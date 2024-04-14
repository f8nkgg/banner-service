package tests

import (
	v1 "banner/internal/controller/http/v1"
	"banner/internal/entity"
	"context"
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
	"net/http/httptest"
	"strings"
)

func (s *APITestSuite) TestUpdateBanner_Success() {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	v1.RegisterRoutes(router, s.handler)
	r := s.Require()
	s.createTestBanner()
	defer s.deleteTestBanner()
	requestBody := `{
		"tag_ids": [7, 8, 9]
	}`
	req, _ := http.NewRequest("PATCH", "/banner/1", strings.NewReader(requestBody))
	req.Header.Set("Content-type", "application/json")
	req.Header.Set("token", "admin_token")
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	r.Equal(http.StatusNoContent, resp.Result().StatusCode)
}
func (s *APITestSuite) TestUpdateBanner_Unauthorized() {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	v1.RegisterRoutes(router, s.handler)
	r := s.Require()

	req, _ := http.NewRequest("PATCH", "/banner/1", nil)
	req.Header.Set("Content-type", "application/json")
	req.Header.Set("token", "")
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	r.Equal(http.StatusUnauthorized, resp.Result().StatusCode)
}
func (s *APITestSuite) TestUpdateBanner_Forbidden() {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	v1.RegisterRoutes(router, s.handler)
	r := s.Require()

	req, _ := http.NewRequest("PATCH", "/banner/1", nil)
	req.Header.Set("Content-type", "application/json")
	req.Header.Set("token", "user_token")
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	r.Equal(http.StatusForbidden, resp.Result().StatusCode)
}

func (s *APITestSuite) TestUpdateBanner_BadRequest() {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	v1.RegisterRoutes(router, s.handler)
	r := s.Require()

	req, _ := http.NewRequest("PATCH", "/banner/abc", nil)
	req.Header.Set("Content-type", "application/json")
	req.Header.Set("token", "admin_token")
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	r.Equal(http.StatusBadRequest, resp.Result().StatusCode)
}

func (s *APITestSuite) TestUpdateBanner_NotFound() {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	v1.RegisterRoutes(router, s.handler)
	r := s.Require()
	requestBody := `{
		"tag_ids": [7, 8, 9]
	}`
	req, _ := http.NewRequest("PATCH", "/banner/9999", strings.NewReader(requestBody))
	req.Header.Set("Content-type", "application/json")
	req.Header.Set("token", "admin_token")
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	r.Equal(http.StatusNotFound, resp.Result().StatusCode)
}

func (s *APITestSuite) TestUpdateBanner_InternalServerError() {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	mockService := &MockBannerService{
		UpdateFunc: func(ctx context.Context, banner *entity.BannerUpdate) error {
			return errors.New("internal server error")
		},
	}
	v1.RegisterRoutes(router, v1.NewBannerController(mockService, s.logger))
	r := s.Require()
	requestBody := `{
		"tag_ids": [7, 8, 9]
	}`
	req, _ := http.NewRequest("PATCH", "/banner/1", strings.NewReader(requestBody))
	req.Header.Set("Content-type", "application/json")
	req.Header.Set("token", "admin_token")
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	r.Equal(http.StatusInternalServerError, resp.Result().StatusCode)
}
