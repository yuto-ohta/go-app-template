package errors

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/labstack/echo/v4"
	"net/http"
)

type responseErrorMessage struct {
	status  int
	message string
}

type responseErrorMessageJSON struct {
	Status  int
	Message string
}

func (a responseErrorMessage) MarshalJSON() ([]byte, error) {
	value, err := json.Marshal(&responseErrorMessageJSON{
		Status:  a.status,
		Message: a.message,
	})
	return value, err
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

func newResponseErrorMessage(status int, message string) *responseErrorMessage {
	return &responseErrorMessage{
		status:  status,
		message: message,
	}
}
