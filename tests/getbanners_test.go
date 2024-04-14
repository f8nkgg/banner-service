package tests

import (
	v1 "banner/internal/controller/http/v1"
	"banner/internal/entity"
	"context"
	"encoding/json"
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
	"net/http/httptest"
)

func (s *APITestSuite) TestGetBanners_Success() {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	v1.RegisterRoutes(router, s.handler)
	r := s.Require()
	s.createTestBanner()
	defer s.deleteTestBanner()
	req, _ := http.NewRequest("GET", "/banner?feature_id=123&tag_id=4&limit=10&offset=0", nil)
	req.Header.Set("Content-type", "application/json")
	req.Header.Set("token", "admin_token")
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	r.Equal(http.StatusOK, resp.Result().StatusCode)
	var banners []entity.FilteredBanner
	err := json.Unmarshal(resp.Body.Bytes(), &banners)
	s.NoError(err)
	r.NotEmpty(banners)
}
func (s *APITestSuite) TestGetBanners_Unauthorized() {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	v1.RegisterRoutes(router, s.handler)
	r := s.Require()

	req, _ := http.NewRequest("GET", "/banner", nil)
	req.Header.Set("Content-type", "application/json")
	req.Header.Set("token", "")
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	r.Equal(http.StatusUnauthorized, resp.Result().StatusCode)
}

func (s *APITestSuite) TestGetBanners_Forbidden() {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	v1.RegisterRoutes(router, s.handler)
	r := s.Require()

	req, _ := http.NewRequest("GET", "/banner", nil)
	req.Header.Set("Content-type", "application/json")
	req.Header.Set("token", "user_token")
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	r.Equal(http.StatusForbidden, resp.Result().StatusCode)
}

func (s *APITestSuite) TestGetBanners_InternalServerError() {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	mockService := &MockBannerService{
		GetBannersFunc: func(ctx context.Context, featureID, tagID, limit, offset *int32) ([]*entity.FilteredBanner, error) {
			return nil, errors.New("internal server error")
		},
	}
	v1.RegisterRoutes(router, v1.NewBannerController(mockService, s.logger))
	r := s.Require()

	req, _ := http.NewRequest("GET", "/banner", nil)
	req.Header.Set("Content-type", "application/json")
	req.Header.Set("token", "admin_token")
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	r.Equal(http.StatusInternalServerError, resp.Result().StatusCode)
}
