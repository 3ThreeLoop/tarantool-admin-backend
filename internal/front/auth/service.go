package auth

import (
	"tarantool-admin-api/pkg/responses"

	"github.com/jmoiron/sqlx"
)

type AuthServiceCreator interface {
	Login(username string, password string) (*LoginResponse, *responses.ErrorResponse)
	Register(register_req RegisterRequest) (*RegisterResponse, *responses.ErrorResponse)
}

type AuthService struct {
	DBPool   *sqlx.DB
	AuthRepo *AuthRepoImpl
}

func NewAuthService(db_pool *sqlx.DB) *AuthService {
	return &AuthService{
		DBPool:   db_pool,
		AuthRepo: NewAuthRepoImpl(db_pool),
	}
}

func (au *AuthService) Login(username string, password string) (*LoginResponse, *responses.ErrorResponse) {
	return au.AuthRepo.Login(username, password)
}

func (au *AuthService) Register(register_req RegisterRequest) (*RegisterResponse, *responses.ErrorResponse) {
	return au.AuthRepo.Register(register_req)
}
