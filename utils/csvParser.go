package utils

import (
	"bufio"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"time"
)

func CsvToMap(filePath string) (string, []string, error) {
	var csvMap map[string][]string
	var headerArray []string
	csvMap = make(map[string][]string)

	// read the header of csv
	csvFile, err := os.Open(filePath)
	if err != nil {
		fmt.Println("Unable to find the file")
		return "", nil, err
	}
	reader := csv.NewReader(bufio.NewReader(csvFile))

	// read the first line of csv
	line, _ := reader.Read()
	for _, head := range line {
		headerArray = append(headerArray, head)
	}

	// read further lines
	for {
		line, error := reader.Read()
		if error == io.EOF {
			break
		} else if error != nil {
			log.Print("her" + error.Error())
			break
		}
		for index, value := range line {
			csvMap[headerArray[index]] = append(csvMap[headerArray[index]], value)
		}
	}

	mapToString, _ := json.Marshal(&csvMap)

	// save to json
	// create file name
	fileNameS := strings.Split(filePath, string(os.PathSeparator))
	fileName := "./json/" + (fileNameS[len(fileNameS)-1]) + string(time.Now().Format(time.RFC3339)) + ".json"

	file, err := os.OpenFile(fileName, os.O_CREATE|os.O_WRONLY, 0666)
	file.Write(mapToString)
	defer file.Close()

	// return the file name and header array
	return fileName, headerArray, nil

}
