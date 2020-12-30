package mock

import (
	"fmt"
)

type AppErrorMock struct {
	err      error
	status   int
	fileName string
	line     int
}

func (e *AppErrorMock) Error() string {
	return fmt.Sprintf("Error: %v,\nLocation: %v:%v", e.err.Error(), e.fileName, e.line)
}

func (e *AppErrorMock) GetHttpStatus() int {
	return e.status
}

func NewAppErrorMock(err error, status int, fileName string, line int) *AppErrorMock {
	return &AppErrorMock{
		err:      err,
		status:   status,
		fileName: fileName,
		line:     line,
	}
}
