package app

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"reflect"
	"tgsd96/onend/models"

	"github.com/rs/cors"

	"github.com/go-http-utils/logger"
	"github.com/jinzhu/gorm"
	"github.com/julienschmidt/httprouter"
)

var SECRET_KEY = "ifhsailh839r3hrroq3rj!92439jiofaj"

type app struct {
	Router *httprouter.Router
	DB     *gorm.DB
	Server *http.Server
}

var App *app

func Create_app(port string) *app {

	defaultRouter := httprouter.New()

	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost", "http://localhost:8000"},
		AllowCredentials: true,
		Debug:            true,
		AllowedMethods:   []string{"GET", "POST", "PUT"},
	})

	return &app{
		Router: defaultRouter,
		DB:     &gorm.DB{},
		Server: &http.Server{
			Addr:    port,
			Handler: c.Handler(logger.Handler(defaultRouter, os.Stdout, logger.DevLoggerType)),
		},
	}
}

func (ap *app) RunServer(dbConfig string) {
	var err error
	ap.DB, err = gorm.Open("postgres", dbConfig)
	if err != nil {
		fmt.Println(err.Error())
		panic("Failed")
	}
	// add your models here
	models.AddToModel(&models.Users{})
	models.AddToModel(&models.Ledger{})
	models.AddToModel(&models.Master{})
	models.AddToModel(&models.File{})
	models.AddToModel(&models.MasterView{})
	models.AddToModel(&models.Purchase{})
	models.AddToModel(&models.ErrorPurchases{})
	// automigrate
	for _, model := range models.Models {
		if err := ap.DB.AutoMigrate(model).Error; err != nil {
			fmt.Println(err)
		} else {
			fmt.Println("Auto Migrating the model:", reflect.TypeOf(model).Name(), "...")
		}
	}

	err = ap.Server.ListenAndServe()

	if err != nil {
		log.Fatal(err.Error())
	}

}
