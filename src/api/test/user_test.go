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

type statusOKCheckParamUser struct {
	base         statusOKCheckParamBase
	expectedBody dto.UserResDto
}

/**************************************
	ユーザー取得
**************************************/
func TestUserController_GetUser_正常系(t *testing.T) {
	// setup
	var params = []statusOKCheckParamUser{
		{base: statusOKCheckParamBase{
			title:        "正常にユーザーが取得できること",
			requestParam: requestParam{httpMethod: http.MethodGet, path: "/users/1", body: nil}},
			expectedBody: dto.UserResDto{Id: 1, Name: "まるお"}},
	}

	// check
	doStatusOKCheck__UserController(t, params, _checkExisting, false)
}

func TestUserController_GetUser_異常系(t *testing.T) {
	// setup
	var params = []errorCheckParam{
		{
			"存在しないuserIdのとき、404になること",
			[]requestParam{{httpMethod: http.MethodGet, path: "/users/9999", body: nil}},
			http.StatusNotFound,
			message.UserNotFound,
		},
		{
			"userIdが数字ではないとき、400になること",
			[]requestParam{{httpMethod: http.MethodGet, path: "/users/hogehoge", body: nil}},
			http.StatusBadRequest,
			message.InvalidUserId,
		},
		{
			"userIdが0以下のとき、400になること",
			[]requestParam{
				{httpMethod: http.MethodGet, path: "/users/0", body: nil},
				{httpMethod: http.MethodGet, path: "/users/-1", body: nil},
			},
			http.StatusBadRequest,
			message.InvalidUserId,
		},
	}

	// check
	doErrorCheck(t, params)
}

