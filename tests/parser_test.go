package utils_test

import (
	"testing"
	"tgsd96/onend/utils"
)

func TestCsvToMap(t *testing.T) {
	var Filename = "colgate.csv"
	FileOuput, headers, err := utils.CsvToMap(Filename)

	if err != nil {
		t.Fatal(err.Error())
	}
	t.Log(FileOuput)
	t.Log(headers)

}
