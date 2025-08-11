package database

import (
	"fmt"
	"os"
	custom_log "tarantool-admin-api/pkg/logs"
	types "tarantool-admin-api/pkg/model"
	"tarantool-admin-api/pkg/postgres"
	tarantool_utils "tarantool-admin-api/pkg/tarantool"
	"tarantool-admin-api/pkg/utils"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type Database struct {
	ID        uint64     `json:"-" db:"id"`
	UserID    uint64     `json:"user_id" db:"user_id"`
	DBUUID    string     `json:"db_uuid" db:"db_uuid"`
	DBName    string     `json:"db_name" db:"db_name"`
	Host      string     `json:"host" db:"host"`
	Port      uint64     `json:"port" db:"port"`
	Username  string     `json:"username" db:"username"`
	Password  string     `json:"password" db:"password"`
	IsActive  bool       `json:"is_active" db:"is_active"`
	CreatedBy uint64     `json:"-" db:"created_by"`
	CreatedAt time.Time  `json:"-" db:"created_at"`
	UpdatedBy *uint64    `json:"-" db:"updated_by"`
	UpdatedAt *time.Time `json:"-" db:"updated_at"`
	DeletedBy *uint64    `json:"-" db:"deleted_by"`
	DeletedAt *time.Time `json:"-" db:"deleted_at"`
}

type DatabaseResponse struct {
	Database Database `json:"database"`
}

type DatabaseNewRequest struct {
	DBName   string `json:"db_name" validate:"required"`
	Host     string `json:"host" validate:"required"`
	Port     uint64 `json:"port" validate:"required"`
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

func (db *DatabaseNewRequest) bind(c *fiber.Ctx, v *utils.Validator) error {
	if err := c.BodyParser(db); err != nil {
		custom_log.NewCustomLog("add_db_failed", err.Error(), "error")
		return fmt.Errorf(utils.Translate("invalid_body", nil, c))
	}

	if err := v.Validate(db, c); err != nil {
		custom_log.NewCustomLog("add_db_failed", err.Error(), "error")
		return err
	}

	return nil
}

type DatabaseNewModel struct {
	ID        uint64    `db:"id"`
	UserID    uint64    `db:"user_id"`
	DBUUID    string    `db:"db_uuid"`
	DBName    string    `db:"db_name"`
	Host      string    `db:"host"`
	Port      int       `db:"port"`
	Username  string    `db:"username"`
	Password  string    `db:"password"`
	IsActive  bool      `db:"is_active"`
	CreatedBy int       `db:"created_by"`
	CreatedAt time.Time `db:"created_at"`
}

func (db *DatabaseNewModel) new(db_new_req DatabaseNewRequest, us_ctx *types.UserContext, conn *sqlx.DB) error {
	// get sequence next value
	id, err := postgres.GetSeqNextVal("tbl_users_databases_id_seq", conn)
	if err != nil {
		return fmt.Errorf("error get seq : %w", err)
	}

	// generate new uuid
	uuid, err := uuid.NewV7()
	if err != nil {
		return fmt.Errorf("error generate new uuid : %w", err)
	}

	// get current os time
	time_zone := os.Getenv("APP_TIMEZONE")
	location, err := time.LoadLocation(time_zone)
	if err != nil {
		return fmt.Errorf("error load location : %w", err)
	}
	now := time.Now().In(location)

	db.ID = uint64(*id)
	db.UserID = uint64(us_ctx.Id)
	db.DBUUID = uuid.String()
	db.DBName = db_new_req.DBName
	db.Host = db_new_req.Host
	db.Port = int(db_new_req.Port)
	db.Username = db_new_req.Username
	db.Password = db_new_req.Password
	db.IsActive = true
	db.CreatedBy = us_ctx.Id
	db.CreatedAt = now

	return nil
}

type SpaceFormatField struct {
	Name       string `msgpack:"name" json:"name"`
	Type       string `msgpack:"type" json:"type"`
	IsNullable bool   `msgpack:"is_nullable" json:"is_nullable"`
}

type TarantoolSpace struct {
	ID         uint32                 `msgpack:"0" json:"-"`
	Owner      uint32                 `msgpack:"1" json:"-"`
	Name       string                 `msgpack:"2" json:"name"`
	Engine     string                 `msgpack:"3" json:"-"`
	FieldCount uint32                 `msgpack:"4" json:"-"`
	Flags      map[string]interface{} `msgpack:"5" json:"-"`
	Format     []SpaceFormatField     `msgpack:"6" json:"format"`
}

type DatabaseDetailResponse struct {
	DatabaseDetail DatabaseDetail `json:"database_detail"`
}

type DatabaseDetail struct {
	DBName string           `json:"db_name"`
	DBUUID string           `json:"db_uuid"`
	Spaces []TarantoolSpace `json:"spaces"`
}

type DatabaseQueryRequest struct {
	Query string `json:"query" validate:"required"`
}

func (db *DatabaseQueryRequest) bind(c *fiber.Ctx, v *utils.Validator) error {
	if err := c.BodyParser(db); err != nil {
		custom_log.NewCustomLog("query_db_failed", err.Error(), "error")
		return fmt.Errorf(utils.Translate("invalid_body", nil, c))
	}

	if err := v.Validate(db, c); err != nil {
		custom_log.NewCustomLog("query_db_failed", err.Error(), "error")
		return err
	}

	return nil
}

type DatabaseQueryResultResponse struct {
	QueryResult tarantool_utils.QueryResult `json:"query_result"`
}
