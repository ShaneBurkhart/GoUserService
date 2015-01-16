package controllers

import (
	"fmt"
	"net/http"
)

func LoginController(resp http.ResponseWriter, req *http.Request) {
	fmt.Fprint(resp, "LoginController")
}
