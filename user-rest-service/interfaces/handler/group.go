package handler

import (
	"net/http"

	"github.com/gorilla/context"

	"github.com/paypay3/kakeibo-rest-api-ddd/user-rest-service/apierrors"
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
	ctx, ok := context.GetOk(r, config.Env.RequestCtx.UserID)
	if !ok {
		presenter.ErrorJSON(w, apierrors.NewInternalServerError(apierrors.NewErrorString("Internal Server Error")))
		return
	}

	ctxUserID, ok := ctx.(string)
	if !ok {
		presenter.ErrorJSON(w, apierrors.NewInternalServerError(apierrors.NewErrorString("Internal Server Error")))
		return
	}

	in := input.AuthenticatedUser{UserID: ctxUserID}

	out, err := h.groupUsecase.FetchGroupList(&in)
	if err != nil {
		presenter.ErrorJSON(w, err)
		return
	}

	presenter.JSON(w, http.StatusOK, out)
}
