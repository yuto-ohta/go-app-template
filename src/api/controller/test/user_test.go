package test

import (
	"github.com/stretchr/testify/assert"
	"go-app-template/src/config/routes"
	"go-app-template/src/domain"
	"go-app-template/src/errors"
	"go-app-template/src/errors/messages"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestUserController_userID_1でユーザーが取得できること(t *testing.T) {
	// setup
	router := routes.NewRouter()
	req := httptest.NewRequest("GET", "/user/1", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	// actual
	actualCode := rec.Code
	var actualBody domain.User
	if err := actualBody.UnmarshalJSON(rec.Body.Bytes()); err != nil {
		t.Errorf("ResponseBodyがdomain.Userの構造と合致していません, Error: %v, ResponseBody: %v", err.Error(), rec.Body.String())
	}

	// expected
	expectedCode := http.StatusOK
	// TODO: DB周りのテスト環境整備
	expectedBody := *domain.NewUser(*domain.NewUserId(1), "taro")

	// check
	assert.Equal(t, expectedCode, actualCode)
	assert.Equal(t, expectedBody, actualBody)
}

func TestUserController_存在しないuserIDで404が返ること(t *testing.T) {
	// setup
	router := routes.NewRouter()
	req := httptest.NewRequest("GET", "/user/9999", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	// actual
	actualCode := rec.Code
	var actualBody errors.ResponseErrorMessage
	if err := actualBody.UnmarshalJSON(rec.Body.Bytes()); err != nil {
		t.Errorf("ResponseBodyがerrors.ResponseErrorMessageの構造と合致していません, Error: %v, ResponseBody: %v", err.Error(), rec.Body.String())
	}
	actualStatus := actualBody.GetStatus()
	actualMessage := actualBody.GetMessage()

	// expected
	expectedCode := http.StatusNotFound
	expectedStatus := expectedCode
	expectedMessage := messages.UserNotFound.String()

	// check
	assert.Equal(t, expectedCode, actualCode)
	assert.Equal(t, expectedStatus, actualStatus)
	assert.Equal(t, expectedMessage, actualMessage)
}

func TestUserController_userIDが数字ではないときは400が返ること(t *testing.T) {
	// setup
	router := routes.NewRouter()
	req := httptest.NewRequest("GET", "/user/taro", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	// actual
	actualCode := rec.Code
	var actualBody errors.ResponseErrorMessage
	if err := actualBody.UnmarshalJSON(rec.Body.Bytes()); err != nil {
		t.Errorf("ResponseBodyがerrors.ResponseErrorMessageの構造と合致していません, Error: %v, ResponseBody: %v", err.Error(), rec.Body.String())
	}
	actualStatus := actualBody.GetStatus()
	actualMessage := actualBody.GetMessage()

	// expected
	expectedCode := http.StatusBadRequest
	expectedStatus := expectedCode
	expectedMessage := messages.InvalidUserId.String()

	// check
	assert.Equal(t, expectedCode, actualCode)
	assert.Equal(t, expectedStatus, actualStatus)
	assert.Equal(t, expectedMessage, actualMessage)
}
