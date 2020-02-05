package main

import (
	"fmt"
	"log"
	"net/http"
	"time"
	"gitlab.com/sgryczan/seel/api"

	"github.com/gorilla/mux"
)

var Items = []string{}

var (
	versionMajor string
	versionMinor string
	versionPatch string
)


func main() {

	r := mux.NewRouter()
	fmt.Println("Seel v%d.%d.%d", versionMajor, versionMinor, versionPatch)

	err := api.MapController()
	if err != nil {
		log.Print(err)
	}

	r.HandleFunc("/", api.HomeHandler)
	r.HandleFunc("/convert", api.ConvertHandler).Methods("POST")
	r.HandleFunc("/create", api.CreateHandler).Methods("POST")

	sh := http.StripPrefix("/api",
		http.FileServer(http.Dir("./swaggerui/")))
	r.PathPrefix("/api/").Handler(sh)

	srv := &http.Server{
		Handler:      r,
		Addr:         ":8080",
		WriteTimeout: 10 * time.Second,
		ReadTimeout:  10 * time.Second,
	}

	log.Fatal(srv.ListenAndServe())
}
