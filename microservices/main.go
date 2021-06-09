package main

import (
	"os"

	"github.com/aofiee/diablos/routes"
	fiber "github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/requestid"
	"github.com/gofiber/fiber/v2/utils"
	"github.com/joho/godotenv"
)

var (
	app = fiber.New()
)

func init() {
	err := godotenv.Load(".env")
	if err != nil {
		panic("Error loading .env file")
	}
	app.Use(requestid.New())
	app.Use(requestid.New(requestid.Config{
		Header: "Diablos-Service-Header",
		Generator: func() string {
			return utils.UUID()
		},
	}))
	app.Use(logger.New(logger.Config{
		Format:     "${pid} ${status} - ${method} ${path}\n",
		TimeFormat: "02-Jan-2006",
		TimeZone:   "Asia/Bangkok",
	}))
}

func main() {
	app.Use(cors.New(cors.Config{
		AllowOrigins: os.Getenv("ALLOW_ORIGINS"),
		AllowHeaders: "Origin, Content-Type, Accept",
	}))
	//not AuthorizationRequired
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World!")
	})
	app.Post("/login", routes.Auth)
	app.Post("/refresh", routes.RefreshToken)

	//AuthorizationRequired Action
	app.Use(routes.AuthorizationRequired())

	//need AuthorizationRequired
	app.Get("/profile", routes.Profile)
	app.Get("/logout", routes.Logout)
	//end AuthorizationRequired

	err := app.Listen(":" + os.Getenv("APP_PORT"))
	if err != nil {
		panic(err)
	}
}
