package handler

import (
	"encoding/json"
	"net/http"

	"github.com/paypay3/kakeibo-rest-api-ddd/user-rest-service/apierrors"
	"github.com/paypay3/kakeibo-rest-api-ddd/user-rest-service/infrastructure/response"
	"github.com/paypay3/kakeibo-rest-api-ddd/user-rest-service/usecase"
	"github.com/paypay3/kakeibo-rest-api-ddd/user-rest-service/usecase/input"
)

type userHandler struct {
	userUsecase usecase.UserUsecase
}

func NewUserHandler(userUsecase usecase.UserUsecase) *userHandler {
	return &userHandler{
		userUsecase: userUsecase,
	}
}

func (h *userHandler) SignUp(w http.ResponseWriter, r *http.Request) {
	var in input.SignUpUser
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		response.ErrorJSON(w, apierrors.NewBadRequestError(apierrors.NewErrorString("正しいデータを入力してください")))
		return
	}

	out, err := h.userUsecase.SignUp(&in)
	if err != nil {
		response.ErrorJSON(w, err)
		return
	}

	response.JSON(w, http.StatusCreated, out)
}
