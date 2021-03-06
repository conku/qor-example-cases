package main

import (
	"net/http"

	"github.com/conku/configor"
	"github.com/conku/gorm"
	_ "github.com/conku/gorm/dialects/postgres"
	"github.com/conku/media"
	"github.com/conku/media/oss"
	"github.com/conku/oss/s3"
	"github.com/conku/qor-example-cases/config"
	"github.com/fatih/color"
	appkitlog "github.com/theplant/appkit/log"
	"github.com/theplant/appkit/server"
)

type Order struct {
	gorm.Model
	File oss.OSS
}

// run with dummy data
// MODE=data go run main.go
type Config struct {
	AccessID  string `env:"QOR_AWS_ACCESS_KEY_ID"`
	AccessKey string `env:"QOR_AWS_SECRET_ACCESS_KEY"`
	Region    string `env:"QOR_AWS_REGION"`
	Bucket    string `env:"QOR_AWS_BUCKET"`
}

func main() {
	db := config.DB
	appConfig := Config{}
	err := configor.Load(&appConfig)
	if err != nil {
		panic(err)
	}

	oss.Storage = s3.New(&s3.Config{AccessID: appConfig.AccessID, AccessKey: appConfig.AccessKey, Region: appConfig.Region, Bucket: appConfig.Bucket})

	media.RegisterCallbacks(db)

	db.AutoMigrate(&Order{})

	orderR := config.Admin.AddResource(&Order{})
	_ = orderR
	mux := http.NewServeMux()
	config.Admin.MountTo("/admin", mux)
	color.Green("URL: %v", "http://localhost:3000/admin/orders")
	server.ListenAndServe(server.Config{Addr: ":3000"}, appkitlog.Default(), mux)
}
