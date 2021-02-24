package router

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/mux"

	"github.com/paypay3/kakeibo-rest-api-ddd/user-rest-service/config"
	"github.com/paypay3/kakeibo-rest-api-ddd/user-rest-service/infrastructure/auth"
	"github.com/paypay3/kakeibo-rest-api-ddd/user-rest-service/infrastructure/auth/imdb"
	"github.com/paypay3/kakeibo-rest-api-ddd/user-rest-service/infrastructure/externalapi"
	"github.com/paypay3/kakeibo-rest-api-ddd/user-rest-service/infrastructure/externalapi/client"
	"github.com/paypay3/kakeibo-rest-api-ddd/user-rest-service/infrastructure/middleware"
	"github.com/paypay3/kakeibo-rest-api-ddd/user-rest-service/infrastructure/persistence"
	"github.com/paypay3/kakeibo-rest-api-ddd/user-rest-service/infrastructure/persistence/rdb"
	"github.com/paypay3/kakeibo-rest-api-ddd/user-rest-service/interfaces/handler"
	"github.com/paypay3/kakeibo-rest-api-ddd/user-rest-service/usecase"
)

func Run() error {
	redisHandler, err := imdb.NewRedisHandler()
	if err != nil {
		return err
	}
	defer redisHandler.Pool.Close()

	mySQLHandler, err := rdb.NewMySQLHandler()
	if err != nil {
		return err
	}
	defer mySQLHandler.Conn.Close()

	accountApiHandler := client.NewAccountApiHandler()

	userRepository := persistence.NewUserRepository(mySQLHandler)
	sessionStore := auth.NewSessionStore(redisHandler)
	accountApi := externalapi.NewAccountApi(accountApiHandler)
	userUsecase := usecase.NewUserUsecase(userRepository, sessionStore, accountApi)
	userHandler := handler.NewUserHandler(userUsecase)

	router := mux.NewRouter()

	// register middleware
	router.Use(
		middleware.NewCorsMiddlewareFunc(),
		middleware.NewAuthMiddlewareFunc(sessionStore),
	)

	router.HandleFunc("/signup", userHandler.SignUp).Methods(http.MethodPost)
	router.HandleFunc("/login", userHandler.Login).Methods(http.MethodPost)
	router.HandleFunc("/logout", userHandler.Logout).Methods(http.MethodDelete)
	router.HandleFunc("/user", userHandler.FetchUserInfo).Methods(http.MethodGet)

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", config.Env.Server.Port),
		Handler: router,
	}

	errorCh := make(chan error, 1)
	go func() {
		if err := srv.ListenAndServe(); err != http.ErrServerClosed {
			errorCh <- err
		}
	}()

	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, syscall.SIGTERM, syscall.SIGINT)

	select {
	case err := <-errorCh:
		return err
	case s := <-signalCh:
		log.Printf("SIGNAL %s received", s.String())
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		if err := srv.Shutdown(ctx); err != nil {
			return err
		}
	}

	return nil
}
