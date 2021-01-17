package impl

import (
	"fmt"
	"go-app-template/src/api/controller/dto"
	"go-app-template/src/apperror"
	"go-app-template/src/domain"
	"go-app-template/src/domain/repository"
	"go-app-template/src/domain/valueobject"
	"go-app-template/src/usecase/appmodel"
	"net/http"
)

type UserUseCaseImpl struct {
	userRepository repository.UserRepository
}

func NewUserUseCaseImpl(userRepository repository.UserRepository) *UserUseCaseImpl {
	return &UserUseCaseImpl{userRepository: userRepository}
}

/**************************************
	ユーザー取得
**************************************/
func (u UserUseCaseImpl) GetUser(id int) (dto.UserResDto, error) {
	var err error

	//get userId
	var userId *valueobject.UserId
	if userId, err = valueobject.NewUserIdWithId(id); err != nil {
		return dto.UserResDto{}, apperror.NewAppErrorWithStatus(err, http.StatusInternalServerError)
	}

	// find user
	var found domain.User
	if found, err = u.userRepository.FindById(*userId); err != nil {
		return dto.UserResDto{}, apperror.NewAppError(err)
	}

	return *found.ToDto(), nil
}

/**************************************
	ユーザー全件取得
**************************************/
func (u UserUseCaseImpl) GetAllUser(condition appmodel.SearchCondition) (dto.UserPage, error) {
	var err error

	// find all user
	var allUser []domain.User
	if allUser, err = u.userRepository.FindAll(); err != nil {
		return dto.UserPage{}, apperror.NewAppError(err)
	}

	// sort
	if err = domain.Sort(condition.GetOrderBy(), condition.GetOrder(), allUser); err != nil {
		return dto.UserPage{}, apperror.NewAppErrorWithStatus(err, http.StatusInternalServerError)
	}

	// make userPage
	var userPage appmodel.Page
	if userPage, err = makeUserPage(condition.GetPage(), condition.GetLimit(), allUser); err != nil {
		return dto.UserPage{}, apperror.NewAppError(err)
	}

	// make userDtoList
	var dtoList []dto.UserResDto
	if dtoList, err = makeUserDtoListFromPage(userPage); err != nil {
		return dto.UserPage{}, apperror.NewAppError(err)
	}

	// make userPageDto
	pageInfoDto := dto.PageInfo{
		PageNum:     userPage.GetInfo().GetPageNum(),
		LastPageNum: userPage.GetInfo().GetLastPageNum(),
		Limit:       userPage.GetInfo().GetLimit(),
		Offset:      userPage.GetInfo().GetOffset(),
	}
	userPageDto := dto.UserPage{
		Users:    dtoList,
		PageInfo: pageInfoDto,
	}

	return userPageDto, nil
}

/**************************************
	ユーザー登録
**************************************/
func (u UserUseCaseImpl) CreateUser(userDto dto.UserReceiveDto) (dto.UserResDto, error) {
	var err error

	// get user domain
	var newUser *domain.User
	if newUser, err = domain.NewUserBuilder().Name(userDto.Name).Password(userDto.Password).Build(); err != nil {
		return dto.UserResDto{}, apperror.NewAppErrorWithStatus(err, http.StatusInternalServerError)
	}

	// create user
	var created domain.User
	if created, err = u.userRepository.CreateUser(*newUser); err != nil {
		return dto.UserResDto{}, apperror.NewAppErrorWithStatus(err, http.StatusInternalServerError)
	}

	return *created.ToDto(), nil
}

/**************************************
	ユーザー削除
**************************************/
func (u UserUseCaseImpl) DeleteUser(id int) (dto.UserResDto, error) {
	var err error

	// get userId
	var userId *valueobject.UserId
	if userId, err = valueobject.NewUserIdWithId(id); err != nil {
		return dto.UserResDto{}, apperror.NewAppErrorWithStatus(err, http.StatusBadRequest)
	}

	// delete user
	var deleted domain.User
	if deleted, err = u.userRepository.DeleteUser(*userId); err != nil {
		return dto.UserResDto{}, apperror.NewAppError(err)
	}

	return *deleted.ToDto(), nil
}

/**************************************
	ユーザー更新
**************************************/
func (u UserUseCaseImpl) UpdateUser(id int, user dto.UserReceiveDto) (dto.UserResDto, error) {
	var err error

	// get userId
	var userId *valueobject.UserId
	if userId, err = valueobject.NewUserIdWithId(id); err != nil {
		return dto.UserResDto{}, apperror.NewAppErrorWithStatus(err, http.StatusBadRequest)
	}

	// get oldParams
	var oldParams domain.User
	if oldParams, err = u.userRepository.FindById(*userId); err != nil {
		return dto.UserResDto{}, apperror.NewAppError(err)
	}

	// get new userDomain
	var newUser *domain.User
	// ユーザー名 && パスワード 両方変更する場合
	if len(user.Name) != 0 && len(user.Password) != 0 {
		if newUser, err = domain.NewUserBuilder().Name(user.Name).Password(user.Password).Build(); err != nil {
			return dto.UserResDto{}, apperror.NewAppErrorWithStatus(err, http.StatusInternalServerError)
		}
		// ユーザー名のみ変更する場合
	} else if len(user.Name) != 0 {
		if err = oldParams.SetName(user.Name); err != nil {
			return dto.UserResDto{}, apperror.NewAppErrorWithStatus(err, http.StatusInternalServerError)
		}
		newUser = &oldParams
		// パスワードのみ変更する場合
	} else {
		if newUser, err = domain.NewUserBuilder().Name(oldParams.GetName()).Password(user.Password).Build(); err != nil {
			return dto.UserResDto{}, apperror.NewAppErrorWithStatus(err, http.StatusInternalServerError)
		}
	}

	// update user
	var updated domain.User
	if updated, err = u.userRepository.UpdateUser(*userId, *newUser); err != nil {
		return dto.UserResDto{}, apperror.NewAppError(err)
	}

	return *updated.ToDto(), nil
}

/**************************************
	private
**************************************/
func makeUserPage(page int, limit int, target []domain.User) (appmodel.Page, error) {
	var err error

	// convert to Pageable
	pageable := make(appmodel.Pageable, len(target))
	for i, u := range target {
		pageable[i] = u
	}

	// paging
	var userPage *appmodel.Page
	if userPage, err = pageable.GetPage(page, limit); err != nil {
		return appmodel.Page{}, apperror.NewAppError(err)
	}

	return *userPage, nil
}

func makeUserDtoListFromPage(userPage appmodel.Page) ([]dto.UserResDto, error) {
	dtoList := make([]dto.UserResDto, len(userPage.GetList()))
	for i, e := range userPage.GetList() {
		userDomain, ok := e.(domain.User)
		if !ok {
			return []dto.UserResDto{}, apperror.NewAppErrorWithStatus(fmt.Errorf("型アサーションエラー\nfrom: %v\nto: domain.User", e), http.StatusInternalServerError)
		}
		userDto := userDomain.ToDto()
		dtoList[i] = *userDto
	}

	return dtoList, nil
}
