package handler

import (
	"restful-api/internal/front/auth"
	"restful-api/pkg/middlewares"
	"fmt"

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

	// middleware
	middlewares.NewJwtMinddleWare(app, pool)

	app.Get("/hello", func(c *fiber.Ctx) error {
		fmt.Println("hello")
		return nil
	})

	return &FrontService{
		AuthRoute: au,
	}
}

func NewServiceHandlers(app *fiber.App, pool *sqlx.DB) *ServiceHandler {
	front := NewFrontService(app, pool)
	return &ServiceHandler{
		Front: front,
	}
}
