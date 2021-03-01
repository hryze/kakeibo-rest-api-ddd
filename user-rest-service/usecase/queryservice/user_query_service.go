package queryservice

import "github.com/paypay3/kakeibo-rest-api-ddd/user-rest-service/domain/userdomain"

type UserQueryService interface {
	FindLoginUserByUserID(userID userdomain.UserID) (*userdomain.LoginUserWithoutPassword, error)
}
