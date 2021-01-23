package test

import (
	"encoding/json"
	"fmt"
	"github.com/labstack/echo/v4"
	"go-app-template/src/api/controller/dto"
	"go-app-template/src/apperror/message"
	"go-app-template/src/config/db/localdata"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

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
	doStatusOKCheck__ApplicationController(t, params)
}

func TestApplicaionController_Login_異常系(t *testing.T) {
	// setup
	passwords := []string{" ", "　", "\n", "Test 1111", "Test　1111", "Test\n1111", " Test1111", "Test1111　", " Test1111\n"}

	var params = []errorCheckParam{
		{
			"パスワードが間違っているとき、401になること",
			[]requestParam{{httpMethod: http.MethodPost, path: "/login", body: strings.NewReader(`{"id": 1, "password":"hogehoge"}`)}},
			http.StatusUnauthorized,
			message.WrongPassword,
		},
		{
			"存在しないuserIdのとき、404になること",
			[]requestParam{{httpMethod: http.MethodPost, path: "/login", body: strings.NewReader(`{"id": 9999, "password":"Test1111"}`)}},
			http.StatusNotFound,
			message.UserNotFound,
		},
		{
			"userIdが数字ではないとき、400になること",
			[]requestParam{{httpMethod: http.MethodPost, path: "/login", body: strings.NewReader(`{"id": hoge, "password":"Test1111"}`)}},
			http.StatusBadRequest,
			message.StatusBadRequest,
		},
		{
			"userIdが0以下のとき、400になること",
			[]requestParam{
				{httpMethod: http.MethodPost, path: "/login", body: strings.NewReader(`{"id": 0, "password":"Test1111"}`)},
				{httpMethod: http.MethodPost, path: "/login", body: strings.NewReader(`{"id": -1, "password":"Test1111"}`)}},
			http.StatusBadRequest,
			message.StatusBadRequest,
		},
		{
			"passwordが空文字のとき、400になること",
			[]requestParam{{httpMethod: http.MethodPost, path: "/login", body: strings.NewReader(`{"id": 1, "password":""}`)}},
			http.StatusBadRequest,
			message.StatusBadRequest,
		},
		{
			"passwordに半角・全角スペース、改行が含まれているとき、400になること",
			makeInputs(http.MethodPost, "/login", makeBodyList(makeLoginReceiveDtoJsonList(passwords))),
			http.StatusBadRequest,
			message.StatusBadRequest,
		},
	}

	// check
	doErrorCheck(t, params)
}

/**************************************
	private
**************************************/
func doStatusOKCheck__ApplicationController(t *testing.T, params []statusOKCheckParamBase) {
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
		var actualBody dto.LoginResDto
		if err := json.Unmarshal(rec.Body.Bytes(), &actualBody); err != nil {
			t.Errorf("responseにTokenが返ってきていません\nresponse: %v", rec.Body.String())
		}

		// expected
		const expectedCode = http.StatusOK

		// check
		// LoginResDtoに変換でき && ステータスコードが200であればOKとする
		fmt.Println(p.title)
		assert.Equal(t, expectedCode, actualCode)

		// clean
		localdata.InitializeLocalData()
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
