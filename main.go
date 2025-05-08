package main

import (
	"api-mini-shop/config"
	"api-mini-shop/db/postgresql"
	"api-mini-shop/handler"
	"api-mini-shop/pkg/logs"
	"api-mini-shop/pkg/swagger"
	"api-mini-shop/router"
	"fmt"
)

func main() {
	// load environment variable from .env file
	app_configs := config.NewAppConfig()

	// log
	log_level := "info"
	logs.NewLog(log_level)

	// init postgresql database and connection pool
	pool, err := postgresql.ConnectDB()
	if err != nil {
		fmt.Println("Error connect database : ", err)
	}

	// init go fiber framework, cors and handler configuration
	apps := router.New()

	// swagger
	swagger.NewSwagger(apps, app_configs.AppHost, app_configs.AppPort)

	// init router
	handler.NewServiceHandlers(apps, pool)

	// http server
	err = apps.Listen(fmt.Sprintf("%s:%d", app_configs.AppHost, app_configs.AppPort))
	if err != nil {
		fmt.Printf("%v", err)
	}
}
