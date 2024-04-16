package main

import (
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"net/smtp"
	"os"
	"time"
)

func reqBodyToString(r *http.Request) (req string) {
	data := r.Body
	text := make([]byte, 256)
	data.Read(text)
	stringData := string(text)
	for i := 0; i < len(stringData); i++ {
		if stringData[i] == '\u0000' {
			stringData = stringData[:i]
			break
		}
	}
	return stringData
}

func checkFileExistance(filename string) (existance bool) {
	_, file_err := os.Stat(filename)
	return !os.IsNotExist(file_err)
}

func serverLog(location string, message string) {
	str := fmt.Sprintf("[%s] [%s] %s", time.Now().String()[:19], location, message)

	fmt.Println(str)

	file, _ := os.OpenFile(fmt.Sprintf("./logs/%s.txt", time.Now().String()[:10]), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	file.WriteString(str + "\n")
	file.Close()
}

func sendMail(subject string, body string, to string) {
	from := "horize.noreply@gmail.com"
	pass := "byxv owbq sjmi kwkf"

	msg := "From: " + from + "\n" +
		"To: " + to + "\n" +
		"Subject: " + subject + "\n\n" +
		body

	smtp_err := smtp.SendMail("smtp.gmail.com:587",
		smtp.PlainAuth("", from, pass, "smtp.gmail.com"),
		from, []string{to}, []byte(msg))

	if smtp_err != nil {
		serverLog("utilities.go:57", fmt.Sprintf("smtp had an error \"%s\", the receiver \"%s\", the message \"%s\" \"%s\"", smtp_err, to, subject, body))
		return
	}
	serverLog("utilities.go:60", fmt.Sprintf("successfully sent a message to \"%s\", the message \"%s\" \"%s\"", to, subject, body))
}

func randNumber(min int, max int) (num int) {
	return rand.Intn(max-min) + min
}

func fileToByte(filename string) (file_data []byte, err error) {
	file, file_err := os.Open(filename)
	if file_err != nil {
		return nil, file_err
	}
	defer file.Close()

	file_byte, read_err := io.ReadAll(file)
	if read_err != nil {
		return nil, read_err
	}

	return file_byte, nil
}

func checkTokenValidity(token string) (valid bool, login string, err error) {
	go deleteExpiredTokens()

	token_data := AccountToken{}
	token_file, token_err := fileToByte(fmt.Sprintf("./save_data/account_tokens/%s.json", token))
	if token_err != nil {
		return false, "", token_err
	}
	json.Unmarshal(token_file, &token_data)

	if time.Now().Unix() > token_data.Time {
		removal_err := os.Remove(fmt.Sprintf("./save_data/account_tokens/%s.json", token))
		if removal_err != nil {
			serverLog("utilities.go:95", fmt.Sprintf("failed to remove the token file \"%s.json\", the error \"%s\"", token, removal_err))
			return false, "", removal_err
		}
		return false, "", nil
	}

	account_token_json, account_token_json_err := json.Marshal(AccountToken{token_data.Login, time.Now().Unix() + 1200})
	if account_token_json_err != nil {
		serverLog("utilities.go:103", fmt.Sprintf("could not marshal the account token data to json, the error \"%s\"", account_token_json_err))
		return
	}
	account_token_file_err := os.WriteFile(fmt.Sprintf("./save_data/account_tokens/%s.json", token), account_token_json, 0644)
	if account_token_file_err != nil {
		serverLog("utilities.go:108", fmt.Sprintf("could not write the account token file \"%s.json\", the error \"%s\"", token, account_token_file_err))
		return
	}

	return true, string(token_data.Login), nil
}

func deleteExpiredTokens() {
	if checkFileExistance("./save_data/account_tokens/.running_deletion") {
		return
	}

	serverLog("utilities.go:120", "deleting expired tokens will start in 10 seconds")

	time.Sleep(time.Second * 10)

	running_deletion_file, running_file_creation_err := os.Create("./save_data/account_tokens/.running_deletion")
	if running_file_creation_err != nil {
		serverLog("utilities.go:126", fmt.Sprintf("failed to create the file \".running_deletion\", the error \"%s\"", running_file_creation_err))
		return
	}
	running_deletion_file.Close()

	files, _ := os.ReadDir("./save_data/account_tokens")
	for _, file := range files {
		if file.Name() == ".running_deletion" {
			continue
		}
		token_data := AccountToken{}
		token_file, token_err := fileToByte(fmt.Sprintf("./save_data/account_tokens/%s", file.Name()))
		if token_err != nil {
			serverLog("utilities.go:139", fmt.Sprintf("failed to read the token file \"%s\", the error \"%s\"", file.Name(), token_err))
			continue
		}
		if json.Unmarshal(token_file, &token_data) != nil {
			serverLog("utilities.go:143", fmt.Sprintf("failed to unmarshal the token file \"%s\"", file.Name()))
			continue
		}

		if time.Now().Unix() > token_data.Time {
			removal_err := os.Remove(fmt.Sprintf("./save_data/account_tokens/%s", file.Name()))
			if removal_err != nil {
				serverLog("utilities.go:150", fmt.Sprintf("failed to remove the token file \"%s\", the error \"%s\"", file.Name(), removal_err))
				continue
			}
		}
	}

	time.Sleep(time.Second * 10)

	running_file_removal_err := os.Remove("./save_data/account_tokens/.running_deletion")
	if running_file_removal_err != nil {
		serverLog("utilities.go:160", fmt.Sprintf("failed to remove the file \".running_deletion\", the error \"%s\"", running_file_removal_err))
	}

	serverLog("utilities.go:163", "deleting expired tokens has finished")
}
