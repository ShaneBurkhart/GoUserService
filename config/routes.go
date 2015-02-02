package config

import (
	"github.com/ShaneBurkhart/GoUserService/app/controllers"
	"github.com/gorilla/mux"
	"net/http"
)

func SetupRoutes(r *mux.Router) {
	// TODO Test that only works when id is number
	r.HandleFunc("/user/{id:[0-9]+}", controllers.EmailController).Methods("GET")

	r.HandleFunc("/register", controllers.RegisterController).Methods("POST")
	r.HandleFunc("/login", controllers.LoginController).Methods("POST")
}

func Serve(r *mux.Router) {
	http.Handle("/", r)
	http.ListenAndServe(":3000", nil)
}
