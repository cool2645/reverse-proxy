package server

import (
	log "github.com/Sirupsen/logrus"
	"net/http"
	"strings"
	"net/url"
	"net/http/httputil"
	"sync"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"

	"github.com/2645Corp/reverse-proxy/model"
	"github.com/2645Corp/reverse-proxy/config"
)

type handle struct {
	host string
	port string
}

var websites map[string]handle
var mu = sync.Mutex{}

func getWebsiteName(r *http.Request) string {
	host := r.Host
	if i := strings.Index(host, "."); i != -1 {
		host = host[:i]
	} else if i := strings.Index(host, ":"); i != -1 {
		host = host[:i]
	}
	return host
}

func serveHTTP(w http.ResponseWriter, r *http.Request) {
	i := getWebsiteName(r)
	if _, ok := websites[i]; ok {
		remote, err := url.Parse("http://" + websites[i].host + ":" + websites[i].port)
		if err != nil {
			panic(err)
		}
		proxy := httputil.NewSingleHostReverseProxy(remote)
		proxy.ServeHTTP(w, r)
	} else {
		http.Error(w, "403 Forbidden", 403)
	}
}

func StartServer() {
	globCfg := config.GlobCfg

	db, err := gorm.Open("mysql", config.ParseDSN(globCfg))
	if err != nil {
		log.Fatal(err)
	}
	log.Infof("Database connected")
	defer db.Close()

	db.AutoMigrate(&model.User{}, &model.Website{})

	db_websites, err := model.ListWebsites(db)
	if err != nil {
		log.Fatal(err)
	}
	websites = make(map[string]handle)
	for _, v := range db_websites {
		websites[v.Name] = handle{host: v.Host, port: v.Port}
	}

	go func() {
		mux := http.NewServeMux()
		mux.HandleFunc("/", serveHTTP)
		err = http.ListenAndServe(":"+globCfg.PROXY_PORT, mux)
		if err != nil {
			log.Fatalln("ListenAndServe: ", err)
		}
	}()
}

func AddWebsite(name string, handle handle) {
	mu.Lock()
	websites[name] = handle
	mu.Unlock()
}

func DelWebsite(name string) {
	mu.Lock()
	delete(websites, name)
	mu.Unlock()
}
