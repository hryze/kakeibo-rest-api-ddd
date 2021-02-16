package handler

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/paypay3/kakeibo-rest-api-ddd/user-rest-service/apierrors"
	"github.com/paypay3/kakeibo-rest-api-ddd/user-rest-service/config"
	"github.com/paypay3/kakeibo-rest-api-ddd/user-rest-service/interfaces/presenter"
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
		presenter.ErrorJSON(w, apierrors.NewBadRequestError(apierrors.NewErrorString("正しいデータを入力してください")))
		return
	}

	out, err := h.userUsecase.SignUp(&in)
	if err != nil {
		presenter.ErrorJSON(w, err)
		return
	}

	presenter.JSON(w, http.StatusCreated, out)
}

func (h *userHandler) Login(w http.ResponseWriter, r *http.Request) {
	var in input.LoginUser
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		presenter.ErrorJSON(w, apierrors.NewBadRequestError(apierrors.NewErrorString("正しいデータを入力してください")))
		return
	}

	out, err := h.userUsecase.Login(&in)
	if err != nil {
		presenter.ErrorJSON(w, err)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     config.Env.Cookie.Name,
		Value:    out.Cookie.SessionID,
		Expires:  time.Now().Add(config.Env.Cookie.Expiration),
		Domain:   config.Env.Cookie.Domain,
		Secure:   config.Env.Cookie.Secure,
		HttpOnly: true,
	})

	presenter.JSON(w, http.StatusCreated, out)
}
