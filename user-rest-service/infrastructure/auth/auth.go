package auth

import (
	"github.com/paypay3/kakeibo-rest-api-ddd/user-rest-service/apierrors"
	"github.com/paypay3/kakeibo-rest-api-ddd/user-rest-service/config"
	"github.com/paypay3/kakeibo-rest-api-ddd/user-rest-service/domain/userdomain"
	"github.com/paypay3/kakeibo-rest-api-ddd/user-rest-service/infrastructure/auth/imdb"
)

type sessionStore struct {
	*imdb.RedisHandler
}

func NewSessionStore(redisHandler *imdb.RedisHandler) *sessionStore {
	return &sessionStore{redisHandler}
}

func (s *sessionStore) AddSessionID(sessionID string, userID userdomain.UserID) error {
	conn := s.RedisHandler.Pool.Get()
	defer conn.Close()

	expirationS := int(config.Env.Cookie.Expiration.Seconds())

	if _, err := conn.Do("SET", sessionID, userID.Value(), "EX", expirationS); err != nil {
		return apierrors.NewInternalServerError(apierrors.NewErrorString("Internal Server Error"))
	}

	return nil
}
