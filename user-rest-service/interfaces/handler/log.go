package handler

import (
	"log"

	"github.com/paypay3/kakeibo-rest-api-ddd/user-rest-service/apperrors"
	"github.com/paypay3/kakeibo-rest-api-ddd/user-rest-service/config"
)

func ErrorLog(err error) {
	appErr := apperrors.AsAppError(err)

	if config.Env.Log.Debug {
		log.Printf("%+v", appErr)
		return
	}

	if appErr.IsLevelError() || appErr.IsLevelCritical() {
		// Transfer logs to CloudWatch Logs
	}
}
