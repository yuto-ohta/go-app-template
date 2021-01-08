package test

import (
	"encoding/json"
	"fmt"
	"go-app-template/src/api/controller/dto"
	"go-app-template/src/apperror"
	"go-app-template/src/apperror/message"
	"go-app-template/src/apputil"
	"go-app-template/src/config/db/localdata"
	"go-app-template/src/config/route"
	"go-app-template/src/domain"
	"go-app-template/src/domain/valueobject"
	"go-app-template/src/infrastructure"
	"go-app-template/src/usecase/impl"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

var (
	_target      = route.NewRouter()
	_userUseCase = *impl.NewUserUseCaseImpl(infrastructure.NewUserRepositoryImpl())
)

type input struct {
	httpMethod string
	path       string
	body       io.Reader
}

type errorCheckParam struct {
	title           string
	input           []input
	expectedCode    int
	expectedMessage message.Message
}

func TestMain(m *testing.M) {
	// before all
	localdata.InitializeLocalData()

	// run each test
	code := m.Run()

	// after all
	localdata.InitializeLocalData()

	// finish test
	os.Exit(code)
}

func TestUserController_GetUser_正常系(t *testing.T) {
	// setup
	var params = []struct {
		title             string
		input             input
		expectedCode      int
		expectedUserIdInt int
		expectedName      string
	}{
		{
			"正常にユーザーが取得できること",
			input{httpMethod: "GET", path: "/user/1", body: nil},
			http.StatusOK,
			1,
			"まるお",
		},
	}

	for _, p := range params {
		req := httptest.NewRequest(p.input.httpMethod, p.input.path, p.input.body)
		rec := httptest.NewRecorder()
		_target.ServeHTTP(rec, req)

		// actual
		actualCode := rec.Code
		var actualBody domain.User
		_ = actualBody.UnmarshalJSON(rec.Body.Bytes())

		// expected
		id, _ := valueobject.NewUserIdWithId(p.expectedUserIdInt)
		expectedCode := p.expectedCode
		expectedBody, _ := domain.NewUserWithUserId(*id, p.expectedName)

		// check
		fmt.Println(p.title)
		assert.Equal(t, expectedCode, actualCode)
		assert.Equal(t, *expectedBody, actualBody)
	}
}

func TestUserController_GetUser_異常系(t *testing.T) {
	// setup
	var params = []errorCheckParam{
		{
			"存在しないuserIdのとき、404になること",
			[]input{{httpMethod: "GET", path: "/user/9999", body: nil}},
			http.StatusNotFound,
			message.UserNotFound,
		},
		{
			"userIdが数字ではないとき、400になること",
			[]input{{httpMethod: "GET", path: "/user/hogehoge", body: nil}},
			http.StatusBadRequest,
			message.InvalidUserId,
		},
		{
			"userIdがマイナスのとき、400になること",
			[]input{{httpMethod: "GET", path: "/user/-1", body: nil}},
			http.StatusBadRequest,
			message.InvalidUserId,
		},
	}

	// check
	doErrorCheck(t, params)
}

func TestUserController_CreateUser_正常系(t *testing.T) {
	// setup
	var params = []struct {
		title        string
		input        input
		expectedCode int
		expectedName string
	}{
		{
			"正常にユーザーが登録されること",
			input{httpMethod: "GET", path: "/user/new?name=新規ユーザー太郎", body: nil},
			http.StatusOK,
			"新規ユーザー太郎",
		},
		{
			"userNameの両端に半角・全角スペースがあるとき、スペースが取り除かれ、ユーザーが登録されること",
			input{httpMethod: "GET", path: fmt.Sprintf("/user/new?name=%v", apputil.QueryEncoding(" 　 　新規ユーザー太郎 　 　")), body: nil},
			http.StatusOK,
			"新規ユーザー太郎",
		},
	}

	for _, p := range params {
		req := httptest.NewRequest(p.input.httpMethod, p.input.path, p.input.body)
		rec := httptest.NewRecorder()
		_target.ServeHTTP(rec, req)

		// actual
		actualCode := rec.Code
		var actualBody domain.User
		_ = actualBody.UnmarshalJSON(rec.Body.Bytes())
		actualName := actualBody.GetName()

		// expected
		expectedCode := p.expectedCode
		expectedBody, _ := _userUseCase.FindById(actualBody.GetId().GetValue())
		expectedName := p.expectedName

		// check
		fmt.Println(p.title)
		assert.Equal(t, expectedCode, actualCode)
		assert.Equal(t, expectedBody, actualBody)
		assert.Equal(t, expectedName, actualName)

		// clean
		localdata.InitializeLocalData()
	}
}
func TestUserController_CreateUser_異常系(t *testing.T) {
	// setup
	userNames := []string{" ", "　", "\n", "新規ユーザー 太郎", "新規ユーザー　太郎", "新規ユーザー\n太郎", " 新規ユーザー\n太郎　"}

	var params = []errorCheckParam{
		{
			"userNameが空文字のとき、400になること",
			[]input{{httpMethod: "GET", path: "/user/new?name=", body: nil}},
			http.StatusBadRequest,
			message.InvalidUserName,
		},
		{
			"userNameに半角・全角スペース、改行が含まれているとき、400になること",
			makeQueryParamInputs("GET", "/user/new?name=", userNames, nil),
			http.StatusBadRequest,
			message.InvalidUserName,
		},
		{
			"userNameが9文字以上のとき、400になること",
			[]input{{httpMethod: "GET", path: "/user/new?name=123456789", body: nil}},
			http.StatusBadRequest,
			message.InvalidUserName,
		},
	}

	// check
	doErrorCheck(t, params)

	// clean
	localdata.InitializeLocalData()
}

func TestUserController_DeleteUser_正常系(t *testing.T) {
	// setup
	var resCheckParams = []struct {
		title             string
		input             input
		expectedCode      int
		expectedUserIdInt int
		expectedName      string
	}{
		{
			"正常にユーザーが削除できること①_レスポンスチェック",
			input{httpMethod: "DELETE", path: "/user/1", body: nil},
			http.StatusOK,
			1,
			"まるお",
		},
	}

	var recordCheckParams = []errorCheckParam{
		{
			"正常にユーザーが削除できること②_レコードチェック",
			[]input{{httpMethod: "GET", path: "/user/1", body: nil}},
			http.StatusNotFound,
			message.UserNotFound,
		},
	}

	// resCheck
	for _, p := range resCheckParams {
		req := httptest.NewRequest(p.input.httpMethod, p.input.path, p.input.body)
		rec := httptest.NewRecorder()
		_target.ServeHTTP(rec, req)

		// actual
		actualCode := rec.Code
		var actualBody domain.User
		_ = actualBody.UnmarshalJSON(rec.Body.Bytes())

		// expected
		id, _ := valueobject.NewUserIdWithId(p.expectedUserIdInt)
		expectedCode := p.expectedCode
		expectedBody, _ := domain.NewUserWithUserId(*id, p.expectedName)

		// check
		fmt.Println(p.title)

		assert.Equal(t, expectedCode, actualCode)
		assert.Equal(t, *expectedBody, actualBody)
	}

	// recordCheck
	doErrorCheck(t, recordCheckParams)

	// clean
	localdata.InitializeLocalData()
}

func TestUserController_DeleteUser_異常系(t *testing.T) {
	// setup
	var params = []errorCheckParam{
		{
			"存在しないuserIdのとき、404になること",
			[]input{{httpMethod: "DELETE", path: "/user/9999", body: nil}},
			http.StatusNotFound,
			message.UserNotFound,
		},
	}

	// check
	doErrorCheck(t, params)

	// clean
	localdata.InitializeLocalData()
}

func TestUserController_UpdateUser_正常系(t *testing.T) {
	// setup
	var resCheckParams = []struct {
		title             string
		input             input
		expectedCode      int
		expectedUserIdInt int
		expectedName      string
	}{
		{
			"正常にユーザー名が更新できること",
			input{httpMethod: http.MethodPost, path: "/user/1/update", body: strings.NewReader(`{"id":1,"name":"ハルキゲニア"}`)},
			http.StatusOK,
			1,
			"ハルキゲニア",
		},
	}

	for _, p := range resCheckParams {
		req := httptest.NewRequest(p.input.httpMethod, p.input.path, p.input.body)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		_target.ServeHTTP(rec, req)

		// actual
		actualCode := rec.Code
		var actualBody dto.UserDto
		_ = json.Unmarshal(rec.Body.Bytes(), &actualBody)
		actualRecordDomain, _ := _userUseCase.FindById(p.expectedUserIdInt)
		actualRecordDto := actualRecordDomain.ToDto()

		// expected
		expectedCode := p.expectedCode
		expectedBody := &dto.UserDto{
			Id:   p.expectedUserIdInt,
			Name: p.expectedName,
		}

		// check
		fmt.Println(p.title)
		assert.Equal(t, expectedCode, actualCode)
		assert.Equal(t, *expectedBody, actualBody)
		assert.Equal(t, expectedBody, actualRecordDto)

		// clean
		localdata.InitializeLocalData()
	}
}

func TestUserController_UpdateUser_異常系(t *testing.T) {
	// setup
	userNames := []string{" ", "　", "\n", "新規ユーザー 太郎", "新規ユーザー　太郎", "新規ユーザー\n太郎", " 新規ユーザー太郎　"}

	var params = []errorCheckParam{
		{
			"存在しないuserIdのとき、404になること",
			[]input{{httpMethod: http.MethodPost, path: "/user/9999/update", body: strings.NewReader(`{"id":1,"name":"ハルキゲニア"}`)}},
			http.StatusNotFound,
			message.UserNotFound,
		},
		{
			"userNameが空文字のとき、400になること",
			[]input{{httpMethod: http.MethodPost, path: "/user/1/update", body: strings.NewReader(`{"id":1,"name":""}`)}},
			http.StatusBadRequest,
			message.StatusBadRequest,
		},
		{
			"userNameがnilのとき、400になること",
			[]input{{httpMethod: http.MethodPost, path: "/user/1/update", body: strings.NewReader(`{"id":1,"name":}`)}},
			http.StatusBadRequest,
			message.StatusBadRequest,
		},
		{
			"userNameに半角・全角スペース、改行が含まれているとき、400になること",
			makePostInputs("/user/1/update", makeBodyList(makeUserDtoJsonList(userNames))),
			http.StatusBadRequest,
			message.StatusBadRequest,
		},
		{
			"userNameが9文字以上のとき、400になること",
			[]input{{httpMethod: http.MethodPost, path: "/user/1/update", body: strings.NewReader(`{"id":1,"name":"123456789"}`)}},
			http.StatusBadRequest,
			message.StatusBadRequest,
		},
	}

	// check
	doErrorCheck(t, params)

	// clean
	localdata.InitializeLocalData()
}

func doErrorCheck(t *testing.T, params []errorCheckParam) {
	for _, p := range params {
		for _, ip := range p.input {
			req := httptest.NewRequest(ip.httpMethod, ip.path, ip.body)
			if ip.httpMethod == http.MethodPost {
				req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			}
			rec := httptest.NewRecorder()
			_target.ServeHTTP(rec, req)

			// actual
			actualCode := rec.Code
			var actualBody apperror.ResponseErrorMessage
			_ = actualBody.UnmarshalJSON(rec.Body.Bytes())
			actualMessage := actualBody.GetMessage()

			// expected
			expectedCode := p.expectedCode
			expectedMessage := p.expectedMessage.String()

			// check
			fmt.Println(p.title)
			assert.Equal(t, expectedCode, actualCode)
			assert.Equal(t, expectedMessage, actualMessage)
		}
	}
}

func makeQueryParamInputs(httpMethod string, pathBase string, pathParams []string, body io.Reader) []input {
	inputs := make([]input, len(pathParams))
	for i, p := range pathParams {
		input := &input{
			httpMethod: httpMethod,
			path:       fmt.Sprintf("%v%v", pathBase, apputil.QueryEncoding(p)),
			body:       body,
		}
		inputs[i] = *input
	}
	return inputs
}

func makePostInputs(path string, body []io.Reader) []input {
	inputs := make([]input, len(body))
	for i, p := range body {
		input := &input{
			httpMethod: http.MethodPost,
			path:       path,
			body:       p,
		}
		inputs[i] = *input
	}
	return inputs
}

func makeBodyList(jsonList [][]byte) []io.Reader {
	list := make([]io.Reader, len(jsonList))
	for i, p := range jsonList {
		list[i] = strings.NewReader(string(p))
	}
	return list
}

func makeUserDtoJsonList(userNames []string) [][]byte {
	list := make([][]byte, len(userNames))
	for i, n := range userNames {
		userDto := dto.UserDto{
			Id:   1,
			Name: n,
		}
		j, _ := json.Marshal(userDto)
		list[i] = j
	}
	return list
}
