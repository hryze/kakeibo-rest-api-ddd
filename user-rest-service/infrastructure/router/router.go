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
	"github.com/paypay3/kakeibo-rest-api-ddd/user-rest-service/infrastructure/persistence/query"
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

	sessionStore := auth.NewSessionStore(redisHandler)
	accountApi := externalapi.NewAccountApi(accountApiHandler)

	userRepository := persistence.NewUserRepository(mySQLHandler)
	userQueryService := query.NewUserQueryService(mySQLHandler)
	userUsecase := usecase.NewUserUsecase(userRepository, userQueryService, sessionStore, accountApi)
	userHandler := handler.NewUserHandler(userUsecase)

	groupRepository := persistence.NewGroupRepository(mySQLHandler)
	groupQueryService := query.NewGroupQueryServiceImpl(mySQLHandler)
	groupUsecase := usecase.NewGroupUsecase(groupRepository, groupQueryService, accountApi, userRepository)
	groupHandler := handler.NewGroupHandler(groupUsecase)

	router := mux.NewRouter()

	// Register auth middleware.
	router.Use(middleware.NewAuthMiddlewareFunc(sessionStore))

	router.HandleFunc("/signup", userHandler.SignUp).Methods(http.MethodPost)
	router.HandleFunc("/login", userHandler.Login).Methods(http.MethodPost)
	router.HandleFunc("/logout", userHandler.Logout).Methods(http.MethodDelete)
	router.HandleFunc("/user", userHandler.FetchLoginUser).Methods(http.MethodGet)
	router.HandleFunc("/groups", groupHandler.FetchGroupList).Methods(http.MethodGet)
	router.HandleFunc("/groups", groupHandler.StoreGroup).Methods(http.MethodPost)
	router.HandleFunc("/groups/{group_id:[0-9]+}", groupHandler.UpdateGroupName).Methods(http.MethodPut)
	router.HandleFunc("/groups/{group_id:[0-9]+}/users", groupHandler.StoreGroupUnapprovedUser).Methods(http.MethodPost)
	router.HandleFunc("/groups/{group_id:[0-9]+}/users", groupHandler.DeleteGroupApprovedUser).Methods(http.MethodDelete)
	router.HandleFunc("/groups/{group_id:[0-9]+}/users/approved", groupHandler.StoreGroupApprovedUser).Methods(http.MethodPost)

	// Apply cors middleware to top-level router.
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", config.Env.Server.Port),
		Handler: middleware.NewCorsMiddlewareFunc()(router),
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
