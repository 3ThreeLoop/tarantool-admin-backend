package database

import (
	types "tarantool-admin-api/pkg/model"
	"tarantool-admin-api/pkg/responses"

	"github.com/jmoiron/sqlx"
)

type DatabaseServiceCreator interface {
	Create(new_db_req DatabaseNewRequest) (*DatabaseResponse, *responses.ErrorResponse)
	GetDBDetail(db_uuid string) (*DatabaseDetailResponse, *responses.ErrorResponse)
	Query(db_uuid string, db_query_req DatabaseQueryRequest) (*DatabaseQueryResultResponse, *responses.ErrorWithDetailResponse)
}

type DatabaseService struct {
	DBPool       *sqlx.DB
	DatabaseRepo *DatabaseRepoImpl
	UserContext  *types.UserContext
}

func NewDatabaseService(us_ctx *types.UserContext, db_pool *sqlx.DB) *DatabaseService {
	return &DatabaseService{
		DBPool:       db_pool,
		DatabaseRepo: NewDatabaseRepoImpl(us_ctx, db_pool),
		UserContext:  us_ctx,
	}
}

func (db *DatabaseService) Create(new_db_req DatabaseNewRequest) (*DatabaseResponse, *responses.ErrorResponse) {
	return db.DatabaseRepo.Create(new_db_req)
}

func (db *DatabaseService) GetDBDetail(db_uuid string) (*DatabaseDetailResponse, *responses.ErrorResponse) {
	return db.DatabaseRepo.GetDBDetail(db_uuid)
}

func (db *DatabaseService) Query(db_uuid string, db_query_req DatabaseQueryRequest) (*DatabaseQueryResultResponse, *responses.ErrorWithDetailResponse) {
	return db.DatabaseRepo.Query(db_uuid, db_query_req)
}
