package main

import (
	"net/http"
	"os"
)

func init_post_handlers() {
	http.HandleFunc("/post_login", loginHandler)
	http.HandleFunc("/post_register", registerHandler)
	http.HandleFunc("/post_logout", logoutHandler)

	http.HandleFunc("/post_gmail_code", gmailCodeHandler)

	http.HandleFunc("/post_get_cars", getCarsHandler)
	http.HandleFunc("/post_add_car", addCarHandler)
	http.HandleFunc("/post_remove_car", removeCarHandler)
}

func init_get_handlers() {
	fs := FileServerWith404(http.Dir("./pages"))
	http.Handle("/", http.StripPrefix("/", fs))
}

func FileServerWith404(root http.FileSystem) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		f, err := root.Open(r.URL.Path)
		if err != nil && os.IsNotExist(err) {
			http.ServeFile(w, r, "./.go_to_404.html")
			return
		}
		if err == nil {
			f.Close()
		}
		http.FileServer(root).ServeHTTP(w, r)
	})
}
