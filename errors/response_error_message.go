package errors

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/labstack/echo/v4"
	"net/http"
)

type ResponseErrorMessage struct {
	status  int
	message string
}

type responseErrorMessageJSON struct {
	Status  int
	Message string
}

func ResponseErrorJSON(c echo.Context, err error, errMessage string) error {
	var httpStatus int
	var appErr *AppError
	// httpStatusを取得する
	if errors.As(err, &appErr) {
		httpStatus = appErr.GetHttpStatus()
	} else {
		// appErr以外はすべて想定外とする
		httpStatus = http.StatusInternalServerError
	}

	// ログにエラー内容を書き出し
	fmt.Print("\n")
	fmt.Println("-----------------------------------")
	fmt.Println(err.Error())
	fmt.Print("\n")

	// Response
	return c.JSON(httpStatus, newResponseErrorMessage(httpStatus, errMessage))
}

func (r ResponseErrorMessage) MarshalJSON() ([]byte, error) {
	value, err := json.Marshal(&responseErrorMessageJSON{
		Status:  r.status,
		Message: r.message,
	})
	return value, err
}

func (r *ResponseErrorMessage) UnmarshalJSON(b []byte) error {
	var responseErrorMessageJSON responseErrorMessageJSON
	if err := json.Unmarshal(b, &responseErrorMessageJSON); err != nil {
		return err
	}

	r.status = responseErrorMessageJSON.Status
	r.message = responseErrorMessageJSON.Message

	return nil
}

func (r ResponseErrorMessage) GetStatus() int {
	return r.status
}

func (r ResponseErrorMessage) GetMessage() string {
	return r.message
}

func newResponseErrorMessage(status int, message string) *ResponseErrorMessage {
	return &ResponseErrorMessage{
		status:  status,
		message: message,
	}
}
