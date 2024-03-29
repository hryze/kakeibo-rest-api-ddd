package gateway

import (
	"github.com/paypay3/kakeibo-rest-api-ddd/user-rest-service/domain/groupdomain"
	"github.com/paypay3/kakeibo-rest-api-ddd/user-rest-service/domain/userdomain"
)

type AccountApi interface {
	PostInitStandardBudgets(userID userdomain.UserID) error
	PostInitGroupStandardBudgets(groupID groupdomain.GroupID) error
}
