package actions

import (
	"tgsd96/onend/app"
	"tgsd96/onend/models"
)

func CreateLegder(ledger *models.Ledger) {

	var purchase models.Purchase
	if ledger.Amount > 0 && ledger.CustID != 0 {
		purchase.CustID = ledger.CustID
		purchase.Date = ledger.Date
		purchase.Amount = ledger.Amount
		purchase.InterfaceCode = ledger.Type
	}

	if app.App.DB.NewRecord(&ledger) {
		app.App.DB.Create(&ledger)

	}

}
