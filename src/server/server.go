package server

import (
	"errors"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/logger"
	"github.com/gofiber/fiber/v3/middleware/static"

	"github.com/kylerequez/go-upload-example/src/handlers"
	"github.com/kylerequez/go-upload-example/src/utils"
)

func InitServer() {
	if err := utils.LoadEnvVariables(); err != nil {
		panic(err)
	}

	server := fiber.New(
		fiber.Config{
			AppName: "go-upload-example",
		},
	)

	server.Use("/", static.New("./src/statics"))
	server.Use("/uploads", static.New("./uploads"))

	server.Use(logger.New(
		logger.Config{
			Format: "[${ip}]:${port} ${status} - ${method} ${path}\n",
		},
	))

	if err := handlers.InitHandlers(server); err != nil {
		panic(err)
	}

	port := utils.GetEnv("PORT")
	if port == "" {
		panic(errors.New("app port is empty"))
	}
	server.Listen(":" + port)
}
