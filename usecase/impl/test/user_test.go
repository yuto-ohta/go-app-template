package test

import (
	"github.com/stretchr/testify/assert"
	"go-app-template/domain"
	"go-app-template/infrastructure"
	"go-app-template/usecase/impl"
	"gorm.io/gorm"
	"testing"
)

func TestUserUseCaseImpl_FindById_userIdでユーザーが返ること(t *testing.T) {
	userId := domain.NewUserId(1)
	target := impl.NewUserUseCaseImpl(infrastructure.NewUserRepositoryImpl())
	actual, err := target.FindById(*userId)

	// 予想外のエラーが発生した場合はFail
	if err != nil {
		t.Errorf("Unknown Error: %v", err.Error())
	}

	// TODO: Mockに差し替えるかテスト用DBを用意する
	expected := *domain.NewUser(*userId, "taro")
	assert.Equal(t, expected, actual)
}

func TestUserUseCaseImpl_FindById_存在しないuserIdでRecordNotFoundが返ること(t *testing.T) {
	userId := domain.NewUserId(9999)
	target := impl.NewUserUseCaseImpl(infrastructure.NewUserRepositoryImpl())
	user, actualErr := target.FindById(*userId)

	// RecordNotFoundが返るはずなので、エラーがnilの場合はFail
	if actualErr == nil {
		t.Errorf("Unknown Error: %v, user: %v", "actualErrがなぜかnilになっています", user)
	}

	expectedErr := gorm.ErrRecordNotFound
	assert.Equal(t, expectedErr, actualErr)
}
