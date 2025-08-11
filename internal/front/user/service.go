package user

import (
	types "tarantool-admin-api/pkg/model"
	"tarantool-admin-api/pkg/responses"

	"github.com/jmoiron/sqlx"
)

type UserServiceCreator interface {
	Info() (*UserInfoResponse, *responses.ErrorResponse)
}

type UserService struct {
	DBPool      *sqlx.DB
	UserRepo    *UserRepoImpl
	UserContext *types.UserContext
}

func NewUserService(us_ctx *types.UserContext, db_pool *sqlx.DB) *UserService {
	return &UserService{
		DBPool:      db_pool,
		UserRepo:    NewUserRepoImpl(us_ctx, db_pool),
		UserContext: us_ctx,
	}
}

func (u *UserService) Info() (*UserInfoResponse, *responses.ErrorResponse) {
	return u.UserRepo.Info()
}
