package main

import (
	"log"

	"github.com/25mordad/mytheresa-challenge/config"
	"github.com/25mordad/mytheresa-challenge/internal/server"
)

func main() {
	conf := config.GetConfig()
	log.Println("Server is running" + conf.Port)
	server.Run(conf)
}