/**************************************
	ユーザー全件取得
**************************************/
func TestUserController_GetAllUser_正常系(t *testing.T) {
	// setup
	params := []struct {
		base             statusOKCheckParamBase
		expectedBody     dto.UserResDto
		expectedPageInfo dto.PageInfo
		expectedLen      int
	}{
		{
			base: statusOKCheckParamBase{
				title:        "正常にユーザーが取得できること",
				requestParam: requestParam{httpMethod: http.MethodGet, path: "/users", body: nil}},
			expectedBody:     dto.UserResDto{Id: 1, Name: "まるお"},
			expectedPageInfo: dto.PageInfo{PageNum: 1, LastPageNum: 1, Limit: 10, Offset: 0},
			expectedLen:      10,
		},
		{
			base: statusOKCheckParamBase{
				title:        "limitのみを指定して、limit件数分、正常にユーザーが取得できること",
				requestParam: requestParam{httpMethod: http.MethodGet, path: "/users?limit=5", body: nil}},
			expectedBody:     dto.UserResDto{Id: 1, Name: "まるお"},
			expectedPageInfo: dto.PageInfo{PageNum: 1, LastPageNum: 2, Limit: 5, Offset: 0},
			expectedLen:      5,
		},
		{
			base: statusOKCheckParamBase{
				title:        "page, limitを指定して、正常にユーザーが取得できること①",
				requestParam: requestParam{httpMethod: http.MethodGet, path: "/users?page=1&limit=3", body: nil}},
			expectedBody:     dto.UserResDto{Id: 1, Name: "まるお"},
			expectedPageInfo: dto.PageInfo{PageNum: 1, LastPageNum: 4, Limit: 3, Offset: 0},
			expectedLen:      3,
		},
		{
			base: statusOKCheckParamBase{
				title:        "page, limitを指定して、正常にユーザーが取得できること②",
				requestParam: requestParam{httpMethod: http.MethodGet, path: "/users?page=3&limit=3", body: nil}},
			expectedBody:     dto.UserResDto{Id: 7, Name: "べらぼう太郎"},
			expectedPageInfo: dto.PageInfo{PageNum: 3, LastPageNum: 4, Limit: 3, Offset: 6},
			expectedLen:      3,
		},
		{
			base: statusOKCheckParamBase{
				title:        "page, limitを指定して、正常にユーザーが取得できること③",
				requestParam: requestParam{httpMethod: http.MethodGet, path: "/users?page=4&limit=3", body: nil}},
			expectedBody:     dto.UserResDto{Id: 10, Name: "先生"},
			expectedPageInfo: dto.PageInfo{PageNum: 4, LastPageNum: 4, Limit: 3, Offset: 9},
			expectedLen:      1,
		},
		{
			base: statusOKCheckParamBase{
				title:        "最終ページを超過したpageを指定して、最終ページが正常に取得できること",
				requestParam: requestParam{httpMethod: http.MethodGet, path: "/users?page=10&limit=9", body: nil}},
			expectedBody:     dto.UserResDto{Id: 10, Name: "先生"},
			expectedPageInfo: dto.PageInfo{PageNum: 2, LastPageNum: 2, Limit: 9, Offset: 9},
			expectedLen:      1,
		},
		{
			base: statusOKCheckParamBase{
				title:        "orderBy=idのみを指定して、id・ASCで、正常にユーザーが取得できること",
				requestParam: requestParam{httpMethod: http.MethodGet, path: "/users?orderBy=id", body: nil}},
			expectedBody:     dto.UserResDto{Id: 1, Name: "まるお"},
			expectedPageInfo: dto.PageInfo{PageNum: 1, LastPageNum: 1, Limit: 10, Offset: 0},
			expectedLen:      10,
		},
		{
			base: statusOKCheckParamBase{
				title:        "orderBy=nameのみを指定して、userName・ASCで、正常にユーザーが取得できること",
				requestParam: requestParam{httpMethod: http.MethodGet, path: "/users?orderBy=name", body: nil}},
			expectedBody:     dto.UserResDto{Id: 7, Name: "べらぼう太郎"},
			expectedPageInfo: dto.PageInfo{PageNum: 1, LastPageNum: 1, Limit: 10, Offset: 0},
			expectedLen:      10,
		},
		{
			base: statusOKCheckParamBase{
				title:        "order=ASCのみを指定して、id・ASCで、正常にユーザーが取得できること",
				requestParam: requestParam{httpMethod: http.MethodGet, path: "/users?order=ASC", body: nil}},
			expectedBody:     dto.UserResDto{Id: 1, Name: "まるお"},
			expectedPageInfo: dto.PageInfo{PageNum: 1, LastPageNum: 1, Limit: 10, Offset: 0},
			expectedLen:      10,
		},
		{
			base: statusOKCheckParamBase{
				title:        "order=DESCのみを指定して、id・DESCで、正常にユーザーが取得できること",
				requestParam: requestParam{httpMethod: http.MethodGet, path: "/users?order=DESC", body: nil}},
			expectedBody:     dto.UserResDto{Id: 10, Name: "先生"},
			expectedPageInfo: dto.PageInfo{PageNum: 1, LastPageNum: 1, Limit: 10, Offset: 0},
			expectedLen:      10,
		},
		{
			base: statusOKCheckParamBase{
				title:        "orderBy=name, order=ASCを指定して、正常にユーザーが取得できること",
				requestParam: requestParam{httpMethod: http.MethodGet, path: "/users?orderBy=name&order=ASC", body: nil}},
			expectedBody:     dto.UserResDto{Id: 7, Name: "べらぼう太郎"},
			expectedPageInfo: dto.PageInfo{PageNum: 1, LastPageNum: 1, Limit: 10, Offset: 0},
			expectedLen:      10,
		},
		{
			base: statusOKCheckParamBase{
				title:        "orderBy=name, order=DESCを指定して、正常にユーザーが取得できること",
				requestParam: requestParam{httpMethod: http.MethodGet, path: "/users?orderBy=name&order=DESC", body: nil}},
			expectedBody:     dto.UserResDto{Id: 4, Name: "腕時計両腕ちゃん"},
			expectedPageInfo: dto.PageInfo{PageNum: 1, LastPageNum: 1, Limit: 10, Offset: 0},
			expectedLen:      10,
		},
		{
			base: statusOKCheckParamBase{
				title:        "orderBy,order,page,limitを指定して、正常にユーザーが取得できること",
				requestParam: requestParam{httpMethod: http.MethodGet, path: "/users?orderBy=name&order=DESC&page=3&limit=2", body: nil}},
			expectedBody:     dto.UserResDto{Id: 2, Name: "トマト君"},
			expectedPageInfo: dto.PageInfo{PageNum: 3, LastPageNum: 5, Limit: 2, Offset: 4},
			expectedLen:      2,
		},
	}

	for _, p := range params {
		req := httptest.NewRequest(p.base.requestParam.httpMethod, p.base.requestParam.path, p.base.requestParam.body)
		rec := httptest.NewRecorder()
		_target.ServeHTTP(rec, req)

		// actual
		actualCode := rec.Code
		var actualBody dto.UserPage
		_ = json.Unmarshal(rec.Body.Bytes(), &actualBody)
		actualPageInfo := actualBody.PageInfo
		actualLen := len(actualBody.Users)
		actualBodyFirst := actualBody.Users[0]

		// expected
		expectedCode := http.StatusOK
		expectedPageInfo := p.expectedPageInfo
		expectedLen := p.expectedLen
		expectedBodyFirst := p.expectedBody

		// check
		fmt.Println(p.base.title)
		assert.Equal(t, expectedCode, actualCode)
		assert.Equal(t, expectedPageInfo, actualPageInfo)
		assert.Equal(t, expectedLen, actualLen)
		assert.Equal(t, expectedBodyFirst, actualBodyFirst)
	}
}

