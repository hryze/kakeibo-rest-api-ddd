package usecase

import (
	"github.com/paypay3/kakeibo-rest-api-ddd/user-rest-service/apierrors"
	"github.com/paypay3/kakeibo-rest-api-ddd/user-rest-service/domain/userdomain"
	"github.com/paypay3/kakeibo-rest-api-ddd/user-rest-service/interfaces/presenter"
	"github.com/paypay3/kakeibo-rest-api-ddd/user-rest-service/usecase/input"
	"github.com/paypay3/kakeibo-rest-api-ddd/user-rest-service/usecase/output"
	"github.com/paypay3/kakeibo-rest-api-ddd/user-rest-service/usecase/queryservice"
)

type GroupUsecase interface {
	FetchGroupList(in *input.AuthenticatedUser) (*output.GroupList, error)
}

type groupUsecase struct {
	groupQueryService queryservice.GroupQueryService
}

func NewGroupUsecase(groupQueryService queryservice.GroupQueryService) *groupUsecase {
	return &groupUsecase{
		groupQueryService: groupQueryService,
	}
}

func (u *groupUsecase) FetchGroupList(in *input.AuthenticatedUser) (*output.GroupList, error) {
	userID, err := userdomain.NewUserID(in.UserID)
	if err != nil {
		return nil, apierrors.NewBadRequestError(&presenter.UserValidationError{UserID: "ユーザーIDが正しくありません"})
	}

	groupList, err := u.groupQueryService.FetchGroupList(userID)
	if err != nil {
		return nil, err
	}

	return groupList, nil
}
