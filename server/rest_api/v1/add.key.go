package v1

import (
	"encoding/json"
	"fmt"
	"github.com/Atluss/TestTaskElma/lib"
	"github.com/Atluss/TestTaskElma/lib/api"
	"github.com/Atluss/TestTaskElma/lib/config"
	"github.com/Atluss/TestTaskElma/lib/data"
	"log"
	"net/http"
)

func V1AddKey(set *config.Setup) error {
	login := &v1AddKey{
		Url: fmt.Sprintf("/%s/addkey", api.V1Api),
		WbS: fmt.Sprintf("ws://%s:%s/st_cpu", set.Config.Address, set.Config.Port),
	}
	set.Route.HandleFunc(login.Url, login.Request)

	return nil
}

type v1AKreplay struct {
	Status      int
	Description string
	ServerWSAdr string
}

func (obj *v1AKreplay) Encode(w http.ResponseWriter) error {
	err := json.NewEncoder(w).Encode(&obj)
	if !lib.LogOnError(err, "error: can't encode ReplayStatus") {
		return err
	}
	return nil
}

type v1AddKey struct {
	api.RequestRest
	WbS string
	Url string
}

func (obj *v1AddKey) Request(w http.ResponseWriter, r *http.Request) {
	api.SetDefaultHeadersV1API(w)
	rep := v1AKreplay{
		Status: http.StatusOK,
	}

	key := r.FormValue("key")
	if key == "" {
		rep.Status = http.StatusBadRequest
		rep.Description = http.StatusText(http.StatusBadRequest)
		lib.LogOnError(rep.Encode(w), "warning")
		return
	}

	dKey := data.Keys{
		Key:    key,
		Status: 0,
	}

	log.Printf("%+v", dKey)

	rep.ServerWSAdr = obj.WbS
	lib.LogOnError(rep.Encode(w), "warning")
}
