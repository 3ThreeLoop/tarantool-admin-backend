package user

import (
	"github.com/gofiber/fiber/v2"
	"github.com/jmoiron/sqlx"
)

type UserRoute struct {
	App         *fiber.App
	DBPool      *sqlx.DB
	UserHandler *UserHandler
}

func NewRoute(db_pool *sqlx.DB, app *fiber.App) *UserRoute {
	return &UserRoute{
		App:         app,
		DBPool:      db_pool,
		UserHandler: NewUserHandler(db_pool),
	}
}

func (u *UserRoute) RegisterUserRoute() *UserRoute {
	user := u.App.Group("/api/v1/front/user")

	user.Get("/info", u.UserHandler.Info)

	return u
}
