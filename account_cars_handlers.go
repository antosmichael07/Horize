package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

func getCarsHandler(w http.ResponseWriter, r *http.Request) {
	request := reqBodyToString(r)

	valid, login, err := checkTokenValidity(request)
	if err != nil {
		serverLog("account_cars_handlers.go:15", fmt.Sprintf("failed to check token validity, the error \"%s\"", err))
		w.Write([]byte("failed to check token validity"))
		return
	}
	if !valid {
		serverLog("account_cars_handlers.go:20", fmt.Sprintf("invalid token \"%s\"", request))
		w.Write([]byte("invalid token"))
		return
	}

	cars, cars_err := getCarsData(login)
	if cars_err != nil {
		serverLog("account_cars_handlers.go:27", fmt.Sprintf("failed to get cars data, the error \"%s\"", cars_err))
		w.Write([]byte("failed to get cars data"))
		return
	}

	cars_json, cars_json_err := json.Marshal(cars)
	if cars_json_err != nil {
		serverLog("account_cars_handlers.go:34", fmt.Sprintf("failed to marshal cars data, the error \"%s\"", cars_json_err))
		w.Write([]byte("failed to marshal cars data"))
		return
	}

	w.Write(cars_json)
}

func addCarHandler(w http.ResponseWriter, r *http.Request) {
	request := reqBodyToString(r)

	pipe_count := 0
	for i := 0; i < len(request); i++ {
		if request[i] == '|' {
			pipe_count++
		}
	}
	if pipe_count != 1 {
		serverLog("account_cars_handlers.go:52", fmt.Sprintf("invalid request \"%s\"", request))
		w.Write([]byte("invalid request"))
		return
	}

	data := strings.Split(request, "|")

	valid, login, err := checkTokenValidity(data[0])
	if err != nil {
		serverLog("account_cars_handlers.go:61", fmt.Sprintf("failed to check token validity, the error \"%s\"", err))
		w.Write([]byte("failed to check token validity"))
		return
	}
	if !valid {
		serverLog("account_cars_handlers.go:66", fmt.Sprintf("invalid token \"%s\"", data[0]))
		w.Write([]byte("invalid token"))
		return
	}

	is_car_valid := false
	for i := 0; i < len(car_list); i++ {
		if data[1] == car_list[i] {
			is_car_valid = true
			break
		}
	}
	if !is_car_valid {
		serverLog("account_cars_handlers.go:79", fmt.Sprintf("invalid car \"%s\"", data[1]))
		w.Write([]byte("invalid car"))
		return
	}

	cars, err := getCarsData(login)
	if err != nil {
		serverLog("account_cars_handlers.go:86", fmt.Sprintf("failed to get cars data, the error \"%s\"", err))
		w.Write([]byte("failed to get cars data"))
		return
	}
	car_count := 0
	for i := 0; i < len(cars.Cars); i++ {
		if cars.Cars[i][:len(cars.Cars[i])-2] == data[1] {
			car_count++
		}
	}
	data[1] = fmt.Sprintf("%s.%d", data[1], car_count)

	add_err := addCarToAccount(login, data[1])
	if add_err != nil {
		serverLog("account_cars_handlers.go:100", fmt.Sprintf("failed to add car to account, the error \"%s\"", add_err))
		w.Write([]byte("failed to add car to account"))
		return
	}

	w.Write([]byte("car added to account"))
}

func removeCarHandler(w http.ResponseWriter, r *http.Request) {
	request := reqBodyToString(r)

	pipe_count := 0
	for i := 0; i < len(request); i++ {
		if request[i] == '|' {
			pipe_count++
		}
	}
	if pipe_count != 1 {
		serverLog("account_cars_handlers.go:118", fmt.Sprintf("invalid request \"%s\"", request))
		w.Write([]byte("invalid request"))
		return
	}

	data := strings.Split(request, "|")

	valid, login, err := checkTokenValidity(data[0])
	if err != nil {
		serverLog("account_cars_handlers.go:127", fmt.Sprintf("failed to check token validity, the error \"%s\"", err))
		w.Write([]byte("failed to check token validity"))
		return
	}
	if !valid {
		serverLog("account_cars_handlers.go:132", fmt.Sprintf("invalid token \"%s\"", data[0]))
		w.Write([]byte("invalid token"))
		return
	}

	account_cars, account_cars_err := getCarsData(login)
	if account_cars_err != nil {
		serverLog("account_cars_handlers.go:139", fmt.Sprintf("failed to get cars data, the error \"%s\"", account_cars_err))
		w.Write([]byte("failed to get cars data"))
		return
	}
	is_car_valid := false
	for i := 0; i < len(account_cars.Cars); i++ {
		if account_cars.Cars[i] == data[1] {
			is_car_valid = true
			break
		}
	}
	if !is_car_valid {
		serverLog("account_cars_handlers.go:151", fmt.Sprintf("invalid car \"%s\"", data[1]))
		w.Write([]byte("invalid car"))
		return
	}

	remove_err := removeCarFromAccount(login, data[1])
	if remove_err != nil {
		serverLog("account_cars_handlers.go:158", fmt.Sprintf("failed to remove car from account, the error \"%s\"", remove_err))
		w.Write([]byte("failed to remove car from account"))
		return
	}

	w.Write([]byte("car removed from account"))
}
