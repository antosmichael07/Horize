package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

func getCarsHandler(w http.ResponseWriter, r *http.Request) {

}

func getCarsData(login string) {
	if !checkFileExistance(fmt.Sprintf("./save_data/account_cars/%s.json", login)) {
		cars_json, cars_json_err := json.Marshal(AccountCars{Cars: []string{}, Modifications: []string{}, Paint: []string{}})
		if cars_json_err != nil {
			serverLog("account_cars.go:18", fmt.Sprintf("failed to marshal to file \"account_cars/%s.json\", the error \"%s\"", login, cars_json_err))
			return
		}
		cars_json_file_err := os.WriteFile(fmt.Sprintf("./save_data/account_cars/%s.json", login), cars_json, 0644)
		if cars_json_file_err != nil {
			serverLog("account_cars.go:23", fmt.Sprintf("failed to create \"account_cars/%s.json\", the error \"%s\"", login, cars_json_file_err))
			return
		}
	}
}

func addCarToAccount(login string) {

}

func removeCarFromAccount(login string) {

}
