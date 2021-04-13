package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/context"
	"github.com/gorilla/mux"

	"github.com/paypay3/kakeibo-rest-api-ddd/user-rest-service/apierrors"
	"github.com/paypay3/kakeibo-rest-api-ddd/user-rest-service/apperrors"
	"github.com/paypay3/kakeibo-rest-api-ddd/user-rest-service/config"
	"github.com/paypay3/kakeibo-rest-api-ddd/user-rest-service/interfaces/presenter"
	"github.com/paypay3/kakeibo-rest-api-ddd/user-rest-service/usecase"
	"github.com/paypay3/kakeibo-rest-api-ddd/user-rest-service/usecase/input"
)

type groupHandler struct {
	groupUsecase usecase.GroupUsecase
}

func NewGroupHandler(groupUsecase usecase.GroupUsecase) *groupHandler {
	return &groupHandler{
		groupUsecase: groupUsecase,
	}
}

func (h *groupHandler) FetchGroupList(w http.ResponseWriter, r *http.Request) {
	in, err := getUserIDOfContext(r)
	if err != nil {
		presenter.ErrorJSON(w, err)
		return
	}

	out, err := h.groupUsecase.FetchGroupList(in)
	if err != nil {
		presenter.ErrorJSON(w, err)
		return
	}

	presenter.JSON(w, http.StatusOK, out)
}

func (h *groupHandler) StoreGroup(w http.ResponseWriter, r *http.Request) {
	authenticatedUser, err := getUserIDOfContext(r)
	if err != nil {
		presenter.ErrorJSON(w, err)
		return
	}

	var group input.Group
	if err := json.NewDecoder(r.Body).Decode(&group); err != nil {
		presenter.ErrorJSON(w, apierrors.NewBadRequestError(apierrors.NewErrorString("正しいデータを入力してください")))
		return
	}

	out, err := h.groupUsecase.StoreGroup(authenticatedUser, &group)
	if err != nil {
		presenter.ErrorJSON(w, err)
		return
	}

	presenter.JSON(w, http.StatusCreated, out)
}

func (h *groupHandler) UpdateGroupName(w http.ResponseWriter, r *http.Request) {
	var group input.Group
	if err := json.NewDecoder(r.Body).Decode(&group); err != nil {
		presenter.ErrorJSON(w, apierrors.NewBadRequestError(apierrors.NewErrorString("正しいデータを入力してください")))
		return
	}

	groupID, err := strconv.Atoi(mux.Vars(r)["group_id"])
	if err != nil {
		presenter.ErrorJSON(w, apierrors.NewBadRequestError(apierrors.NewErrorString("グループIDを正しく指定してください")))
		return
	}

	group.GroupID = groupID

	out, err := h.groupUsecase.UpdateGroupName(&group)
	if err != nil {
		presenter.ErrorJSON(w, err)
		return
	}

	presenter.JSON(w, http.StatusOK, out)
}

func (h *groupHandler) StoreGroupUnapprovedUser(w http.ResponseWriter, r *http.Request) {
	var unapprovedUser input.UnapprovedUser
	if err := json.NewDecoder(r.Body).Decode(&unapprovedUser); err != nil {
		presenter.ErrorJSON(w, apierrors.NewBadRequestError(apierrors.NewErrorString("正しいデータを入力してください")))
		return
	}

	var group input.Group
	groupID, err := strconv.Atoi(mux.Vars(r)["group_id"])
	if err != nil {
		presenter.ErrorJSON(w, apierrors.NewBadRequestError(apierrors.NewErrorString("グループIDを正しく指定してください")))
		return
	}

	group.GroupID = groupID

	out, err := h.groupUsecase.StoreGroupUnapprovedUser(&unapprovedUser, &group)
	if err != nil {
		presenter.ErrorJSON(w, err)
		return
	}

	presenter.JSON(w, http.StatusCreated, out)
}

func (h *groupHandler) DeleteGroupApprovedUser(w http.ResponseWriter, r *http.Request) {
	authenticatedUser, err := getUserIDOfContext(r)
	if err != nil {
		presenter.ErrorJSON(w, err)
		return
	}

	groupID, err := strconv.Atoi(mux.Vars(r)["group_id"])
	if err != nil {
		presenter.ErrorJSON(w, apierrors.NewBadRequestError(apierrors.NewErrorString("グループIDを正しく指定してください")))
		return
	}

	group := input.Group{GroupID: groupID}

	if err := h.groupUsecase.DeleteGroupApprovedUser(authenticatedUser, &group); err != nil {
		presenter.ErrorJSON(w, err)
		return
	}

	presenter.JSON(w, http.StatusOK, presenter.NewSuccessString("グループを退会しました"))
}

func (h *groupHandler) StoreGroupApprovedUser(w http.ResponseWriter, r *http.Request) {
	authenticatedUser, err := getUserIDOfContext(r)
	if err != nil {
		err = apperrors.Wrap(err)
		context.Set(r, config.Env.ContextKey.AppError, err)
		presenter.ErrorJSONV2(w, err)
		return
	}

	groupID, err := strconv.Atoi(mux.Vars(r)["group_id"])
	if err != nil {
		err = apperrors.InvalidParameter.SetInfoMessage(apperrors.NewErrorString("グループIDを正しく指定してください")).Wrap(err)
		context.Set(r, config.Env.ContextKey.AppError, err)
		presenter.ErrorJSONV2(w, err)
		return
	}

	group := input.Group{GroupID: groupID}

	out, err := h.groupUsecase.StoreGroupApprovedUser(authenticatedUser, &group)
	if err != nil {
		err = apperrors.Wrap(err, "StoreGroupApprovedUser handler failed ")
		context.Set(r, config.Env.ContextKey.AppError, err)
		presenter.ErrorJSONV2(w, err)
		return
	}

	presenter.JSON(w, http.StatusCreated, out)
}
