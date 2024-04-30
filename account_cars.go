package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

func getCarsHandler(w http.ResponseWriter, r *http.Request) {

}

func getCarsData(login string) (cars AccountCars, err error) {
	if !checkFileExistance(fmt.Sprintf("./save_data/account_cars/%s.json", login)) {
		cars_json, cars_json_err := json.Marshal(AccountCars{Cars: []string{}, Modifications: []string{}, Paint: []string{}})
		if cars_json_err != nil {
			serverLog("account_cars.go:18", fmt.Sprintf("failed to marshal to file \"account_cars/%s.json\", the error \"%s\"", login, cars_json_err))
			return AccountCars{Cars: []string{}, Modifications: []string{}, Paint: []string{}}, cars_json_err
		}

		cars_json_file_err := os.WriteFile(fmt.Sprintf("./save_data/account_cars/%s.json", login), cars_json, 0644)
		if cars_json_file_err != nil {
			serverLog("account_cars.go:24", fmt.Sprintf("failed to create \"account_cars/%s.json\", the error \"%s\"", login, cars_json_file_err))
			return AccountCars{Cars: []string{}, Modifications: []string{}, Paint: []string{}}, cars_json_file_err
		}

		serverLog("account_cars.go:28", fmt.Sprintf("created \"account_cars/%s.json\"", login))
		return AccountCars{Cars: []string{}, Modifications: []string{}, Paint: []string{}}, nil
	}

	cars_json, cars_json_err := fileToByte(fmt.Sprintf("./save_data/account_cars/%s.json", login))
	if cars_json_err != nil {
		serverLog("account_cars.go:34", fmt.Sprintf("failed to read \"account_cars/%s.json\", the error \"%s\"", login, cars_json_err))
		return AccountCars{Cars: []string{}, Modifications: []string{}, Paint: []string{}}, cars_json_err
	}

	cars = AccountCars{}
	cars_json_err = json.Unmarshal(cars_json, &cars)
	if cars_json_err != nil {
		serverLog("account_cars.go:41", fmt.Sprintf("failed to unmarshal \"account_cars/%s.json\", the error \"%s\"", login, cars_json_err))
		return AccountCars{Cars: []string{}, Modifications: []string{}, Paint: []string{}}, cars_json_err
	}

	serverLog("account_cars.go:45", fmt.Sprintf("successfully read \"account_cars/%s.json\"", login))
	return cars, nil
}

func addCarToAccount(login string, car string) (err error) {
	cars, cars_err := getCarsData(login)
	if cars_err != nil {
		serverLog("account_cars.go:52", fmt.Sprintf("failed to get cars data, the error \"%s\"", cars_err))
		return cars_err
	}

	cars.Cars = append(cars.Cars, car)
	cars.Modifications = append(cars.Modifications, "")
	cars.Paint = append(cars.Paint, "")

	cars_json, cars_json_err := json.Marshal(cars)
	if cars_json_err != nil {
		serverLog("account_cars.go:62", fmt.Sprintf("failed to marshal to file \"account_cars/%s.json\", the error \"%s\"", login, cars_json_err))
		return cars_json_err
	}
	write_err := os.WriteFile(fmt.Sprintf("./save_data/account_cars/%s.json", login), cars_json, 0644)
	if write_err != nil {
		serverLog("account_cars.go:67", fmt.Sprintf("failed to write to file \"account_cars/%s.json\", the error \"%s\"", login, write_err))
		return write_err
	}

	return nil
}

func removeCarFromAccount(login string, car string) {

}
