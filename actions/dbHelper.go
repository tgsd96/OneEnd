package actions

import (
	"tgsd96/onend/app"
)

type CustID struct {
	Nextval int64 `json:"nextval"`
}

func GenereateCustID() int64 {

	// execute db

	var nextVal CustID

	app.App.DB.Raw("Select nextval('seq_cust_id')").Scan(&nextVal)

	return nextVal.Nextval
}
