package apperror

import (
	"errors"
	"fmt"
	"go-app-template/src/config"
	"runtime"
	"strings"
)

const (
	_notEvaluated = -1
)

type AppError struct {
	err      error
	status   int
	fileName string
	line     int
}

func (e *AppError) Error() string {
	wrappedError := e.err

	// 内側がappErrorである限り、ループする
	var appErr *AppError
	if errors.As(wrappedError, &appErr) {
		return appErr.Error()
	}

	// 内側が一番下の普通のerrorになったときに、プリントする
	return fmt.Sprintf("Error: %v,\nLocation: %v:%v", wrappedError.Error(), e.fileName, e.line)
}

func (e *AppError) ErrorWithoutLocation() string {
	wrappedError := e.err

	// 内側がappErrorである限り、ループする
	var appErr *AppError
	if errors.As(wrappedError, &appErr) {
		return appErr.ErrorWithoutLocation()
	}

	// 内側が一番下の普通のerrorになったときに、プリントする
	return fmt.Sprintf("Error: %v", wrappedError.Error())
}

func (e AppError) Is(target error) bool {
	_, ok := target.(*AppError)
	return ok
}

func (e AppError) Unwrap() error {
	return e.err
}

func NewAppError(err error) *AppError {
	fileName, line := getCallerData(2)
	return &AppError{
		err:      err,
		status:   _notEvaluated,
		fileName: fileName,
		line:     line,
	}
}

func NewAppErrorWithStatus(err error, status int) *AppError {
	fileName, line := getCallerData(2)
	return &AppError{
		err:      err,
		status:   status,
		fileName: fileName,
		line:     line,
	}
}

func (e AppError) GetHttpStatus() int {
	return e.status
}

func (e AppError) isStatusEvaluated() bool {
	return e.status != _notEvaluated
}

func getCallerData(skip int) (string, int) {
	_, fileAbsPath, line, ok := runtime.Caller(skip)
	if !ok {
		return "???", 0
	}

	projectName := config.GetConfig()["project_name"].(string)
	fileProjectPath := fileAbsPath[strings.Index(fileAbsPath, projectName):]

	return fileProjectPath, line
}
