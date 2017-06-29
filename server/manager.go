package server

import (
	"net/http"
	log "github.com/Sirupsen/logrus"
	"text/template"

	"github.com/2645Corp/reverse-proxy/config"
	"fmt"
)

func index(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("tmpl/index.html")
	if err != nil {
		log.Fatal(err)
	}
	t.Execute(w, nil)
}

func addWebsite(w http.ResponseWriter, r *http.Request) {
	AddWebsite("homestead", handle{host: "127.0.0.1", port: "8000"})
	fmt.Fprint(w, "OK")
}

func StartManager() {
	globCfg := config.GlobCfg

	go func() {
		mux := http.NewServeMux()
		mux.HandleFunc("/", index)
		mux.HandleFunc("/add", addWebsite)
		err := http.ListenAndServe(":"+globCfg.MANAGE_PORT, mux)
		if err != nil {
			log.Fatalln("ListenAndServe: ", err)
		}
	}()
}
