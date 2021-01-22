package test

import (
	"errors"
	"fmt"
	"github.com/labstack/echo/v4"
	"go-app-template/src/apperror"
	"go-app-template/src/apperror/message"
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

	"github.com/stretchr/testify/assert"
)

var (
	_target      = route.NewRouter()
	_userUseCase = *impl.NewUserUseCaseImpl(infrastructure.NewUserRepositoryImpl())
)

const (
	_doNothing recordCheckPattern = iota + 1
	_checkExisting
	_checkNotExiting
)

type recordCheckPattern int

type requestParam struct {
	httpMethod string
	path       string
	body       io.Reader
}

/*
	※各テストファイルにて、expectedBodyを加える
	ex)
		type statusOKCheckParamUser struct {
			base         statusOKCheckParamBase
			expectedBody dto.UserResDto
		}
*/
type statusOKCheckParamBase struct {
	title        string
	requestParam requestParam
}

type errorCheckParam struct {
	title           string
	requestParams   []requestParam
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

/**************************************
	private
**************************************/
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

func doErrorCheck(t *testing.T, params []errorCheckParam) {
	for _, p := range params {
		for _, ip := range p.requestParams {
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

func makeInputs(method string, path string, body []io.Reader) []requestParam {
	inputs := make([]requestParam, len(body))
	for i, p := range body {
		input := &requestParam{
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

func makeSameStringList(length int, str string) []string {
	res := make([]string, length)
	for i := 0; i < length; i++ {
		res[i] = str
	}
	return res
}
