package main

import (
	"bytes"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"net/http"
	"net/mail"
	"os"
	"regexp"
	"strings"
	"time"

	"golang.org/x/crypto/argon2"
)

func loginHandler(w http.ResponseWriter, r *http.Request) {
	request := reqBodyToString(r)

	validity := checkLoginValidity(request)
	if validity != "0" {
		w.Write([]byte(validity))
		serverLog("login.go:24", fmt.Sprintf("an account could not created because of \"%s\"", validity))
		return
	}

	data := strings.Split(request, "|")

	password_file, password_file_err := fileToByte(fmt.Sprintf("./save_data/account_passwords/%s-pass", data[0]))
	if password_file_err != nil {
		serverLog("login.go:32", fmt.Sprintf("could not read the file \"%s-pass\", the error \"%s\"", data[0], password_file_err))
		w.Write([]byte("text-danger\\Server could not read account data, try again later"))
		return
	}
	salt_file, salt_file_err := fileToByte(fmt.Sprintf("./save_data/account_passwords/%s-salt", data[0]))
	if salt_file_err != nil {
		serverLog("login.go:38", fmt.Sprintf("could not read the file \"%s-salt\", the error \"%s\"", data[0], salt_file_err))
		w.Write([]byte("text-danger\\Server could not read account data, try again later"))
		return
	}

	if !bytes.Equal(password_file, argon2.IDKey([]byte(data[1]), salt_file, 1, 64*1024, 4, 32)) {
		w.Write([]byte("text-danger\\Account doesn't exist or the password is wrong"))
		return
	}

	files, _ := os.ReadDir("./save_data/account_tokens")
	for _, file := range files {
		token_data := AccountToken{}
		token_file, token_err := fileToByte(fmt.Sprintf("./save_data/account_tokens/%s", file.Name()))
		if token_err != nil {
			serverLog("login.go:53", fmt.Sprintf("failed to read the token file \"%s\", the error \"%s\"", file.Name(), token_err))
			continue
		}
		if json.Unmarshal(token_file, &token_data) != nil {
			serverLog("login.go:57", fmt.Sprintf("failed to unmarshal the token file \"%s\"", file.Name()))
			continue
		}

		if data[0] == token_data.Login {
			removal_err := os.Remove(fmt.Sprintf("./save_data/account_tokens/%s", file.Name()))
			if removal_err != nil {
				serverLog("login.go:64", fmt.Sprintf("failed to remove the token file \"%s\", the error \"%s\"", file.Name(), removal_err))
				continue
			}
		}
	}

	account_token := ""
	for i := 0; i < 32; i++ {
		account_token = fmt.Sprintf("%s%d", account_token, randNumber(0, 9))
	}

	account_token_json, account_token_json_err := json.Marshal(AccountToken{data[0], time.Now().Unix() + 1200})
	if account_token_json_err != nil {
		serverLog("login.go:77", fmt.Sprintf("could not marshal the account token data to json, the error \"%s\"", account_token_json_err))
		w.Write([]byte("text-danger\\Server could not create account token, try again later"))
		return
	}
	account_token_file_err := os.WriteFile(fmt.Sprintf("./save_data/account_tokens/%s.json", account_token), account_token_json, 0644)
	if account_token_file_err != nil {
		serverLog("login.go:83", fmt.Sprintf("could not write the account token file \"%s.json\", the error \"%s\"", account_token, account_token_file_err))
		w.Write([]byte("text-danger\\Server could not create account token, try again later"))
		return
	}

	serverLog("login.go:88", fmt.Sprintf("an account was just logged in with the login \"%s\"", data[0]))
	w.Write([]byte(fmt.Sprintf("text-success\\Account logged in successfully\\%s", account_token)))
}

