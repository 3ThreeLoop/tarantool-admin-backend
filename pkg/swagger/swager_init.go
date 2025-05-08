package swagger

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/swagger"
)

func NewSwagger(apps *fiber.App, host string, port int) {
	apps.Get("/swagger/*", swagger.New(swagger.Config{
		URL: fmt.Sprintf("http://%s:%d/swagger/doc.json", host, port), // The url pointing to API definition
	}))
}
