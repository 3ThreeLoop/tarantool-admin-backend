package database

import (
	"fmt"
	"net/http"
	response "tarantool-admin-api/pkg/http/response"
	custom_log "tarantool-admin-api/pkg/logs"
	types "tarantool-admin-api/pkg/model"
	"tarantool-admin-api/pkg/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/jmoiron/sqlx"
)

type DatabaseHandler struct {
	DBPool          *sqlx.DB
	DatabaseService func(c *fiber.Ctx) *DatabaseService
}

func NewDatabaseHandler(db_pool *sqlx.DB) *DatabaseHandler {
	return &DatabaseHandler{
		DBPool: db_pool,
		DatabaseService: func(c *fiber.Ctx) *DatabaseService {
			user_context := c.Locals("UserContext")

			var us_ctx types.UserContext
			if context_map, ok := user_context.(types.UserContext); ok {
				us_ctx = context_map
			} else {
				custom_log.NewCustomLog("user_context_failed", "Failed to cast UserContext to map[string]interface{}", "warn")
				us_ctx = types.UserContext{}
			}

			return NewDatabaseService(&us_ctx, db_pool)
		},
	}
}

func (db *DatabaseHandler) Create(c *fiber.Ctx) error {
	var db_new_req DatabaseNewRequest
	v := utils.NewValidator()

	if err := db_new_req.bind(c, v); err != nil {
		return c.Status(http.StatusBadRequest).JSON(
			response.NewResponseError(
				utils.Translate("add_db_failed", nil, c),
				-2001,
				err,
			),
		)
	}

	resp, err := db.DatabaseService(c).Create(db_new_req)
	if err != nil {
		fmt.Println("hello error")
		return c.Status(http.StatusBadRequest).JSON(
			response.NewResponseError(
				utils.Translate(err.MessageID, nil, c),
				-2001,
				fmt.Errorf(utils.Translate(err.Err.Error(), nil, c)),
			),
		)
	}

	return c.Status(http.StatusOK).JSON(
		response.NewResponse(
			utils.Translate("add_db_success", nil, c),
			2000,
			resp,
		),
	)
}

func (db *DatabaseHandler) GetDBDetail(c *fiber.Ctx) error {
	uuid := c.Params("db_uuid")

	resp, err := db.DatabaseService(c).GetDBDetail(uuid)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(
			response.NewResponseError(
				utils.Translate(err.MessageID, nil, c),
				-2000,
				fmt.Errorf(utils.Translate(err.Err.Error(), nil, c)),
			),
		)
	}

	return c.Status(http.StatusOK).JSON(
		response.NewResponse(
			utils.Translate("db_detail_show_failed", nil, c),
			2000,
			resp,
		),
	)
}

func (db *DatabaseHandler) Query(c *fiber.Ctx) error {
	db_uuid := c.Params("db_uuid")

	var db_query_req DatabaseQueryRequest
	v := utils.NewValidator()
	if err := db_query_req.bind(c, v); err != nil {
		return c.Status(http.StatusBadRequest).JSON(
			response.NewResponseError(
				utils.Translate("query_db_failed", nil, c),
				-2005,
				err,
			),
		)
	}

	query_resp, err := db.DatabaseService(c).Query(db_uuid, db_query_req)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(
			response.NewResponseErrorWithDetail(
				utils.Translate(err.MessageID, nil, c),
				-2005,
				fmt.Errorf(utils.Translate(err.Err.Error(), nil, c)),
				err.Detail,
			),
		)
	}

	return c.Status(http.StatusOK).JSON(
		response.NewResponse(
			utils.Translate("query_db_success", nil, c),
			2005,
			query_resp,
		),
	)
}
