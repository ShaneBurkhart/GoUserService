package config

import (
	"github.com/ShaneBurkhart/GoUserService/app/controllers"
	"github.com/gorilla/mux"
	"net/http"
)

func SetupRoutes(r *mux.Router) {
	r.HandleFunc("/user/{id}", controllers.EmailController).Methods("GET")

	r.HandleFunc("/register", controllers.RegisterController).Methods("POST")
	r.HandleFunc("/login", controllers.LoginController).Methods("POST")
}

func Serve(r *mux.Router) {
	http.Handle("/", r)
	http.ListenAndServe(":3000", nil)
}
