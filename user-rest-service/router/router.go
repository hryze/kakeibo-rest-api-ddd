package router

import (
	"context"
	"errors"
	"flag"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/rs/cors"

	"github.com/paypay3/kakeibo-rest-api-ddd/user-rest-service/handler"
	"github.com/paypay3/kakeibo-rest-api-ddd/user-rest-service/infrastructure/externalapi"
	"github.com/paypay3/kakeibo-rest-api-ddd/user-rest-service/infrastructure/externalapi/client"
	"github.com/paypay3/kakeibo-rest-api-ddd/user-rest-service/infrastructure/persistence"
	"github.com/paypay3/kakeibo-rest-api-ddd/user-rest-service/infrastructure/persistence/db"
	"github.com/paypay3/kakeibo-rest-api-ddd/user-rest-service/usecase"
)

func Run() error {
	isLocal := flag.Bool("local", false, "Please specify -local flag")
	flag.Parse()

	if *isLocal {
		if err := godotenv.Load("../../.env"); err != nil {
			return err
		}
	}

	if len(os.Getenv("ALLOWED_ORIGIN")) == 0 ||
		len(os.Getenv("COOKIE_DOMAIN")) == 0 ||
		len(os.Getenv("ACCOUNT_HOST")) == 0 ||
		len(os.Getenv("MYSQL_DSN")) == 0 ||
		len(os.Getenv("REDIS_DSN")) == 0 {
		return errors.New("environment variable not defined")
	}

	redisHandler, err := db.NewRedisHandler()
	if err != nil {
		return err
	}

	mySQLHandler, err := db.NewMySQLHandler()
	if err != nil {
		return err
	}

	accountApiHandler := client.NewAccountApiHandler()

	userRepository := persistence.NewUserRepository(redisHandler, mySQLHandler)
	accountApi := externalapi.NewAccountApi(accountApiHandler)
	userUsecase := usecase.NewUserUsecase(userRepository, accountApi)
	userHandler := handler.NewUserHandler(userUsecase)

	router := mux.NewRouter()
	router.HandleFunc("/signup", userHandler.SignUp).Methods(http.MethodPost)

	allowedOrigin := os.Getenv("ALLOWED_ORIGIN")

	corsWrapper := cors.New(cors.Options{
		AllowedOrigins:   []string{allowedOrigin},
		AllowedMethods:   []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete, "OPTIONS"},
		AllowedHeaders:   []string{"Origin", "Content-Type", "Accept", "Accept-Language"},
		AllowCredentials: true,
	})

	srv := &http.Server{
		Addr:    ":8080",
		Handler: corsWrapper.Handler(router),
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