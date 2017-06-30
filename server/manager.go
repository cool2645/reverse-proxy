package server

import (
	"net/http"
	log "github.com/Sirupsen/logrus"
	"strconv"
	"fmt"
	"text/template"
	"github.com/astaxie/beego/session"

	"github.com/2645Corp/reverse-proxy/config"
	"github.com/2645Corp/reverse-proxy/model"
)

var globalSessions *session.Manager

func init() {
	globalSessions, _ = session.NewManager("memory", &session.ManagerConfig{CookieName: "gosessionid", EnableSetCookie: true, Gclifetime: 3600 })
	go globalSessions.GC()
}

func index(w http.ResponseWriter, r *http.Request) {
	sess, _ := globalSessions.SessionStart(w, r)
	defer sess.SessionRelease(w)
	if user_id, ok := sess.Get("user_id").(uint); ok {
		if r.Method == "GET" {
			t, err := template.ParseFiles("tmpl/index.html")
			if err != nil {
				log.Fatal(err)
			}
			nickname := sess.Get("nickname")
			t.Execute(w, nickname)
		} else if r.Method == "POST" {
			r.ParseForm()
			log.Infof("Handling request of adding reverse proxy [" + r.Form.Get("name") + "] for port " + r.Form.Get("port") + " of " + r.Form.Get("host") +  " by user " + fmt.Sprint(user_id))
			model.AddWebsite(db, r.Form.Get("name"), r.Form.Get("host"), r.Form.Get("port"), uint(user_id))
			addWebsite(r.Form.Get("name"), handle{host:r.Form.Get("host"), port:r.Form.Get("port")})
			http.Redirect(w, r, "/list", 302)
		}
	} else {
		http.Redirect(w, r, "/auth", 302)
	}
}

func auth(w http.ResponseWriter, r *http.Request) {
	sess, _ := globalSessions.SessionStart(w, r)
	defer sess.SessionRelease(w)
	if user_id := sess.Get("user_id"); user_id != nil {
		sess.Delete("user_id")
		sess.Delete("nickname")
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
			if ret, user := model.CheckUser(db, r.Form.Get("email"), r.Form.Get("password")); ret {
				sess.Set("user_id", user.ID)
				sess.Set("nickname", user.Nickname)
				http.Redirect(w, r, "/", 302)
			} else {
				t, err := template.ParseFiles("tmpl/auth.html")
				if err != nil {
					log.Fatal(err)
				}
				t.Execute(w, "用户不存在！")
			}
		} else if r.Form.Get("smt") == "注册" {
			if r.Form.Get("password") != r.Form.Get("password_confirm") {
				t, err := template.ParseFiles("tmpl/auth.html")
				if err != nil {
					log.Fatal(err)
				}
				t.Execute(w, "两次输入的密码不一致！")
			}
			model.AddUser(db, r.Form.Get("nickname"), r.Form.Get("email"), r.Form.Get("password"))
			t, err := template.ParseFiles("tmpl/auth.html")
			if err != nil {
				log.Fatal(err)
			}
			t.Execute(w, "注册成功！")
		}
	}
}

func list(w http.ResponseWriter, r *http.Request) {
	sess, _ := globalSessions.SessionStart(w, r)
	defer sess.SessionRelease(w)
	if user_id, ok := sess.Get("user_id").(uint); ok {
		if r.Method == "GET" {
			db_websites := model.ListUserWebsites(db, user_id)
			t, err := template.ParseFiles("tmpl/list.html")
			if err != nil {
				log.Fatal(err)
			}
			t.Execute(w, db_websites)
		} else if r.Method == "POST" {
			r.ParseForm()
			id, _ := strconv.ParseUint(r.Form.Get("id"), 10, 32)
			if ok, website := model.DelWebsite(db, uint(id) , user_id); ok {
				log.Infof("Handling request of deleting reverse proxy [" + website.Name + "] by user " + fmt.Sprint(user_id))
				delWebsite(website.Name)
			}
			http.Redirect(w, r, "/list", 302)
		}
	} else {
		http.Redirect(w, r, "/auth", 302)
	}
}

func StartManager() {
	globCfg := config.GlobCfg

	go func() {
		mux := http.NewServeMux()
		mux.HandleFunc("/", index)
		mux.HandleFunc("/auth", auth)
		mux.HandleFunc("/list", list)
		log.Infof("Starting manager on port " + globCfg.MANAGE_PORT)
		err := http.ListenAndServe(":"+globCfg.MANAGE_PORT, mux)
		if err != nil {
			log.Fatalln("ListenAndServe: ", err)
		}
	}()
}
