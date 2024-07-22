package handlers

import (
	"github.com/a-h/templ"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/adaptor"

	"github.com/kylerequez/go-upload-example/src/db"
	"github.com/kylerequez/go-upload-example/src/repositories"
)

func InitHandlers(server *fiber.App) error {
	if err := db.ConnectDB(); err != nil {
		return err
	}

	if err := db.PingDB(); err != nil {
		return err
	}

	ur := repositories.NewUploadRepository(db.DB, "uploads")
	uh := NewUploadHandler(ur)
	uh.InitRoutes(server)

	return nil
}

func Render(
	c fiber.Ctx,
	component templ.Component,
	options ...func(*templ.ComponentHandler),
) error {
	componentHandler := templ.Handler(component)
	for _, o := range options {
		o(componentHandler)
	}
	return adaptor.HTTPHandler(componentHandler)(c)
}
