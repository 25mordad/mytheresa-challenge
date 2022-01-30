package server

import (
	"net/http"
	"time"

	"github.com/25mordad/mytheresa-challenge/config"
	"github.com/25mordad/mytheresa-challenge/internal/controller"
	"github.com/gorilla/mux"
	cache "github.com/patrickmn/go-cache"
)

//Routing is
func Routing(r *mux.Router, c config.Configuration) {

	var Cache = cache.New(18*time.Hour, 18*time.Hour)
	h := controller.NewBaseHandler(Cache)

	r.Handle("/", http.HandlerFunc(controller.Home)).Methods("GET")
	r.Handle("/api/v0/products", http.HandlerFunc(h.ImportProduct)).Methods("PUT")
	r.Handle("/api/v0/products", http.HandlerFunc(h.FilterProduct)).Methods("GET")

}
