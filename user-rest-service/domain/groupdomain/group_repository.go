package groupdomain

import "github.com/paypay3/kakeibo-rest-api-ddd/user-rest-service/domain/userdomain"

type Repository interface {
	StoreGroupAndApprovedUser(group *Group, userID userdomain.UserID) (*Group, error)
	DeleteGroupAndApprovedUser(group *Group) error
}
