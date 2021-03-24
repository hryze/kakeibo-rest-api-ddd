package groupdomain

import "github.com/paypay3/kakeibo-rest-api-ddd/user-rest-service/domain/userdomain"

type Repository interface {
	StoreGroupAndApprovedUser(groupName GroupName, userID userdomain.UserID) (*Group, error)
	DeleteGroupAndApprovedUser(groupID GroupID) error
}
