// Package webServer server for web pages
package webServer

import (
	"github.com/Atluss/TestTaskElma/pkg/v1/api"
	"github.com/Atluss/TestTaskElma/pkg/v1/auth"
	"github.com/Atluss/TestTaskElma/pkg/v1/config"
	"net/http"
)

// AddPage add static page to web server
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
	Page   string // path to static html page
	Secure bool   // if it true allow to this page only admins
}

// Request setup new page
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
