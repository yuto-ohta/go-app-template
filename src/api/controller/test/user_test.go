package test

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"go-app-template/src/apperror"
	"go-app-template/src/apperror/messages"
	"go-app-template/src/config"
	"go-app-template/src/config/db/localdata"
	"go-app-template/src/config/routes"
	"go-app-template/src/domain"
	"go-app-template/src/domain/valueobject"
	"go-app-template/src/infrastructure"
	"go-app-template/src/usecase"
	"go-app-template/src/usecase/impl"
	appUtils "go-app-template/src/utils"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

var (
	userUseCase usecase.UserUseCase
)

func TestMain(m *testing.M) {
	// before all
	config.LoadConfig()
	localdata.InitializeLocalData()
	userRepository := infrastructure.NewUserRepositoryImpl()
	userUseCase = *impl.NewUserUseCaseImpl(userRepository)

	// run each test
	code := m.Run()

	// after all

	// finish test
	localdata.InitializeLocalData()
	os.Exit(code)
}

func TestUserController_GetUser_正常にユーザーが取得できること(t *testing.T) {
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
	expectedBody := *domain.NewUserWithUserId(*valueobject.NewUserIdWithId(1), "まるお")

	// check
	assert.Equal(t, expectedCode, actualCode)
	assert.Equal(t, expectedBody, actualBody)
}

func TestUserController_GetUser_userIdに紐づくユーザーがいない場合_404が返ること(t *testing.T) {
	// setup
	router := routes.NewRouter()
	req := httptest.NewRequest("GET", "/user/9999", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	// actual
	actualCode := rec.Code
	var actualBody apperror.ResponseErrorMessage
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

func TestUserController_GetUser_userIDが数字ではないとき_400が返ること(t *testing.T) {
	// setup
	router := routes.NewRouter()
	req := httptest.NewRequest("GET", "/user/taro", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	// actual
	actualCode := rec.Code
	var actualBody apperror.ResponseErrorMessage
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

func TestUserController_CreateUser_正常にユーザーが登録されること(t *testing.T) {
	// setup
	userNameParam := "新規ユーザー太郎"
	router := routes.NewRouter()
	req := httptest.NewRequest("GET", fmt.Sprintf("/user/new?name=%v", userNameParam), nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	// actual
	var actualCode int
	var actualBody domain.User
	actualCode = rec.Code
	if err := actualBody.UnmarshalJSON(rec.Body.Bytes()); err != nil {
		t.Errorf("ResponseBodyがdomain.Userの構造と合致していません, Error: %v, ResponseBody: %v", err.Error(), rec.Body.String())
	}

	// expected
	var expectedCode int
	var expectedBody domain.User
	expectedCode = http.StatusOK
	expectedBody, _ = userUseCase.FindById(actualBody.GetId())

	// check
	assert.Equal(t, expectedCode, actualCode)
	assert.Equal(t, expectedBody, actualBody)
}

func TestUserController_CreateUser_userNameが存在しない場合_400エラーが返ること(t *testing.T) {
	// setup
	userNameParam := ""
	router := routes.NewRouter()
	req := httptest.NewRequest("GET", fmt.Sprintf("/user/new?name=%v", userNameParam), nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	// actual
	actualCode := rec.Code
	var actualBody apperror.ResponseErrorMessage
	if err := actualBody.UnmarshalJSON(rec.Body.Bytes()); err != nil {
		t.Errorf("ResponseBodyがerrors.ResponseErrorMessageの構造と合致していません, Error: %v, ResponseBody: %v", err.Error(), rec.Body.String())
	}
	actualStatus := actualBody.GetStatus()
	actualMessage := actualBody.GetMessage()

	// expected
	expectedCode := http.StatusBadRequest
	expectedStatus := expectedCode
	expectedMessage := messages.InvalidUserName.String()

	// check
	assert.Equal(t, expectedCode, actualCode)
	assert.Equal(t, expectedStatus, actualStatus)
	assert.Equal(t, expectedMessage, actualMessage)
}

func TestUserController_CreateUser_userNameに半角_全角スペース_改行が含まれている場合_400エラーが返ること(t *testing.T) {
	// setup
	router := routes.NewRouter()
	userNameParams := []string{" ", "　", "\n", "新規ユーザー 太郎", "新規ユーザー　太郎", "新規ユーザー\n太郎"}

	// expected
	expectedCode := http.StatusBadRequest
	expectedStatus := expectedCode
	expectedMessage := messages.InvalidUserName.String()

	for _, param := range userNameParams {
		param = appUtils.QueryEncoding(param)
		req := httptest.NewRequest("GET", fmt.Sprintf("/user/new?name=%v", param), nil)
		rec := httptest.NewRecorder()
		router.ServeHTTP(rec, req)

		// actual
		actualCode := rec.Code
		var actualBody apperror.ResponseErrorMessage
		if err := actualBody.UnmarshalJSON(rec.Body.Bytes()); err != nil {
			t.Errorf("ResponseBodyがerrors.ResponseErrorMessageの構造と合致していません, Error: %v, ResponseBody: %v", err.Error(), rec.Body.String())
		}
		actualStatus := actualBody.GetStatus()
		actualMessage := actualBody.GetMessage()

		// check
		assert.Equal(t, expectedCode, actualCode)
		assert.Equal(t, expectedStatus, actualStatus)
		assert.Equal(t, expectedMessage, actualMessage)
	}
}
