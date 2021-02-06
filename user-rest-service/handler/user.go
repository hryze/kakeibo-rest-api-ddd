package handler

import (
	"encoding/json"
	"github.com/paypay3/kakeibo-rest-api-ddd/user-rest-service/errors"
	"net/http"

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
		errors.ErrorResponseByJSON(w, errors.NewInternalServerError(errors.NewErrorString("Internal Server Error")))
		return
	}

	out, err := h.userUsecase.SignUp(&in)
	if err != nil {
		errors.ErrorResponseByJSON(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(out); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
