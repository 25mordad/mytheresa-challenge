package server

import (
	"log"
	"net/http"

	"github.com/25mordad/mytheresa-challenge/config"

	"github.com/gorilla/mux"
)

//Run is going to run the server
func Run(c config.Configuration) {
	r := mux.NewRouter()

	//routing
	Routing(r, c)

	// return "Server Starts"
	log.Fatal(http.ListenAndServe(c.Port, r))

}
