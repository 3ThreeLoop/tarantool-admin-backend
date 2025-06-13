package database

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"
	custom_log "tarantool-admin-api/pkg/logs"
	types "tarantool-admin-api/pkg/model"
	"tarantool-admin-api/pkg/responses"
	tarantool_utils "tarantool-admin-api/pkg/tarantool"

	"github.com/jmoiron/sqlx"
	"github.com/tarantool/go-tarantool/v2"
	"github.com/tarantool/go-tarantool/v2/pool"
	"github.com/vmihailenco/msgpack/v5"
)

type DatabaseRepo interface {
	Create(new_db_req DatabaseNewRequest) (*DatabaseResponse, *responses.ErrorResponse)
	GetDBDetail(db_uuid string) (*DatabaseDetailResponse, *responses.ErrorResponse)
	Query(db_uuid string, db_query_req DatabaseQueryRequest) (*DatabaseQueryResultResponse, *responses.ErrorWithDetailResponse)
}

type DatabaseRepoImpl struct {
	DBPool      *sqlx.DB
	UserContext *types.UserContext
}

func NewDatabaseRepoImpl(us_ctx *types.UserContext, db_pool *sqlx.DB) *DatabaseRepoImpl {
	return &DatabaseRepoImpl{
		DBPool:      db_pool,
		UserContext: us_ctx,
	}
}

func (db *DatabaseRepoImpl) ShowOne(db_uuid string) (*DatabaseResponse, *responses.ErrorResponse) {
	// prepare query
	query := `
		SELECT 
			id, user_id, db_uuid, db_name, host, port, username, password, is_active, 
			created_by, created_at, updated_by, updated_at, deleted_by, deleted_at
		FROM tbl_users_databases
		WHERE deleted_at IS NULL
		AND db_uuid = $1
	`

	var database Database

	// execute query
	err := db.DBPool.Get(&database, query, db_uuid)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			custom_log.NewCustomLog("db_show_failed", err.Error(), "error")
			err_msg := &responses.ErrorResponse{}
			return nil, err_msg.NewErrorResponse("db_show_failed", fmt.Errorf("no_db_found"))
		}
		custom_log.NewCustomLog("db_show_failed", err.Error(), "error")
		err_msg := &responses.ErrorResponse{}
		return nil, err_msg.NewErrorResponse("db_show_failed", fmt.Errorf("get_db_error"))
	}

	return &DatabaseResponse{
		Database: database,
	}, nil
}

func (db *DatabaseRepoImpl) Create(new_db_req DatabaseNewRequest) (*DatabaseResponse, *responses.ErrorResponse) {
	var database_new_model DatabaseNewModel

	// test connect the creating database
	if err := tarantool_utils.TestTarantoolConnection(
		new_db_req.Host,
		int(new_db_req.Port),
		new_db_req.Username,
		new_db_req.Password,
	); err != nil {
		custom_log.NewCustomLog("add_db_failed", err.Error(), "error")
		err_msg := &responses.ErrorResponse{}
		return nil, err_msg.NewErrorResponse("add_db_failed", fmt.Errorf("invalid_connection_settings"))
	}

	// create insert model
	if err := database_new_model.new(new_db_req, db.UserContext, db.DBPool); err != nil {
		custom_log.NewCustomLog("add_db_failed", err.Error(), "error")
		err_msg := &responses.ErrorResponse{}
		return nil, err_msg.NewErrorResponse("add_db_failed", fmt.Errorf("invalid_info_to_add_db"))
	}

	// prepare query
	query := `
		INSERT INTO tbl_users_databases (
			id, user_id, db_uuid, db_name, host, port, username, password,
			is_active, created_by, created_at
		) VALUES (
			:id, :user_id, :db_uuid, :db_name, :host, :port, :username, :password,
			:is_active, :created_by, :created_at 
		)
	`

	// execute request
	_, err := db.DBPool.NamedExec(query, database_new_model)
	if err != nil {
		custom_log.NewCustomLog("add_db_failed", err.Error(), "error")
		err_msg := &responses.ErrorResponse{}
		return nil, err_msg.NewErrorResponse("add_db_failed", fmt.Errorf("error_add_db"))
	}

	return db.ShowOne(database_new_model.DBUUID)
}

