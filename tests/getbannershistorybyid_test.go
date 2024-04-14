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
	"strings"
)

func (s *APITestSuite) TestGetBannersHistoryByID_Success() {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	v1.RegisterRoutes(router, s.handler)
	r := s.Require()

	s.createTestBanner()
	defer s.deleteTestBanner()
	req, _ := http.NewRequest("GET", "/banner/history/1", nil)
	req.Header.Set("Content-type", "application/json")
	req.Header.Set("token", "admin_token")
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	updateReqBody := `{
		  "content": {
    "title": "some_title12",
    "text": "some_text11",
    "url": "some_url13"
  }
	}`
	updateReq, _ := http.NewRequest("PATCH", "/banner/1", strings.NewReader(updateReqBody))
	updateReq.Header.Set("Content-Type", "application/json")
	updateReq.Header.Set("token", "admin_token")
	updateResp := httptest.NewRecorder()
	router.ServeHTTP(updateResp, updateReq)
	r.Equal(http.StatusNoContent, updateResp.Result().StatusCode)

	r.Equal(http.StatusOK, resp.Result().StatusCode)
	var banners []entity.FilteredBanner
	err := json.Unmarshal(resp.Body.Bytes(), &banners)
	s.NoError(err)
	r.NotEmpty(banners)
}
func (s *APITestSuite) TestGetBannersHistoryByID_Unauthorized() {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	v1.RegisterRoutes(router, s.handler)
	r := s.Require()

	req, _ := http.NewRequest("GET", "/banner/history/1", nil)
	req.Header.Set("Content-type", "application/json")
	req.Header.Set("token", "")
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)
	r.Equal(http.StatusUnauthorized, resp.Result().StatusCode)
}
func (s *APITestSuite) TestGetBannersHistoryByID_Forbidden() {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	v1.RegisterRoutes(router, s.handler)
	r := s.Require()

	req, _ := http.NewRequest("GET", "/banner/history/1", nil)
	req.Header.Set("Content-type", "application/json")
	req.Header.Set("token", "user_token")
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)
	r.Equal(http.StatusForbidden, resp.Result().StatusCode)
}

func (s *APITestSuite) TestGetBannersHistoryByID_BadRequest() {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	v1.RegisterRoutes(router, s.handler)
	r := s.Require()

	req, _ := http.NewRequest("GET", "/banner/history/abc", nil)
	req.Header.Set("Content-type", "application/json")
	req.Header.Set("token", "admin_token")
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	r.Equal(http.StatusBadRequest, resp.Result().StatusCode)
}

func (s *APITestSuite) TestGetBannersHistoryByID_InternalServerError() {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	mockService := &MockBannerService{
		GetBannersHistoryByIDFunc: func(ctx context.Context, id int32) ([]*entity.BannerHistoryItem, error) {
			return nil, errors.New("internal server error")
		},
	}
	v1.RegisterRoutes(router, v1.NewBannerController(mockService, s.logger))
	r := s.Require()
	req, _ := http.NewRequest("GET", "/banner/history/1", nil)
	req.Header.Set("Content-type", "application/json")
	req.Header.Set("token", "admin_token")
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	r.Equal(http.StatusInternalServerError, resp.Result().StatusCode)
}
