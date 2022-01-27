package server

import (
	"net/http"

	"github.com/gorilla/mux"
	// cache "github.com/patrickmn/go-cache"
	"github.com/25mordad/mytheresa-challenge/config"
	"github.com/25mordad/mytheresa-challenge/internal/controller"
)

//Routing is
func Routing(r *mux.Router, c config.Configuration) {

	// var Cache = cache.New(18*time.Hour, 18*time.Hour)
	// h := controller.NewBaseHandler(Cache)

	r.Handle("/", http.HandlerFunc(controller.Home)).Methods("GET")
	// r.Handle("/cars", http.HandlerFunc(h.Test)).Methods("PUT")

}
