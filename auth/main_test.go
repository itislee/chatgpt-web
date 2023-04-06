package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestHandleCallback(t *testing.T) {
	db, _ := NewMySQLDatabase("username:password@tcp(localhost:3306)/dbname")

	router := gin.New()
	router.GET("/callback", func(c *gin.Context) { handleCallback(c, db) })

	// 正确的请求参数
	qs := "code=valid_code&state=valid_state"
	req, _ := http.NewRequest("GET", "/callback?"+qs, nil)
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	if resp.Code != http.StatusOK {
		t.Errorf("Expected status 200, but got %d", resp.Code)
	}

	// 不正确的请求参数
	qs = "code=invalid_code&state=invalid_state"
	req, _ = http.NewRequest("GET", "/callback?"+qs, nil)
	resp = httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	if resp.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, but got %d", resp.Code)
	}
}

func TestCheckAccessToken(t *testing.T) {
	db, _ := NewMySQLDatabase("username:password@tcp(localhost:3306)/dbname")

	router := gin.New()
	router.GET("/check", func(c *gin.Context) { checkAccessToken(c, db) })

	// 带有正确的Cookie
	req, _ := http.NewRequest("GET", "/check", nil)
	req.AddCookie(&http.Cookie{Name: "access_token", Value: "valid_access_token"})
	req.AddCookie(&http.Cookie{Name: "open_id", Value: "valid_open_id"})
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	if resp.Code != http.StatusOK {
		t.Errorf("Expected status 200, but got %d", resp.Code)
	}

	// 带有错误的Cookie
	req, _ = http.NewRequest("GET", "/check", nil)
	req.AddCookie(&http.Cookie{Name: "access_token", Value: "invalid_access_token"})
	req.AddCookie(&http.Cookie{Name: "open_id", Value: "invalid_open_id"})
	resp = httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	if resp.Code != http.StatusUnauthorized {
		t.Errorf("Expected status 401, but got %d", resp.Code)
	}
}
