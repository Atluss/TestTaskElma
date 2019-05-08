// Package restApi login endpoint
package restApi

import (
	"fmt"
	"github.com/Atluss/TestTaskElma/pkg/v1"
	"github.com/Atluss/TestTaskElma/pkg/v1/api"
	"github.com/Atluss/TestTaskElma/pkg/v1/auth"
	"github.com/Atluss/TestTaskElma/pkg/v1/config"
	"net/http"
)

var login = "admin"
var pass = "admin"

type requestV1Login struct {
	Login string `json:"Login"`
	Pass  string `json:"Pass"`
}

// V1Login new login logic for API
func V1Login(set *config.Setup, secure bool) error {
	login := &v1Login{
		Url:    fmt.Sprintf("/%s/login", api.V1Api),
		Secure: secure,
	}
	set.Route.HandleFunc(login.Url, login.Request)
	return nil
}

type v1Login struct {
	api.RequestRest
	Url    string
	Secure bool
}

// Request page for login
func (obj *v1Login) Request(w http.ResponseWriter, r *http.Request) {
	api.SetDefaultHeadersV1API(w)
	rep := api.ReplayStatus{
		Status: http.StatusOK,
	}

	session := auth.GetSession(r)
	if obj.Secure && !auth.CheckAuth(session) {
		rep.Status = http.StatusForbidden
		rep.Description = http.StatusText(http.StatusForbidden)
		v1.LogOnError(rep.Encode(w), "warning")
		return
	}

	req := requestV1Login{
		Login: r.FormValue("login"),
		Pass:  r.FormValue("pass"),
	}
	if req.Login == "" || req.Pass == "" {
		rep.Status = http.StatusBadRequest
		rep.Description = http.StatusText(http.StatusBadRequest)
		v1.LogOnError(rep.Encode(w), "warning")
		return
	}

	if req.Login != login || req.Pass != pass {
		rep.Status = http.StatusForbidden
		rep.Description = http.StatusText(http.StatusForbidden)
		v1.LogOnError(rep.Encode(w), "warning")
		return
	} else {
		auth.SaveSessionLogin(session, w, r)
		v1.LogOnError(rep.Encode(w), "warning")
	}
}
