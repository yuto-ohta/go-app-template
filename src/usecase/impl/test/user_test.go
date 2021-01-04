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
	"gorm.io/gorm"
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

func TestUserUseCaseImpl_FindById_userIdでユーザーが返ること(t *testing.T) {
	// setup
	userId := valueobject.NewUserIdWithId(1)
	target := impl.NewUserUseCaseImpl(infrastructure.NewUserRepositoryImpl())

	// actual
	actual, err := target.FindById(*userId)
	if err != nil {
		t.Errorf(fmt.Sprintf("ユーザー取得にエラーが発生しています, エラー: %v", err))
	}

	// expected
	// TODO: DB周りのテスト環境整備
	expected := *domain.NewUserWithUserId(*userId, "まるお")

	// check
	assert.Equal(t, expected, actual)
}

func TestUserUseCaseImpl_FindById_存在しないuserIdでRecordNotFoundが返ること(t *testing.T) {
	// setup
	userId := valueobject.NewUserIdWithId(9999)
	target := impl.NewUserUseCaseImpl(infrastructure.NewUserRepositoryImpl())

	// actual
	var (
		actualAppErr    apperror.AppError
		actualErrStatus int
	)

	_, actualErr := target.FindById(*userId)
	if actualErr == nil {
		t.Error("エラーが発生していません。RecordNotFoundが返るはず")
	}
	var appErr *apperror.AppError
	if errors.As(actualErr, &appErr) {
		actualErrStatus = appErr.GetHttpStatus()
		actualAppErr = *appErr
	} else {
		t.Error("エラーがAppErrorになっていません")
	}

	// expected
	expectedAppErr := apperror.NewAppErrorWithStatus(gorm.ErrRecordNotFound, http.StatusNotFound)
	expectedErrStatus := expectedAppErr.GetHttpStatus()

	// check
	assert.Equal(t, expectedErrStatus, actualErrStatus)
	assert.Equal(t, expectedAppErr.ErrorWithoutLocation(), actualAppErr.ErrorWithoutLocation())
}

func TestUserUseCaseImpl_CreateUser_正常にユーザーが登録できること(t *testing.T) {
	// setup
	target := impl.NewUserUseCaseImpl(infrastructure.NewUserRepositoryImpl())
	userName := "新規ユーザー太郎"
	userDomain, _ := domain.NewUser(userName)

	// actual
	var actualCreatedUser domain.User
	var err error
	actualCreatedUser, err = target.CreateUser(*userDomain)
	if err != nil {
		t.Errorf("ユーザー登録に失敗しています, Error: %v", err.Error())
	}

	// expected
	expectedCreatedUser, _ := target.FindById(actualCreatedUser.GetId())

	// check
	assert.Equal(t, expectedCreatedUser, actualCreatedUser)
}

func TestUserUseCaseImpl_CreateUser_すでにuserIdがある場合_登録できないこと(t *testing.T) {
	// setup
	target := impl.NewUserUseCaseImpl(infrastructure.NewUserRepositoryImpl())
	userId := valueobject.NewUserIdWithId(9999)
	userName := "新規ユーザー太郎"
	userDomain := domain.NewUserWithUserId(*userId, userName)

	// actual
	var actualAppErr apperror.AppError
	var actualErrStatus int

	_, err := target.CreateUser(*userDomain)
	if err == nil {
		t.Error("エラーが発生していません")
	}
	var appErr *apperror.AppError
	if errors.As(err, &appErr) {
		actualErrStatus = appErr.GetHttpStatus()
		actualAppErr = *appErr
	} else {
		t.Error("エラーがAppErrorになっていません")
	}

	// expected
	expectedAppErr := apperror.NewAppErrorWithStatus(fmt.Errorf("未登録のユーザーにuserIdが割り当てられています, user: %v", *userDomain), http.StatusInternalServerError)
	expectedErrStatus := expectedAppErr.GetHttpStatus()

	// check
	assert.Equal(t, expectedErrStatus, actualErrStatus)
	assert.Equal(t, expectedAppErr.ErrorWithoutLocation(), actualAppErr.ErrorWithoutLocation())
}
