package groupdomain

import "github.com/paypay3/kakeibo-rest-api-ddd/user-rest-service/domain/userdomain"

type Repository interface {
	StoreGroupAndApprovedUser(group *Group, userID userdomain.UserID) (*Group, error)
	DeleteGroupAndApprovedUser(group *Group) error
	UpdateGroupName(group *Group) error
	StoreUnapprovedUser(unapprovedUser *UnapprovedUser) error
	DeleteApprovedUser(approvedUser *ApprovedUser) error
	StoreApprovedUser(approvedUser *ApprovedUser) error
	FindGroupByID(groupID *GroupID) (*Group, error)
	FindApprovedUser(groupID GroupID, userID userdomain.UserID) (*ApprovedUser, error)
	FindUnapprovedUser(groupID GroupID, userID userdomain.UserID) (*UnapprovedUser, error)
	FetchApprovedUserIDList(groupID GroupID) ([]userdomain.UserID, error)
}
