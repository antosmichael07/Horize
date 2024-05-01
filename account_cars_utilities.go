package main

import (
	"encoding/json"
	"fmt"
	"os"
)

func getCarsData(login string) (cars AccountCars, err error) {
	if !checkFileExistance(fmt.Sprintf("./save_data/account_cars/%s.json", login)) {
		cars_json, cars_json_err := json.Marshal(AccountCars{Cars: []string{}, Modifications: []string{}, Paint: []string{}})
		if cars_json_err != nil {
			serverLog("account_cars_utilities.go:13", fmt.Sprintf("failed to marshal to file \"account_cars/%s.json\", the error \"%s\"", login, cars_json_err))
			return AccountCars{Cars: []string{}, Modifications: []string{}, Paint: []string{}}, cars_json_err
		}

		cars_json_file_err := os.WriteFile(fmt.Sprintf("./save_data/account_cars/%s.json", login), cars_json, 0644)
		if cars_json_file_err != nil {
			serverLog("account_cars_utilities.go:19", fmt.Sprintf("failed to create \"account_cars/%s.json\", the error \"%s\"", login, cars_json_file_err))
			return AccountCars{Cars: []string{}, Modifications: []string{}, Paint: []string{}}, cars_json_file_err
		}

		serverLog("account_cars_utilities.go:23", fmt.Sprintf("created \"account_cars/%s.json\"", login))
		return AccountCars{Cars: []string{}, Modifications: []string{}, Paint: []string{}}, nil
	}

	cars_json, cars_json_err := fileToByte(fmt.Sprintf("./save_data/account_cars/%s.json", login))
	if cars_json_err != nil {
		serverLog("account_cars_utilities.go:29", fmt.Sprintf("failed to read \"account_cars/%s.json\", the error \"%s\"", login, cars_json_err))
		return AccountCars{Cars: []string{}, Modifications: []string{}, Paint: []string{}}, cars_json_err
	}

	cars = AccountCars{}
	cars_json_err = json.Unmarshal(cars_json, &cars)
	if cars_json_err != nil {
		serverLog("account_cars_utilities.go:36", fmt.Sprintf("failed to unmarshal \"account_cars/%s.json\", the error \"%s\"", login, cars_json_err))
		return AccountCars{Cars: []string{}, Modifications: []string{}, Paint: []string{}}, cars_json_err
	}

	serverLog("account_cars_utilities.go:40", fmt.Sprintf("successfully read \"account_cars/%s.json\"", login))
	return cars, nil
}

func addCarToAccount(login string, car string) (err error) {
	cars, cars_err := getCarsData(login)
	if cars_err != nil {
		serverLog("account_cars_utilities.go:47", fmt.Sprintf("failed to get cars data, the error \"%s\"", cars_err))
		return cars_err
	}

	cars.Cars = append(cars.Cars, car)
	cars.Modifications = append(cars.Modifications, "")
	cars.Paint = append(cars.Paint, "")

	cars_json, cars_json_err := json.Marshal(cars)
	if cars_json_err != nil {
		serverLog("account_cars_utilities.go:57", fmt.Sprintf("failed to marshal to file \"account_cars/%s.json\", the error \"%s\"", login, cars_json_err))
		return cars_json_err
	}
	write_err := os.WriteFile(fmt.Sprintf("./save_data/account_cars/%s.json", login), cars_json, 0644)
	if write_err != nil {
		serverLog("account_cars_utilities.go:62", fmt.Sprintf("failed to write to file \"account_cars/%s.json\", the error \"%s\"", login, write_err))
		return write_err
	}

	return nil
}

func removeCarFromAccount(login string, car string) (err error) {
	cars, cars_err := getCarsData(login)
	if cars_err != nil {
		serverLog("account_cars_utilities.go:72", fmt.Sprintf("failed to get cars data, the error \"%s\"", cars_err))
		return cars_err
	}

	for i, v := range cars.Cars {
		if v == car {
			cars.Cars = append(cars.Cars[:i], cars.Cars[i+1:]...)
			cars.Modifications = append(cars.Modifications[:i], cars.Modifications[i+1:]...)
			cars.Paint = append(cars.Paint[:i], cars.Paint[i+1:]...)
			break
		}
	}

	cars_json, cars_json_err := json.Marshal(cars)
	if cars_json_err != nil {
		serverLog("account_cars_utilities.go:87", fmt.Sprintf("failed to marshal to file \"account_cars/%s.json\", the error \"%s\"", login, cars_json_err))
		return cars_json_err
	}
	write_err := os.WriteFile(fmt.Sprintf("./save_data/account_cars/%s.json", login), cars_json, 0644)
	if write_err != nil {
		serverLog("account_cars_utilities.go:92", fmt.Sprintf("failed to write to file \"account_cars/%s.json\", the error \"%s\"", login, write_err))
		return write_err
	}

	return nil
}