func checkLoginValidity(str string) string {
	pipeCount := 0
	for i := 0; i < len(str); i++ {
		if str[i] == '|' {
			pipeCount++
		}
	}
	if pipeCount != 1 {
		return "text-danger\\Invalid data (cannot contain | or the data sent was too long)"
	}

	data := strings.Split(str, "|")

	if strings.ReplaceAll(data[0], " ", "") == "" {
		return "text-danger\\Invalid login (cannot be empty)"
	}

	login_check, login_check_err := regexp.MatchString("^[a-zA-Z0-9]*$", data[0])
	if login_check_err != nil {
		return "text-danger\\Server could not check the login validity, try again later"
	} else if !login_check {
		return "text-danger\\Invalid login (can only contain letters and numbers)"
	}

	if len(data[0]) > 16 {
		return "text-danger\\Invalid login (cannot be longer than 16 characters)"
	}

	if !checkFileExistance(fmt.Sprintf("./save_data/accounts/%s.json", data[0])) {
		return "text-danger\\Account doesn't exist or the password is wrong"
	}

	return "0"
}

func gmailCodeHandler(w http.ResponseWriter, r *http.Request) {
	file_infos, file_infos_err := os.ReadDir("./save_data/gmail_codes/")
	if file_infos_err != nil {
		serverLog("login.go:130", fmt.Sprintf("could not read directory: %s", file_infos_err))
	} else {
		for _, file_info := range file_infos {
			gmail_code := GmailCode{}
			gmail_code_json, gmail_code_json_err := fileToByte(fmt.Sprintf("./save_data/gmail_codes/%s", file_info.Name()))
			if gmail_code_json_err != nil {
				serverLog("login.go:136", fmt.Sprintf("could not read the gmail code file \"%s\", the error \"%s\"", file_info.Name(), gmail_code_json_err))
				continue
			}
			if json.Unmarshal(gmail_code_json, &gmail_code) != nil {
				serverLog("login.go:140", fmt.Sprintf("could not unmarshal the gmail code file \"%s\", the error \"%s\"", file_info.Name(), gmail_code_json_err))
				continue
			}
			if time.Now().Unix() > gmail_code.Time {
				removal_err := os.Remove(fmt.Sprintf("./save_data/gmail_codes/%s", file_info.Name()))
				if removal_err != nil {
					serverLog("login.go:146", fmt.Sprintf("could not delete the gmail code file \"%s\", the error \"%s\"", file_info.Name(), removal_err))
					continue
				}
				serverLog("login.go:149", fmt.Sprintf("the gmail code file \"%s\" was deleted because it expired", file_info.Name()))
			}
		}
	}

	gmail := reqBodyToString(r)
	gmail_validity := checkGmailValidity(gmail)
	if gmail_validity != "0" {
		w.Write([]byte(gmail_validity))
		return
	}

	code := ""
	for i := 0; i < 6; i++ {
		code = fmt.Sprintf("%s%d", code, randNumber(0, 9))
	}

	gmail_code := GmailCode{time.Now().Unix() + 180, code}
	json_data, json_err := json.Marshal(gmail_code)
	if json_err != nil {
		serverLog("login.go:169", fmt.Sprintf("could not marshal the gmail code data to json, the gmail \"%s\", the error \"%s\"", gmail, json_err))
		w.Write([]byte("text-danger\\Server could not send the code, try again later"))
		return
	}
	gmail_code_file_err := os.WriteFile(fmt.Sprintf("./save_data/gmail_codes/%s.json", gmail), json_data, 0644)
	if gmail_code_file_err != nil {
		serverLog("login.go:175", fmt.Sprintf("could not write the gmail code file \"%s.json\", the error \"%s\"", gmail, gmail_code_file_err))
		w.Write([]byte("text-danger\\Server could not send the code, try again later"))
		return
	}

	sendMail("Horize Account Verification", fmt.Sprintf("You have 3 minutes to verify your account, your verification code is: %s", code), gmail)
	w.Write([]byte("text-success\\Code sent successfully"))
}

func checkGmailValidity(gmail string) string {
	_, mail_err := mail.ParseAddress(gmail)
	if mail_err != nil {
		return "text-danger\\Invalid gmail syntax"
	} else if gmail[len(gmail)-10:] != "@gmail.com" {
		return "text-danger\\Invalid gmail (must be a gmail account)"
	}

	if checkFileExistance(fmt.Sprintf("./save_data/existing_gmails/%s", gmail)) {
		return "text-danger\\Gmail is already in use"
	}

	return "0"
}

