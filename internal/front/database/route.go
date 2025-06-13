package database

import (
	"github.com/gofiber/fiber/v2"
	"github.com/jmoiron/sqlx"
)

type DatabaseRoute struct {
	App             *fiber.App
	DBPool          *sqlx.DB
	DatabaseHandler *DatabaseHandler
}

func NewRoute(db_pool *sqlx.DB, app *fiber.App) *DatabaseRoute {
	return &DatabaseRoute{
		App:             app,
		DBPool:          db_pool,
		DatabaseHandler: NewDatabaseHandler(db_pool),
	}
}

func (db *DatabaseRoute) RegisterAuthRoute() *DatabaseRoute {
	database := db.App.Group("/api/v1/front/database")

	database.Post("/", db.DatabaseHandler.Create)
	database.Get("/:db_uuid/detail", db.DatabaseHandler.GetDBDetail)
	database.Post("/:db_uuid/query", db.DatabaseHandler.Query)

	return db
}
