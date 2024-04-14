package tests

import (
	v1 "banner/internal/controller/http/v1"
	"context"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"net/http/httptest"
)

func (s *APITestSuite) TestDeleteBanner_Success() {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	v1.RegisterRoutes(router, s.handler)
	r := s.Require()

	s.createTestBanner()
	defer s.deleteTestBanner()
	req, _ := http.NewRequest("DELETE", fmt.Sprintf("/banner/%d", 1), nil)
	req.Header.Set("Content-type", "application/json")
	req.Header.Set("token", "admin_token")
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	r.Equal(http.StatusNoContent, resp.Result().StatusCode)
}

func (s *APITestSuite) TestDeleteBanner_Unauthorized() {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	v1.RegisterRoutes(router, s.handler)
	r := s.Require()

	req, _ := http.NewRequest("DELETE", "/banner/1", nil)
	req.Header.Set("Content-type", "application/json")
	req.Header.Set("token", "")
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	r.Equal(http.StatusUnauthorized, resp.Result().StatusCode)
}
func (s *APITestSuite) TestDeleteBanner_Forbidden() {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	v1.RegisterRoutes(router, s.handler)
	r := s.Require()

	req, _ := http.NewRequest("DELETE", "/banner/1", nil)
	req.Header.Set("Content-type", "application/json")
	req.Header.Set("token", "user_token")
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	r.Equal(http.StatusForbidden, resp.Result().StatusCode)
}
func (s *APITestSuite) TestDeleteBanner_BadRequest() {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	v1.RegisterRoutes(router, s.handler)
	r := s.Require()

	req, _ := http.NewRequest("DELETE", "/banner/abc", nil)
	req.Header.Set("Content-type", "application/json")
	req.Header.Set("token", "admin_token")
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	r.Equal(http.StatusBadRequest, resp.Result().StatusCode)
}

func (s *APITestSuite) TestDeleteBanner_NotFound() {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	v1.RegisterRoutes(router, s.handler)
	r := s.Require()
	req, _ := http.NewRequest("DELETE", "/banner/9999", nil)
	req.Header.Set("Content-type", "application/json")
	req.Header.Set("token", "admin_token")
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	r.Equal(http.StatusNotFound, resp.Result().StatusCode)
}

func (s *APITestSuite) TestDeleteBanner_InternalServerError() {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	mockService := &MockBannerService{
		DeleteFunc: func(ctx context.Context, id int32) error {
			return errors.New("internal server error")
		},
	}
	v1.RegisterRoutes(router, v1.NewBannerController(mockService, s.logger))
	r := s.Require()

	req, _ := http.NewRequest("DELETE", "/banner/1", nil)
	req.Header.Set("Content-type", "application/json")
	req.Header.Set("token", "admin_token")
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	r.Equal(http.StatusInternalServerError, resp.Result().StatusCode)
}
