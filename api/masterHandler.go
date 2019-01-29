package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"tgsd96/onend/actions"
	"tgsd96/onend/app"
	"tgsd96/onend/models"
	"tgsd96/onend/utils"

	"github.com/julienschmidt/httprouter"
)

const MASTER_VIEW = "master_view"

type MasterWithBalance struct {
	CustID  int64  `json:"cust_id"`
	Balance int64  `json:"balance"`
	Name    string `json:"name"`
	Area    string `json:"area"`
}

// get masters list

func GetMastersList(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {

	// sql for getting balance
	getSQL := ` select B.*, A.name,A.area from master_view A Inner join(
		select A.cust_id, sum(B.amount) as balance from master_view A left outer join ledgers B ON A.cust_id = B.cust_id group by A.cust_id) B
		on A.cust_id = B.cust_id order by A.name`

	// define result interface
	var result []MasterWithBalance

	// execute query

	app.App.DB.Raw(getSQL).Scan(&result)

	// marshal the result and return

	msg, _ := json.Marshal(&result)

	w.Write(msg)
}

type LedgerWithName struct {
	models.Ledger
	Name string `json:"name"`
}

type LedgerDetail struct {
	Main     MasterWithBalance `json:"main, omitempty"`
	Accounts []models.Master   `json:"accounts"`
	Ledgers  []LedgerWithName  `json:"ledgers"`
}

// Get Master Details

func GetMasterDetails(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

	// get the cust id
	cust_id := ps.ByName("cust_id")

	var allMerged []models.Master

	var allLedgers []LedgerWithName

	// get all the masters for this cust_id

	app.App.DB.Order("created_at desc").Where("cust_id = ?", cust_id).Find(&allMerged)

	getLedgerWithNameSQL := `select A.name, B.* from master_view A, ledgers B where A.cust_id = B.cust_id and B.cust_id = ? order by B.date desc`

	app.App.DB.Raw(getLedgerWithNameSQL, cust_id).Scan(&allLedgers)
	// getSQL := ` select B.*, A.name,A.area from master_view A Inner join(
	// 	select A.cust_id, sum(B.amount) as balance from master_view A left outer join ledgers B ON A.cust_id = B.cust_id group by A.cust_id) B
	// 	on A.cust_id = B.cust_id and A.cust_id = ?`

	// var mainAccount MasterWithBalance

	// app.App.DB.Raw(getSQL, cust_id).Scan(&mainAccount)

	result := LedgerDetail{
		Accounts: allMerged,
		Ledgers:  allLedgers,
	}

	msg, _ := json.Marshal(result)

	w.Write(msg)

}

// Merge Request struct

type MergeTwoIDs struct {
	OldID int64 `json:"old_id"`
	NewID int64 `json:"new_id"`
}

func PostMergeCustID(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {

	// parse the body of json

	decoder := json.NewDecoder(r.Body)

	var newMergeRequest MergeTwoIDs

	// decode
	err := decoder.Decode(&newMergeRequest)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Wrong format"))
		return
	}

	// update the cust id
	var master models.Master
	app.App.DB.Model(&master).Where("cust_id = ?", newMergeRequest.OldID).Update("cust_id", newMergeRequest.NewID)

	w.Write([]byte("Sucessful"))
}

// Update Balance to mobile server

const serverURL = "http://localhost:8000"

func GetUpdateMobileServer(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {

	// get the masters
	getSQL := ` select B.*, A.name,A.area from master_view A Inner join(
		select A.cust_id, sum(B.amount) as balance from master_view A left outer join ledgers B ON A.cust_id = B.cust_id group by A.cust_id) B
		on A.cust_id = B.cust_id`

	// define result interface
	var result []MasterWithBalance

	// execute query

	app.App.DB.Raw(getSQL).Scan(&result)

	msg, _ := json.Marshal(result)

	// request to mobile server
	_, err := http.Post(serverURL+"/master", "application/json", bytes.NewBuffer(msg))

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	w.Write([]byte("successful"))
}

// CRUD

func PostCreateMaster(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {

	// parse the master from req
	var master models.Master

	// decode using util
	err := utils.BodyToInteface(r.Body, &master)

	if err != nil {
		fmt.Println(err)
		return
	}

	if master.CustID == 0 {
		master.CustID = actions.GenereateCustID()

	}
	app.App.DB.Create(&master)

}

func PutUpdateMaster(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {

	var master models.Master

	err := utils.BodyToInteface(r.Body, &master)

	if err != nil {
		fmt.Println(err)
		return
	}

	app.App.DB.Model(&master).Updates(master)

	w.Write([]byte("Successful"))
}

func DeleteMaster(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {

	var master models.Master

	err := utils.BodyToInteface(r.Body, &master)

	if err != nil {
		fmt.Println(err)
		return
	}

	app.App.DB.Delete(&master)

	w.Write([]byte("Successful"))

}

func SearchMaster(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {

	// get search query

	searchName := strings.ToUpper(r.URL.Query().Get("search"))

	// search using like
	var masters []models.Master
	app.App.DB.Where("name like ?", "%"+searchName+"%").Find(&masters)

	msg, _ := json.Marshal(&masters)

	w.Write(msg)

}
