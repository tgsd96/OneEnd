package utils

import (
	"encoding/json"
	"io"
)

func BodyToInteface(body io.ReadCloser, v interface{}) error {

	// b, _ := ioutil.ReadAll(body)
	// fmt.Println(string(b))

	decoder := json.NewDecoder(body)

	err := decoder.Decode(v)

	if err != nil {
		return err
	}
	return nil

}
