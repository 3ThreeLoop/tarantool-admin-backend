package handler

import (
	"api-mini-shop/internal/front/auth"

	"github.com/gofiber/fiber/v2"
	"github.com/jmoiron/sqlx"
)

// group all the module factories
type ServiceHandler struct {
	Front *FrontService
}

// register modules route here
type FrontService struct {
	AuthRoute *auth.AuthRoute
}

func NewFrontService(app *fiber.App, pool *sqlx.DB) *FrontService {
	// register auth route
	au := auth.NewRoute(pool, app).RegisterAuthRoute()

	return &FrontService{
		AuthRoute: au,
	}
}

func NewServiceHandlers(app *fiber.App, pool *sqlx.DB) *ServiceHandler {
	return &ServiceHandler{
		Front: NewFrontService(app, pool),
	}
}
