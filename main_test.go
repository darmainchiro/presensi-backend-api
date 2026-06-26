package main

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

// Skenario 1: Test Login dengan Email yang Benar
func TestLoginAPI_Success(t *testing.T) {
	// 1. Setup Router dalam Mode Test
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	
	// Kita inisialisasi AppHandler kosong karena API Login saat ini tidak butuh ke Database
	appHandler := &AppHandler{} 
	r.POST("/api/login", appHandler.LoginAPI)

	// 2. Buat Request Palsu (Mock Form Data)
	formData := url.Values{}
	formData.Set("email", "aji@example.com")
	
	req, _ := http.NewRequest(http.MethodPost, "/api/login", strings.NewReader(formData.Encode()))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	// 3. Tangkap Response menggunakan HTTP Recorder (Tanpa menyalakan server)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	// 4. Validasi Hasil (Assertion)
	// Kita berharap Status HTTP adalah 200 OK
	assert.Equal(t, http.StatusOK, w.Code)
	// Kita berharap di dalam response JSON terdapat kata "token"
	assert.Contains(t, w.Body.String(), "token")
}

// Skenario 2: Test Login dengan Email yang Salah
func TestLoginAPI_Failed(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	appHandler := &AppHandler{}
	r.POST("/api/login", appHandler.LoginAPI)

	// Memasukkan email yang salah
	formData := url.Values{}
	formData.Set("email", "salah@example.com")
	
	req, _ := http.NewRequest(http.MethodPost, "/api/login", strings.NewReader(formData.Encode()))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	// Kita berharap Status HTTP adalah 401 Unauthorized
	assert.Equal(t, http.StatusUnauthorized, w.Code)
	// Kita berharap response mengeluarkan pesan error
	assert.Contains(t, w.Body.String(), "Email tidak terdaftar")
}