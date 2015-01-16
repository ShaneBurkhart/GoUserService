package main

import (
	"github.com/ShaneBurkhart/GoUserService/config"
	"github.com/gorilla/mux"
	"log"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	err := config.SetupDB()
	if err != nil {
		log.Fatal(err)
		return
	}
	defer config.CloseDB()

	if err := config.VerifyDB(); err != nil {
		log.Fatal(err)
		return
	}

	r := mux.NewRouter()
	config.SetupRoutes(r)
	config.Serve(r)
}
