package main

import (
	"github.com/BurntSushi/toml"
	log "github.com/Sirupsen/logrus"

	"github.com/2645Corp/reverse-proxy/config"
	"github.com/2645Corp/reverse-proxy/server"
	"time"
)

func main() {

	log.Warnf("Initializing program")
	log.Warnf("Loading config file")
	_, err := toml.DecodeFile("config.toml", &config.GlobCfg)
	if err != nil {
		panic(err)
	}

	log.Warnf("Starting reverse proxy server")
	server.StartServer()

	log.Warnf("Starting manage server")
	server.StartManager()

	for {
		time.Sleep(time.Second)
	}
}
