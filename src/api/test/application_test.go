package test

import (
	"encoding/json"
	"fmt"
	"go-app-template/src/api/controller/dto"
	"go-app-template/src/apperror/message"
	"go-app-template/src/config/db/localdata"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"

	"github.com/stretchr/testify/assert"
)

/**************************************
	ログイン
**************************************/

func TestApplicaionController_Login_正常系(t *testing.T) {
	// setup
	var params = []statusOKCheckParamBase{
		{
			title:        "正常にログインできること",
			requestParam: requestParam{httpMethod: http.MethodPost, path: "/login", body: strings.NewReader(`{"id": 1, "password":"Test1111"}`)},
		},
	}

	// check
	doStatusOKCheck__ApplicationController(t, params, false)
}

func TestApplicaionController_Login_異常系(t *testing.T) {
	// setup
	passwords := []string{" ", "　", "\n", "Test 1111", "Test　1111", "Test\n1111", " Test1111", "Test1111　", " Test1111\n"}

	var params = []errorCheckParam{
		{
			title:           "パスワードが間違っているとき、401になること",
			requestParams:   []requestParam{{httpMethod: http.MethodPost, path: "/login", body: strings.NewReader(`{"id": 1, "password":"hogehoge"}`)}},
			expectedCode:    http.StatusUnauthorized,
			expectedMessage: message.WrongPassword,
		},
		{
			title:           "存在しないuserIdのとき、404になること",
			requestParams:   []requestParam{{httpMethod: http.MethodPost, path: "/login", body: strings.NewReader(`{"id": 9999, "password":"Test1111"}`)}},
			expectedCode:    http.StatusNotFound,
			expectedMessage: message.UserNotFound,
		},
		{
			title:           "userIdが数字ではないとき、400になること",
			requestParams:   []requestParam{{httpMethod: http.MethodPost, path: "/login", body: strings.NewReader(`{"id": hoge, "password":"Test1111"}`)}},
			expectedCode:    http.StatusBadRequest,
			expectedMessage: message.StatusBadRequest,
		},
		{
			title: "userIdが0以下のとき、400になること",
			requestParams: []requestParam{
				{httpMethod: http.MethodPost, path: "/login", body: strings.NewReader(`{"id": 0, "password":"Test1111"}`)},
				{httpMethod: http.MethodPost, path: "/login", body: strings.NewReader(`{"id": -1, "password":"Test1111"}`)}},
			expectedCode:    http.StatusBadRequest,
			expectedMessage: message.StatusBadRequest,
		},
		{
			title:           "passwordが空文字のとき、400になること",
			requestParams:   []requestParam{{httpMethod: http.MethodPost, path: "/login", body: strings.NewReader(`{"id": 1, "password":""}`)}},
			expectedCode:    http.StatusBadRequest,
			expectedMessage: message.StatusBadRequest,
		},
		{
			title:           "passwordに半角・全角スペース、改行が含まれているとき、400になること",
			requestParams:   makeInputs(http.MethodPost, "/login", makeBodyList(makeLoginReceiveDtoJsonList(passwords))),
			expectedCode:    http.StatusBadRequest,
			expectedMessage: message.StatusBadRequest,
		},
	}

	// check
	doErrorCheck(t, params)
}

/**************************************
	ログアウト
**************************************/

func TestApplicaionController_Logout_正常系(t *testing.T) {
	// setup
	var params = []statusOKCheckParamBase{
		{
			title:        "正常にログアウトできること",
			requestParam: requestParam{httpMethod: http.MethodGet, path: "/logout", body: nil},
		},
	}

	// check
	doStatusOKCheck__ApplicationController(t, params, false)
}

/**************************************
	private
**************************************/

func doStatusOKCheck__ApplicationController(t *testing.T, params []statusOKCheckParamBase, doCleanData bool) {
	for _, p := range params {
		// setup
		req := httptest.NewRequest(p.requestParam.httpMethod, p.requestParam.path, p.requestParam.body)
		if p.requestParam.httpMethod == http.MethodPost || p.requestParam.httpMethod == http.MethodPut {
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		}
		rec := httptest.NewRecorder()
		_target.ServeHTTP(rec, req)

		// actual
		actualCode := rec.Code

		// expected
		const expectedCode = http.StatusOK

		// check
		// ステータスコードが200であればOKとする
		fmt.Println(p.title)
		assert.Equal(t, expectedCode, actualCode)

		// clean
		if doCleanData {
			localdata.InitializeLocalData()
		}
	}
}

func makeLoginReceiveDtoJsonList(passwords []string) [][]byte {
	list := make([][]byte, len(passwords))
	for i, p := range passwords {
		userDto := dto.LoginReceiveDto{
			Id:       1,
			Password: p,
		}
		b, _ := json.Marshal(userDto)
		list[i] = b
	}
	return list
}
