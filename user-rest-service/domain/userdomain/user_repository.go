package userdomain

import "github.com/paypay3/kakeibo-rest-api-ddd/user-rest-service/domain/vo"

type Repository interface {
	FindSignUpUserByUserID(userID UserID) (*SignUpUser, error)
	FindSignUpUserByEmail(email vo.Email) (*SignUpUser, error)
	CreateSignUpUser(user *SignUpUser) error
	DeleteSignUpUser(signUpUser *SignUpUser) error
	FindLoginUserByEmail(email vo.Email) (*LoginUser, error)
	AddSessionID(sessionID string, userID UserID) error
}
