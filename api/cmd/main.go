package main

import (
	"log"

	"github.com/aofiee/diablos/diablosutils"
	"github.com/aofiee/diablos/routes"
	fiber "github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/requestid"
	"github.com/gofiber/fiber/v2/utils"
)

var (
	config diablosutils.Config
)

func Setup() *fiber.App {
	var err error
	config, err = diablosutils.LoadConfig("../")
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	app := fiber.New()
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

	app.Use(cors.New(cors.Config{
		AllowOrigins: config.AllowOrigins,
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
	app.Delete("/logout", routes.Logout)
	//end AuthorizationRequired
	/*
		for _, r := range app.Stack() {
			for _, v := range r {
				if v.Path != "/" && v.Method != "HEAD" {
					fmt.Printf("%v \t %s\n", v.Method, v.Path)
				}
			}
		}
	*/
	return app
}

func main() {
	app := Setup()
	err := app.Listen(":" + config.AppPort)
	if err != nil {
		panic(err)
	}
}