func (db *DatabaseRepoImpl) GetDBDetail(db_uuid string) (*DatabaseDetailResponse, *responses.ErrorResponse) {
	// get database info
	db_resp, err_resp := db.ShowOne(db_uuid)
	if err_resp != nil {
		return nil, err_resp
	}

	// connect database to get data
	conn, err := tarantool_utils.ConnectTarantool(
		db_resp.Database.Host,
		int(db_resp.Database.Port),
		db_resp.Database.Username,
		db_resp.Database.Password,
	)
	if err != nil {
		custom_log.NewCustomLog("db_detail_show_failed", err.Error(), "error")
		err_msg := &responses.ErrorResponse{}
		return nil, err_msg.NewErrorResponse("db_detail_show_failed", fmt.Errorf("failed_connect_to_target_db"))
	}

	// close connection after function end
	defer conn.Close()

	// execute request to select all space
	var resp interface{}

	resp, err = conn.Do(
		tarantool.NewSelectRequest("_vspace").
			Index("primary").
			Iterator(tarantool.IterAll).
			Limit(1000),
		pool.ANY,
	).Get()
	if err != nil {
		custom_log.NewCustomLog("db_detail_show_failed", err.Error(), "error")
		err_msg := &responses.ErrorResponse{}
		return nil, err_msg.NewErrorResponse("db_detail_show_failed", fmt.Errorf("failed_to_get_db_detail"))
	}

	// cast the result to []interface{}
	rawRows, ok := resp.([]interface{})
	if !ok {
		custom_log.NewCustomLog("db_detail_show_failed", err.Error(), "error")
		err_msg := &responses.ErrorResponse{}
		return nil, err_msg.NewErrorResponse("db_detail_show_failed", fmt.Errorf("unexpected_data_format"))
	}

	// loop through and filter system spaces
	var spaces []TarantoolSpace
	for _, row := range rawRows {
		// cast row to []interface{}
		// and validate of each row have atleast 7 field
		// to prevent incomplete data
		fields, ok := row.([]interface{})
		if !ok || len(fields) < 7 {
			continue
		}

		// filter out system spaces (name is at index 2)
		name, ok := fields[2].(string)
		if !ok || strings.HasPrefix(name, "_") {
			continue
		}

		// marshal the row into binary msgpack, and unmarshal it back to struct
		// this can handle complex nested structure
		b, _ := msgpack.Marshal(fields)
		var space TarantoolSpace
		if err := msgpack.Unmarshal(b, &space); err != nil {
			continue
		}

		// append to go variable
		spaces = append(spaces, space)
	}

	return &DatabaseDetailResponse{
		DatabaseDetail: DatabaseDetail{
			DBName: db_resp.Database.DBName,
			DBUUID: db_resp.Database.DBUUID,
			Spaces: spaces,
		},
	}, nil
}

func (db *DatabaseRepoImpl) Query(db_uuid string, db_query_req DatabaseQueryRequest) (*DatabaseQueryResultResponse, *responses.ErrorWithDetailResponse) {
	// get database info
	db_resp, err_resp := db.ShowOne(db_uuid)
	if err_resp != nil {
		err_msg := &responses.ErrorWithDetailResponse{}
		return nil, err_msg.NewErrorResponse(err_resp.MessageID, err_resp.Err, fmt.Errorf(""))
	}

	// connect database to get data
	conn, err := tarantool_utils.ConnectTarantool(
		db_resp.Database.Host,
		int(db_resp.Database.Port),
		db_resp.Database.Username,
		db_resp.Database.Password,
	)
	if err != nil {
		custom_log.NewCustomLog("query_db_failed", err.Error(), "error")
		err_msg := &responses.ErrorWithDetailResponse{}
		return nil, err_msg.NewErrorResponse("query_db_failed", fmt.Errorf("failed_connect_to_target_db"), err)
	}

	// close connection after function end
	defer conn.Close()

	query := fmt.Sprintf(`
		return (function()
			local sql = [[%s]]
			local res, err_msg = box.execute(sql)

			if res == nil then
				return {
					metadata = {},
					rows = {},
					error = tostring(err_msg or "SQL execution returned nil")
				}
			end

			-- process rows
			local processed_rows = {}
			for _, tuple in ipairs(res.rows or {}) do
				local row = {}
				for i = 1, #tuple do
					local val = tuple[i]
					if val == nil then
						row[i] = "null"
					elseif type(val) == 'cdata' or type(val) == 'userdata' then
						row[i] = tostring(val)
					else
						row[i] = val
					end
				end
				table.insert(processed_rows, row)
			end

			res.rows = processed_rows
			return res
		end)()
	`, db_query_req.Query)

	// run query string
	row_result, err := conn.Do(tarantool.NewEvalRequest(query), pool.ANY).Get()
	if err != nil {
		custom_log.NewCustomLog("query_db_failed", err.Error(), "error")
		err_msg := &responses.ErrorWithDetailResponse{}
		return nil, err_msg.NewErrorResponse("query_db_failed", fmt.Errorf("failed_to_query_db"), err)
	}

	// map the result into go struct
	result, err := tarantool_utils.MapToDetailedQueryResult(row_result)
	if err != nil {
		custom_log.NewCustomLog("query_db_failed", err.Error(), "error")
		err_msg := &responses.ErrorWithDetailResponse{}
		return nil, err_msg.NewErrorResponse("query_db_failed", fmt.Errorf("failed_to_query_db"), err)
	}

	return &DatabaseQueryResultResponse{
		QueryResult: *result,
	}, nil
}
