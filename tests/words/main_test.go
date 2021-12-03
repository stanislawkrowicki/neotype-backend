package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"neotype-backend/pkg/words"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetRandomWords(t *testing.T) {
	var goodNumbers = [...]int{12, 36, 100, 999}
	var badNumbers = [...]int{0, -2, -15, -999}
	var notNumbers = [...]string{"jade", "10f", "9asf", "_", "10.2", "10,2"}

	for _, num := range goodNumbers {
		request, _ := http.NewRequest("GET", fmt.Sprintf("/words/%d", num), nil)
		response := httptest.NewRecorder()
		router().ServeHTTP(response, request)
		assert.Equal(t, 200, response.Code, "200 OK was expected for proper numbers")
	}

	for _, num := range badNumbers {
		request, _ := http.NewRequest("GET", fmt.Sprintf("/words/%d", num), nil)
		response := httptest.NewRecorder()
		router().ServeHTTP(response, request)
		assert.Equal(t, 400, response.Code, "400 Bad Request was expected for bad number.")
	}

	for _, num := range notNumbers {
		request, _ := http.NewRequest("GET", fmt.Sprintf("/words/%s", num), nil)
		response := httptest.NewRecorder()
		router().ServeHTTP(response, request)
		assert.Equal(t, 400, response.Code, "400 Bad Request was expected for not numbers.")
	}
}

func router() *gin.Engine {
	r := gin.New()
	r.GET("/words/:count", words.GetRandomWords)
	return r
}
