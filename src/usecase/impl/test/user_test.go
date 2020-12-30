package test

import (
	"errors"
	"fmt"
	"github.com/stretchr/testify/assert"
	"go-app-template/src/config"
	"go-app-template/src/domain"
	appErrors "go-app-template/src/errors"
	"go-app-template/src/errors/test/mock"
	"go-app-template/src/infrastructure"
	"go-app-template/src/usecase/impl"
	"gorm.io/gorm"
	"net/http"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	// before all
	config.LoadConfig()

	// run each test
	code := m.Run()

	// after all

	// finish test
	os.Exit(code)
}

func TestUserUseCaseImpl_FindById_userIdでユーザーが返ること(t *testing.T) {
	// setup
	userId := domain.NewUserId(1)
	target := impl.NewUserUseCaseImpl(infrastructure.NewUserRepositoryImpl())

	// actual
	actual, err := target.FindById(*userId)
	if err != nil {
		t.Errorf(fmt.Sprintf("ユーザー取得にエラーが発生しています, エラー: %v", err))
	}

	// expected
	// TODO: DB周りのテスト環境整備
	expected := *domain.NewUser(*userId, "taro")

	// check
	assert.Equal(t, expected, actual)
}

func TestUserUseCaseImpl_FindById_存在しないuserIdでRecordNotFoundが返ること(t *testing.T) {
	// setup
	userId := domain.NewUserId(9999)
	target := impl.NewUserUseCaseImpl(infrastructure.NewUserRepositoryImpl())

	// actual
	_, actualErr := target.FindById(*userId)
	var actualErrStatus int

	var appErr *appErrors.AppError
	if actualErr == nil {
		t.Errorf("エラーが発生していません。RecordNotFoundが返るはず")
	}
	if errors.As(actualErr, &appErr) {
		actualErrStatus = appErr.GetHttpStatus()
	} else {
		t.Errorf("エラーがAppErrorになっていません")
	}

	// expected
	const FilePath = "go-app-template/infrastructure/user.go"
	const Line = 31
	expectedErr := mock.NewAppErrorMock(gorm.ErrRecordNotFound, http.StatusNotFound, FilePath, Line)
	expectedErrStatus := expectedErr.GetHttpStatus()

	// check
	assert.Equal(t, expectedErrStatus, actualErrStatus)
	assert.Equal(t, expectedErr.Error(), actualErr.Error())
}
