// logout endpoint
package v1

import (
	"fmt"
	"github.com/Atluss/TestTaskElma/lib"
	"github.com/Atluss/TestTaskElma/lib/api"
	"github.com/Atluss/TestTaskElma/lib/auth"
	"github.com/Atluss/TestTaskElma/lib/config"
	"net/http"
)

func V1Logout(set *config.Setup) error {

	login := &v1Logout{
		Url: fmt.Sprintf("/%s/logout", api.V1Api),
	}
	set.Route.HandleFunc(login.Url, login.Request)
	return nil
}

type v1Logout struct {
	api.RequestRest
	Url    string
	Secure bool
}

func (obj *v1Logout) Request(w http.ResponseWriter, r *http.Request) {

	api.SetDefaultHeadersV1API(w)
	auth.SaveSessionLogout(auth.GetSession(r), w, r)
	rep := api.ReplayStatus{
		Status: http.StatusOK,
	}
	lib.LogOnError(rep.Encode(w), "warning")
}
