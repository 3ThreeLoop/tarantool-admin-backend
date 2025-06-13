package handler

import (
	"tarantool-admin-api/internal/front/auth"
	"tarantool-admin-api/internal/front/database"
	"tarantool-admin-api/pkg/middlewares"

	"github.com/gofiber/fiber/v2"
	"github.com/jmoiron/sqlx"
)

// group all the module factories
type ServiceHandler struct {
	Front *FrontService
}

// register modules route here
type FrontService struct {
	AuthRoute     *auth.AuthRoute
	DatabaseRoute *database.DatabaseRoute
}

func NewFrontService(app *fiber.App, pool *sqlx.DB) *FrontService {
	// register auth route
	au := auth.NewRoute(pool, app).RegisterAuthRoute()

	// middleware
	middlewares.NewJwtMinddleWare(app, pool)

	// register database route
	db := database.NewRoute(pool, app).RegisterAuthRoute()

	return &FrontService{
		AuthRoute:     au,
		DatabaseRoute: db,
	}
}

func NewServiceHandlers(app *fiber.App, pool *sqlx.DB) *ServiceHandler {
	front := NewFrontService(app, pool)
	return &ServiceHandler{
		Front: front,
	}
}
