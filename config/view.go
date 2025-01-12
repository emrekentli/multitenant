package config

import (
	"github.com/gofiber/template/django/v3"
	"github.com/markbates/pkger"
	"os"
	"path/filepath"
)

const (
	LiveReloadScript = `<script type="text/javascript" src="https://livejs.com/live.js"></script>`
)

type TemplateConfig struct {
	TemplateEngine *django.Engine
	Path           string `env:"TEMPLATE_PATH" env-default:"resources/views"`
	Extension      string `env:"TEMPLATE_EXTENSION" env-default:".django"`
}
type ViewConfig struct {
	Template TemplateConfig
}

func (v *ViewConfig) Load() {
	v.Template.Path = os.Getenv("TEMPLATE_PATH")
	v.Template.Extension = os.Getenv("TEMPLATE_EXTENSION")
	path, _ := os.Getwd()
	path = MakeDir(filepath.Join(path, v.Template.Path))

	engine := django.NewFileSystem(pkger.Dir(path), v.Template.Extension)

	if os.Getenv("ENV") == "development" {
		engine.Reload(true)
		engine.AddFunc(
			"liveReloadTag", func() string {
				return LiveReloadScript
			},
		)
	} else {
		engine.AddFunc(
			"liveReloadTag", func() string {
				return ""
			},
		)
	}

	v.Template.TemplateEngine = engine
}
