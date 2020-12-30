package errors

import (
	"fmt"
	"go-app-template/src/config"
	"runtime"
	"strings"
)

type AppError struct {
	err      error
	status   int
	fileName string
	line     int
}

func (e *AppError) Error() string {
	return fmt.Sprintf("Error: %v,\nLocation: %v:%v", e.err.Error(), e.fileName, e.line)
}

func (e *AppError) GetHttpStatus() int {
	return e.status
}

func NewAppError(err error, status int) *AppError {
	fileName, line := getCallerData(2)
	return &AppError{
		err:      err,
		status:   status,
		fileName: fileName,
		line:     line,
	}
}

func getCallerData(skip int) (string, int) {
	_, fileAbsPath, line, ok := runtime.Caller(skip)
	if !ok {
		return "???", 0
	}

	projectName := config.Properties.ProjectName
	fileProjectPath := fileAbsPath[strings.Index(fileAbsPath, projectName):]

	return fileProjectPath, line
}
