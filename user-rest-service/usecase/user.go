package usecase

import (
	"golang.org/x/xerrors"

	"github.com/paypay3/kakeibo-rest-api-ddd/user-rest-service/apierrors"
	"github.com/paypay3/kakeibo-rest-api-ddd/user-rest-service/domain/userdomain"
	"github.com/paypay3/kakeibo-rest-api-ddd/user-rest-service/domain/vo"
	"github.com/paypay3/kakeibo-rest-api-ddd/user-rest-service/usecase/input"
	"github.com/paypay3/kakeibo-rest-api-ddd/user-rest-service/usecase/interfaces"
	"github.com/paypay3/kakeibo-rest-api-ddd/user-rest-service/usecase/output"
)

type UserUsecase interface {
	SignUp(in *input.SignUpUser) (*output.SignUpUser, error)
}

type userUsecase struct {
	userRepository userdomain.Repository
	accountApi     interfaces.AccountApi
}

func NewUserUsecase(userRepository userdomain.Repository, accountApi interfaces.AccountApi) *userUsecase {
	return &userUsecase{
		userRepository: userRepository,
		accountApi:     accountApi,
	}
}

func (u *userUsecase) SignUp(in *input.SignUpUser) (*output.SignUpUser, error) {
	var userValidationError apierrors.UserValidationError

	userID, err := vo.NewUserID(in.UserID)
	if err != nil {
		userValidationError.UserID = "ユーザーIDを正しく入力してください"
	}

	email, err := vo.NewEmail(in.Email)
	if err != nil {
		userValidationError.Email = "メールアドレスを正しく入力してください"
	}

	password, err := vo.NewPassword(in.Password)
	if err != nil {
		if xerrors.Is(err, apierrors.ErrInvalidPassword) {
			userValidationError.Password = "パスワードを正しく入力してください"
		} else {
			return nil, apierrors.NewInternalServerError(apierrors.NewErrorString("Internal Server Error"))
		}
	}

	signUpUser, err := userdomain.NewSignUpUser(userID, in.Name, email, password)
	if err != nil {
		userValidationError.Name = "名前を正しく入力してください"
	}

	if userValidationError.UserID != "" ||
		userValidationError.Name != "" ||
		userValidationError.Email != "" ||
		userValidationError.Password != "" {
		return nil, apierrors.NewBadRequestError(&userValidationError)
	}

	if err := checkForUniqueUser(u, signUpUser); err != nil {
		var userConflictError *apierrors.UserConflictError
		if xerrors.As(err, &userConflictError) {
			return nil, apierrors.NewConflictError(userConflictError)
		}

		return nil, apierrors.NewInternalServerError(apierrors.NewErrorString("Internal Server Error"))
	}

	if err := u.userRepository.CreateSignUpUser(signUpUser); err != nil {
		return nil, apierrors.NewInternalServerError(apierrors.NewErrorString("Internal Server Error"))
	}

	if err := u.accountApi.PostInitStandardBudgets(signUpUser.UserID().Value()); err != nil {
		if err := u.userRepository.DeleteSignUpUser(signUpUser); err != nil {
			return nil, apierrors.NewInternalServerError(apierrors.NewErrorString("Internal Server Error"))
		}

		return nil, apierrors.NewInternalServerError(apierrors.NewErrorString("Internal Server Error"))
	}

	return &output.SignUpUser{
		UserID: signUpUser.UserID().Value(),
		Name:   signUpUser.Name(),
		Email:  signUpUser.Email().Value(),
	}, nil
}

func checkForUniqueUser(u *userUsecase, signUpUser *userdomain.SignUpUser) error {
	_, errUserID := u.userRepository.FindSignUpUserByUserID(signUpUser.UserID().Value())
	if errUserID != nil && !xerrors.Is(errUserID, apierrors.ErrUserNotFound) {
		return errUserID
	}

	_, errEmail := u.userRepository.FindSignUpUserByEmail(signUpUser.Email().Value())
	if errEmail != nil && !xerrors.Is(errEmail, apierrors.ErrUserNotFound) {
		return errEmail
	}

	existsUserByUserID := !xerrors.Is(errUserID, apierrors.ErrUserNotFound)
	existsUserByEmail := !xerrors.Is(errEmail, apierrors.ErrUserNotFound)

	if !existsUserByUserID && !existsUserByEmail {
		return nil
	}

	var userConflictError apierrors.UserConflictError

	if existsUserByUserID {
		userConflictError.UserID = "このユーザーIDは既に利用されています"
	}

	if existsUserByEmail {
		userConflictError.Email = "このメールアドレスは既に利用されています"
	}

	return &userConflictError
}