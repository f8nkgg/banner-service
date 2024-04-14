package tests

import (
	v1 "banner/internal/controller/http/v1"
	"context"
	"errors"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
	"net/http/httptest"
)

func (s *APITestSuite) TestBannerGet_Success() {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	v1.RegisterRoutes(router, s.handler)
	r := s.Require()
	s.createTestBanner()
	defer s.deleteTestBanner()
	req, _ := http.NewRequest("GET", "/user_banner?tag_id=4&feature_id=123", nil)
	req.Header.Set("Content-type", "application/json")
	req.Header.Set("token", "user_token")
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)
	r.Equal(http.StatusOK, resp.Result().StatusCode)
}
func (s *APITestSuite) TestBannerGet_BadRequest() {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	v1.RegisterRoutes(router, s.handler)
	r := s.Require()
	req, _ := http.NewRequest("GET", "/user_banner?tag_id=stashge", nil)
	req.Header.Set("Content-type", "application/json")
	req.Header.Set("token", "user_token")
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)
	r.Equal(http.StatusBadRequest, resp.Result().StatusCode)
	responseBody, err := io.ReadAll(resp.Body)
	s.NoError(err)
	r.Equal("{\"message\":\"Некорректные данные tagId\"}", string(responseBody))
}
func (s *APITestSuite) TestBannerGet_Unauthorized() {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	v1.RegisterRoutes(router, s.handler)
	r := s.Require()
	req, _ := http.NewRequest("GET", "/user_banner?tag_id=4&feature_id=123", nil)
	req.Header.Set("Content-type", "application/json")
	req.Header.Set("token", "")
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)
	r.Equal(http.StatusUnauthorized, resp.Result().StatusCode)
}

func (s *APITestSuite) TestBannerGet_NotFound() {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	v1.RegisterRoutes(router, s.handler)
	r := s.Require()
	req, _ := http.NewRequest("GET", "/user_banner?tag_id=9999&feature_id=9999", nil)
	req.Header.Set("Content-type", "application/json")
	req.Header.Set("token", "user_token")
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)
	r.Equal(http.StatusNotFound, resp.Result().StatusCode)
}

func (s *APITestSuite) TestBannerGet_InternalServerError() {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	mockService := &MockBannerService{
		GetForUserFunc: func(ctx context.Context, tagID, featureID int32, isActiveParam, lastRevision bool) (map[string]interface{}, error) {
			return nil, errors.New("internal server error")
		},
	}
	v1.RegisterRoutes(router, v1.NewBannerController(mockService, s.logger))
	r := s.Require()
	s.createTestBanner()
	defer s.deleteTestBanner()
	req, _ := http.NewRequest("GET", "/user_banner?tag_id=4&feature_id=124", nil)
	req.Header.Set("Content-type", "application/json")
	req.Header.Set("token", "user_token")
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)
	r.Equal(http.StatusInternalServerError, resp.Result().StatusCode)
	responseBody, err := io.ReadAll(resp.Body)
	s.NoError(err)
	r.Equal("{\"error\":\"Внутренняя ошибка сервера\"}", string(responseBody))
}
