package web_server

import (
	"github.com/Atluss/TestTaskElma/lib/api"
	"github.com/Atluss/TestTaskElma/lib/auth"
	"github.com/Atluss/TestTaskElma/lib/config"
	"net/http"
)

// AddPage add static page to webserver
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
	api.HeadRequest
	Setup  *config.Setup
	Url    string
	Page   string
	Secure bool
}

func (obj *pageSt) Request(w http.ResponseWriter, r *http.Request) {

	api.SetDefaultHeadersHttp(w)

	if obj.Secure {
		if !auth.CheckAuth(auth.GetSession(r)) {
			http.ServeFile(w, r, "http_files/login.html")
			return
		}
	}

	http.ServeFile(w, r, obj.Page)
}
