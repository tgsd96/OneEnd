package api

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"tgsd96/onend/actions"
	"tgsd96/onend/app"
	"tgsd96/onend/models"
	"tgsd96/onend/utils"

	"github.com/julienschmidt/httprouter"
)

type FileResponse struct {
	Filename string   `json:"file_name"`
	Headers  []string `json:"keys"`
}

type FileRequest struct {
	Filename string            `json:"file_name"`
	Keys     map[string]string `json:"keys"`
}

func PostFileUpload(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

	// parse the form
	r.ParseMultipartForm(32 << 20)

	file, handler, err := r.FormFile("upload")
	if err != nil {
		fmt.Println(err.Error())
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("File not found"))
		return
	}
	defer file.Close()

	//save the file
	// fmt.Fprintf(w, "%v", handler.Header)
	f, err := os.OpenFile("./uploads/"+handler.Filename, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		fmt.Println(err)
		// w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}
	defer f.Close()
	io.Copy(f, file)

	// save to database
	company := ps.ByName("company")
	var File models.File

	// save the file as json
	csvFileName, headers, err := utils.CsvToMap("./uploads/" + handler.Filename)

	if err != nil {
		log.Fatalf("Error while processing the file: %s", err.Error())
		w.Write([]byte("Not a csv file"))
	}

	File.Filename = csvFileName
	File.Location = "./upload"
	File.Company = company

	app.App.DB.Create(&File)

	// marshal headers to json

	var fileResponse FileResponse

	fileResponse.Filename = csvFileName
	fileResponse.Headers = headers

	headerString, _ := json.Marshal(&fileResponse)

	w.Write(headerString)
}

func PutCsvToDB(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

	// parse the body to gain insights

	decoder := json.NewDecoder(r.Body)

	company := ps.ByName("company")

	// parse the body to file request
	var fileRequest FileRequest
	err := decoder.Decode(&fileRequest)

	if err != nil {
		log.Fatalf(err.Error())
	}

	// call the csvToDb and get results
	s_cnt, e_cnt, err := actions.SaveCsvToDB("./json/"+fileRequest.Filename, fileRequest.Keys, company)

	if err != nil {
		fmt.Println("Error:", err.Error())
		w.Write([]byte("failed"))
		return
	}

	response := map[string]interface{}{
		"success": s_cnt,
		"errors":  e_cnt,
	}

	msg, _ := json.Marshal(&response)

	w.Write(msg)

}
