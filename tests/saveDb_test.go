package utils_test

import (
	"testing"
	"tgsd96/onend/actions"
)

func TestSaveToDB(t *testing.T) {
	var testFileName = "../json/test.csv2019-01-08T16:05:04+05:30.json"
	keys := map[string]string{"name": "Name", "amount": "Amount"}

	_, _, err := actions.SaveCsvToDB(testFileName, keys, "COL")

	if err != nil {
		t.Fatal(err.Error())
	}
}
