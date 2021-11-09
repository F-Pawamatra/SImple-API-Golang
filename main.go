package main

import (
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/login", Login)
	http.HandleFunc("/home", Home)

	log.Fatal(http.ListenAndServe(":8080", nil))
}