func registerHandler(w http.ResponseWriter, r *http.Request) {
	request := reqBodyToString(r)

	validity := checkRegisterValidity(request)
	if validity != "0" {
		w.Write([]byte(validity))
		serverLog("login.go:205", fmt.Sprintf("an account could not created because of \"%s\"", validity))
		return
	}

	data := strings.Split(request, "|")

	salt := make([]byte, 16)
	_, salt_err := rand.Read(salt)
	if salt_err != nil {
		serverLog("login.go:214", fmt.Sprintf("could not generate salt for the account password, the error \"%s\"", salt_err))
		w.Write([]byte("text-danger\\Server could not generate salt for the account password, try again later"))
		return
	}
	hashed_password := argon2.IDKey([]byte(data[1]), salt, 1, 64*1024, 4, 32)

	json_data, json_err := json.Marshal(Account{data[2], data[3], ""})
	if json_err != nil {
		serverLog("login.go:222", fmt.Sprintf("could not marshal the account data to json, the error \"%s\"", json_err))
		w.Write([]byte("text-danger\\Server could not create the account, try again later"))
		return
	}
	account_json_file_err := os.WriteFile(fmt.Sprintf("./save_data/accounts/%s.json", data[0]), json_data, 0644)
	if account_json_file_err != nil {
		serverLog("login.go:228", fmt.Sprintf("could not write the account file \"%s.json\", the error \"%s\"", data[0], account_json_file_err))
		w.Write([]byte("text-danger\\Server could not create the account, try again later"))
		return
	}
	gmail_file_err := os.WriteFile(fmt.Sprintf("./save_data/existing_gmails/%s", data[3]), []byte(""), 0644)
	if gmail_file_err != nil {
		serverLog("login.go:234", fmt.Sprintf("could not write the gmail file \"%s.json\", the error \"%s\"", data[3], gmail_file_err))
		account_file_remove_err := os.Remove(fmt.Sprintf("./save_data/accounts/%s.json", data[0]))
		if account_file_remove_err != nil {
			serverLog("login.go:237", fmt.Sprintf("could not remove the account file \"%s.json\", the error \"%s\"", data[0], account_file_remove_err))
		}
		w.Write([]byte("text-danger\\Server could not create the account, try again later"))
		return
	}
	password_file_err := os.WriteFile(fmt.Sprintf("./save_data/account_passwords/%s-pass", data[0]), hashed_password, 0644)
	if password_file_err != nil {
		serverLog("login.go:244", fmt.Sprintf("could not write the password file \"%s-pass\", the error \"%s\"", data[0], password_file_err))
		account_file_remove_err := os.Remove(fmt.Sprintf("./save_data/accounts/%s.json", data[0]))
		if account_file_remove_err != nil {
			serverLog("login.go:247", fmt.Sprintf("could not remove the account file \"%s.json\", the error \"%s\"", data[0], account_file_remove_err))
		}
		gmail_file_remove_err := os.Remove(fmt.Sprintf("./save_data/existing_gmails/%s", data[3]))
		if gmail_file_remove_err != nil {
			serverLog("login.go:251", fmt.Sprintf("could not remove the gmail file \"%s.json\", the error \"%s\"", data[3], gmail_file_remove_err))
		}
		w.Write([]byte("text-danger\\Server could not create the account, try again later"))
		return
	}
	salt_file_err := os.WriteFile(fmt.Sprintf("./save_data/account_passwords/%s-salt", data[0]), salt, 0644)
	if salt_file_err != nil {
		serverLog("login.go:258", fmt.Sprintf("could not write the salt file \"%s-salt\", the error \"%s\"", data[0], salt_file_err))
		account_file_remove_err := os.Remove(fmt.Sprintf("./save_data/accounts/%s.json", data[0]))
		if account_file_remove_err != nil {
			serverLog("login.go:261", fmt.Sprintf("could not remove the account file \"%s.json\", the error \"%s\"", data[0], account_file_remove_err))
		}
		gmail_file_remove_err := os.Remove(fmt.Sprintf("./save_data/existing_gmails/%s", data[3]))
		if gmail_file_remove_err != nil {
			serverLog("login.go:265", fmt.Sprintf("could not remove the gmail file \"%s.json\", the error \"%s\"", data[3], gmail_file_remove_err))
		}
		password_file_remove_err := os.Remove(fmt.Sprintf("./save_data/account_passwords/%s-pass", data[0]))
		if password_file_remove_err != nil {
			serverLog("login.go:269", fmt.Sprintf("could not remove the password file \"%s-pass\", the error \"%s\"", data[0], password_file_remove_err))
		}
		w.Write([]byte("text-danger\\Server could not create the account, try again later"))
		return
	}

	serverLog("login.go:275", fmt.Sprintf("an account was just created with the login \"%s\"", data[0]))
	w.Write([]byte("text-success\\Account created successfully, you can now login"))
}

