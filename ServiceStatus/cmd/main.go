package main

import (
	"log"
	"net/http"
	"service_status/pkg/handler"
	"service_status/pkg/middleware"

	"github.com/gorilla/mux"
)

func main() {
	r := mux.NewRouter()
	allhandler := handler.NewHandler()

	r.HandleFunc("/startNewTick", allhandler.NewTick).Methods(http.MethodPost)
	r.Use(middleware.Jwt)
	r.HandleFunc("/responseData", allhandler.GetResult).Methods(http.MethodPost)

	log.Fatal(http.ListenAndServe(":8089", r))
}
