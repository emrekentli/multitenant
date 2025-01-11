package config

import (
	"github.com/gofiber/template/django/v3"
	"github.com/markbates/pkger"
	"html/template"
	"os"
	"path/filepath"
)

const (
	LiveReloadScript = `<script type="text/javascript" src="https://livejs.com/live.js"></script>`
)

type TemplateConfig struct {
	TemplateEngine *django.Engine
	Path           string `yaml:"path" env-default:"resources/views"`
	Extension      string `yaml:"extension" env-default:".django"`
}
type ViewConfig struct {
	Template TemplateConfig `yaml:"template"`
}

func (v *ViewConfig) Load(path string) {
	path = MakeDir(filepath.Join(path, v.Template.Path))
	engine := django.NewFileSystem(pkger.Dir(path), v.Template.Extension)

	if os.Getenv("ENV") == "development" {
		engine.Reload(true)
		engine.AddFunc(
			"liveReloadTag", func() template.HTML {
				return LiveReloadScript
			},
		)
	} else {
		engine.AddFunc(
			"liveReloadTag", func() template.HTML {
				return ""
			},
		)
	}
	v.Template.TemplateEngine = engine
}
