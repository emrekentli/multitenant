package config

import (
	"encoding/json"
	"fmt"
	"github.com/emrekentli/multitenant-boilerplate/app/router"
	"github.com/emrekentli/multitenant-boilerplate/src/util/rest"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/log"
	"github.com/gofiber/template/django/v3"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
)

type ServerConfig struct {
	*fiber.App
	TemplateEngine *django.Engine
	Name           string
	Env            string
	Port           string
	PublicPath     string
	StoragePath    string
	LogPath        string
}

func (s *ServerConfig) Setup() error {
	path, _ := os.Getwd()
	s.StoragePath = MakeDir(filepath.Join(path, s.StoragePath))
	s.LogPath = MakeDir(filepath.Join(path, s.LogPath))
	s.Name = os.Getenv("APP_NAME")
	s.Env = os.Getenv("ENV")
	s.Port = os.Getenv("PORT")
	s.PublicPath = os.Getenv("PUBLIC_PATH")
	s.StoragePath = os.Getenv("STORAGE_PATH")
	s.LogPath = os.Getenv("LOG_PATH")
	s.App = fiber.New(fiber.Config{
		Views:             s.TemplateEngine,
		ServerHeader:      s.Name,
		ReduceMemoryUsage: true,
		AppName:           s.Name,
		ErrorHandler:      ErrorHandler,
		JSONDecoder:       json.Unmarshal,
		JSONEncoder:       json.Marshal,
	})
	router.LoadRoutes(s.App)
	return nil
}
func (s *ServerConfig) ServeWithGraceFullShutdown(addr ...string) error {
	go func() {
		if err := s.Listen(s.Port); err != nil {
			log.Fatal(fmt.Sprintf("Error starting server: %v", err))
		}
	}()

	c := make(chan os.Signal, 1) // Create channel to signify a signal being sent
	signal.Notify(c,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGABRT,
		syscall.SIGQUIT,
	) // When an interrupt is sent, notify the channel
	<-c // This blocks the main thread until an interrupt is received
	fmt.Println("I'm shutting down")
	return s.Shutdown()
}
func ErrorHandler(c fiber.Ctx, err error) error {
	log.Error(fmt.Sprintf("Error: %v", err))
	code, message := rest.Error(err)
	return rest.ErrorRes(c, code, message)
}

func MakeDir(path string) string {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		if mkdirErr := os.MkdirAll(path, os.ModePerm); mkdirErr != nil {
			log.Error(fmt.Sprintf("Error creating directory: %v", mkdirErr))
			return ""
		}
	}
	return path
}
