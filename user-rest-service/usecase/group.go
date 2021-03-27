package usecase

import (
	"golang.org/x/xerrors"

	"github.com/paypay3/kakeibo-rest-api-ddd/user-rest-service/apierrors"
	"github.com/paypay3/kakeibo-rest-api-ddd/user-rest-service/domain/groupdomain"
	"github.com/paypay3/kakeibo-rest-api-ddd/user-rest-service/domain/userdomain"
	"github.com/paypay3/kakeibo-rest-api-ddd/user-rest-service/interfaces/presenter"
	"github.com/paypay3/kakeibo-rest-api-ddd/user-rest-service/usecase/gateway"
	"github.com/paypay3/kakeibo-rest-api-ddd/user-rest-service/usecase/input"
	"github.com/paypay3/kakeibo-rest-api-ddd/user-rest-service/usecase/output"
	"github.com/paypay3/kakeibo-rest-api-ddd/user-rest-service/usecase/queryservice"
)

type GroupUsecase interface {
	FetchGroupList(in *input.AuthenticatedUser) (*output.GroupList, error)
	StoreGroup(authenticatedUser *input.AuthenticatedUser, group *input.Group) (*output.Group, error)
	UpdateGroupName(group *input.Group) (*output.Group, error)
}

type groupUsecase struct {
	groupRepository   groupdomain.Repository
	groupQueryService queryservice.GroupQueryService
	accountApi        gateway.AccountApi
}

func NewGroupUsecase(groupRepository groupdomain.Repository, groupQueryService queryservice.GroupQueryService, accountApi gateway.AccountApi) *groupUsecase {
	return &groupUsecase{
		groupRepository:   groupRepository,
		groupQueryService: groupQueryService,
		accountApi:        accountApi,
	}
}

func (u *groupUsecase) FetchGroupList(in *input.AuthenticatedUser) (*output.GroupList, error) {
	groupList, err := u.groupQueryService.FetchGroupList(in.UserID)
	if err != nil {
		return nil, err
	}

	return groupList, nil
}

func (u *groupUsecase) StoreGroup(authenticatedUser *input.AuthenticatedUser, groupInput *input.Group) (*output.Group, error) {
	userID, err := userdomain.NewUserID(authenticatedUser.UserID)
	if err != nil {
		return nil, apierrors.NewBadRequestError(&presenter.UserValidationError{UserID: "ユーザーIDを正しく入力してください"})
	}

	groupName, err := groupdomain.NewGroupName(groupInput.GroupName)
	if err != nil {
		if xerrors.Is(err, groupdomain.ErrCharacterCountGroupName) {
			return nil, apierrors.NewBadRequestError(apierrors.NewErrorString("グループ名は1文字以上、20文字以内で入力してください"))
		}

		if xerrors.Is(err, groupdomain.ErrPrefixSpaceGroupName) {
			return nil, apierrors.NewBadRequestError(apierrors.NewErrorString("文字列先頭に空白がないか確認してください"))
		}

		if xerrors.Is(err, groupdomain.ErrSuffixSpaceGroupName) {
			return nil, apierrors.NewBadRequestError(apierrors.NewErrorString("文字列末尾に空白がないか確認してください"))
		}

		return nil, apierrors.NewInternalServerError(apierrors.NewErrorString("Internal Server Error"))
	}

	group := groupdomain.NewGroupWithoutID(groupName)

	group, err = u.groupRepository.StoreGroupAndApprovedUser(group, userID)
	if err != nil {
		return nil, err
	}

	groupID, err := group.ID()
	if err != nil {
		return nil, apierrors.NewInternalServerError(apierrors.NewErrorString("Internal Server Error"))
	}

	if err := u.accountApi.PostInitGroupStandardBudgets(groupID); err != nil {
		if err := u.groupRepository.DeleteGroupAndApprovedUser(group); err != nil {
			return nil, err
		}

		return nil, err
	}

	return &output.Group{
		GroupID:   groupID.Value(),
		GroupName: group.GroupName().Value(),
	}, nil
}

func (u *groupUsecase) UpdateGroupName(groupInput *input.Group) (*output.Group, error) {
	groupID, err := groupdomain.NewGroupID(groupInput.GroupID)
	if err != nil {
		return nil, apierrors.NewBadRequestError(apierrors.NewErrorString("グループIDは1以上の整数で指定してください"))
	}

	group, err := u.groupRepository.FindGroupByID(&groupID)
	if err != nil {
		return nil, err
	}

	groupName, err := groupdomain.NewGroupName(groupInput.GroupName)
	if err != nil {
		if xerrors.Is(err, groupdomain.ErrCharacterCountGroupName) {
			return nil, apierrors.NewBadRequestError(apierrors.NewErrorString("グループ名は1文字以上、20文字以内で入力してください"))
		}

		if xerrors.Is(err, groupdomain.ErrPrefixSpaceGroupName) {
			return nil, apierrors.NewBadRequestError(apierrors.NewErrorString("文字列先頭に空白がないか確認してください"))
		}

		if xerrors.Is(err, groupdomain.ErrSuffixSpaceGroupName) {
			return nil, apierrors.NewBadRequestError(apierrors.NewErrorString("文字列末尾に空白がないか確認してください"))
		}

		return nil, apierrors.NewInternalServerError(apierrors.NewErrorString("Internal Server Error"))
	}

	group.UpdateGroupName(groupName)

	if err := u.groupRepository.UpdateGroupName(group); err != nil {
		return nil, err
	}

	return &output.Group{
		GroupID:   groupID.Value(),
		GroupName: groupName.Value(),
	}, nil
}
