package main

import (
	"log"
	"net/http"
	"server/handlers"
)

func main() {
	router := handlers.HandleRequests()
	log.Fatal(http.ListenAndServe(":8081", router))
}