func TestUserController_GetAllUser_異常系(t *testing.T) {
	// setup
	var params = []errorCheckParam{
		{
			"pageのみ指定され、limitがないとき、400になること",
			[]requestParam{
				{httpMethod: http.MethodGet, path: "/users?page=5", body: nil},
			},
			http.StatusBadRequest,
			message.StatusBadRequest,
		},
		{
			"page, limitが数字ではないとき、400になること",
			[]requestParam{
				{httpMethod: http.MethodGet, path: "/users?limit=hoge", body: nil},
				{httpMethod: http.MethodGet, path: "/users?page=hoge&limit=5", body: nil},
			},
			http.StatusBadRequest,
			message.StatusBadRequest,
		},
		{
			"page, limitが0以下のとき、400になること",
			[]requestParam{
				{httpMethod: http.MethodGet, path: "/users?limit=0", body: nil},
				{httpMethod: http.MethodGet, path: "/users?limit=-1", body: nil},
				{httpMethod: http.MethodGet, path: "/users?page=0&limit=2", body: nil},
				{httpMethod: http.MethodGet, path: "/users?page=-1&limit=2", body: nil},
			},
			http.StatusBadRequest,
			message.StatusBadRequest,
		},
		{
			"orderByに存在しないColumnが指定されているとき, 400になること",
			[]requestParam{
				{httpMethod: http.MethodGet, path: "/users?orderBy=hoge", body: nil},
			},
			http.StatusBadRequest,
			message.StatusBadRequest,
		},
		{
			"orderにASC・desc以外が指定されているとき, 400になること",
			[]requestParam{
				{httpMethod: http.MethodGet, path: "/users?order=hoge", body: nil},
			},
			http.StatusBadRequest,
			message.StatusBadRequest,
		},
	}

	// check
	doErrorCheck(t, params)
}

