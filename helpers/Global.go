package helpers

import (
	"encoding/json"
	"fmt"
	"github.com/teris-io/shortid"
	"io/ioutil"
	"log"
	"net/http"
)

type ResponseObject struct {
	RandomNumber int `json:"random_number"`
}

func InBetween(i int, min int, max int) bool {
	if (i >= min) && (i <= max) {
		return true
	} else {
		return false
	}
}
func MapNumber() (int, error) {
	random, err := CallEndpoint("https://codechallenge.boohma.com/random")
	var responseObject ResponseObject
	json.Unmarshal(random, &responseObject)
	fmt.Printf("API Response as struct %+v\n", responseObject)
	if err != nil {
		return 0, err
	}
	numberNumber := responseObject.RandomNumber
	number := 0
	if InBetween(numberNumber, 80, 100) {
		number = 5
	} else if InBetween(numberNumber, 60, 79) {
		number = 4
	} else if InBetween(numberNumber, 40, 69) {
		number = 3
	} else if InBetween(numberNumber, 20, 39) {
		number = 2
	} else if InBetween(numberNumber, 1, 19) {
		number = 1
	}
	return number, err
}
func CallEndpoint(endpoint string) ([]byte, error) {
	response, _ := http.Get(endpoint)
	responseData, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}
	return responseData, err

}

func Includes(s []int, e int) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

func New() (string, error) {
	return shortid.Generate()
}
