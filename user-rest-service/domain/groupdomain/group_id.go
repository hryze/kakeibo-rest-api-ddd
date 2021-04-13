package groupdomain

import "github.com/paypay3/kakeibo-rest-api-ddd/user-rest-service/apperrors"

type GroupID int

func NewGroupID(id int) (GroupID, error) {
	if id < 1 {
		return 0, apperrors.Errorf("group id must be an integer greater than or equal to 1: %d", id)
	}

	return GroupID(id), nil
}

func (i GroupID) Value() int {
	return int(i)
}
