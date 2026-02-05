package user

import (
	"context"

	"gitlab.com/skeleton-golang/internal/pkg/constants"
)

type Service interface {
	GetAllUser(ctx context.Context, req GetUsersReq) (res constants.DefaultResponse, err error)
}