func checkRegisterValidity(str string) string {
	pipeCount := 0
	for i := 0; i < len(str); i++ {
		if str[i] == '|' {
			pipeCount++
		}
	}
	if pipeCount != 4 {
		return "text-danger\\Invalid data (cannot contain | or the data sent was too long)"
	}

	data := strings.Split(str, "|")

	if strings.ReplaceAll(data[0], " ", "") == "" {
		return "text-danger\\Invalid login (cannot be empty)"
	}

	login_check, login_check_err := regexp.MatchString("^[a-zA-Z0-9]*$", data[0])
	if login_check_err != nil {
		return "text-danger\\Server could not check the login validity, try again later"
	} else if !login_check {
		return "text-danger\\Invalid login (can only contain letters and numbers)"
	}

	if checkFileExistance(fmt.Sprintf("./save_data/accounts/%s.json", data[0])) {
		return "text-danger\\Account already exists"
	}

	if len(data[0]) > 16 {
		return "text-danger\\Invalid login (cannot be longer than 16 characters)"
	}

	gmail_validity := checkGmailValidity(data[3])
	if gmail_validity != "0" {
		return gmail_validity
	}

	if !checkFileExistance(fmt.Sprintf("./save_data/gmail_codes/%s.json", data[3])) {
		return "text-danger\\Gmail doesn't have a verification code"
	}
	gmail_code := GmailCode{}
	gmail_code_json, gmail_code_json_err := fileToByte(fmt.Sprintf("./save_data/gmail_codes/%s.json", data[3]))
	if gmail_code_json_err != nil {
		return "text-danger\\Server could not check the gmail code, try again later"
	}
	json.Unmarshal(gmail_code_json, &gmail_code)
	if time.Now().Unix() > gmail_code.Time {
		return "text-danger\\Gmail verification code expired"
	}
	if data[4] != gmail_code.Code {
		return "text-danger\\Invalid gmail verification code or the data sent was too long"
	}

	return "0"
}

func logoutHandler(w http.ResponseWriter, r *http.Request) {
	request := reqBodyToString(r)

	validity, _, err := checkTokenValidity(request)
	if err != nil {
		serverLog("login.go:340", fmt.Sprintf("could not check validity of \"%s.json\", the error \"%s\"", request, err))
		w.Write([]byte("text-danger\\Server could not log you out, try again later"))
		return
	}
	if !validity {
		serverLog("login.go:345", fmt.Sprintf("someone tried to log out but the token was invalid, the token \"%s\"", request))
	}

	removal_err := os.Remove(fmt.Sprintf("./save_data/account_tokens/%s.json", request))
	if removal_err != nil {
		serverLog("login.go:350", fmt.Sprintf("could not remove \"%s.json\", the error \"%s\"", request, err))
		w.Write([]byte("text-danger\\Server could not log you out, try again later"))
		return
	}

	serverLog("login.go:355", fmt.Sprintf("An account just logged out with the token \"%s\"", request))
	w.Write([]byte("text-success\\Account logged out"))
}
