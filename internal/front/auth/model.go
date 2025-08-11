package auth

import (
	"fmt"
	"os"
	"strings"
	custom_log "tarantool-admin-api/pkg/logs"
	"tarantool-admin-api/pkg/postgres"
	"tarantool-admin-api/pkg/responses"
	"tarantool-admin-api/pkg/utils"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type LoginRequest struct {
	UserName string `json:"user_name" validate:"required"`
	Password string `json:"password" validate:"required"`
}

func (au *LoginRequest) Bind(c *fiber.Ctx, v *utils.Validator) error {
	if err := c.BodyParser(au); err != nil {
		custom_log.NewCustomLog("login_failed", err.Error(), "error")
		return fmt.Errorf(utils.Translate("invalid_body", nil, c))
	}

	if err := v.Validate(au, c); err != nil {
		custom_log.NewCustomLog("login_failed", err.Error(), "error")
		return err
	}

	return nil
}

type LoginResponse struct {
	Auth Auth `json:"auth"`
}

type Auth struct {
	Token     string `json:"token"`
	TokenType string `json:"token_type"`
}

type User struct {
	UserUUID uuid.UUID `json:"user_uuid" db:"user_uuid"`
}
type UserInfo struct {
	ID           int    `json:"id" db:"id"`
	UserUUID     string `json:"user_uuid" db:"user_uuid"`
	UserName     string `json:"user_name" db:"user_name"`
	LoginSession string `json:"login_session" db:"login_session"`
	StatusID     int    `json:"status_id" db:"status_id"`
}

type RegisterRequest struct {
	FirstName       string `json:"first_name" validate:"required,min=2,max=100"`
	LastName        string `json:"last_name" validate:"required,min=2,max=100"`
	UserName        string `json:"user_name" validate:"required,min=3,max=50,alphanum"`
	Password        string `json:"password" validate:"required,min=6,max=100"`
	ConfirmPassword string `json:"confirm_password" validate:"required,min=6,max=100"`
	Email           string `json:"email" validate:"required,email"`
	ProfilePhoto    string `json:"profile_photo" validate:"omitempty"`
}

func (au *RegisterRequest) Bind(c *fiber.Ctx, v *utils.Validator) error {
	if err := c.BodyParser(au); err != nil {
		custom_log.NewCustomLog("register_failed", err.Error(), "error")
		return fmt.Errorf(utils.Translate("invalid_body", nil, c))
	}

	if err := v.Validate(au, c); err != nil {
		custom_log.NewCustomLog("register_failed", err.Error(), "error")
		return err
	}

	if au.Password != au.ConfirmPassword {
		return fmt.Errorf("confirm_pass_and_pass_dont_match")
	}

	return nil
}

type RegisterModel struct {
	ID           uint64    `db:"id" json:"-"`
	UserUUID     string    `db:"user_uuid" json:"user_uuid"`
	FirstName    string    `db:"first_name" json:"first_name"`
	LastName     string    `db:"last_name" json:"last_name"`
	Username     string    `db:"user_name" json:"user_name"`
	Password     string    `db:"password" json:"password"`
	Email        string    `db:"email" json:"email"`
	ProfilePhoto string    `db:"profile_photo" json:"profile_photo"`
	StatusID     uint64    `db:"status_id" json:"status_id"`
	Order        uint64    `db:"order" json:"order"`
	CreatedBy    uint64    `db:"created_by" json:"created_by"`
	CreatedAt    time.Time `db:"created_at" json:"created_at"`
}

func (au *RegisterModel) New(register_req RegisterRequest, conn *sqlx.DB) *responses.ErrorWithDetailResponse {
	// check user exists
	username := strings.TrimSpace(strings.ToUpper(register_req.UserName))
	is_username, err := postgres.IsExists("tbl_users", "user_name", username, conn)

	fmt.Println()
	if err != nil {
		err_msg := &responses.ErrorWithDetailResponse{}
		return err_msg.NewErrorResponse("register_failed", fmt.Errorf("cannot_check_username"), err)
	} else {
		if is_username {
			err_msg := &responses.ErrorWithDetailResponse{}
			return err_msg.NewErrorResponse("register_failed", fmt.Errorf("username_already_exists"), fmt.Errorf("register username is already exists"))
		}
	}

	// assign default profile photo if none provided
	if register_req.ProfilePhoto == "" {
		register_req.ProfilePhoto = "user1.png"
	}

	// generate new UUID
	uuid, _ := uuid.NewV7()

	// get current OS time
	time_zone := os.Getenv("APP_TIMEZONE")
	location, err := time.LoadLocation(time_zone)
	if err != nil {
		err_msg := &responses.ErrorWithDetailResponse{}
		return err_msg.NewErrorResponse("register_failed", fmt.Errorf("technical_error"), err)
	}
	now := time.Now().In(location)

	// get next sequence ID
	id, err := postgres.GetSeqNextVal("tbl_users_id_seq", conn)
	if err != nil {
		err_msg := &responses.ErrorWithDetailResponse{}
		return err_msg.NewErrorResponse("register_failed", fmt.Errorf("technical_error"), err)
	}

	// assign values to the RegisterModel
	au.ID = uint64(*id)
	au.UserUUID = uuid.String()
	au.FirstName = register_req.FirstName
	au.LastName = register_req.LastName
	au.Username = username
	au.Password = register_req.Password
	au.Email = register_req.Email
	au.ProfilePhoto = register_req.ProfilePhoto
	au.StatusID = 1
	au.Order = uint64(*id)
	au.CreatedBy = uint64(*id)
	au.CreatedAt = now

	return nil
}

type RegisterResponse struct {
	UserInfo RegisterModel `json:"user_info"`
}
