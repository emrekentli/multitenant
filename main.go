package main

import (
	"flag"
	"github.com/emrekentli/multitenant-boilerplate/app"
	"github.com/emrekentli/multitenant-boilerplate/migrate"
	"github.com/emrekentli/multitenant-boilerplate/rest/routes"
	"log"
)

func main() {
	configFile := flag.String("config", "config.yml", "User Config file from user")
	flag.Parse()
	app.Load(*configFile)
	err := migrate.RunMigrations(app.Http.Database.DB)
	if err != nil {
		return
	}
	routes.LoadRoutes(app.Http.Server.App)
	app.Http.Route404()
	log.Fatal(app.Http.Server.ServeWithGraceFullShutdown())

}
