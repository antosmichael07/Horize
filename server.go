package main

import (
	"log"
	"net/http"
)

func main() {
	init_post_handlers()
	init_get_handlers()
	log.Fatal(http.ListenAndServe(":8080", nil))
}
