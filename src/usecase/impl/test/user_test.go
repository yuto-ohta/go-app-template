package test

import (
	"errors"
	"fmt"
	"go-app-template/src/apperror"
	"go-app-template/src/config/db/localdata"
	"go-app-template/src/domain"
	"go-app-template/src/domain/valueobject"
	"go-app-template/src/infrastructure"
	"go-app-template/src/usecase/impl"
	"net/http"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	_target = impl.NewUserUseCaseImpl(infrastructure.NewUserRepositoryImpl())
)

func TestMain(m *testing.M) {
	// before all
	localdata.InitializeLocalData()

	// run each test
	code := m.Run()

	// after all

	// finish test
	localdata.InitializeLocalData()
	os.Exit(code)
}

func TestUserUseCaseImpl_TestFindById_正常系(t *testing.T) {
	// setup
	var params = []struct {
		inputUserIdInt   int
		expectedUserName string
	}{
		{1, "まるお"},
	}

	for _, p := range params {
		// actual
		actual, _ := _target.FindById(p.inputUserIdInt)

		// expected
		userId, _ := valueobject.NewUserIdWithId(p.inputUserIdInt)
		expected, _ := domain.NewUserWithUserId(*userId, p.expectedUserName)

		// check
		assert.Equal(t, *expected, actual)
	}
}

func TestUserUseCaseImpl_FindById_異常系(t *testing.T) {
	// setup
	var params = []struct {
		title              string
		inputUserIdInt     int
		expectedErrStr     string
		expectedHttpStatus int
	}{
		{"存在しないuserIdはNG", 9999, "Error: record not found", http.StatusNotFound},
	}

	for _, p := range params {
		// actual
		var (
			actualAppErr    *apperror.AppError
			actualErrStr    string
			actualErrStatus int
		)
		_, actualErr := _target.FindById(p.inputUserIdInt)
		if errors.As(actualErr, &actualAppErr) {
			actualErrStr = actualAppErr.ErrorWithoutLocation()
			actualErrStatus = actualAppErr.GetHttpStatus()
		}

		// expected
		expectedErrStr := p.expectedErrStr
		expectedErrStatus := p.expectedHttpStatus

		// check
		fmt.Println(p.title)
		assert.Equal(t, expectedErrStatus, actualErrStatus)
		assert.Equal(t, expectedErrStr, actualErrStr)
	}
}

func TestUserUseCaseImpl_CreateUser_正常系(t *testing.T) {
	// setup
	var params = []struct {
		inputUserName string
	}{
		{"新規ユーザー太郎"},
	}

	for _, p := range params {
		// actual
		actualCreatedUser, _ := _target.CreateUser(p.inputUserName)

		// expected
		expectedCreatedUser, _ := _target.FindById(actualCreatedUser.GetId().GetValue())

		// check
		assert.Equal(t, expectedCreatedUser, actualCreatedUser)
	}
}

func TestUserUseCaseImpl_CreateUser_異常系(t *testing.T) {
	// setup
	var params = []struct {
		title              string
		inputUserName      string
		expectedErrStr     string
		expectedHttpStatus int
	}{
		{"９文字以上のユーザー名はNG", "123456789", "Error: Key: '' Error:Field validation for '' failed on the 'max' tag", http.StatusBadRequest},
	}

	for _, p := range params {
		// actual
		var (
			actualAppErr    *apperror.AppError
			actualErrStr    string
			actualErrStatus int
		)
		_, actualErr := _target.CreateUser(p.inputUserName)
		if errors.As(actualErr, &actualAppErr) {
			actualErrStr = actualAppErr.ErrorWithoutLocation()
			actualErrStatus = actualAppErr.GetHttpStatus()
		}

		// expected
		expectedErrStr := p.expectedErrStr
		expectedErrStatus := p.expectedHttpStatus

		// check
		fmt.Println(p.title)
		assert.Equal(t, expectedErrStatus, actualErrStatus)
		assert.Equal(t, expectedErrStr, actualErrStr)
	}
}
