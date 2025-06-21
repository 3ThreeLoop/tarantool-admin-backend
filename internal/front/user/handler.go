package user

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

type UserHandler struct {
	DBPool      *sqlx.DB
	UserService func(c *fiber.Ctx) *UserService
}

func NewUserHandler(db_pool *sqlx.DB) *UserHandler {
	return &UserHandler{
		DBPool: db_pool,
		UserService: func(c *fiber.Ctx) *UserService {
			user_context := c.Locals("UserContext")

			var us_ctx types.UserContext
			if context_map, ok := user_context.(types.UserContext); ok {
				us_ctx = context_map
			} else {
				custom_log.NewCustomLog("user_context_failed", "Failed to cast UserContext to map[string]interface{}", "warn")
				us_ctx = types.UserContext{}
			}

			return NewUserService(&us_ctx, db_pool)
		},
	}
}

func (u *UserHandler) Info(c *fiber.Ctx) error {
	user_info_resp, err := u.UserService(c).Info()
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(
			response.NewResponseError(
				utils.Translate(err.MessageID, nil, c),
				-3000,
				fmt.Errorf(utils.Translate(err.Err.Error(), nil, c)),
			),
		)
	}

	return c.Status(http.StatusOK).JSON(
		response.NewResponse(
			utils.Translate("userinfo_show_success", nil, c),
			3000,
			user_info_resp,
		),
	)
}
