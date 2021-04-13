package middleware

import (
	"log"
	"net/http"

	"github.com/gorilla/context"

	"github.com/paypay3/kakeibo-rest-api-ddd/user-rest-service/apperrors"
	"github.com/paypay3/kakeibo-rest-api-ddd/user-rest-service/config"
)

func NewLoggingMiddlewareFunc() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			next.ServeHTTP(w, r)

			ctx, ok := context.GetOk(r, config.Env.ContextKey.AppError)
			if !ok {
				log.Print("failed to get appError context")
				return
			}

			err, ok := ctx.(error)
			if !ok {
				log.Print("failed to type assertion for error context")
				return
			}

			appErr := apperrors.AsAppError(err)

			if config.Env.Log.Debug {
				log.Printf("%+v", appErr)
				return
			}

			if appErr.IsLevelError() || appErr.IsLevelCritical() {
				// Transfer logs to CloudWatch Logs
			}
		})
	}
}
