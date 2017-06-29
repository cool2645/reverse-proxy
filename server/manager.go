package server

import (
	"net/http"
	log "github.com/Sirupsen/logrus"
	"text/template"
	"fmt"
	"github.com/astaxie/beego/session"

	"github.com/2645Corp/reverse-proxy/config"
)

var globalSessions *session.Manager

func init() {
	globalSessions, _ = session.NewManager("memory", &session.ManagerConfig{CookieName: "gosessionid", EnableSetCookie: true, Gclifetime: 3600 })
	go globalSessions.GC()
}

func index(w http.ResponseWriter, r *http.Request) {
	sess, _ := globalSessions.SessionStart(w, r)
	defer sess.SessionRelease(w)
	if r.Method == "GET" {
		if nickname := sess.Get("nickname"); nickname != nil {
			t, err := template.ParseFiles("tmpl/index.html")
			if err != nil {
				log.Fatal(err)
			}
			t.Execute(w, nil)
		} else {
			http.Redirect(w, r, "/auth", 302)
		}
	}
}

func auth(w http.ResponseWriter, r *http.Request) {
	sess, _ := globalSessions.SessionStart(w, r)
	defer sess.SessionRelease(w)
	if nickname := sess.Get("nickname"); nickname != nil {
		http.Redirect(w, r, "/", 302)
	}
	if r.Method == "GET" {
		t, err := template.ParseFiles("tmpl/auth.html")
		if err != nil {
			log.Fatal(err)
		}
		t.Execute(w, nil)
	} else if r.Method == "POST" {
		r.ParseForm()
		if r.Form.Get("smt") == "登录" {
			fmt.Fprint(w, "登录")
		} else if r.Form.Get("smt") == "注册" {
			fmt.Fprint(w, "注册")
		}
	}
}

func list(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("tmpl/list.html")
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
		mux.HandleFunc("/auth", auth)
		mux.HandleFunc("/list", list)
		mux.HandleFunc("/add", addWebsite)
		err := http.ListenAndServe(":"+globCfg.MANAGE_PORT, mux)
		if err != nil {
			log.Fatalln("ListenAndServe: ", err)
		}
	}()
}
