// login endpoint
package v1

import (
	"fmt"
	"github.com/Atluss/TestTaskElma/lib"
	"github.com/Atluss/TestTaskElma/lib/api"
	"github.com/Atluss/TestTaskElma/lib/auth"
	"github.com/Atluss/TestTaskElma/lib/config"
	"net/http"
)

var login = "admin"
var pass = "admin"

type requestV1Login struct {
	Login string `json:"Login"`
	Pass  string `json:"Pass"`
}

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

func (obj *v1Login) Request(w http.ResponseWriter, r *http.Request) {

	api.SetDefaultHeadersV1API(w)
	rep := api.ReplayStatus{
		Status: http.StatusOK,
	}

	session := auth.GetSession(r)

	if obj.Secure && !auth.CheckAuth(session) {
		rep.Status = http.StatusForbidden
		rep.Description = http.StatusText(http.StatusForbidden)
		lib.LogOnError(rep.Encode(w), "warning")
		return
	}

	req := requestV1Login{
		Login: r.FormValue("login"),
		Pass:  r.FormValue("pass"),
	}
	if req.Login == "" || req.Pass == "" {
		rep.Status = http.StatusBadRequest
		rep.Description = http.StatusText(http.StatusBadRequest)
		lib.LogOnError(rep.Encode(w), "warning")
		return
	}

	if req.Login != login || req.Pass != pass {
		rep.Status = http.StatusForbidden
		rep.Description = http.StatusText(http.StatusForbidden)
		lib.LogOnError(rep.Encode(w), "warning")
		return
	} else {
		auth.SaveSessionLogin(session, w, r)
		lib.LogOnError(rep.Encode(w), "warning")
	}
}
