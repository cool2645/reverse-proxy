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

func logHttp(r *http.Request, website string) {
	log.Infof("[%s] Request From %s, Request URI %s, Header %+v\n", website, r.Header.Get("Origin"), r.RequestURI, r.Header)
}

func index(w http.ResponseWriter, r *http.Request) {
	logHttp(r, "[manager]")
	sess, _ := globalSessions.SessionStart(w, r)
	defer sess.SessionRelease(w)
	if user_id, ok := sess.Get("user_id").(uint); ok {
		if r.Method == "GET" {
			t, err := template.ParseFiles("tmpl/index.html")
			if err != nil {
				log.Fatal(err)
			}
			m := make(map[string]string)
			m["nickname"] = sess.Get("nickname").(string)
			t.Execute(w, m)
		} else if r.Method == "POST" {
			r.ParseForm()
			log.Warnf("Handling request of adding reverse proxy [" + r.Form.Get("name") + "] for port " + r.Form.Get("port") + " of " + r.Form.Get("host") +  " by user " + fmt.Sprint(user_id))
			if ok, _ = model.AddWebsite(db, r.Form.Get("name"), r.Form.Get("host"), r.Form.Get("port"), uint(user_id)); ok {
				addWebsite(r.Form.Get("name"), handle{host:r.Form.Get("host"), port:r.Form.Get("port")})
			} else {
				t, err := template.ParseFiles("tmpl/index.html")
				if err != nil {
					log.Fatal(err)
				}
				m := make(map[string]string)
				m["nickname"] = sess.Get("nickname").(string)
				m["error"] = "此域名已存在"
				t.Execute(w, m)
			}
			http.Redirect(w, r, "/list", 302)
		}
	} else {
		http.Redirect(w, r, "/auth", 302)
	}
}

func auth(w http.ResponseWriter, r *http.Request) {
	logHttp(r, "[manager]")
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
			if ok, _ := model.AddUser(db, r.Form.Get("nickname"), r.Form.Get("email"), r.Form.Get("password")); ok {
				t, err := template.ParseFiles("tmpl/auth.html")
				if err != nil {
					log.Fatal(err)
				}
				t.Execute(w, "注册成功！")
			} else {
				t, err := template.ParseFiles("tmpl/auth.html")
				if err != nil {
					log.Fatal(err)
				}
				t.Execute(w, "该邮箱已被注册！")
			}

		}
	}
}

func list(w http.ResponseWriter, r *http.Request) {
	logHttp(r, "[manager]")
	sess, _ := globalSessions.SessionStart(w, r)
	defer sess.SessionRelease(w)
	if user_id, ok := sess.Get("user_id").(uint); ok {
		if r.Method == "GET" {
			db_websites := model.ListUserWebsites(db, user_id)
			t, err := template.ParseFiles("tmpl/list.html")
			if err != nil {
				log.Fatal(err)
			}
			m := make(map[string]interface{})
			m["websites"] = db_websites
			m["addr"] = config.GlobCfg.DOMAIN
			t.Execute(w, m)
		} else if r.Method == "POST" {
			r.ParseForm()
			id, _ := strconv.ParseUint(r.Form.Get("id"), 10, 32)
			if ok, website := model.DelWebsite(db, uint(id) , user_id); ok {
				log.Warnf("Handling request of deleting reverse proxy [" + website.Name + "] by user " + fmt.Sprint(user_id))
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
		log.Warnf("Starting manager on port " + globCfg.MANAGE_PORT)
		err := http.ListenAndServe(":"+globCfg.MANAGE_PORT, mux)
		if err != nil {
			log.Fatalln("ListenAndServe: ", err)
		}
	}()
}