/**************************************
	ユーザー登録
**************************************/
func TestUserController_CreateUser_正常系(t *testing.T) {
	// setup
	const initializedLocalDataRecordCounts = 10
	firstExpectedUserIdInt := initializedLocalDataRecordCounts + 1
	var paramsUserController = []statusOKCheckParamUser{
		{base: statusOKCheckParamBase{
			title:        "ユーザー登録①__正常にユーザーが登録されること",
			requestParam: requestParam{httpMethod: http.MethodPost, path: "/users/new", body: strings.NewReader(`{"name":"新規ユーザー太郎", "password":"NewUser1"}`)}},
			expectedBody: dto.UserResDto{Id: firstExpectedUserIdInt, Name: "新規ユーザー太郎"},
		},
		{base: statusOKCheckParamBase{
			title:        "ユーザー登録②__userIdに既存の値が指定されているときにも、正常にユーザーが登録されること（userIdが無視されること）",
			requestParam: requestParam{httpMethod: http.MethodPost, path: "/users/new", body: strings.NewReader(`{"id":1,"name":"二番煎次郎", "password":"NewUser2"}`)}},
			expectedBody: dto.UserResDto{Id: firstExpectedUserIdInt + 1, Name: "二番煎次郎"},
		},
	}

	var paramsAppController = []statusOKCheckParamBase{
		{
			title:        "ログイン確認①__正常にログインできること",
			requestParam: requestParam{httpMethod: http.MethodPost, path: "/login", body: strings.NewReader(`{"id": 11, "password":"NewUser1"}`)},
		},
		{
			title:        "ログイン確認②__正常にログインできること",
			requestParam: requestParam{httpMethod: http.MethodPost, path: "/login", body: strings.NewReader(`{"id": 12, "password":"NewUser2"}`)},
		},
	}

	// check
	// 2ユーザーを登録する
	doStatusOKCheck__UserController(t, paramsUserController, _checkExisting, false)
	// 2ユーザーでログインする
	doStatusOKCheck__ApplicationController(t, paramsAppController, false)

	// clean
	localdata.InitializeLocalData()
}
func TestUserController_CreateUser_異常系(t *testing.T) {
	// setup
	userNames1 := []string{" ", "　", "\n", "新規ユーザー 太郎", "新規ユーザー　太郎", "新規ユーザー\n太郎", " 新規ユーザー太郎", "新規ユーザー太郎　", " 新規ユーザー太郎\n"}
	passwords1 := makeSameStringList(9, "Test1111")

	userNames2 := makeSameStringList(9, "新規ユーザー太郎")
	passwords2 := []string{" ", "　", "\n", "Test 1111", "Test　1111", "Test\n1111", " Test1111", "Test1111　", " Test1111\n"}

	userNames3 := makeSameStringList(3, "新規ユーザー太郎")
	passwords3 := []string{"TEST1111", "test1111", "Testaaaa"}

	userNames4 := makeSameStringList(3, "新規ユーザー太郎")
	passwords4 := []string{"Test1111あ", "Test1111漢字", "Test1111(*"}

	var params = []errorCheckParam{
		{
			"userNameが空文字のとき、400になること",
			[]requestParam{{httpMethod: http.MethodPost, path: "/users/new", body: strings.NewReader(`{"name":"", "password":"Test1111"}`)}},
			http.StatusBadRequest,
			message.StatusBadRequest,
		},
		{
			"userNameに半角・全角スペース、改行が含まれているとき、400になること",
			makeInputs(http.MethodPost, "/users/new", makeBodyList(makeUserReceiveDtoJsonList(userNames1, passwords1))),
			http.StatusBadRequest,
			message.StatusBadRequest,
		},
		{
			"userNameが9文字以上のとき、400になること",
			[]requestParam{{httpMethod: http.MethodPost, path: "/users/new", body: strings.NewReader(`{"name":"123456789", "password":"Test1111"}`)}},
			http.StatusBadRequest,
			message.StatusBadRequest,
		},
		{
			"passwordが空文字のとき、400になること",
			[]requestParam{{httpMethod: http.MethodPost, path: "/users/new", body: strings.NewReader(`{"name":"test", "password":""}`)}},
			http.StatusBadRequest,
			message.StatusBadRequest,
		},
		{
			"passwordに半角・全角スペース、改行が含まれているとき、400になること",
			makeInputs(http.MethodPost, "/users/new", makeBodyList(makeUserReceiveDtoJsonList(userNames2, passwords2))),
			http.StatusBadRequest,
			message.StatusBadRequest,
		},
		{
			"passwordが7文字以下のとき、400になること",
			[]requestParam{{httpMethod: http.MethodPost, path: "/users/new", body: strings.NewReader(`{"name":"新規ユーザー太郎", "password":"Test111"}`)}},
			http.StatusBadRequest,
			message.StatusBadRequest,
		},
		{
			"passwordが英数大文字小文字を最低１つずつ含んでいないとき、400になること",
			makeInputs(http.MethodPost, "/users/new", makeBodyList(makeUserReceiveDtoJsonList(userNames3, passwords3))),
			http.StatusBadRequest,
			message.StatusBadRequest,
		},
		{
			"passwordが英数大文字小文字以外の文字を含んでいるとき、400になること",
			makeInputs(http.MethodPost, "/users/new", makeBodyList(makeUserReceiveDtoJsonList(userNames4, passwords4))),
			http.StatusBadRequest,
			message.StatusBadRequest,
		},
	}

	// check
	doErrorCheck(t, params)
}

/**************************************
	ユーザー削除
**************************************/
func TestUserController_DeleteUser_正常系(t *testing.T) {
	// setup
	var params = []statusOKCheckParamUser{
		{base: statusOKCheckParamBase{
			title:        "正常にユーザーが削除できること_レスポンスチェック",
			requestParam: requestParam{httpMethod: http.MethodDelete, path: "/users/1", body: nil}},
			expectedBody: dto.UserResDto{Id: 1, Name: "まるお"}},
	}

	// check
	doStatusOKCheck__UserController(t, params, _checkNotExiting, true)
}

func TestUserController_DeleteUser_異常系(t *testing.T) {
	// setup
	var params = []errorCheckParam{
		{
			"存在しないuserIdのとき、404になること",
			[]requestParam{{httpMethod: http.MethodDelete, path: "/users/9999", body: nil}},
			http.StatusNotFound,
			message.UserNotFound,
		},
	}

	// check
	doErrorCheck(t, params)
}

