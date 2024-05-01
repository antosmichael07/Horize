package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
)

func getCarsData(login string) (cars AccountCars, err error) {
	if !checkFileExistance(fmt.Sprintf("./save_data/account_cars/%s.json", login)) {
		cars_json, cars_json_err := json.Marshal(AccountCars{Cars: []string{}, Modifications: [][]string{}, Paint: [][]string{}})
		if cars_json_err != nil {
			serverLog("account_cars_utilities.go:14", fmt.Sprintf("failed to marshal to file \"account_cars/%s.json\", the error \"%s\"", login, cars_json_err))
			return AccountCars{Cars: []string{}, Modifications: [][]string{}, Paint: [][]string{}}, cars_json_err
		}

		cars_json_file_err := os.WriteFile(fmt.Sprintf("./save_data/account_cars/%s.json", login), cars_json, 0644)
		if cars_json_file_err != nil {
			serverLog("account_cars_utilities.go:20", fmt.Sprintf("failed to create \"account_cars/%s.json\", the error \"%s\"", login, cars_json_file_err))
			return AccountCars{Cars: []string{}, Modifications: [][]string{}, Paint: [][]string{}}, cars_json_file_err
		}

		serverLog("account_cars_utilities.go:24", fmt.Sprintf("created \"account_cars/%s.json\"", login))
		return AccountCars{Cars: []string{}, Modifications: [][]string{}, Paint: [][]string{}}, nil
	}

	cars_json, cars_json_err := fileToByte(fmt.Sprintf("./save_data/account_cars/%s.json", login))
	if cars_json_err != nil {
		serverLog("account_cars_utilities.go:30", fmt.Sprintf("failed to read \"account_cars/%s.json\", the error \"%s\"", login, cars_json_err))
		return AccountCars{Cars: []string{}, Modifications: [][]string{}, Paint: [][]string{}}, cars_json_err
	}

	cars = AccountCars{}
	cars_json_err = json.Unmarshal(cars_json, &cars)
	if cars_json_err != nil {
		serverLog("account_cars_utilities.go:37", fmt.Sprintf("failed to unmarshal \"account_cars/%s.json\", the error \"%s\"", login, cars_json_err))
		return AccountCars{Cars: []string{}, Modifications: [][]string{}, Paint: [][]string{}}, cars_json_err
	}

	return cars, nil
}

func addCarToAccount(login string, car string) (err error) {
	cars, cars_err := getCarsData(login)
	if cars_err != nil {
		serverLog("account_cars_utilities.go:47", fmt.Sprintf("failed to get cars data, the error \"%s\"", cars_err))
		return cars_err
	}

	cars.Cars = append(cars.Cars, car)
	cars.Modifications = append(cars.Modifications, []string{})
	cars.Paint = append(cars.Paint, []string{})

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

	serverLog("account_cars_utilities.go:66", fmt.Sprintf("added the car \"%s\" to the account \"%s\"", car, login))
	return nil
}

func removeCarFromAccount(login string, car string) (err error) {
	cars, cars_err := getCarsData(login)
	if cars_err != nil {
		serverLog("account_cars_utilities.go:73", fmt.Sprintf("failed to get cars data, the error \"%s\"", cars_err))
		return cars_err
	}

	car_index := 0
	for i, v := range cars.Cars {
		if v == car {
			cars.Cars = append(cars.Cars[:i], cars.Cars[i+1:]...)
			cars.Modifications = append(cars.Modifications[:i], cars.Modifications[i+1:]...)
			cars.Paint = append(cars.Paint[:i], cars.Paint[i+1:]...)
			car_index = i
			break
		}
	}
	for i := car_index; i < len(cars.Cars); i++ {
		if cars.Cars[i][:len(cars.Cars[i])-2] == car[:len(car)-2] {
			car_number, car_number_err := strconv.Atoi(string(cars.Cars[i][len(cars.Cars[i])-1]))
			if car_number_err != nil {
				serverLog("account_cars_utilities.go:91", fmt.Sprintf("failed to convert the car number to an integer, the error \"%s\"", car_number_err))
				return car_number_err
			}
			cars.Cars[i] = fmt.Sprintf("%s%d", cars.Cars[i][:len(cars.Cars[i])-1], car_number-1)
		}
	}

	cars_json, cars_json_err := json.Marshal(cars)
	if cars_json_err != nil {
		serverLog("account_cars_utilities.go:100", fmt.Sprintf("failed to marshal to file \"account_cars/%s.json\", the error \"%s\"", login, cars_json_err))
		return cars_json_err
	}
	write_err := os.WriteFile(fmt.Sprintf("./save_data/account_cars/%s.json", login), cars_json, 0644)
	if write_err != nil {
		serverLog("account_cars_utilities.go:105", fmt.Sprintf("failed to write to file \"account_cars/%s.json\", the error \"%s\"", login, write_err))
		return write_err
	}

	serverLog("account_cars_utilities.go:109", fmt.Sprintf("removed the car \"%s\" from the account \"%s\"", car, login))
	return nil
}
