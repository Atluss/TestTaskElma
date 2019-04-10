package web_serve

import (
	"github.com/Atluss/TestTaskElma/lib/config"
	"net/http"
)

func AddPage(page, url string, secure bool, set *config.Setup) {

	login := &pageSt{
		Setup:  set,
		Url:    url,
		Page:   page,
		Secure: secure,
	}

	set.Route.HandleFunc(login.Url, login.Request)
}

type pageSt struct {
	RequestHttp
	Setup  *config.Setup
	Url    string
	Page   string
	Secure bool
}

func (obj *pageSt) Request(w http.ResponseWriter, r *http.Request) {

	setDefaultHeadersHttp(w)

	if obj.Secure {
		session, _ := store.Get(r, sessionName)
		if auth, ok := session.Values["authenticated"].(bool); !ok || !auth {
			http.ServeFile(w, r, "http.files/login.html")
			return
		}
	}

	http.ServeFile(w, r, obj.Page)
}
