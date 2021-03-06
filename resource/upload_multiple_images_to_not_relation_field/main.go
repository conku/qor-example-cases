package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/conku/admin"
	"github.com/conku/configor"
	"github.com/conku/gorm"
	_ "github.com/conku/gorm/dialects/postgres"
	"github.com/conku/media"
	"github.com/conku/media/oss"
	"github.com/conku/oss/s3"
	"github.com/conku/qor"
	"github.com/conku/qor-example-cases/config"
	"github.com/fatih/color"
	appkitlog "github.com/theplant/appkit/log"
	"github.com/theplant/appkit/server"
)

type Order struct {
	gorm.Model
	Name   string
	Images []*Image `gorm:"-"`
}

type Image struct {
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
	config := Config{}
	err := configor.Load(&config)
	if err != nil {
		panic(err)
	}

	db.LogMode(true)
	media.RegisterCallbacks(db)

	if os.Getenv("DATA") != "" {
		db.DropTable(&Order{})
	}

	db.AutoMigrate(&Order{})

	if os.Getenv("DATA") != "" {
		order := &Order{}
		err = db.Create(order).Error
		if err != nil {
			panic(err)
		}
	}

	oss.Storage = s3.New(&s3.Config{AccessID: config.AccessID, AccessKey: config.AccessKey, Region: config.Region, Bucket: config.Bucket})

	adm := admin.New(&admin.AdminConfig{DB: db})
	orderR := adm.AddResource(&Order{})

	orderR.Meta(&admin.Meta{Name: "Images", Type: "collection_edit"})

	oldSaveHandler := orderR.SaveHandler
	orderR.SaveHandler = func(v interface{}, ctx *qor.Context) (err error) {
		ord := v.(*Order)
		fmt.Println("len(ord.Images) = ", len(ord.Images))
		return oldSaveHandler(v, ctx)
	}

	mux := http.NewServeMux()
	adm.MountTo("/admin", mux)
	color.Green("URL: %v", "http://localhost:3000/admin/orders")
	server.ListenAndServe(server.Config{Addr: ":3000"}, appkitlog.Default(), mux)
}
