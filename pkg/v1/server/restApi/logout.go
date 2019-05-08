// Package restApi logout endpoint
package restApi

import (
	"fmt"
	v1 "github.com/Atluss/TestTaskElma/pkg/v1"
	"github.com/Atluss/TestTaskElma/pkg/v1/api"
	"github.com/Atluss/TestTaskElma/pkg/v1/auth"
	"github.com/Atluss/TestTaskElma/pkg/v1/config"
	"net/http"
)

// V1Logout logout logic
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

// Request page for logout
func (obj *v1Logout) Request(w http.ResponseWriter, r *http.Request) {
	api.SetDefaultHeadersV1API(w)
	auth.SaveSessionLogout(auth.GetSession(r), w, r)
	rep := api.ReplayStatus{
		Status: http.StatusOK,
	}
	v1.LogOnError(rep.Encode(w), "warning")
}
