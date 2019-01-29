package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"tgsd96/onend/app"
	"tgsd96/onend/models"

	"github.com/julienschmidt/httprouter"
)

type ErrorObj struct {
	Account        models.ErrorPurchases `json:"account"`
	Recommendation []models.Master       `json:"recommendations"`
}

type mergeRequest struct {
	CustID int64  `json:"cust_id"`
	Name   string `json:"name"`
}

func GetErrors(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {

	// get error purchases
	var errorArray []models.ErrorPurchases
	var masterTemp []models.Master
	var errorReturn []ErrorObj

	app.App.DB.Find(&errorArray)

	// recommended sql
	recommendSQL := ` SELECT distinct cust_id, name from masters where levenshtein(name,?)<5`

	for _, errors := range errorArray {

		// get recommended masters
		app.App.DB.Raw(recommendSQL, errors.Name).Scan(&masterTemp)
		tempError := ErrorObj{
			Account:        errors,
			Recommendation: masterTemp,
		}
		errorReturn = append(errorReturn, tempError)
	}

	// unmarshal the erroArray and send

	msg, err := json.Marshal(&errorReturn)

	if err != nil {
		fmt.Println("Error while unmarshaling" + err.Error())
	}

	w.Write(msg)

}

func PostErrors(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {

	newSql := `
		Insert into masters(cust_id, name, interface_code)
		SELECT nextval('seq_cust_id'), A.name, A.interface_code from error_purchases A
		where A.name=?
	`

	mergeSql := `
		Insert into masters(cust_id, name, interface_code)
		SELECT ?, A.name, A.interface_code from error_purchases A
		where A.name = ?
	`

	moveToMasterSql := `
		Insert into Purchases(cust_id, bill_no,created_at, amount, interface_code)
		SELECT A.cust_id, B.bill_no, B.created_at, B.Amount, B.interface_code 
		FROM Masters A, ERROR_PURCHASES B where A.name = B.name
	`

	deleteFromErrorSql := `
		Delete from error_purchases where name in (select name from Masters)
	`
	removeEntrySQl := `
		Delete from error_purchases where name = ?
	`

	// parse the request for data
	decoder := json.NewDecoder(r.Body)

	var request []mergeRequest

	err := decoder.Decode(&request)

	if err != nil {
		fmt.Println(err.Error())
		w.Write([]byte("Errors"))
	}

	for _, obj := range request {

		if obj.CustID == -1 {
			app.App.DB.Exec(newSql, obj.Name)
		} else if obj.CustID == 0 {
			app.App.DB.Exec(removeEntrySQl, obj.Name)
		} else {
			app.App.DB.Exec(mergeSql, obj.CustID, obj.Name)
		}
	}

	// execute insert and remove
	app.App.DB.Exec(moveToMasterSql)
	app.App.DB.Exec(deleteFromErrorSql)

	w.Write([]byte("Successful"))

}
