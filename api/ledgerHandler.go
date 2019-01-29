package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"tgsd96/onend/app"
	"tgsd96/onend/models"
	"tgsd96/onend/utils"
	"time"

	"github.com/julienschmidt/httprouter"
)

// Sync ledgers from api

// const serverURL = "http://localhost:8000"

func GetLedgerFromMobile(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {

	// make request

	response, err := http.Get(serverURL + "/download")

	if err != nil {
		log.Fatal(err.Error())
		return
	} else {
		decoder := json.NewDecoder(response.Body)

		// decode the response and add to ledger
		var ledgerResult []models.Ledger

		err := decoder.Decode(&ledgerResult)

		if err != nil {
			log.Fatal(err)
			return
		}

		for _, ledger := range ledgerResult {

			fmt.Println(ledger.CustID)
			var tempLedger models.Ledger

			tempLedger.Amount = -ledger.Amount
			tempLedger.Comments = ledger.Comments
			tempLedger.CreatedAt = ledger.CreatedAt
			tempLedger.CustID = ledger.CustID
			tempLedger.Date = ledger.Date
			tempLedger.DeletedAt = ledger.DeletedAt
			tempLedger.Type = ledger.Type
			tempLedger.UpdatedAt = ledger.UpdatedAt
			tempLedger.UserID = ledger.UserID

			fmt.Println(tempLedger)
			// save to db
			if err := app.App.DB.Create(&tempLedger).Error; err != nil {
				fmt.Println(err)
			}
		}
	}

	w.Write([]byte("Succesful"))
}

// create new ledger

func PostCreateLedger(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {

	// parse the body for ledger
	var tempLedger models.Ledger

	err := utils.BodyToInteface(r.Body, &tempLedger)

	if err != nil {
		fmt.Println(err.Error())
		return
	}

	// create ledger
	app.App.DB.Create(&tempLedger)

	w.Write([]byte("Successful"))
}

func PostCreateLedgers(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {

	var ledgerArray []models.Ledger
	err := utils.BodyToInteface(r.Body, &ledgerArray)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	for _, ledger := range ledgerArray {
		if ledger.Amount > 0 && ledger.CustID != 0 {
			var purchase models.Purchase
			purchase.CustID = ledger.CustID
			purchase.Date = ledger.Date
			purchase.Amount = ledger.Amount
			purchase.InterfaceCode = ledger.Type
			app.App.DB.Create(&purchase)
		}
		app.App.DB.Create(&ledger)
	}
	w.Write([]byte("Successful"))
}

func PutUpdateLedger(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {

	var tempLedger models.Ledger

	err := utils.BodyToInteface(r.Body, &tempLedger)

	if err != nil {
		fmt.Println(err.Error())
		return
	}

	app.App.DB.Model(&tempLedger).Updates(tempLedger)

	w.Write([]byte("Successful"))
}

func PostUpdateOrCreateLedger(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var ledgerArray []models.Ledger

	err := utils.BodyToInteface(r.Body, &ledgerArray)

	if err != nil {
		fmt.Println(err.Error())
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Wrong format"))
		return
	}

	for _, ledger := range ledgerArray {
		if app.App.DB.NewRecord(&ledger) {
			app.App.DB.Create(&ledger)
		} else {
			app.App.DB.Model(&ledger).Updates(ledger)
		}
	}

	w.Write([]byte("Successful"))
}

func GetLedgerFromID(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var tempLedger models.Ledger

	ID, _ := strconv.Atoi(ps.ByName("id"))

	// query db for the ledger
	app.App.DB.Where("id = ?", ID).First(&tempLedger)

	msg, _ := json.Marshal(tempLedger)

	w.Write(msg)
}

type LedgerRequest struct {
	StartDate *time.Time `json:"start_date"`
	EndDate   *time.Time `json:"end_date"`
	CustID    int64      `json:"cust_id"`
}

func PostViewAllLedgers(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {

	// check if dates are given
	var lRequest LedgerRequest

	err := utils.BodyToInteface(r.Body, &lRequest)

	if err != nil {
		fmt.Println(err.Error())
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Unknown format"))
	}

	fmt.Printf("%v \n", lRequest.CustID)

	var allLedgers []models.Ledger

	if lRequest.CustID == 0 {
		app.App.DB.Where("date>=? and date<=?", lRequest.StartDate, lRequest.EndDate).Order("date desc").Find(&allLedgers)
	} else {
		app.App.DB.Where("date>=? and date<=? and cust_id=?", lRequest.StartDate, lRequest.EndDate, lRequest.CustID).Order("date desc").Find(&allLedgers)
	}

	msg, _ := json.Marshal(allLedgers)
	w.Write(msg)

}

func GetCreateLedgerAndID(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var newLedger models.Ledger

	app.App.DB.Create(&newLedger)

	msg, _ := json.Marshal(newLedger)

	w.Write(msg)
}