/**************************************
	ユーザー更新
**************************************/
func TestUserController_UpdateUser_正常系(t *testing.T) {
	// setup
	var paramsUserController = []statusOKCheckParamUser{
		{base: statusOKCheckParamBase{
			title:        "正常にユーザー名が更新できること",
			requestParam: requestParam{httpMethod: http.MethodPut, path: "/users/1/update", body: strings.NewReader(`{"name":"ハルキゲニア"}`)}},
			expectedBody: dto.UserResDto{Id: 1, Name: "ハルキゲニア"},
		},
		{base: statusOKCheckParamBase{
			title:        "正常にパスワードが更新できること",
			requestParam: requestParam{httpMethod: http.MethodPut, path: "/users/1/update", body: strings.NewReader(`{"password":"HogeHoge123"}`)}},
			expectedBody: dto.UserResDto{Id: 1, Name: "ハルキゲニア"},
		},
		{base: statusOKCheckParamBase{
			title:        "正常にユーザー名とパスワードが更新できること",
			requestParam: requestParam{httpMethod: http.MethodPut, path: "/users/2/update", body: strings.NewReader(`{"name":"バオバブの木", "password":"BaoBab123"}`)}},
			expectedBody: dto.UserResDto{Id: 2, Name: "バオバブの木"},
		},
		{base: statusOKCheckParamBase{
			title:        "ボディのuserIdにパスと異なる値が指定されているときにも、正常にパスで指定したユーザーが更新されること（ボディのuserIdが無視されること）",
			requestParam: requestParam{httpMethod: http.MethodPut, path: "/users/3/update", body: strings.NewReader(`{"id":4,"name":"ビクトリア3世"}`)}},
			expectedBody: dto.UserResDto{Id: 3, Name: "ビクトリア3世"},
		},
	}
	var paramsAppController = []statusOKCheckParamBase{
		{
			title:        "ログイン確認①__パスワードを変更した後に、正常にログインできること",
			requestParam: requestParam{httpMethod: http.MethodPost, path: "/login", body: strings.NewReader(`{"id": 1, "password":"HogeHoge123"}`)},
		},
		{
			title:        "ログイン確認②__パスワードを変更した後に、正常にログインできること",
			requestParam: requestParam{httpMethod: http.MethodPost, path: "/login", body: strings.NewReader(`{"id": 2, "password":"BaoBab123"}`)},
		},
	}

	// check
	// 一通り更新系チェックをする
	doStatusOKCheck__UserController(t, paramsUserController, _checkExisting, false)
	// パスワードを変更したユーザーにて、ログイン確認する
	doStatusOKCheck__ApplicationController(t, paramsAppController, false)

	// clean
	localdata.InitializeLocalData()
}

