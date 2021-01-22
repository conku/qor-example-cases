package main

import (
	"net/http"

	"github.com/conku/gorm"
	_ "github.com/conku/gorm/dialects/postgres"
	"github.com/conku/qor-example-cases/config"
	"github.com/fatih/color"
	appkitlog "github.com/theplant/appkit/log"
	"github.com/theplant/appkit/server"
)

// Your model definition

type Order struct {
	gorm.Model
	Num   string
	State string
	Price float64
}

func main() {
	var (
		DB    = config.DB
		Admin = config.Admin
	)

	// Your logic

	mux := http.NewServeMux()
	Admin.MountTo("/admin", mux)
	color.Green("URL: %v", "http://localhost:3000/admin/orders")
	server.ListenAndServe(server.Config{Addr: ":3000"}, appkitlog.Default(), mux)
}
