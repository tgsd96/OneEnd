package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"tgsd96/onend/actions"
	"time"

	"github.com/jszwec/csvutil"

	"github.com/julienschmidt/httprouter"
)

type QueryResult struct {
	CustID int64  `json:"cust_id"`
	Col    string `json:"col"`
	Marg   string `json:"marg"`
	Total  string `json:"total"`
	Name   string `json:"name"`
	Area   string `json:"area"`
}

// api to get the lists
func GETListForDates(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {

	// get the dates from query
	var startDate = "1950-01-01"
	var endDate = "2050-01-01"

	if sDate := r.URL.Query().Get("startDate"); sDate != "" {
		startDate = sDate
	}

	if eDate := r.URL.Query().Get("endDate"); eDate != "" {
		endDate = eDate
	}

	// query for result

	msg := actions.ExecuteListSQL(startDate, endDate)

	w.Write(msg)
}

func DownloadListFile(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {

	// get the dates
	var startDate = "1950-01-01"
	var endDate = "2050-01-01"

	if sDate := r.URL.Query().Get("startDate"); sDate != "" {
		startDate = sDate
	}

	if eDate := r.URL.Query().Get("endDate"); eDate != "" {
		endDate = eDate
	}

	// get json string from function
	msg := actions.ExecuteListSQL(startDate, endDate)

	var queryResult []QueryResult

	err := json.Unmarshal(msg, &queryResult)

	if err != nil {
		fmt.Print(err)
	}

	// marshal to csv
	b, err := csvutil.Marshal(queryResult)

	if err != nil {
		fmt.Print(err)
	}
	// open a file
	filename := "list_" + time.Now().Format(time.RFC3339) + ".csv"

	file, _ := os.OpenFile("./csv/"+filename, os.O_CREATE|os.O_WRONLY, 0666)

	// write csv to file
	file.Write(b)

	w.Header().Set("Content-Disposition", "attachment; filename="+filename)
	http.ServeFile(w, r, "./csv/"+filename)

}
