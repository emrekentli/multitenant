package server

import (
	"app/app/middlewares"
	"app/app/router"
	"app/config"
	"app/src/general/util/exitcode"
	"app/src/general/util/rest"
	"github.com/goccy/go-json"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/log"
	"os"
)

type FiberServer struct {
	*fiber.App
}

func New() *FiberServer {
	server := &FiberServer{
		App: fiber.New(fiber.Config{
			AppName:      config.AppName,
			ErrorHandler: ErrorHandler,
			JSONDecoder:  json.Unmarshal,
			JSONEncoder:  json.Marshal,
			ViewsLayout:  "layouts/main.html",
		}),
	}
	middlewares.RegisterBuiltInMiddlewares(server.App)
	router.RegisterFiberRoutes(server.App)
	err := server.Listen(config.Port)
	if err != nil {
		os.Exit(exitcode.ServerStartError)
	}

	return server
}

func ErrorHandler(c fiber.Ctx, err error) error {
	log.Error("Error: ", err)
	code, message := rest.Error(err)
	return rest.ErrorRes(c, code, message)
}
