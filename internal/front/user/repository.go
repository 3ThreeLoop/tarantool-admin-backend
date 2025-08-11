package user

import (
	"database/sql"
	"errors"
	"fmt"
	custom_log "tarantool-admin-api/pkg/logs"
	types "tarantool-admin-api/pkg/model"
	"tarantool-admin-api/pkg/responses"

	"github.com/jmoiron/sqlx"
)

type UserRepo interface {
	Info() (*UserInfoResponse, *responses.ErrorResponse)
}

type UserRepoImpl struct {
	DBPool      *sqlx.DB
	UserContext *types.UserContext
}

func NewUserRepoImpl(us_ctx *types.UserContext, db_pool *sqlx.DB) *UserRepoImpl {
	return &UserRepoImpl{
		DBPool:      db_pool,
		UserContext: us_ctx,
	}
}

func (u *UserRepoImpl) Info() (*UserInfoResponse, *responses.ErrorResponse) {
	// prepare query
	query := `
		SELECT
			id, user_uuid, first_name, last_name, user_name, password, email, 
			login_session, profile_photo, status_id, "order", created_by, created_at,
			updated_by, updated_at, deleted_by, deleted_at
		FROM tbl_users
		WHERE deleted_at IS NULL AND id = $1
	`

	// execute request
	var user User
	err := u.DBPool.Get(&user, query, u.UserContext.Id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			custom_log.NewCustomLog("userinfo_show_failed", err.Error(), "error")
			err_msg := &responses.ErrorResponse{}
			return nil, err_msg.NewErrorResponse("userinfo_show_failed", fmt.Errorf("no_user_found"))
		}
		custom_log.NewCustomLog("userinfo_show_failed", err.Error(), "error")
		err_msg := &responses.ErrorResponse{}
		return nil, err_msg.NewErrorResponse("userinfo_show_failed", fmt.Errorf("get_user_error"))
	}

	// get user databases
	user_databases, err := u.getUserDatabases(user.ID)
	if err != nil {
		custom_log.NewCustomLog("userinfo_show_failed", err.Error(), "error")
		err_msg := &responses.ErrorResponse{}
		return nil, err_msg.NewErrorResponse("userinfo_show_failed", fmt.Errorf("error_get_user_database"))
	}

	user.UserDatabases = user_databases

	return &UserInfoResponse{
		UserInfo: user,
	}, nil

}

func (u *UserRepoImpl) getUserDatabases(user_id int) ([]UserDatabase, error) {
	// prepare query
	query := `
		SELECT
			db_uuid, db_name, host, port
		FROM tbl_users_databases
		WHERE deleted_at IS NULL
		AND user_id = $1
	`

	// execute query
	var databases []UserDatabase
	if err := u.DBPool.Select(&databases, query, user_id); err != nil {
		return nil, err
	}

	if databases == nil {
		databases = []UserDatabase{}
	}

	return databases, nil
}
