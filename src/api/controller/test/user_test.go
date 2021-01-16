package test

import (
	"encoding/json"
	"errors"
	"fmt"
	"go-app-template/src/api/controller/dto"
	"go-app-template/src/apperror"
	"go-app-template/src/apperror/message"
	"go-app-template/src/apputil"
	"go-app-template/src/config/db/localdata"
	"go-app-template/src/config/route"
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

type statusOKCheckParam struct {
	title             string
	input             input
	expectedUserIdInt int
	expectedName      string
}

type errorCheckParam struct {
	title           string
	input           []input
	expectedCode    int
	expectedMessage message.Message
}

type RecordCheckPattern int

const (
	DoNothing RecordCheckPattern = iota + 1
	ExistingCheck
	NotExistingCheck
)

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
	var params = []statusOKCheckParam{
		{

			"正常にユーザーが取得できること",
			input{httpMethod: http.MethodGet, path: "/users/1", body: nil},
			1,
			"まるお",
		},
	}

	// check
	doStatusOKCheck(t, params, ExistingCheck)
}

func TestUserController_GetUser_異常系(t *testing.T) {
	// setup
	var params = []errorCheckParam{
		{
			"存在しないuserIdのとき、404になること",
			[]input{{httpMethod: http.MethodGet, path: "/users/9999", body: nil}},
			http.StatusNotFound,
			message.UserNotFound,
		},
		{
			"userIdが数字ではないとき、400になること",
			[]input{{httpMethod: http.MethodGet, path: "/users/hogehoge", body: nil}},
			http.StatusBadRequest,
			message.InvalidUserId,
		},
		{
			"userIdが0以下のとき、400になること",
			[]input{
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

func TestUserController_GetAllUser_正常系(t *testing.T) {
	// setup
	params := []struct {
		base             statusOKCheckParam
		expectedPageInfo dto.PageInfo
		expectedLen      int
	}{
		{statusOKCheckParam{
			"正常にユーザーが全件取得できること",
			input{httpMethod: http.MethodGet, path: "/users", body: nil},
			1,
			"まるお"},
			dto.PageInfo{
				PageNum:     1,
				LastPageNum: 1,
				Limit:       10,
				Offset:      0,
			},
			10,
		},
		{statusOKCheckParam{
			"limitのみを指定して、limit件数分、正常にユーザーが取得できること",
			input{httpMethod: http.MethodGet, path: "/users?limit=5", body: nil},
			1,
			"まるお"},
			dto.PageInfo{
				PageNum:     1,
				LastPageNum: 2,
				Limit:       5,
				Offset:      0,
			},
			5,
		},
		{statusOKCheckParam{
			"page, limitを指定して、正常にユーザーが取得できること①",
			input{httpMethod: http.MethodGet, path: "/users?page=1&limit=3", body: nil},
			1,
			"まるお"},
			dto.PageInfo{
				PageNum:     1,
				LastPageNum: 4,
				Limit:       3,
				Offset:      0,
			},
			3,
		},
		{statusOKCheckParam{
			"page, limitを指定して、正常にユーザーが取得できること②",
			input{httpMethod: http.MethodGet, path: "/users?page=3&limit=3", body: nil},
			7,
			"べらぼう太郎"},
			dto.PageInfo{
				PageNum:     3,
				LastPageNum: 4,
				Limit:       3,
				Offset:      6,
			},
			3,
		},
		{statusOKCheckParam{
			"page, limitを指定して、正常にユーザーが取得できること③",
			input{httpMethod: http.MethodGet, path: "/users?page=4&limit=3", body: nil},
			10,
			"先生"},
			dto.PageInfo{
				PageNum:     4,
				LastPageNum: 4,
				Limit:       3,
				Offset:      9,
			},
			1,
		},
		{statusOKCheckParam{
			"最終ページを超過したpageを指定して、最終ページが正常に取得できること",
			input{httpMethod: http.MethodGet, path: "/users?page=10&limit=9", body: nil},
			10,
			"先生"},
			dto.PageInfo{
				PageNum:     2,
				LastPageNum: 2,
				Limit:       9,
				Offset:      9,
			},
			1,
		},
		{statusOKCheckParam{
			"orderBy=idのみを指定して、id・ASCで、正常にユーザーが取得できること",
			input{httpMethod: http.MethodGet, path: "/users?orderBy=id", body: nil},
			1,
			"まるお"},
			dto.PageInfo{
				PageNum:     1,
				LastPageNum: 1,
				Limit:       10,
				Offset:      0,
			},
			10,
		},
		{statusOKCheckParam{
			"orderBy=nameのみを指定して、userName・ASCで、正常にユーザーが取得できること",
			input{httpMethod: http.MethodGet, path: "/users?orderBy=name", body: nil},
			7,
			"べらぼう太郎"},
			dto.PageInfo{
				PageNum:     1,
				LastPageNum: 1,
				Limit:       10,
				Offset:      0,
			},
			10,
		},
		{statusOKCheckParam{
			"order=ASCのみを指定して、id・ASCで、正常にユーザーが取得できること",
			input{httpMethod: http.MethodGet, path: "/users?order=ASC", body: nil},
			1,
			"まるお"},
			dto.PageInfo{
				PageNum:     1,
				LastPageNum: 1,
				Limit:       10,
				Offset:      0,
			},
			10,
		},
		{statusOKCheckParam{
			"order=DESCのみを指定して、id・DESCで、正常にユーザーが取得できること",
			input{httpMethod: http.MethodGet, path: "/users?order=DESC", body: nil},
			10,
			"先生"},
			dto.PageInfo{
				PageNum:     1,
				LastPageNum: 1,
				Limit:       10,
				Offset:      0,
			},
			10,
		},
		{statusOKCheckParam{
			"orderBy=name, order=ASCを指定して、正常にユーザーが取得できること",
			input{httpMethod: http.MethodGet, path: "/users?orderBy=name&order=ASC", body: nil},
			7,
			"べらぼう太郎"},
			dto.PageInfo{
				PageNum:     1,
				LastPageNum: 1,
				Limit:       10,
				Offset:      0,
			},
			10,
		},
		{statusOKCheckParam{
			"orderBy=name, order=DESCを指定して、正常にユーザーが取得できること",
			input{httpMethod: http.MethodGet, path: "/users?orderBy=name&order=DESC", body: nil},
			4,
			"腕時計両腕ちゃん"},
			dto.PageInfo{
				PageNum:     1,
				LastPageNum: 1,
				Limit:       10,
				Offset:      0,
			},
			10,
		},
		{statusOKCheckParam{
			"orderBy,order,page,limitを指定して、正常にユーザーが取得できること",
			input{httpMethod: http.MethodGet, path: "/users?orderBy=name&order=DESC&page=3&limit=2", body: nil},
			2,
			"トマト君"},
			dto.PageInfo{
				PageNum:     3,
				LastPageNum: 5,
				Limit:       2,
				Offset:      4,
			},
			2,
		},
	}

	for _, p := range params {
		req := httptest.NewRequest(p.base.input.httpMethod, p.base.input.path, p.base.input.body)
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
		expectedBodyFirst := dto.UserReceiveDto{
			Id:   p.base.expectedUserIdInt,
			Name: p.base.expectedName,
		}

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
			[]input{
				{httpMethod: http.MethodGet, path: "/users?page=5", body: nil},
			},
			http.StatusBadRequest,
			message.StatusBadRequest,
		},
		{
			"page, limitが数字ではないとき、400になること",
			[]input{
				{httpMethod: http.MethodGet, path: "/users?limit=hoge", body: nil},
				{httpMethod: http.MethodGet, path: "/users?page=hoge&limit=5", body: nil},
			},
			http.StatusBadRequest,
			message.StatusBadRequest,
		},
		{
			"page, limitが0以下のとき、400になること",
			[]input{
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
			[]input{
				{httpMethod: http.MethodGet, path: "/users?orderBy=hoge", body: nil},
			},
			http.StatusBadRequest,
			message.StatusBadRequest,
		},
		{
			"orderにASC・desc以外が指定されているとき, 400になること",
			[]input{
				{httpMethod: http.MethodGet, path: "/users?order=hoge", body: nil},
			},
			http.StatusBadRequest,
			message.StatusBadRequest,
		},
	}

	// check
	doErrorCheck(t, params)
}

func TestUserController_CreateUser_正常系(t *testing.T) {
	// setup
	const initializedLocalDataRecordCounts = 10
	expectedUserIdInt := initializedLocalDataRecordCounts + 1
	var params = []statusOKCheckParam{
		{
			"正常にユーザーが登録されること",
			input{httpMethod: http.MethodPost, path: "/users/new", body: strings.NewReader(`{"name":"新規ユーザー太郎"}`)},
			expectedUserIdInt,
			"新規ユーザー太郎",
		},
		{
			"userIdに既存の値が指定されているときにも、正常にユーザーが登録されること（userIdが無視されること）",
			input{httpMethod: http.MethodPost, path: "/users/new", body: strings.NewReader(`{"id":1,"name":"新規ユーザー太郎"}`)},
			expectedUserIdInt,
			"新規ユーザー太郎",
		},
	}

	// check
	doStatusOKCheck(t, params, ExistingCheck)
}
func TestUserController_CreateUser_異常系(t *testing.T) {
	// setup
	userNames := []string{" ", "　", "\n", "新規ユーザー 太郎", "新規ユーザー　太郎", "新規ユーザー\n太郎", " 新規ユーザー太郎", "新規ユーザー太郎　", " 新規ユーザー\n太郎　"}

	var params = []errorCheckParam{
		{
			"userNameが空文字のとき、400になること",
			[]input{{httpMethod: http.MethodPost, path: "/users/new", body: strings.NewReader(`{"name":""}`)}},
			http.StatusBadRequest,
			message.StatusBadRequest,
		},
		{
			"userNameに半角・全角スペース、改行が含まれているとき、400になること",
			makeInputs(http.MethodPost, "/users/new", makeBodyList(makeUserDtoJsonList(userNames))),
			http.StatusBadRequest,
			message.StatusBadRequest,
		},
		{
			"userNameが9文字以上のとき、400になること",
			[]input{{httpMethod: http.MethodPost, path: "/users/new", body: strings.NewReader(`{"name":"123456789"}`)}},
			http.StatusBadRequest,
			message.StatusBadRequest,
		},
	}

	// check
	doErrorCheck(t, params)
}

func TestUserController_DeleteUser_正常系(t *testing.T) {
	// setup
	var params = []statusOKCheckParam{
		{
			"正常にユーザーが削除できること①_レスポンスチェック",
			input{httpMethod: http.MethodDelete, path: "/users/1", body: nil},
			1,
			"まるお",
		},
	}

	// check
	doStatusOKCheck(t, params, NotExistingCheck)
}

func TestUserController_DeleteUser_異常系(t *testing.T) {
	// setup
	var params = []errorCheckParam{
		{
			"存在しないuserIdのとき、404になること",
			[]input{{httpMethod: http.MethodDelete, path: "/users/9999", body: nil}},
			http.StatusNotFound,
			message.UserNotFound,
		},
	}

	// check
	doErrorCheck(t, params)
}

func TestUserController_UpdateUser_正常系(t *testing.T) {
	// setup
	var params = []statusOKCheckParam{
		{
			"正常にユーザー名が更新できること",
			input{httpMethod: http.MethodPut, path: "/users/1/update", body: strings.NewReader(`{"name":"ハルキゲニア"}`)},
			1,
			"ハルキゲニア",
		},
		{
			"ボディのuserIdにパスと異なる値が指定されているときにも、正常にパスで指定したユーザーが更新されること（ボディのuserIdが無視されること）",
			input{httpMethod: http.MethodPut, path: "/users/2/update", body: strings.NewReader(`{"id":3,"name":"ビクトリア3世"}`)},
			2,
			"ビクトリア3世",
		},
	}

	// check
	doStatusOKCheck(t, params, ExistingCheck)
}

func TestUserController_UpdateUser_異常系(t *testing.T) {
	// setup
	userNames := []string{" ", "　", "\n", "新規ユーザー 太郎", "新規ユーザー　太郎", "新規ユーザー\n太郎", " 新規ユーザー太郎", "新規ユーザー太郎　", " 新規ユーザー\n太郎　"}

	var params = []errorCheckParam{
		{
			"存在しないuserIdのとき、404になること",
			[]input{{httpMethod: http.MethodPut, path: "/users/9999/update", body: strings.NewReader(`{"name":"ハルキゲニア"}`)}},
			http.StatusNotFound,
			message.UserNotFound,
		},
		{
			"userNameが空文字のとき、400になること",
			[]input{{httpMethod: http.MethodPut, path: "/users/1/update", body: strings.NewReader(`{"name":""}`)}},
			http.StatusBadRequest,
			message.StatusBadRequest,
		},
		{
			"userNameがnilのとき、400になること",
			[]input{{httpMethod: http.MethodPut, path: "/users/1/update", body: strings.NewReader(`{"name":}`)}},
			http.StatusBadRequest,
			message.StatusBadRequest,
		},
		{
			"userNameに半角・全角スペース、改行が含まれているとき、400になること",
			makeInputs(http.MethodPut, "/users/1/update", makeBodyList(makeUserDtoJsonList(userNames))),
			http.StatusBadRequest,
			message.StatusBadRequest,
		},
		{
			"userNameが9文字以上のとき、400になること",
			[]input{{httpMethod: http.MethodPut, path: "/users/1/update", body: strings.NewReader(`{"name":"123456789"}`)}},
			http.StatusBadRequest,
			message.StatusBadRequest,
		},
	}

	// check
	doErrorCheck(t, params)
}

func doStatusOKCheck(t *testing.T, params []statusOKCheckParam, recordCheckPattern RecordCheckPattern) {
	for _, p := range params {
		req := httptest.NewRequest(p.input.httpMethod, p.input.path, p.input.body)
		if p.input.httpMethod == http.MethodPost || p.input.httpMethod == http.MethodPut {
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		}
		rec := httptest.NewRecorder()
		_target.ServeHTTP(rec, req)

		// actual
		actualCode := rec.Code
		var actualBody dto.UserReceiveDto
		_ = json.Unmarshal(rec.Body.Bytes(), &actualBody)

		// expected
		expectedCode := http.StatusOK
		expectedBody := &dto.UserReceiveDto{
			Id:   p.expectedUserIdInt,
			Name: p.expectedName,
		}

		// check①
		fmt.Println(p.title)
		assert.Equal(t, expectedCode, actualCode)
		assert.Equal(t, *expectedBody, actualBody)

		// check② レコード存在チェック
		switch recordCheckPattern {
		case DoNothing:
		case ExistingCheck:
			actualRecord, _ := _userUseCase.GetUser(p.expectedUserIdInt)
			assert.Equal(t, *expectedBody, actualRecord)
		case NotExistingCheck:
			doRecordNotExistingCheck(t, p.expectedUserIdInt)
		}

		// clean
		localdata.InitializeLocalData()
	}
}

func doErrorCheck(t *testing.T, params []errorCheckParam) {
	for _, p := range params {
		for _, ip := range p.input {
			req := httptest.NewRequest(ip.httpMethod, ip.path, ip.body)
			if ip.httpMethod == http.MethodPost || ip.httpMethod == http.MethodPut {
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

			// clean
			localdata.InitializeLocalData()
		}
	}
}

func doRecordNotExistingCheck(t *testing.T, recordId int) {
	// actual
	_, actualErr := _userUseCase.GetUser(recordId)
	var actualErrCode int
	var actualErrMessage string
	var appErr *apperror.AppError
	if errors.As(actualErr, &appErr) {
		actualErrCode = int(appErr.GetHttpStatus())
		actualErrMessage = appErr.ErrorWithoutLocation()
	}

	// expected
	const expectedErrCode = http.StatusNotFound
	const expectedErrMessage = "Error: record not found"

	// check
	assert.Equal(t, expectedErrCode, actualErrCode)
	assert.Equal(t, expectedErrMessage, actualErrMessage)
}

func makePathParamInputs(httpMethod string, pathBase string, pathParams []string, body io.Reader) []input {
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

func makeInputs(method string, path string, body []io.Reader) []input {
	inputs := make([]input, len(body))
	for i, p := range body {
		input := &input{
			httpMethod: method,
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
		userDto := dto.UserReceiveDto{
			Id:   1,
			Name: n,
		}
		j, _ := json.Marshal(userDto)
		list[i] = j
	}
	return list
}