func TestUserController_UpdateUser_異常系(t *testing.T) {
	// setup
	userNames1 := []string{" ", "　", "\n", "新規ユーザー 太郎", "新規ユーザー　太郎", "新規ユーザー\n太郎", " 新規ユーザー太郎", "新規ユーザー太郎　", " 新規ユーザー太郎\n"}
	passwords1 := makeSameStringList(9, "Test1111")

	userNames2 := makeSameStringList(9, "新規ユーザー太郎")
	passwords2 := []string{" ", "　", "\n", "Test 1111", "Test　1111", "Test\n1111", " Test1111", "Test1111　", " Test1111\n"}

	userNames3 := makeSameStringList(3, "新規ユーザー太郎")
	passwords3 := []string{"TEST1111", "test1111", "Testaaaa"}

	userNames4 := makeSameStringList(3, "新規ユーザー太郎")
	passwords4 := []string{"Test1111あ", "Test1111漢字", "Test1111(*"}

	var params = []errorCheckParam{
		{
			"存在しないuserIdのとき、404になること",
			[]requestParam{{httpMethod: http.MethodPut, path: "/users/9999/update", body: strings.NewReader(`{"name":"ハルキゲニア"}`)}},
			http.StatusNotFound,
			message.UserNotFound,
		},
		{
			"userNameが空文字のとき、400になること",
			[]requestParam{{httpMethod: http.MethodPut, path: "/users/1/update", body: strings.NewReader(`{"name":""}`)}},
			http.StatusBadRequest,
			message.StatusBadRequest,
		},
		{
			"userNameがnilのとき、400になること",
			[]requestParam{{httpMethod: http.MethodPut, path: "/users/1/update", body: strings.NewReader(`{"name":}`)}},
			http.StatusBadRequest,
			message.StatusBadRequest,
		},
		{
			"userNameに半角・全角スペース、改行が含まれているとき、400になること",
			makeInputs(http.MethodPut, "/users/1/update", makeBodyList(makeUserReceiveDtoJsonList(userNames1, passwords1))),
			http.StatusBadRequest,
			message.StatusBadRequest,
		},
		{
			"userNameが9文字以上のとき、400になること",
			[]requestParam{{httpMethod: http.MethodPut, path: "/users/1/update", body: strings.NewReader(`{"name":"123456789"}`)}},
			http.StatusBadRequest,
			message.StatusBadRequest,
		},
		{
			"passwordが空文字のとき、400になること",
			[]requestParam{{httpMethod: http.MethodPut, path: "/users/1/update", body: strings.NewReader(`{"password":""}`)}},
			http.StatusBadRequest,
			message.StatusBadRequest,
		},
		{
			"passwordがnilのとき、400になること",
			[]requestParam{{httpMethod: http.MethodPut, path: "/users/1/update", body: strings.NewReader(`{"password":}`)}},
			http.StatusBadRequest,
			message.StatusBadRequest,
		},
		{
			"passwordに半角・全角スペース、改行が含まれているとき、400になること",
			makeInputs(http.MethodPut, "/users/1/update", makeBodyList(makeUserReceiveDtoJsonList(userNames2, passwords2))),
			http.StatusBadRequest,
			message.StatusBadRequest,
		},
		{
			"passwordが7文字以下のとき、400になること",
			[]requestParam{{httpMethod: http.MethodPut, path: "/users/1/update", body: strings.NewReader(`{"password":"Test111"}`)}},
			http.StatusBadRequest,
			message.StatusBadRequest,
		},
		{
			"passwordが英数大文字小文字を最低１つずつ含んでいないとき、400になること",
			makeInputs(http.MethodPut, "/users/1/update", makeBodyList(makeUserReceiveDtoJsonList(userNames3, passwords3))),
			http.StatusBadRequest,
			message.StatusBadRequest,
		},
		{
			"passwordが英数大文字小文字以外の文字を含んでいるとき、400になること",
			makeInputs(http.MethodPut, "/users/1/update", makeBodyList(makeUserReceiveDtoJsonList(userNames4, passwords4))),
			http.StatusBadRequest,
			message.StatusBadRequest,
		},
		{
			"ユーザー名もパスワードも指定されていないとき、400になること",
			[]requestParam{{httpMethod: http.MethodPut, path: "/users/1/update", body: strings.NewReader(`{}`)}},
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
func doStatusOKCheck__UserController(t *testing.T, params []statusOKCheckParamUser, recordCheckPattern recordCheckPattern, doCleanData bool) {
	for _, p := range params {
		// setup
		req := httptest.NewRequest(p.base.requestParam.httpMethod, p.base.requestParam.path, p.base.requestParam.body)
		if p.base.requestParam.httpMethod == http.MethodPost || p.base.requestParam.httpMethod == http.MethodPut {
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		}
		rec := httptest.NewRecorder()
		_target.ServeHTTP(rec, req)

		// actual
		actualCode := rec.Code
		var actualBody dto.UserResDto
		_ = json.Unmarshal(rec.Body.Bytes(), &actualBody)

		// expected
		const expectedCode = http.StatusOK
		expectedBody := p.expectedBody
		expectedRecordId := expectedBody.Id

		// check①
		fmt.Println(p.base.title)
		assert.Equal(t, expectedCode, actualCode)
		assert.Equal(t, expectedBody, actualBody)

		// check② レコード存在チェック
		switch recordCheckPattern {
		case _doNothing:
		case _checkExisting:
			actualRecord, _ := _userUseCase.GetUser(expectedRecordId)
			assert.Equal(t, expectedBody, actualRecord)
		case _checkNotExiting:
			doRecordNotExistingCheck(t, expectedRecordId)
		}

		// clean
		if doCleanData {
			localdata.InitializeLocalData()
		}
	}
}

func makeUserReceiveDtoJsonList(userNames []string, passwords []string) [][]byte {
	list := make([][]byte, len(userNames))
	for i, n := range userNames {
		userDto := dto.UserReceiveDto{
			Id:       1,
			Name:     n,
			Password: passwords[i],
		}
		b, _ := json.Marshal(userDto)
		list[i] = b
	}
	return list
}
