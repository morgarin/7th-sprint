package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCafeCount(t *testing.T) {
	handler := http.HandlerFunc(mainHandle)

	testCases := []struct {
		count string
		want  int
	}{
		{"/cafe?city=moscow&count=0", 0},
		{"/cafe?city=moscow&count=1", 1},
		{"/cafe?city=moscow&count=2", 2},
		{"/cafe?city=moscow&count=100", 5},
	}

	for _, tc := range testCases {
		response := httptest.NewRecorder()
		req := httptest.NewRequest("GET", tc.count, nil)

		handler.ServeHTTP(response, req)

		require.Equal(t, http.StatusOK, response.Code)

		resp := response.Body.String()

		respnum := strings.Count(resp, ",") + 1

		if resp == "" {
			respnum = 0
		}

		assert.Equal(t, tc.want, respnum, "Request: %s", tc.count)
	}
}

func TestCafeSearch(t *testing.T) {
	handler := http.HandlerFunc(mainHandle)

	testCases := []struct {
		search    string
		wantCount int
	}{
		{"фасоль", 0},
		{"кофе", 2},
		{"вилка", 1},
	}
	for _, v := range testCases {
		response := httptest.NewRecorder()
		request := fmt.Sprintf("/cafe?city=moscow&search=%s", v.search)
		req := httptest.NewRequest("GET", request, nil)

		handler.ServeHTTP(response, req)

		require.Equal(t, http.StatusOK, response.Code)

		resp := response.Body.String()
		resp = strings.ToLower(resp)
		respStrs := strings.Split(resp, ",")
		respNum := 0
		for _, z := range respStrs {
			if strings.Contains(z, v.search) {
				respNum++
			}
		}
		lenResp := len(respStrs)
		if resp == "" {
			lenResp = 0
		}
		//Quantity check
		assert.Equal(t, v.wantCount, lenResp, "Request: %s", request)
		//Search contains check
		assert.Equal(t, v.wantCount, respNum, "Request: %s", request)
	}
}

func TestCafeNegative(t *testing.T) {
	handler := http.HandlerFunc(mainHandle)

	requests := []struct {
		request string
		status  int
		message string
	}{
		{"/cafe", http.StatusBadRequest, "unknown city"},
		{"/cafe?city=omsk", http.StatusBadRequest, "unknown city"},
		{"/cafe?city=tula&count=na", http.StatusBadRequest, "incorrect count"},
	}
	for _, v := range requests {
		response := httptest.NewRecorder()
		req := httptest.NewRequest("GET", v.request, nil)
		handler.ServeHTTP(response, req)

		assert.Equal(t, v.status, response.Code)
		assert.Equal(t, v.message, strings.TrimSpace(response.Body.String()))
	}
}

func TestCafeWhenOk(t *testing.T) {
	handler := http.HandlerFunc(mainHandle)

	requests := []string{
		"/cafe?count=2&city=moscow",
		"/cafe?city=tula",
		"/cafe?city=moscow&search=ложка",
	}
	for _, v := range requests {
		response := httptest.NewRecorder()
		req := httptest.NewRequest("GET", v, nil)

		handler.ServeHTTP(response, req)

		assert.Equal(t, http.StatusOK, response.Code)
	}
}
