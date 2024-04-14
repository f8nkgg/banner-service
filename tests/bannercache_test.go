package tests

import (
	v1 "banner/internal/controller/http/v1"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
)

func (s *APITestSuite) TestBannerCache() {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	v1.RegisterRoutes(router, s.handler)
	r := s.Require()
	s.createTestBanner()
	defer s.deleteTestBanner()
	// Запрос на получение баннера в первый раз
	req, _ := http.NewRequest("GET", "/user_banner?tag_id=4&feature_id=123", nil)
	req.Header.Set("Content-type", "application/json")
	req.Header.Set("token", "user_token")
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)
	r.Equal(http.StatusOK, resp.Result().StatusCode)

	// Проверяем, что баннер был получен из базы данных
	responseBody, err := io.ReadAll(resp.Body)
	s.NoError(err)
	r.Equal("{\"text\":\"some_text3\",\"title\":\"some_title\",\"url\":\"some_url2\"}", string(responseBody))

	// Обновление баннера
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

	// Запрос на получение баннера во второй раз
	req2, _ := http.NewRequest("GET", "/user_banner?tag_id=4&feature_id=123", nil)
	req2.Header.Set("Content-type", "application/json")
	req2.Header.Set("token", "user_token")
	resp2 := httptest.NewRecorder()
	router.ServeHTTP(resp2, req2)
	r.Equal(http.StatusOK, resp2.Result().StatusCode)

	// Проверяем, что баннер был получен из кэша
	responseBody2, err := io.ReadAll(resp2.Body)
	s.NoError(err)
	r.Equal("{\"text\":\"some_text3\",\"title\":\"some_title\",\"url\":\"some_url2\"}", string(responseBody2))

	// Запрос на получение баннера из бд
	req3, _ := http.NewRequest("GET", "/user_banner?tag_id=4&feature_id=123&use_last_revision=true", nil)
	req3.Header.Set("Content-type", "application/json")
	req3.Header.Set("token", "user_token")
	resp3 := httptest.NewRecorder()
	router.ServeHTTP(resp3, req3)
	r.Equal(http.StatusOK, resp3.Result().StatusCode)
	responseBody3, err := io.ReadAll(resp3.Body)
	s.NoError(err)
	r.Equal("{\"text\":\"some_text11\",\"title\":\"some_title12\",\"url\":\"some_url13\"}", string(responseBody3))

}
