package main

import (
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"net/smtp"
	"os"
	"time"
)

func reqBodyToString(r *http.Request) string {
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

func checkFileExistance(filename string) bool {
	_, file_err := os.Stat(filename)
	return !os.IsNotExist(file_err)
}

func serverLog(message string) {
	str := fmt.Sprintf("[%s] %s", time.Now().String()[:19], message)

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
		serverLog(fmt.Sprintf("smtp had an error \"%s\", the receiver \"%s\", the message \"%s\" \"%s\"", smtp_err, to, subject, body))
		return
	}
	serverLog(fmt.Sprintf("successfully sent a message to \"%s\", the message \"%s\" \"%s\"", to, subject, body))
}

func randNumber(min int, max int) int {
	return rand.Intn(max-min) + min
}

func fileToByte(filename string) ([]byte, error) {
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
