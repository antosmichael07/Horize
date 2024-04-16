package main

import (
	"net/http"
)

func init_post_handlers() {
	http.HandleFunc("/post_login", loginHandler)
	http.HandleFunc("/post_register", registerHandler)
	http.HandleFunc("/post_logout", logoutHandler)

	http.HandleFunc("/post_gmail_code", gmailCodeHandler)
}

func init_get_handlers() {
	http.Handle("/", http.FileServer(http.Dir("./pages")))
}
