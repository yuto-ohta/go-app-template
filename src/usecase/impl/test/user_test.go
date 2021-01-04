package test

import (
	"errors"
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
	target := impl.NewUserUseCaseImpl(infrastructure.NewUserRepositoryImpl())
	userIdInt := 1
	userName := "まるお"

	// actual
	actual, _ := target.FindById(userIdInt)

	// expected
	userId := valueobject.NewUserIdWithId(userIdInt)
	expected := *domain.NewUserWithUserId(*userId, userName)

	// check
	assert.Equal(t, expected, actual)
}

func TestUserUseCaseImpl_FindById_存在しないuserIdでRecordNotFoundが返ること(t *testing.T) {
	// setup
	target := impl.NewUserUseCaseImpl(infrastructure.NewUserRepositoryImpl())
	userIdInt := 9999

	// actual
	var actualAppErr apperror.AppError
	var actualErrStatus int

	_, actualErr := target.FindById(userIdInt)
	var appErr *apperror.AppError
	if errors.As(actualErr, &appErr) {
		actualErrStatus = appErr.GetHttpStatus()
		actualAppErr = *appErr
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

	// actual
	actualCreatedUser, _ := target.CreateUser(userName)

	// expected
	expectedCreatedUser, _ := target.FindById(actualCreatedUser.GetId().GetValue())

	// check
	assert.Equal(t, expectedCreatedUser, actualCreatedUser)
}
