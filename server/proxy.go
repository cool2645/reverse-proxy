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
var db *gorm.DB

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
	var i string
	if strings.Count(r.Host, config.GlobCfg.DOMAIN) > 0 {
		i = getWebsiteName(r)
	} else if ok, website := model.FindWebsiteByDomain(db, r.Host); ok {
		i = website.Name
	}
	if _, ok := websites[i]; ok {
		remote, err := url.Parse("http://" + websites[i].host + ":" + websites[i].port)
		if err != nil {
			panic(err)
		}
		proxy := httputil.NewSingleHostReverseProxy(remote)
		r.Host = websites[i].host + ":" + websites[i].port
		r.URL.Host = websites[i].host + ":" + websites[i].port
		r.Header.Set("Host", websites[i].host+":"+websites[i].port)
		logHttp(r, i)
		proxy.ServeHTTP(w, r)
	} else {
		//redirect to the default
		http.Redirect(w, r, "http://www."+config.GlobCfg.DOMAIN, 302)
	}
}

func StartServer() {
	globCfg := config.GlobCfg

	var err error
	log.Warnf("Connecting to database")
	db, err = gorm.Open("mysql", config.ParseDSN(globCfg))
	if err != nil {
		log.Fatal(err)
	}
	log.Warnf("Database connected")

	db.AutoMigrate(&model.User{}, &model.Website{}, &model.Domain{})

	db_websites := model.ListWebsites(db)
	if err != nil {
		log.Fatal(err)
	}
	websites = make(map[string]handle)
	for _, v := range db_websites {
		websites[v.Name] = handle{host: v.Host, port: v.Port}
		log.Warnf("Loaded reverse proxy [" + v.Name + "] for port " + v.Port + " of " + v.Host)
	}

	go func() {
		mux := http.NewServeMux()
		mux.HandleFunc("/", serveHTTP)
		log.Warnf("Starting reverse proxy on port " + globCfg.PROXY_PORT)
		err = http.ListenAndServe(":"+globCfg.PROXY_PORT, mux)
		if err != nil {
			log.Fatalln("ListenAndServe: ", err)
		}
	}()
}

func addWebsite(name string, handle handle) {
	mu.Lock()
	websites[name] = handle
	log.Warnf("Added reverse proxy [" + name + "] for port " + handle.port + " of " + handle.host)
	mu.Unlock()
}

func delWebsite(name string) {
	mu.Lock()
	delete(websites, name)
	log.Warnf("Terminated reverse proxy [" + name + "]")
	mu.Unlock()
}
