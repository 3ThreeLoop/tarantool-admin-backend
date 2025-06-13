package main

import (
	"tarantool-admin-api/configs"
	"tarantool-admin-api/db/postgresql"
	"tarantool-admin-api/handler"
	"tarantool-admin-api/pkg/logs"
	"tarantool-admin-api/pkg/redis"
	"tarantool-admin-api/pkg/swagger"
	"tarantool-admin-api/router"
	"fmt"
)

// @title       Mini Shop API
// @version     1.0.0
// @description Professional API documentation for the Mini Shop backend
// @BasePath    /api/v1

// @schemes     http
func main() {
	// load environment variable from .env file
	app_configs := configs.NewAppConfig()

	// log
	log_level := "info"
	logs.NewLog(log_level)

	// init postgresql database and connection pool
	pool, err := postgresql.ConnectDB()
	if err != nil {
		fmt.Println("Error connect database : ", err)
	}

	// init redis
	_ = redis.NewRedis()

	// init go fiber framework, cors and handler configuration
	apps := router.New()

	// swagger
	swagger.Setup(apps, app_configs.AppHost, app_configs.AppPort)

	// init router
	handler.NewServiceHandlers(apps, pool)

	// http server
	err = apps.Listen(fmt.Sprintf("%s:%d", app_configs.AppHost, app_configs.AppPort))
	if err != nil {
		fmt.Printf("%v", err)
	}
}
