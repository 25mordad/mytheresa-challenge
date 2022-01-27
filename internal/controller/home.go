package controller

import (
	"fmt"
	"log"

	"net/http"
)

//Home is
func Home(w http.ResponseWriter, r *http.Request) {
	log.Println("github.com/25mordad/mytheresa-challenge")
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprint(w, "github.com/25mordad/mytheresa-challenge")

}
