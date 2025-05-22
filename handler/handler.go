package handler

import (
	"api-mini-shop/internal/front/auth"
	user "api-mini-shop/internal/front/user"
	"api-mini-shop/pkg/middlewares"

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
	UserRoute *user.UserRoute
}

func NewFrontService(app *fiber.App, pool *sqlx.DB) *FrontService {
	// register auth route
	au := auth.NewRoute(pool, app).RegisterAuthRoute()

	// middleware
	middlewares.NewJwtMinddleWare(app, pool)

	user := user.NewUserRoute(app, pool).RegisterUserRoute()

	// app.Get("/hello", func(c *fiber.Ctx) error {
	// 	fmt.Println("hello")
	// 	return nil
	// })

	return &FrontService{
		AuthRoute: au,
		UserRoute: user,
	}
}

func NewServiceHandlers(app *fiber.App, pool *sqlx.DB) *ServiceHandler {
	front := NewFrontService(app, pool)
	return &ServiceHandler{
		Front: front,
	}
}
