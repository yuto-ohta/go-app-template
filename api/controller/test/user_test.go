package test

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"go-app-template/config"
	"go-app-template/domain"
	"gorm.io/gorm"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
)

func TestUserController_userID_1でユーザーが取得できること(t *testing.T) {
	router := config.NewRouter()
	req := httptest.NewRequest("GET", "/user/1", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	expectedCode := http.StatusOK
	expectedBody := *domain.NewUser(*domain.NewUserId(1), "taro")

	actualCode := rec.Code
	var actualBody domain.User
	if err := actualBody.UnmarshalJSON(rec.Body.Bytes()); err != nil {
		t.Errorf("UnknownError: err %v, ResponceBody %v", err.Error(), rec.Body.String())
	}

	assert.Equal(t, expectedCode, actualCode)
	assert.Equal(t, expectedBody, actualBody)
}

func TestUserController_存在しないuserIDで404が返ること(t *testing.T) {
	router := config.NewRouter()
	req := httptest.NewRequest("GET", "/user/9999", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	expectedCode := http.StatusNotFound
	expectedBody := gorm.ErrRecordNotFound.Error()

	actualCode := rec.Code
	var actualBody string
	if err := json.Unmarshal(rec.Body.Bytes(), &actualBody); err != nil {
		t.Errorf("UnknownError: err %v, ResponceBody %v", err.Error(), rec.Body.String())
	}

	assert.Equal(t, expectedCode, actualCode)
	assert.Equal(t, expectedBody, actualBody)
}

func TestUserController_userIDが数字ではないときは400が返ること(t *testing.T) {
	router := config.NewRouter()
	req := httptest.NewRequest("GET", "/user/taro", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	expectedCode := http.StatusBadRequest
	_, err := strconv.Atoi("taro")
	expectedBody := err.Error()

	actualCode := rec.Code
	var actualBody string
	if err := json.Unmarshal(rec.Body.Bytes(), &actualBody); err != nil {
		t.Errorf("UnknownError: err %v, ResponceBody %v", err.Error(), rec.Body.String())
	}

	assert.Equal(t, expectedCode, actualCode)
	assert.Equal(t, expectedBody, actualBody)
}
