package queryservice

import "github.com/paypay3/kakeibo-rest-api-ddd/user-rest-service/usecase/output"

type GroupQueryService interface {
	FetchGroupList(userID string) (*output.GroupList, error)
}
