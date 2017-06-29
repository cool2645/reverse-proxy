package main

import (
	"github.com/BurntSushi/toml"

	"github.com/2645Corp/reverse-proxy/config"
	"github.com/2645Corp/reverse-proxy/server"
	"time"
)

func main() {
	_, err := toml.DecodeFile("config.toml", &config.GlobCfg)
	if err != nil {
		panic(err)
	}

	server.StartServer()
	server.StartManager()

	for {
		time.Sleep(time.Second)
	}
}
