package main

import (
	"github.com/emrekentli/multitenant-boilerplate/app"
	"log"
)

func main() {
	app.Load()
	log.Fatal(app.Http.Server.ServeWithGraceFullShutdown())
}
