package main

import (
	"log"
	"net/http"

	"app/routes"
)

func main() {
	r := routes.NewRouter()
	log.Fatal(http.ListenAndServe("0.0.0.0:8081", r))
}
