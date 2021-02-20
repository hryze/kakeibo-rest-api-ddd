package sessionstore

import "github.com/paypay3/kakeibo-rest-api-ddd/user-rest-service/domain/userdomain"

type SessionStore interface {
	StoreLoginInfo(sessionID string, userID userdomain.UserID) error
	DeleteLoginInfo(sessionID string) error
}
