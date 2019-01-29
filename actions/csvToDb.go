package actions

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"
	"tgsd96/onend/app"
	"tgsd96/onend/models"
	"time"

	_ "github.com/jinzhu/gorm/dialects/postgres"
)

func SaveCsvToDB(filename string, keys map[string]string, company string) (int, int, error) {

	var successCount = 0
	var errorCount = 0

	// open csv file
	file, _ := os.Open(filename)
	defer file.Close()

	var nameCol = keys["name"]
	var amountCol = keys["amount"]

	// fmt.Println(nameCol)

	var csvToMap map[string][]string

	fileReader, _ := ioutil.ReadAll(file)
	// unmarshal the map
	err := json.Unmarshal(fileReader, &csvToMap)

	if err != nil {
		log.Fatal(err.Error())
		return 0, 0, err
	}

	// fmt.Printf("%v", csvToMap)

	// extract the names and amount
	amountArray := csvToMap[amountCol]
	nameArray := csvToMap[nameCol]

	// fmt.Println(nameArray)
	// fmt.Println(amountArray)

	// traverse the arrays and save to db
	for index := range amountArray {
		// get cust_id for the name
		accountName := strings.ToUpper(nameArray[index])
		accountAmount, _ := strconv.Atoi(amountArray[index])

		// fmt.Printf("The name is %s", accountName)

		var master models.Master

		// if app.App.DB == nil {
		// app.App = app.Create_app(":8080")
		// app.App.DB, _ = gorm.Open("postgres", "host=localhost dbname=one-end sslmode=disable")
		// }

		if err := app.App.DB.Where("name = ?", accountName).First(&master).Error; err != nil {

			// log.Fatal(err)
			// fmt.Println("Name not found")
			errorCount++
			errorResult := &models.ErrorPurchases{
				Name:          accountName,
				Amount:        int64(accountAmount),
				InterfaceCode: strings.ToUpper(company),
			}
			//store error purchase
			app.App.DB.Create(&errorResult)
		} else {

			// test
			// fmt.Println("Master found")

			// create a new entry
			createdAt := time.Now()
			purchase := &models.Purchase{
				CustID:        master.CustID,
				Amount:        int64(accountAmount),
				InterfaceCode: strings.ToUpper(company),
				Date:          &createdAt,
			}

			// make a ledger entry
			ledger := &models.Ledger{
				CustID: master.CustID,
				Amount: int64(accountAmount),
				Type:   "Bill: " + company,
				UserID: 0,
			}

			// save to db
			app.App.DB.Create(&purchase)
			app.App.DB.Create(&ledger)
			successCount++
		}

	}

	return successCount, errorCount, nil

}
