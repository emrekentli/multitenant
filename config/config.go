package config

import (
	"fmt"
	"github.com/emrekentli/multitenant-boilerplate/config/database"
	_ "github.com/joho/godotenv/autoload"
	"github.com/oarkflow/log"
	"os"
	"path/filepath"
)

type AppConfig struct {
	Mail       Mail
	Hash       Hash
	View       ViewConfig
	Cache      CacheConfig
	Database   database.DatabaseConfig
	Session    SessionConfig
	JwtSecrets JwtSecrets
	Server     ServerConfig
	Log        LogConfig
	Token      Token
}

func (cfg *AppConfig) Setup() {
	cfg.View.Load()
	cfg.Server.TemplateEngine = cfg.View.Template.TemplateEngine
	err := cfg.Server.Setup()
	if err != nil {
		return
	}
	cfg.Mail.View = &cfg.View
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
	pathWd, _ := os.Getwd()
	return &log.FileWriter{
		Filename:     filepath.Join(MakeDir(filepath.Join(pathWd, path)), fmt.Sprintf("%s.log", level)),
		EnsureFolder: true,
		TimeFormat:   cfg.Log.TimeFormat,
	}
}

func (cfg *AppConfig) LoadComponents() {
	cfg.PrepareLog()
	_ = cfg.Database.InitDB()
	_ = cfg.Session.Setup()
	cfg.Cache.Setup()
}
