package config

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
	"github.com/oarkflow/log"
)

type AppConfig struct {
	Mail       Mail
	Hash       Hash
	View       ViewConfig     `yaml:"view"`
	Cache      CacheConfig    `yaml:"cache"`
	Database   DatabaseConfig `yaml:"database"`
	Session    SessionConfig  `yaml:"session"`
	JwtSecrets JwtSecrets     `yaml:"jwt"`
	Server     ServerConfig   `yaml:"server"`
	Log        LogConfig      `yaml:"log"`
	Token      Token          `yaml:"token"`
	ConfigFile string
}

func (cfg *AppConfig) Setup() {
	if err := godotenv.Load(); err != nil {
		log.Warn()
	}

	if err := cleanenv.ReadConfig(cfg.ConfigFile, cfg); err != nil {
		log.Panic()
	}

	cfg.Server.LoadPath()
	cfg.View.Load(cfg.Server.Path)
	cfg.Mail.View = &cfg.View
	cfg.Server.TemplateEngine = cfg.View.Template.TemplateEngine
	cfg.Server.Setup()
	cfg.LoadComponents()
}

func (cfg *AppConfig) PrepareLog() {
	writer := &log.MultiWriter{
		InfoWriter:  cfg.createLogWriter("INFO", cfg.Log.InfoLevel.Path),
		WarnWriter:  cfg.createLogWriter("WARN", cfg.Log.WarnLevel.Path),
		ErrorWriter: cfg.createLogWriter("ERROR", cfg.Log.ErrorLevel.Path),
	}
	if cfg.Log.ConsoleLog.Show {
		writer.ConsoleWriter = &log.IOWriter{Writer: os.Stderr}
		writer.ConsoleLevel = log.PanicLevel
	}
	log.DefaultLogger = log.Logger{
		TimeField:  cfg.Log.TimeField,
		TimeFormat: cfg.Log.TimeFormat,
		Writer:     writer,
	}
}

func (cfg *AppConfig) createLogWriter(level string, path string) *log.FileWriter {
	return &log.FileWriter{
		Filename:     filepath.Join(MakeDir(filepath.Join(cfg.Server.Path, path)), fmt.Sprintf("%s.log", level)),
		EnsureFolder: true,
		TimeFormat:   cfg.Log.TimeFormat,
	}
}

func (cfg *AppConfig) Route404() {
	cfg.Server.Use(func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusNotFound).Render("404", fiber.Map{
			"Title": "Page Not Found",
		})
	})
}

func (cfg *AppConfig) LoadComponents() {
	cfg.LoadStatic()
	cfg.PrepareLog()
	_ = cfg.Database.InitDB()
	_ = cfg.Session.Setup()
	cfg.Cache.Setup()
}

func (cfg *AppConfig) LoadStatic() {
	cfg.Server.Static("/", filepath.Join(cfg.Server.Path, cfg.Server.PublicPath), fiber.Static{
		Compress:      true,
		ByteRange:     true,
		CacheDuration: 24 * time.Hour,
	})
}
