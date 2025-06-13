package auth

import (
	"fmt"
	"net/http"
	response "tarantool-admin-api/pkg/http/response"
	"tarantool-admin-api/pkg/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/jmoiron/sqlx"
)

type AuthHandler struct {
	DBPool      *sqlx.DB
	AuthService *AuthService
}

func NewAuthHandler(db_pool *sqlx.DB) *AuthHandler {
	return &AuthHandler{
		DBPool:      db_pool,
		AuthService: NewAuthService(db_pool),
	}
}

// @Summary      Login
// @Description  Authenticates a user and returns a token
// @Tags         Admin/Auth
// @Accept       json
// @Produce      json
// @Param        user  body      auth.LoginRequest  true  "Credentials to use"
// @Success      200   {object}  auth.LoginResponse
// @Failure      400   {object}  utils.Error
// @Failure      401   {object}  utils.Error
// @Router       /front/auth [post]
func (au *AuthHandler) Login(c *fiber.Ctx) error {
	var login_request LoginRequest
	v := utils.NewValidator()

	if err := login_request.Bind(c, v); err != nil {
		return c.Status(http.StatusBadRequest).JSON(
			response.NewResponseError(
				utils.Translate("login_failed", nil, c),
				-1000,
				err,
			),
		)
	}

	resp, err := au.AuthService.Login(login_request.UserName, login_request.Password)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(
			response.NewResponseError(
				utils.Translate(err.MessageID, nil, c),
				-1000,
				fmt.Errorf(utils.Translate(err.Err.Error(), nil, c)),
			),
		)
	}

	return c.Status(http.StatusOK).JSON(
		response.NewResponse(
			utils.Translate("login_success", nil, c),
			1000,
			resp,
		),
	)
}

func (au *AuthHandler) Register(c *fiber.Ctx) error {
	var register_req RegisterRequest
	v := utils.NewValidator()

	if err := register_req.Bind(c, v); err != nil {
		return c.Status(http.StatusBadRequest).JSON(
			response.NewResponseError(
				utils.Translate("register_failed", nil, c),
				-1001,
				err,
			),
		)
	}

	register_resp, err := au.AuthService.Register(register_req)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(
			response.NewResponseError(
				utils.Translate(err.MessageID, nil, c),
				-1001,
				fmt.Errorf(utils.Translate(err.Err.Error(), nil, c)),
			),
		)
	}

	return c.Status(http.StatusBadRequest).JSON(
		response.NewResponse(
			utils.Translate("register_success", nil, c),
			1001,
			register_resp,
		),
	)
}
