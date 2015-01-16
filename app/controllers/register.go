package controllers

import (
	"fmt"
	"net/http"
)

func RegisterController(resp http.ResponseWriter, req *http.Request) {
	fmt.Fprint(resp, "RegisterController")
}
