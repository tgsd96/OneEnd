package main

import (
	"fmt"
	"net/http"
	"tgsd96/onend/api"
	"tgsd96/onend/app"

	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/julienschmidt/httprouter"
)

func rootHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	fmt.Fprint(w, "Hi")
}

func main() {
	app.App = app.Create_app(":8080")
	app.App.Router.GET("/", rootHandler)
	app.App.Router.POST("/api/fileupload/:company", api.PostFileUpload)
	app.App.Router.PUT("/api/savefile/:company", api.PutCsvToDB)

	// get the query
	app.App.Router.GET("/api/download", api.GETListForDates)
	app.App.Router.GET("/api/file", api.DownloadListFile)

	// error apis
	app.App.Router.GET("/api/errors", api.GetErrors)
	app.App.Router.POST("/api/errors", api.PostErrors)

	// masters api
	app.App.Router.GET("/api/master", api.GetMastersList)
	// get master detail
	app.App.Router.GET("/api/master/:cust_id", api.GetMasterDetails)
	app.App.Router.POST("/api/master", api.PostCreateMaster)
	app.App.Router.PUT("/api/master", api.PutUpdateMaster)

	// search masters
	app.App.Router.GET("/api/search", api.SearchMaster)

	// merge cusids
	app.App.Router.POST("/api/merge", api.PostMergeCustID)

	// run sync to server
	app.App.Router.GET("/api/sync", api.GetUpdateMobileServer)

	// update ledgers
	app.App.Router.GET("/api/syncledgers", api.GetLedgerFromMobile)

	// api to create a newLedger template
	app.App.Router.GET("/api/ledger", api.GetCreateLedgerAndID)

	app.App.Router.POST("/api/ledger", api.PostCreateLedger)
	app.App.Router.PUT("/api/ledger", api.PutUpdateLedger)
	app.App.Router.GET("/api/ledgers/:id", api.GetLedgerFromID)
	app.App.Router.POST("/api/mledgers", api.PostCreateLedgers)
	app.App.Router.PUT("/api/ledgers", api.PostUpdateOrCreateLedger)
	app.App.Router.POST("/api/ledgers", api.PostViewAllLedgers)
	app.App.RunServer("host=localhost dbname=one-end sslmode=disable")
}
