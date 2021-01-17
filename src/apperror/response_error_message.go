package apperror

import (
	"encoding/json"
	"errors"
	"fmt"
	"go-app-template/src/apperror/message"
	"net/http"
	"runtime"

	"github.com/labstack/echo/v4"
)

type ResponseErrorMessage struct {
	status  HttpStatus
	message string
}

type responseErrorMessageJSON struct {
	Status  int
	Message string
}

func ResponseErrorJSON(c echo.Context, err error, errMessage message.Message) error {
	var (
		httpStatus HttpStatus
		appErr     *AppError
	)

	// httpStatusを取得する
	if errors.As(err, &appErr) && appErr.GetHttpStatus().isEvaluated() {
		httpStatus = appErr.GetHttpStatus()
	} else {
		// 予期せぬエラー
		httpStatus = http.StatusInternalServerError
		errMessage = message.SystemError
	}

	// ログにエラー内容を書き出し
	fmt.Print("\n")
	fmt.Println("-----------------------------------")
	_, fileAbsPath, line, ok := runtime.Caller(1)
	if !ok {
		fileAbsPath = "???"
		line = 0
	}
	fmt.Printf("Error has been responsed at..., Location: %v:%v\n", fileAbsPath, line)
	fmt.Println(err.Error())
	fmt.Print("\n")

	// Response
	return c.JSON(int(httpStatus), newResponseErrorMessage(httpStatus, errMessage.String()))
}

func newResponseErrorMessage(status HttpStatus, message string) *ResponseErrorMessage {
	return &ResponseErrorMessage{
		status:  status,
		message: message,
	}
}

/**************************************
	JSON
**************************************/
func (r ResponseErrorMessage) MarshalJSON() ([]byte, error) {
	value, err := json.Marshal(&responseErrorMessageJSON{
		Status:  int(r.status),
		Message: r.message,
	})
	return value, err
}

func (r *ResponseErrorMessage) UnmarshalJSON(b []byte) error {
	var responseErrorMessageJSON responseErrorMessageJSON
	if err := json.Unmarshal(b, &responseErrorMessageJSON); err != nil {
		return err
	}

	r.status = HttpStatus(responseErrorMessageJSON.Status)
	r.message = responseErrorMessageJSON.Message

	return nil
}

/**************************************
	Getter & Setter
**************************************/
func (r ResponseErrorMessage) GetStatus() HttpStatus {
	return r.status
}

func (r ResponseErrorMessage) GetMessage() string {
	return r.message
}

/**************************************
	private
**************************************/
