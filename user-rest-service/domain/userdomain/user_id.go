package userdomain

import (
	"strings"
	"unicode/utf8"

	"github.com/paypay3/kakeibo-rest-api-ddd/user-rest-service/apperrors"
)

type UserID string

const (
	minUserIDLength = 1
	maxUserIDLength = 10
)

func NewUserID(userID string) (UserID, error) {
	if n := utf8.RuneCountInString(userID); n < minUserIDLength || n > maxUserIDLength {
		return "", apperrors.Errorf(
			"user id must be %d or more and %d or less: %s",
			minUserIDLength, maxUserIDLength, userID,
		)
	}

	if strings.Contains(userID, " ") ||
		strings.Contains(userID, "　") {
		return "", apperrors.Errorf("user id cannot contain spaces: %s", userID)
	}

	return UserID(userID), nil
}

func (i UserID) Value() string {
	return string(i)
}
