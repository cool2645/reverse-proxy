package main

import (
	"github.com/BurntSushi/toml"
	log "github.com/Sirupsen/logrus"

	"github.com/2645Corp/reverse-proxy/config"
	"github.com/2645Corp/reverse-proxy/server"
	"time"
)

func main() {

	log.Infof("Initializing program")
	log.Infof("Loading config file")
	_, err := toml.DecodeFile("config.toml", &config.GlobCfg)
	if err != nil {
		panic(err)
	}

	log.Infof("Starting reverse proxy server")
	server.StartServer()

	log.Infof("Starting manage server")
	server.StartManager()

	for {
		time.Sleep(time.Second)
	}
}
