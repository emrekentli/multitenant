package main

import (
	"flag"
	"github.com/emrekentli/multitenant-boilerplate/internal/app"
	"github.com/emrekentli/multitenant-boilerplate/internal/rest/routes"
	"github.com/emrekentli/multitenant-boilerplate/migrations"
	"log"
)

func main() {
	configFile := flag.String("config", "config.yml", "User Config file from user")
	flag.Parse()
	app.Load(*configFile)
	err := migrations.RunMigrations(app.Http.Database.DB)
	if err != nil {
		return
	}
	routes.LoadRoutes(app.Http.Server.App)
	app.Http.Route404()
	log.Fatal(app.Http.Server.ServeWithGraceFullShutdown())

}
