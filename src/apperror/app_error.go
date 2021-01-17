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

type HttpStatus int

type AppError struct {
	err      error
	status   HttpStatus
	fileName string
	line     int
}

/**************************************
	ErrorImpl
**************************************/
func (e *AppError) Error() string {
	wrappedError := e.err

	// 内側がappErrorである限り、ループする
	var appErr *AppError
	if errors.As(wrappedError, &appErr) {
		fmt.Printf("Error: wrapping..., Location: %v:%v\n", e.fileName, e.line)
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

/**************************************
	Constructor
**************************************/
func NewAppError(err error) *AppError {
	fileName, line := getCallerData(2)
	return &AppError{
		err:      err,
		status:   _notEvaluated,
		fileName: fileName,
		line:     line,
	}
}

func NewAppErrorWithStatus(err error, status HttpStatus) *AppError {
	fileName, line := getCallerData(2)
	return &AppError{
		err:      err,
		status:   status,
		fileName: fileName,
		line:     line,
	}
}

/**************************************
	Getter & Setter
**************************************/
func (e AppError) GetHttpStatus() HttpStatus {
	// httpStatusがない && 内側がappErrorである限り, ループする
	var appErr *AppError
	if !e.status.isEvaluated() && errors.As(e.err, &appErr) {
		return appErr.GetHttpStatus()
	}

	return e.status
}

/**************************************
	private
**************************************/
func (s HttpStatus) isEvaluated() bool {
	return s != _notEvaluated
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
